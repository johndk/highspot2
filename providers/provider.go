package providers

import (
	"highspot2/providers/http"
	"io"
	"log"
)

type Provider interface {
	Fetch(fetch func(resource interface{}) error) error
}

func fetch(url string) (*io.PipeReader, error) {
	piper, pipew := io.Pipe()
	client := http.NewClient(pipew, url)

	go func() {
		defer pipew.Close()
		err := client.Read()
		if err != nil {
			log.Fatalf("Error encountered. %v", err)
			return
		}
	}()

	return piper, nil
}
