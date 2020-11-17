package brokers

import (
	"encoding/json"
	"errors"
	"fmt"
	"highspot/data/validation"
	"highspot2/providers"
	"highspot2/resources"
	"log"
	"os"
	"regexp"
)

type Ingester struct {
	host    string // input file host
	path    string // output file path
	db      resources.DataStore
	counter int
}

func NewIngestor(host, path string) *Ingester {
	ingestor := Ingester{
		host: host,
		path: path,
		db:   resources.NewDataStore(),
	}
	return &ingestor
}

func (self *Ingester) DoIngest() error {
	err := self.ingestUsers()
	if err != nil {
		return err
	}

	err = self.ingestSongs()
	if err != nil {
		return err
	}

	err = self.ingestPlayLists()
	if err != nil {
		return err
	}

	err = self.ingestChanges()
	if err != nil {
		return err
	}

	err = self.produceOutputPlayList()
	if err != nil {
		return err
	}

	return nil
}

func (self *Ingester) ingestUsers() error {
	usersProvider := providers.NewUsersProvider(self.host + "/users.json")
	self.counter = 0
	err := usersProvider.Fetch(self.fetchUsers)
	if err != nil {
		return err
	}
	log.Printf("Ingested %v Users", self.counter)
	return nil
}

func (self *Ingester) ingestSongs() error {
	songsProvider := providers.NewSongsProvider(self.host + "/songs.json")
	self.counter = 0
	err := songsProvider.Fetch(self.fetchSongs)
	if err != nil {
		return err
	}
	log.Printf("Ingested %v Songs", self.counter)
	return nil
}

func (self *Ingester) ingestPlayLists() error {
	playListsProvider := providers.NewPlayListsProvider(self.host + "/playLists.json")
	self.counter = 0
	err := playListsProvider.Fetch(self.fetchPlayLists)
	if err != nil {
		return err
	}
	log.Printf("Ingested %v PlayLists", self.counter)
	return nil
}

func (self *Ingester) ingestChanges() error {
	changesProvider := providers.NewChangesProvider(self.host + "/changes.json")
	self.counter = 0
	err := changesProvider.Fetch(self.fetchChanges)
	if err != nil {
		return err
	}
	log.Printf("Ingested %v Changes", self.counter)
	return nil
}

func (self *Ingester) fetchUsers(resource interface{}) error {
	user, ok := resource.(*resources.User)
	if !ok {
		errors.New("Invalid cast. User expected.")
	}

	err := self.db.AddUser(user)
	if err != nil {
		return err
	}

	self.counter++

	if self.counter%10000 == 0 {
		log.Printf("Ingested %v Users", self.counter)
	}

	// log.Printf("Add User: %v %v", user.ID, user.Name)

	return nil
}

func (self *Ingester) fetchSongs(resource interface{}) error {
	song, ok := resource.(*resources.Song)
	if !ok {
		errors.New("Invalid cast. Song expected.")
	}

	err := self.db.AddSong(song)
	if err != nil {
		return err
	}

	self.counter++

	if self.counter%10000 == 0 {
		log.Printf("Ingested %v Songs", self.counter)
	}

	// log.Printf("Add Song: %v %v %v", song.ID, song.Artist, song.Title)

	return nil
}

func (self *Ingester) fetchPlayLists(resource interface{}) error {
	playList, ok := resource.(*resources.PlayList)
	if !ok {
		errors.New("Invalid cast. PlayList expected.")
	}

	err := self.db.AddPlayList(playList)
	if err != nil {
		return err
	}

	self.counter++

	if self.counter%10000 == 0 {
		log.Printf("Ingested %v PlayLists", self.counter)
	}

	// log.Printf("Add PlayList: %v %v %v", playList.ID, playList.UserID, playList.SongIDs)

	return nil
}

func (self *Ingester) fetchChanges(resource interface{}) error {
	change, ok := resource.(*resources.Change)
	if !ok {
		errors.New("Invalid cast. Change expected.")
	}

	err := self.applyChanges(change)
	if err != nil {
		return err
	}

	self.counter++

	if self.counter%10000 == 0 {
		log.Printf("Ingested %v Changes", self.counter)
	}

	// log.Printf("Add Change: %v %v %v", change.Op, change.Path, change.Value)

	return nil
}

