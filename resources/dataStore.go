package resources

type DataStore interface {
	AddUser(user *User) error
	AddSong(song *Song) error
	AddPlayList(playlist *PlayList) error
	RemovePlayList(playlistID string) error
	AddSongToPlayList(playlistID, songID string) error
	ForeachUser(callback func(user *User) error) error
	ForeachSong(callback func(song *Song) error) error
	ForeachPlayList(callback func(playList *PlayList) error) error
}

func NewDataStore() DataStore {
	return &Storage{
		users:     make(map[string]*User),
		songs:     make(map[string]*Song),
		playLists: make(map[string]*PlayList),
	}
}
