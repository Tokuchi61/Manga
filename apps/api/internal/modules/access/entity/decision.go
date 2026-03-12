package entity

type Decision struct {
	Allowed       bool
	Effect        string
	ReasonCode    string
	Reason        string
	Permission    string
	PolicyVersion int
	SubjectKind   string
}
