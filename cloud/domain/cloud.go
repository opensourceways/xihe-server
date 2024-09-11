package domain

import (
	"errors"
	"fmt"
)

const (
	cloudIdCPU   = "cpu_001"
	cloudIdNPU   = "ascend_001"
	cloudTypeCPU = "cpu"
	cloudTypeNPU = "npu"
)

type CloudConf struct {
	Id            string
	Name          CloudName
	Specs         []CloudSpec
	Images        []CloudImage
	Feature       CloudFeature
	Processor     CloudProcessor
	SingleLimited CloudLimited
	MultiLimited  CloudLimited
	Credit        Credit
}

type CloudImage struct {
	Alias CloudImageAlias
	Image ICloudImage
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

	SingleRemain CloudRemain
	MultiRemain  CloudRemain
}

func (c *Cloud) HasSingleCardIdle() bool {
	return c.SingleRemain.CloudRemain() > 0
}

func (c *Cloud) HasMultiCardsIdle(deduction int) bool {
	return c.MultiRemain.CloudRemain()-deduction >= 0
}

func (c *CloudConf) GetImage(alias string) (ICloudImage, error) {
	for i := range c.Images {
		if alias == c.Images[i].Alias.CloudImageAlias() {
			return c.Images[i].Image, nil
		}
	}

	return nil, fmt.Errorf("%s doesn't exist", alias)
}

func (c *CloudConf) GetSpecDesc(cardsNum int) (CloudSpecDesc, error) {
	for _, spec := range c.Specs {
		if spec.CardsNum.CloudSpecCardsNum() == cardsNum {
			return spec.Desc, nil
		}
	}

	return nil, errors.New("invalid cards number")
}
