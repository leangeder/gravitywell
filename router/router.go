package router

type RouterPath struct {
	verb string
	kind string
}

func Run(verb string, generalConfig *configuration.GeneralConfig) {

	routerPath := &RouterPath {verb: verb, kind: generalConfig.Kind }

	switch *routerPath {
	case RouterPath{ verb: "apply", kind: "cluster"}:
		fmt.Println("Apply cluster")
	case RouterPath{verb: "apply", kind; "application"}:
		fmt.Println("Apply Application")
	default:
		fmt.Println("Route not recognize")
	}
}
