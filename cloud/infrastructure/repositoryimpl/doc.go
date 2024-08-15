package repositoryimpl

type DCloudConf struct {
	Id        string    `bson:"id"         json:"id"`
	Name      string    `bson:"name"       json:"name"`
	Spec      string    `bson:"spec"       json:"spec"`
	Images    []ImageDO `bson:"images"     json:"images"`
	Feature   string    `bson:"feature"    json:"feature"`
	Processor string    `bson:"processor"  json:"processor"`
	Limited   int       `bson:"limited"    json:"limited"`
	Credit    int64     `bson:"credit"     json:"credit"`
}

type ImageDO struct {
	Alias   string `bson:"alias" json:"alias"`
	Image   string `bson:"image" json:"image"`
	Default bool   `bson:"default" json:"default"`
}
