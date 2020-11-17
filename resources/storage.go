package resources

import (
	"errors"
	"fmt"
	"strconv"
)

// An in-memory key/value store
type Storage struct {
	users     map[string]*User
	songs     map[string]*Song
	playLists map[string]*PlayList
}

func (self *Storage) AddUser(user *User) error {
	_, err := strconv.ParseUint(user.ID, 10, 32)
	if err != nil {
		return errors.New(fmt.Sprintf("User ID %v is invalid.", user.ID))
	}

	if _, ok := self.users[user.ID]; ok {
		return errors.New(fmt.Sprintf("Duplicate user ID %v.", user.ID))
	}

	self.users[user.ID] = user

	return nil
}

func (self *Storage) AddSong(song *Song) error {
	_, err := strconv.ParseUint(song.ID, 10, 32)
	if err != nil {
		return errors.New(fmt.Sprintf("Song ID %v is invalid.", song.ID))
	}

	if _, ok := self.songs[song.ID]; ok {
		return errors.New(fmt.Sprintf("Duplicate song ID %v.", song.ID))
	}

	self.songs[song.ID] = song

	return nil
}

func (self *Storage) AddPlayList(playlist *PlayList) error {
	_, err := strconv.ParseUint(playlist.ID, 10, 32)
	if err != nil {
		return errors.New(fmt.Sprintf("Playlist ID %v is invalid.", playlist.ID))
	}

	if _, ok := self.playLists[playlist.ID]; ok {
		return errors.New(fmt.Sprintf("Duplicate playlist ID %v.", playlist.ID))
	}

	_, err = strconv.ParseUint(playlist.UserID, 10, 32)
	if err != nil {
		return errors.New(fmt.Sprintf("User ID %v is invalid.", playlist.UserID))
	}

	if _, ok := self.users[playlist.UserID]; !ok {
		return errors.New(fmt.Sprintf("The user ID %v does not exist.", playlist.UserID))
	}

	for _, songID := range playlist.SongIDs {
		_, err := strconv.ParseUint(songID, 10, 32)
		if err != nil {
			return errors.New(fmt.Sprintf("Song ID %v is invalid.", songID))
		}

		if _, ok := self.songs[songID]; !ok {
			return errors.New(fmt.Sprintf("The song ID %v does not exist.", songID))
		}
	}

	self.playLists[playlist.ID] = playlist

	return nil
}

// Remove a playlist from the storage model
func (self *Storage) RemovePlayList(playlistID string) error {
	_, err := strconv.ParseUint(playlistID, 10, 32)
	if err != nil {
		return errors.New(fmt.Sprintf("Playlist ID %v is invalid.", playlistID))
	}

	_, ok := self.playLists[playlistID]
	if !ok {
		return errors.New(fmt.Sprintf("Playlist ID %v does not exist.", playlistID))
	}

	delete(self.playLists, playlistID)

	return nil
}

// Add a song to a playlist in the storage model
func (self *Storage) AddSongToPlayList(playlistID, songID string) error {
	_, err := strconv.ParseUint(playlistID, 10, 32)
	if err != nil {
		return errors.New(fmt.Sprintf("Playlist ID %v is invalid.", playlistID))
	}

	_, ok := self.playLists[playlistID]
	if !ok {
		return errors.New(fmt.Sprintf("Playlist ID %v does not exist.", playlistID))
	}

	_, err = strconv.ParseUint(songID, 10, 32)
	if err != nil {
		return errors.New(fmt.Sprintf("Song ID %v is invalid.", songID))
	}

	_, ok = self.songs[songID]
	if !ok {
		return errors.New(fmt.Sprintf("Song ID %v does not exist.", songID))
	}

	self.playLists[playlistID].SongIDs = append(self.playLists[playlistID].SongIDs, songID)

	return nil
}

func (self *Storage) ForeachUser(callback func(user *User) error) error {
	for _, user := range self.users {
		err := callback(user)
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *Storage) ForeachSong(callback func(song *Song) error) error {
	for _, song := range self.songs {
		err := callback(song)
		if err != nil {
			return err
		}
	}
	return nil
}

func (self *Storage) ForeachPlayList(callback func(playList *PlayList) error) error {
	for _, playList := range self.playLists {
		err := callback(playList)
		if err != nil {
			return err
		}
	}
	return nil
}
