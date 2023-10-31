package repositoryimpl

type DCloudConf struct {
	Id        string `bson:"id"        json:"id"`
	Name      string `bson:"name"      json:"name"`
	Spec      string `bson:"spec"      json:"spec"`
	Image     string `bson:"image"     json:"image"`
	Feature   string `bson:"feature"   json:"feature"`
	Processor string `bson:"processor" json:"processor"`
	Limited   int    `bson:"limited"   json:"limited"`
	Credit    int64  `bson:"credit"    json:"credit"`
}
