package repositoryimpl

// import (
// 	// 	"errors"
// 	// 	"fmt"

// 	// 	// "github.com/opensourceways/xihe-server/domain"
// 	"github.com/opensourceways/xihe-server/domain"
// 	spacedomain "github.com/opensourceways/xihe-server/space/domain"
// 	// 	"gorm.io/gorm/clause"
// )

// // type datasetAdapter struct {
// // 	relatedDaoImpl
// // }

// func (adapter *datasetAdapter) Getdataset(p *spacedomain.Project, datasetResult []datasetDO) {
// 	if len(datasetResult) == 0 {
// 		return
// 	}
// 	v.err := adapter.db().ClauseBuilders

// 	relatedDatasets := make(domain.RelatedResources, len(datasetResult))

// 	for i, dataset := range datasetResult {
// 		relatedDatasets[i] = domain.ResourceIndex{
// 			Id: dataset.ProjectId,
// 			Owner: spacedomain.Account{
// 				Name: dataset.owner,
// 			},
// 		}
// 	}

// 	p.RelatedDatasets = relatedDatasets
// }
