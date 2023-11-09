package hellocmd

type THelloCmdResponse struct {
	MyModels map[string]string `json:"myModels"`
	MyVariables map[string]string `json:"myVariables"`
	MyEmbeddings map[string]string `json:"myEmbeddings"`
}
