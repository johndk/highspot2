package providers

import (
	"encoding/json"
	"highspot2/resources"
)

type Songs struct {
	url string
}

func NewSongsProvider(url string) Provider {
	return &Songs{
		url: url,
	}
}

func (self *Songs) Fetch(fetchSong func(resource interface{}) error) error {
	piper, err := fetch(self.url)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(piper)
	for decoder.More() {
		var u resources.Song
		err = decoder.Decode(&u)
		if err != nil {
			return err
		}
		fetchSong(&u)
		if err != nil {
			return err
		}
	}

	piper.Close()

	return nil
}
