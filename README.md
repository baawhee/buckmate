![buckmate-logo](docs/docs/assets/logo.png)

Visit [buckmate.org](https://buckmate.org) for more information!

## buckmate - made primarily to deploy static websites to AWS S3, however it can be used to:
* transfer files between **buckets**,
* transfer files between **servers** and **buckets**,
* **replace content** in transfered files according to **yaml configuration**.

1. Define your configuration using `yaml` files
2. Configure your environment with AWS credentials and region details
3. Run `buckmate apply` to swap placeholders and upload your files to desired location

## Testing

Currently only testing in this project is a funny shell test suite. To run them you will need your own AWS buckets (e2e).

> Note this will override / delete content in your buckets!

1. Create two buckets in AWS
2. Replace bucket names placeholders: `cd e2e && ./replace-buckets.sh your_bucket_name_1 your_bucket_name_2`
3. Run test suite: `./e2e.sh`

## Contributing / Code of conduct

Feel free to open issues, or PRs. Be nice.