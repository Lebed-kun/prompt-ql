package hellocmd

type THelloCmdResponse struct {
	MyModels map[string]bool `json:"myModels"`
	MyVariables map[string]bool `json:"myVariables"`
}
