package configuration

type Application struct {
	Template struct {
		Spec struct {
			Applications []struct {
				Application struct {
					Name            string `yaml:"Name"`
					Namespace       string `yaml:"Namespace"`
					CreateNamespace bool   `yaml:"CreateNamespace"`
					Git             string `yaml:"Git"`
					Action          []struct {
						Execute struct {
							Shell   string `yaml:"Shell"`
							Kubectl struct {
								Path    string `yaml:"Path"`
								Type    string `yaml:"Type"`
								Command string `yaml:"Command"`
							} `yaml:"Kubectl"`
						} `yaml:"Execute"`
					} `yaml:"Action"`
				} `yaml:"Deployment"`
			} `yaml:"Deployments"`
		} `yaml:"Spec"`
	} `yaml:"Template"`
}

type ApplicationConfig struct {
	Kind       string      `yaml:"Kind"`
	APIVersion string      `yaml:"APIVersion"`
	Spec       Application `yaml:"Spec"`
}
