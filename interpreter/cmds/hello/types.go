package hellocmd

type THelloCmdResponse struct {
	Models map[string]bool `json:"models"`
	Variables map[string]bool `json:"variables"`
}
