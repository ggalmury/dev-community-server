package dto

type PartyArticleCreateDto struct {
	Poster      string         `json:"poster"`
	Category    string         `json:"category"`
	Title       string         `json:"title"`
	Description *string        `json:"description"`
	TechSkill   []string       `json:"techSkill"`
	Position    map[string]int `json:"position"`
	Process     string         `json:"process"`
	Location    *string        `json:"location"`
	Deadline    string         `json:"deadline"`
	StartDate   string         `json:"startDate"`
	Span        string         `json:"span"`
}
