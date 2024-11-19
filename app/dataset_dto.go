package app

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/utils"
)

func (s datasetService) toDatasetDTO(d *domain.Dataset, dto *DatasetDTO) {
	*dto = DatasetDTO{
		Id:            d.Id,
		Owner:         d.Owner.Account(),
		Name:          d.Name.ResourceName(),
		Protocol:      d.Protocol.ProtocolName(),
		RepoType:      d.RepoType.RepoType(),
		RepoId:        d.RepoId,
		Tags:          d.Tags,
		CreatedAt:     utils.ToDate(d.CreatedAt),
		UpdatedAt:     utils.ToDate(d.UpdatedAt),
		LikeCount:     d.LikeCount,
		DownloadCount: d.DownloadCount,
	}

	if d.Desc != nil {
		dto.Desc = d.Desc.ResourceDesc()
	}

	if d.Title != nil {
		dto.Title = d.Title.ResourceTitle()
	}
}

func (s datasetService) toDatasetSummaryDTO(d *domain.DatasetSummary, dto *DatasetSummaryDTO) {
	*dto = DatasetSummaryDTO{
		Id:            d.Id,
		Owner:         d.Owner.Account(),
		Name:          d.Name.ResourceName(),
		Tags:          d.Tags,
		UpdatedAt:     utils.ToDate(d.UpdatedAt),
		LikeCount:     d.LikeCount,
		DownloadCount: d.DownloadCount,
	}

	if d.Desc != nil {
		dto.Desc = d.Desc.ResourceDesc()
	}

	if d.Title != nil {
		dto.Title = d.Title.ResourceTitle()
	}

}
