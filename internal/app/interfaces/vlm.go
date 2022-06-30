package interfaces

type VoteLinkData struct {
	Email    string
	TaskID   int
	Result   bool
	Checksum string
}

type VoteLinkManager interface {
	Generate(data VoteLinkData) string
	Parse(string) (VoteLinkData, error)
}
