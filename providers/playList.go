package providers

import (
	"encoding/json"
	"highspot2/resources"
)

type PlayLists struct {
	url string
}

func NewPlayListsProvider(url string) Provider {
	return &PlayLists{
		url: url,
	}
}

func (self *PlayLists) Fetch(fetchPlayList func(resource interface{}) error) error {
	piper, err := fetch(self.url)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(piper)
	for decoder.More() {
		var u resources.PlayList
		err = decoder.Decode(&u)
		if err != nil {
			return err
		}
		fetchPlayList(&u)
		if err != nil {
			return err
		}
	}

	piper.Close()

	return nil
}
