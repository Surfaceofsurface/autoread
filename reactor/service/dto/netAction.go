package dto

type LoginReq struct {
	Username string
	Password string
	Pinst    string
}
type LoginRes struct {
	Code       string `json:"code"`
	Text       string `json:"text"`
	Token      string `json:"token"`
	Expiration int    `json:"expiration"`
}
type PendingBook struct {
	Title  string `json:"title"`
	BookID string `json:"bookId"`
}
type GBookPendingRes struct {
	Data []PendingBook `json:"data"`
}
type InData struct {
	ParaIndex int      `json:"paraIndex"`
	Page      int      `json:"page"`
	Children  []InData `json:"children"`
}
type CataLogDTO struct {
	Data      []InData `json:"data"`
	Paragraph int      `json:"paragraph"`
}
type StateDTO struct {
	Time int64 `json:"time"`
}
type ReadGetDTO struct {
	Content   string `json:"content"`
	Page      int    `json:"page"`
	TotalPage int    `json:"totalPage"`
	TrialPage int    `json:"trialPage"`
	LogId     string `json:"logId"`
	Token     string `json:"token"`
}

// impl interface Mark to use Diff()
func (b PendingBook) Mark() string { return b.BookID }
