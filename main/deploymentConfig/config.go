package deploymentConfig

import (
	"buckmate/main/common/util"
)

const InternalBuckmateFilePrefix string = "//buckmate//internal"
const InternalBuckmateVersionMetadataKey string = "buckmate-version"

type Location struct {
	Address string `yaml:"address"`
	Prefix  string `yaml:"prefix"`
}

type FileOptions struct {
	Metadata     map[string]string `yaml:"metadata"`
	CacheControl string            `yaml:"cacheControl"`
}

type Deployment struct {
	Source         Location               `yaml:"source"`
	Target         Location               `yaml:"target"`
	ConfigBoundary string                 `yaml:"configBoundary"`
	KeepPrevious   bool                   `yaml:"keepPrevious"`
	ConfigMap      map[string]string      `yaml:"configMap"`
	FileOptions    map[string]FileOptions `yaml:"fileOptions"`
}

func Load(env string, rootDir string) (Deployment, error) {
	commonPath := rootDir + "/Deployment.yaml"
	commonFile, err := util.LoadYaml(commonPath)
	if err != nil {
		return Deployment{}, err
	}

	envConfig := Deployment{}

	if len(env) > 0 {
		envPath := rootDir + "/" + env + "/Deployment.yaml"
		envFile, err := util.LoadYaml(envPath)
		if err != nil {
			return Deployment{}, err
		}
		err = util.YamlToStruct(envFile, &envConfig)
		if err != nil {
			return Deployment{}, err
		}
	}

	commonConfig := Deployment{}
	commonConfig.ConfigBoundary = "%%%"
	commonConfig.KeepPrevious = false

	err = util.YamlToStruct(commonFile, &commonConfig)
	if err != nil {
		return Deployment{}, err
	}

	err = util.MergeStruct(&commonConfig, envConfig)
	if err != nil {
		return Deployment{}, err
	}

	return commonConfig, nil
}
