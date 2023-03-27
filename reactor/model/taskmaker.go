package model

type TaskMaker struct {
	Username string
	Password string

	SchoolName string `json:"-"`
	SchoolID   string `json:"Pinst"`
	// TaskList   *Tasks `json:"-"`
	PendingBook []Book
}
