package http

import (
	"encoding/json"
	"io/ioutil"
)

func GetJSON(url string, readStruct interface{}) error {
	client := NewDefaultRetryClient()
	res, err := client.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Todo check content-type response
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, &readStruct)
}
