package util

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

func UnmarshalResponseBody(body io.ReadCloser, v interface{}) error {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	err = json.Unmarshal(b, v)
	if err != nil {
		return err
	}

	err = body.Close()
	if err != nil {
		return err
	}

	return nil
}
