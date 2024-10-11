package imageutils

import (
	"fmt"
	"io"
	"net/http"

	_ "image/jpeg"
	_ "image/png"
)

func LoadImageFromURL(URL string) ([]byte, error) {
	response, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("received %v response code", response.StatusCode)
	}
	return io.ReadAll(response.Body)
}
