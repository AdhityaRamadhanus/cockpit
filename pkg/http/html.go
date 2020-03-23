package http

import (
	"io/ioutil"
)

func GetHTML(url string) (string, error) {
	client := NewDefaultRetryClient()
	res, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// Todo check content-type response
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), err
}
