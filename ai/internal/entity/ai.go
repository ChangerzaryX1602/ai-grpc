package entity

type AiRequest struct {
	Question string
	Image    []byte
}
type AiAnswer struct {
	Answer string
	Image  []byte
}
type Ai struct {
	Keys   []string
	Models []string
}
