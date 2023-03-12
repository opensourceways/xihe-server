package cloud

type CloudPodCreateInfo struct {
	PodId        string
	SurvivalTime int64
}

type CloudPod interface {
	Create(*CloudPodCreateInfo) error
}