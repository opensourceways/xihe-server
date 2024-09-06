package domain

import (
	"fmt"
	"sort"
)

const (
	cloudIdCPU   = "cpu_001"
	cloudIdNPU   = "ascend_001"
	cloudTypeCPU = "cpu"
	cloudTypeNPU = "npu"
)

type CloudConf struct {
	Id        string
	Name      CloudName
	Specs     []CloudSpec
	Images    []CloudImage
	Feature   CloudFeature
	Processor CloudProcessor
	Limited   CloudLimited
	Credit    Credit
}

type CloudImage struct {
	Alias   CloudImageAlias
	Image   ICloudImage
	Default bool
}

type CloudSpec struct {
	Desc     CloudSpecDesc
	CardsNum CloudSpecCardsNum
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

func (c *CloudConf) GetImage(alias string) (ICloudImage, error) {
	for i := range c.Images {
		if alias == c.Images[i].Alias.CloudImageAlias() {
			return c.Images[i].Image, nil
		}
	}

	return nil, fmt.Errorf("%s doesn't exist", alias)
}

func (c *CloudConf) MoveDefaultImageToHead() {
	sort.Slice(c.Images, func(i, j int) bool {
		return c.Images[i].Default
	})
}
