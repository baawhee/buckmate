package aws

import (
	"buckmate/main/deploymentConfig"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const InternalBuckmateFilePrefix string = "//buckmate//internal"

type fileWalk chan string

type Bucket struct {
	bucket     string
	uploader   *manager.Uploader
	downloader *manager.Downloader
	paginator  *s3.ListObjectsV2Paginator
	remover    manager.DeleteObjectsAPIClient
	header     s3.HeadObjectAPIClient
}

type UploadOptions struct {
	Prefix      string
	FileOptions map[string]deploymentConfig.FileOptions
	TempDir     string
}

func Init() (c *s3.Client, e error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	client := s3.NewFromConfig(cfg)
	return client, nil
}

func NewBucket(client *s3.Client, location deploymentConfig.Location) (bucket Bucket) {
	uploader := manager.NewUploader(client)
	downloader := manager.NewDownloader(client)
	paginator := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
		Bucket: &location.Address,
	})
	header := s3.HeadObjectAPIClient(client)
	remover := manager.DeleteObjectsAPIClient(client)
	newBucket := Bucket{
		bucket:     location.Address,
		uploader:   uploader,
		downloader: downloader,
		paginator:  paginator,
		header:     header,
		remover:    remover,
	}
	return newBucket
}

func (bucket Bucket) Upload(context context.Context, options UploadOptions) (err error) {
	walker := make(fileWalk)

	go func() {
		if err := filepath.Walk(options.TempDir, walker.Walk); err != nil {
			log.Fatal(err)
		}
		close(walker)
	}()

	for path := range walker {
		rel, err := filepath.Rel(options.TempDir, path)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}

		fileKey := aws.String(filepath.Join(options.Prefix, rel))

		metadataForFile := map[string]string{}
		for k, v := range options.FileOptions[InternalBuckmateFilePrefix].Metadata {
			metadataForFile[k] = v
		}
		for k, v := range options.FileOptions["*"].Metadata {
			metadataForFile[k] = v
		}
		if options.FileOptions[*fileKey].Metadata != nil {
			for k, v := range options.FileOptions[*fileKey].Metadata {
				metadataForFile[k] = v
			}
		}
		result, err := bucket.uploader.Upload(context, &s3.PutObjectInput{
			Bucket:       &bucket.bucket,
			Key:          fileKey,
			Body:         file,
			Metadata:     metadataForFile,
			CacheControl: aws.String(options.FileOptions[*fileKey].CacheControl),
		})

		if err != nil {
			return err
		}

		err = file.Close()

		if err != nil {
			return err
		}
		fmt.Println("Successfully uploaded " + *result.Key)
	}
	if err := os.RemoveAll(options.TempDir); err != nil {
		return err
	}
	return nil
}

type RemoveOptions struct {
	CurrentVersion string
}

func (bucket Bucket) RemovePreviousVersion(context context.Context, options RemoveOptions) (err error) {
	var objectsToRemove []types.ObjectIdentifier
	for bucket.paginator.HasMorePages() {
		page, err := bucket.paginator.NextPage(context)
		if err != nil {
			return err
		}
		for _, obj := range page.Contents {
			if obj.Size > 0 {
				header, err := bucket.header.HeadObject(context, &s3.HeadObjectInput{
					Bucket: &bucket.bucket,
					Key:    obj.Key,
				})
				if err != nil {
					return err
				}

				if header.Metadata[deploymentConfig.InternalBuckmateVersionMetadataKey] != options.CurrentVersion {
					objectsToRemove = append(objectsToRemove, types.ObjectIdentifier{Key: obj.Key})
				}
			}
		}
	}
	if len(objectsToRemove) > 0 {
		_, err := bucket.remover.DeleteObjects(context, &s3.DeleteObjectsInput{
			Bucket: &bucket.bucket,
			Delete: &types.Delete{
				Objects: objectsToRemove,
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (f fileWalk) Walk(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if !info.IsDir() {
		f <- path
	}
	return nil
}

type DownloadOptions struct {
	Prefix  string
	TempDir string
}

func (bucket Bucket) Download(context context.Context, options DownloadOptions) (err error) {
	for bucket.paginator.HasMorePages() {
		page, err := bucket.paginator.NextPage(context)
		if err != nil {
			return err
		}
		for _, obj := range page.Contents {
			if obj.Size > 0 {
				err := downloadToFile(*bucket.downloader, options.TempDir, bucket.bucket, aws.ToString(obj.Key), options.Prefix)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func downloadToFile(downloader manager.Downloader, targetDirectory, bucket, key string, prefix string) (err error) {
	file := filepath.Clean(strings.Replace(filepath.Join(targetDirectory, key), prefix, "", 1))

	if err := os.MkdirAll(filepath.Dir(file), 0775); err != nil {
		return err
	}

	fd, err := os.Create(file)
	if err != nil {
		return err
	}

	defer func() {
		if cerr := fd.Close(); cerr != nil {
			err = cerr
		}
	}()

	if err != nil {
		return err
	}

	if _, err := downloader.Download(context.TODO(), fd, &s3.GetObjectInput{Bucket: &bucket, Key: &key}); err != nil {
		return err
	}
	return nil
}
