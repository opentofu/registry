package v1api

// Platform represents a platform that a provider supports.
type Platform struct {
	OS   string `json:"os"`
	Arch string `json:"arch"`
}
