package providers

import (
	"encoding/json"
	"highspot2/resources"
)

type Changes struct {
	url string
}

func NewChangesProvider(url string) Provider {
	return &Changes{
		url: url,
	}
}

func (self *Changes) Fetch(fetchChange func(resource interface{}) error) error {
	piper, err := fetch(self.url)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(piper)
	for decoder.More() {
		var u resources.Change
		err = decoder.Decode(&u)
		if err != nil {
			return err
		}
		fetchChange(&u)
		if err != nil {
			return err
		}
	}

	piper.Close()

	return nil
}
