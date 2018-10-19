package typeCluster

type ObjectMeta struct {
	Name   string            `yaml:"name"`
	Labels map[string]string `yaml:"labels"`
}

type TypeMeta struct {
	Kind       string `yaml:"kind"`
	APIVersion string `yaml:"apiVersion"`
}

type Cluster struct {
	TypeMeta   `yaml:","`
	ObjectMeta `yaml:"metadata`
	Spec       ClusterSpec `yaml:"spec"`
}

type ClusterSpec struct {
	Replicas *int32               `yaml:"replicas"`
	Template ProviderTemplateSpec `yaml:"template"`
}

type ProviderTemplateSpec struct {
	ObjectMeta
	Spec ProviderSpec
}

type ProviderSpec struct {
	ServiceAccountName `yaml:"serviceAccountName"`
	Providers          []Provider
}

type Provider struct {
	Name        string     `yaml:"name"`
	Project     string     `yaml:"project"`
	Version     string     `yaml:"version"`
	Image       string     `yaml:"image`
	Region      string     `yaml:"region`
	Zone        string     `yaml:"zone`
	Network     string     `yaml:"network"`
	Dashboard   string     `yaml:"dashboard"`
	Monitoring  string     `yaml:"monitoring"`
	Autoscaling string     `yaml:"autoscaling"`
	BootDisk    string     `yaml:"bootDisk"`
	LocalDisk   string     `yaml:"localDisk"`
	NodePools   []NodePool `yaml:"nodepools"`
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
	Name    string `yaml:"name"`
	Replica int32  `yaml:"replica"`
	Type    string `yaml:"type"`
	Image   string `yaml:"image"`
}
