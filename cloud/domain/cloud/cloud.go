package cloud

type CloudPodCreateInfo struct {
	PodId        string
	SurvivalTime int64
	User         string
	CloudType    string
}

type CloudPod interface {
	Create(*CloudPodCreateInfo) error
	Release(podId, cloudType string) error
}
