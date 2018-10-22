package configuration

type ByKind []*GeneralConfig

func (s ByKind) Len() int {
	return len(s)
}
func (s ByKind) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByKind) Less(i, j int) bool {
	return s[i].Kind == "Cluster" && s[j].Kind == "Application"
}
