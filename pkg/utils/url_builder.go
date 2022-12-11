package utils

import (
	"errors"
	"fmt"
	"os"
)

var ErrorURLConnection = errors.New("connection not supported")

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
