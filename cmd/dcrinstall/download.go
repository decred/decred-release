package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func downloadToFile(url, filename string) error {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error: %v", res.Status)
	}

	_, err = io.Copy(f, res.Body)
	if err != nil {
		return err
	}

	return nil
}
