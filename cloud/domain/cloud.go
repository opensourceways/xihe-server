package domain

const (
	cloudIdCPU   = "cpu_001"
	cloudIdNPU   = "ascend_001"
	cloudTypeCPU = "cpu"
	cloudTypeNPU = "npu"
)

type CloudConf struct {
	Id        string
	Name      CloudName
	Spec      CloudSpec
	Image     CloudImage
	Feature   CloudFeature
	Processor CloudProcessor
	Limited   CloudLimited
	Credit    Credit
}

func (c *CloudConf) IsNPU() bool {
	return c.Id == cloudIdNPU
}

type Cloud struct {
	CloudConf

	Remain CloudRemain
}

func (c *Cloud) HasIdle() bool {
	return c.Remain.CloudRemain() > 0
}
