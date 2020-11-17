package resources

type Song struct {
	Resource
	Artist string `json:"artist"`
	Title  string `json:"title"`
}
