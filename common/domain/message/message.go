package message

type MsgNormal struct {
	Type      string            `json:"type"`
	User      string            `json:"user"`
	Desc      string            `json:"desc"`
	Details   map[string]string `json:"details"`
	CreatedAt int64             `json:"created_at"`
}
