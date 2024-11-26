package domain

// DeleteSpaceEvent
type DeleteSpaceEvent struct {
	Time      int64  `json:"time"`
	Owner     string `json:"owner"`
	SpaceId   string `json:"space_id"`
	SpaceName string `json:"space_name"`
	DeletedBy string `json:"deleted_by"`
}
