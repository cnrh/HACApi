package utils

import (
	"errors"
	"fmt"
	"os"
)

// ErrorURLConnection is the error thrown when a not supported URL is passed into BuildConnectionURL().
var ErrorURLConnection = errors.New("connection not supported")

// BuildConnectionURL builds a connection URL for an inputted service.
func BuildConnectionURL(n string) (string, error) {
	var url string

	switch n {
	case "fiber":
		url = fmt.Sprintf("%s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))
	default:
		return "", ErrorURLConnection
	}

	return url, nil
}
