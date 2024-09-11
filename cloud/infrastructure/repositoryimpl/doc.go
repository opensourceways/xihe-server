package repositoryimpl

type DCloudConf struct {
	Id            string    `bson:"id"                json:"id"`
	Name          string    `bson:"name"              json:"name"`
	Specs         []SpecDO  `bson:"specs"             json:"specs"`
	Images        []ImageDO `bson:"images"            json:"images"`
	Feature       string    `bson:"feature"           json:"feature"`
	Processor     string    `bson:"processor"         json:"processor"`
	SingleLimited int       `bson:"single_limited"    json:"single_limited"`
	MultiLimited  int       `bson:"multi_limited"     json:"multi_limited"`
	Credit        int64     `bson:"credit"            json:"credit"`
}

type ImageDO struct {
	Alias string `bson:"alias" json:"alias"`
	Image string `bson:"image" json:"image"`
}

type SpecDO struct {
	Desc     string `bson:"desc" json:"desc"`
	CardsNum int    `bson:"cards_num" json:"cards_num"`
}
