package config

import (
	"time"

	"go.uber.org/zap/zapcore"

	"code.cloudfoundry.org/korifi/tools"
)

type ControllerConfig struct {
	// components
	IncludeKpackImageBuilder bool `yaml:"includeKpackImageBuilder"`
	IncludeJobTaskRunner     bool `yaml:"includeJobTaskRunner"`
	IncludeStatefulsetRunner bool `yaml:"includeStatefulsetRunner"`

	// core controllers
	CFProcessDefaults                CFProcessDefaults  `yaml:"cfProcessDefaults"`
	CFStagingResources               CFStagingResources `yaml:"cfStagingResources"`
	CFRootNamespace                  string             `yaml:"cfRootNamespace"`
	ContainerRegistrySecretNames     []string           `yaml:"containerRegistrySecretNames"`
	TaskTTL                          string             `yaml:"taskTTL"`
	BuilderName                      string             `yaml:"builderName"`
	RunnerName                       string             `yaml:"runnerName"`
	NamespaceLabels                  map[string]string  `yaml:"namespaceLabels"`
	ExtraVCAPApplicationValues       map[string]any     `yaml:"extraVCAPApplicationValues"`
	MaxRetainedPackagesPerApp        int                `yaml:"maxRetainedPackagesPerApp"`
	MaxRetainedBuildsPerApp          int                `yaml:"maxRetainedBuildsPerApp"`
	LogLevel                         zapcore.Level      `yaml:"logLevel"`
	SpaceFinalizerAppDeletionTimeout *int64             `yaml:"spaceFinalizerAppDeletionTimeout"`

	// job-task-runner
	JobTTL                                     string `yaml:"jobTTL"`
	JobTaskRunnerTemporarySetPodSeccompProfile bool   `yaml:"jobTaskRunnerTemporarySetPodSeccompProfile"`

	// statefulset-runner
	StatefulsetRunnerTemporarySetPodSeccompProfile bool `yaml:"statefulsetRunnerTemporarySetPodSeccompProfile"`

	// kpack-image-builder
	ClusterBuilderName        string     `yaml:"clusterBuilderName"`
	BuilderServiceAccount     string     `yaml:"builderServiceAccount"`
	BuilderReadinessTimeout   string     `yaml:"builderReadinessTimeout"`
	ContainerRepositoryPrefix string     `yaml:"containerRepositoryPrefix"`
	ContainerRegistryType     string     `yaml:"containerRegistryType"`
	Networking                Networking `yaml:"networking"`
}

type CFProcessDefaults struct {
	MemoryMB    int64  `yaml:"memoryMB"`
	DiskQuotaMB int64  `yaml:"diskQuotaMB"`
	Timeout     *int64 `yaml:"timeout"`
}

type CFStagingResources struct {
	BuildCacheMB int64 `yaml:"buildCacheMB"`
	DiskMB       int64 `yaml:"diskMB"`
	MemoryMB     int64 `yaml:"memoryMB"`
}

type Networking struct {
	GatewayName      string `yaml:"gatewayName"`
	GatewayNamespace string `yaml:"gatewayNamespace"`
}

const (
	defaultTaskTTL            = 30 * 24 * time.Hour
	defaultTimeout      int64 = 60
	defaultJobTTL             = 24 * time.Hour
	defaultBuildCacheMB       = 2048
)

func LoadFromPath(path string) (*ControllerConfig, error) {
	var config ControllerConfig
	err := tools.LoadConfigInto(&config, path)
	if err != nil {
		return nil, err
	}

	if config.CFProcessDefaults.Timeout == nil {
		config.CFProcessDefaults.Timeout = tools.PtrTo(defaultTimeout)
	}

	if config.SpaceFinalizerAppDeletionTimeout == nil {
		config.SpaceFinalizerAppDeletionTimeout = tools.PtrTo(defaultTimeout)
	}

	if config.CFStagingResources.BuildCacheMB == 0 {
		config.CFStagingResources.BuildCacheMB = defaultBuildCacheMB
	}

	return &config, nil
}

func GetLogLevelFromPath(path string) (zapcore.Level, error) {
	cfg, err := LoadFromPath(path)
	if err != nil {
		return zapcore.InfoLevel, err
	}

	return cfg.LogLevel, nil
}

func (c ControllerConfig) ParseTaskTTL() (time.Duration, error) {
	if c.TaskTTL == "" {
		return defaultTaskTTL, nil
	}

	return tools.ParseDuration(c.TaskTTL)
}

func (c ControllerConfig) ParseBuilderReadinessTimeout() (time.Duration, error) {
	return tools.ParseDuration(c.BuilderReadinessTimeout)
}

func (c ControllerConfig) ParseJobTTL() (time.Duration, error) {
	if c.JobTTL == "" {
		return defaultJobTTL, nil
	}

	return tools.ParseDuration(c.JobTTL)
}
