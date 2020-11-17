package resources

type PlayList struct {
	Resource
	UserID  string   `json:"user_id"`
	SongIDs []string `json:"song_ids"`
}
