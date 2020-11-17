package providers

import (
	"encoding/json"
	"highspot2/resources"
)

type Users struct {
	url string
}

func NewUsersProvider(url string) Provider {
	return &Users{
		url: url,
	}
}

func (self *Users) Fetch(fetchUser func(resource interface{}) error) error {
	piper, err := fetch(self.url)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(piper)
	for decoder.More() {
		var u resources.User
		err = decoder.Decode(&u)
		if err != nil {
			return err
		}
		fetchUser(&u)
		if err != nil {
			return err
		}
	}

	piper.Close()

	return nil
}