func (self *Ingester) applyChanges(change *resources.Change) error {
	//
	// Loop over the changes and apply each change to the mixtape data model
	//
	if change.Op == "add" {
		//
		// Add a new playlist; the playlist should contain at least one song.
		//

		ok, err := regexp.MatchString("/playlists/-", change.Path)
		if err != nil {
			return err
		}
		if ok {
			err = self.addPlaylist(change)
			if err != nil {
				log.Printf("Skipping add playlist. %v", err)
			}
			return nil
		}

		//
		// Add an existing song to an existing playlist
		//

		ok, err = regexp.MatchString("/playlists/[0-9]+/song_ids/-", change.Path)
		if err != nil {
			return err
		}
		if ok {
			err = self.addSongToPlaylist(change)
			if err != nil {
				log.Printf("Skipping add song to playlist. %v", err)
			}
			return nil
		}

		log.Printf("Skipping unknown change Path %v", change.Path)
	} else if change.Op == "remove" {
		//
		// Remove a playlist.
		//

		ok, err := regexp.MatchString("/playlists/[0-9]+", change.Path)
		if err != nil {
			return err
		}
		if ok {
			err = self.removePlaylist(change)
			if err != nil {
				log.Printf("Skipping remove playlist. %v", err)
			}
			return nil
		}

		log.Printf("Skipping unknown change Path %v", change.Path)
	}

	log.Printf("Skipping unknown change Op %v", change.Op)

	return nil
}

func (self *Ingester) addPlaylist(change *resources.Change) error {
	if change.Value == nil {
		return errors.New("Missing playlist value.")
	}

	playlistJSON, err := json.Marshal(change.Value)
	if err != nil {
		return errors.New(fmt.Sprintf("Invalid playlist value. %v", err))
	}

	err = validation.Validate(validation.PatchPlaylistSchema, string(playlistJSON))
	if err != nil {
		return errors.New(fmt.Sprintf("Invalid playlist value. %v", err))
	}

	var playlist resources.PlayList
	err = json.Unmarshal(playlistJSON, &playlist)
	if err != nil {
		return errors.New(fmt.Sprintf("Invalid playlist value. %v", err))
	}

	return self.db.AddPlayList(&playlist)
}

func (self *Ingester) removePlaylist(change *resources.Change) error {
	re := regexp.MustCompile("/playlists/([0-9]+)")
	match := re.FindStringSubmatch(change.Path)
	return self.db.RemovePlayList(match[1])
}

func (self *Ingester) addSongToPlaylist(change *resources.Change) error {
	re := regexp.MustCompile("/playlists/([0-9]+)/song_ids/-")
	if change.Value == nil {
		return errors.New("Missing song ID value.")
	}
	songID, ok := change.Value.(string)
	if !ok {
		return errors.New("Invalid song ID value.")
	}
	match := re.FindStringSubmatch(change.Path)
	return self.db.AddSongToPlayList(match[1], songID)
}

func (self *Ingester) produceOutputPlayList() error {
	f, err := os.OpenFile(self.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(f)

	self.db.ForeachPlayList(func(playlist *resources.PlayList) error {
		return encoder.Encode(playlist)
	})

	return nil
}

func (self *Ingester) produceOutput() error {
	f, err := os.OpenFile(self.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(f)

	_, err = f.WriteString("{\n")
	if err != nil {
		return err
	}

	// Output users

	_, err = f.WriteString("\"users\": [\n")
	if err != nil {
		return err
	}

	i := 0
	self.db.ForeachUser(func(user *resources.User) error {
		if i > 0 {
			f.WriteString(",")
		}
		i++
		return encoder.Encode(user)
	})

	_, err = f.WriteString("],\n")
	if err != nil {
		return err
	}

	// Output playlists

	_, err = f.WriteString("\"playlists\": [\n")
	if err != nil {
		return err
	}

	i = 0
	self.db.ForeachPlayList(func(playlist *resources.PlayList) error {
		if i > 0 {
			f.WriteString(",")
		}
		i++
		return encoder.Encode(playlist)
	})

	_, err = f.WriteString("],\n")
	if err != nil {
		return err
	}

	// Output songs

	_, err = f.WriteString("\"songs\": [\n")
	if err != nil {
		return err
	}

	i = 0
	self.db.ForeachSong(func(song *resources.Song) error {
		if i > 0 {
			f.WriteString(",")
		}
		i++
		return encoder.Encode(song)
	})

	_, err = f.WriteString("]\n")
	if err != nil {
		return err
	}

	_, err = f.WriteString("}\n")
	if err != nil {
		return err
	}

	f.Close()

	return nil
}
