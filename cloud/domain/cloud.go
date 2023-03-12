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

type Cloud struct {
	CloudConf

	Remain CloudRemain
}

func (c *Cloud) HasFree() bool {
	return c.Remain.CloudRemain() > 0
}
