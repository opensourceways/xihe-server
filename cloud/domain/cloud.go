package domain

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

func (c *CloudConf) IsAscend() bool {
	return c.Id == "ascend_001"
}

type Cloud struct {
	CloudConf

	Remain CloudRemain
}

func (c *Cloud) HasIdle() bool {
	return c.Remain.CloudRemain() > 0
}
