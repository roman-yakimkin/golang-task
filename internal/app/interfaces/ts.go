package interfaces

type TokenStorage interface {
	SetParams(map[string]interface{})
	GetTokens() (string, string, error)
	SetTokens(string, string) error
}
