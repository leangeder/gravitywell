package configuration

type ObjectMeta struct {
	Name   string            `yaml:"Name"`
	Labels map[string]string `yaml:"Labels"`
}

type ClusterConfig struct {
	Kind       string      `yaml:"Kind"`
	APIVersion string      `yaml:"APIVersion"`
	Metadata   ObjectMeta  `yaml:"Metadata`
	Spec       ClusterSpec `yaml:"Spec"`
}

type ClusterSpec struct {
	Replicas *int32               `yaml:"Replicas"`
	Template ProviderTemplateSpec `yaml:"Template"`
}

type ProviderTemplateSpec struct {
	ObjectMeta
	Spec ProviderSpec
}

type ProviderSpec struct {
	ServiceAccountName string `yaml:"ServiceAccountName"`
	Providers          []Provider
}

type Provider struct {
	Name        string     `yaml:"Name"`
	Project     string     `yaml:"Project"`
	Version     string     `yaml:"Version"`
	Image       string     `yaml:"Image`
	Region      string     `yaml:"Region`
	Zone        string     `yaml:"Zone`
	Network     string     `yaml:"Network"`
	Dashboard   string     `yaml:"Dashboard"`
	Monitoring  string     `yaml:"Monitoring"`
	Autoscaling string     `yaml:"Autoscaling"`
	BootDisk    string     `yaml:"BootDisk"`
	LocalDisk   string     `yaml:"LocalDisk"`
	NodePools   []NodePool `yaml:"Nodepools"`
	// Command                []string             `yaml:"command`
	// Args                   []string             `yaml:"args`
	// WorkingDir             string               `yaml:"workingDir`
	// Ports                  []ContainerPort      `yaml:"ports`
	// EnvFrom                []EnvFromSource      `yaml:"envFrom`
	// Env                    []EnvVar             `yaml:"env`
	// Resources              ResourceRequirements `yaml:"resources`
	// VolumeMounts           []VolumeMount        `yaml:"volumeMounts`
	// VolumeDevices          []VolumeDevice       `yaml:"volumeDevices`
	// LivenessProbe          *Probe               `yaml:"livenessProbe`
	// ReadinessProbe         *Probe               `yaml:"readinessProbe`
	// Lifecycle              *Lifecycle           `yaml:"lifecycle`
	// TerminationMessagePath string               `yaml:"terminationMessagePath`
}

type NodePool struct {
	Name    string `yaml:"Name"`
	Replica int32  `yaml:"Replica"`
	Type    string `yaml:"Type"`
	Image   string `yaml:"Image"`
}
