package httputil

import (
	"bytes"
	"encoding/json"
	"errors"
	"minik8s/pkg/config"
	"net/http"
)

func Get(URL string, result interface{}) error {
	res, err := http.Get(URL)

	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		err = errors.New("Http get response not ok")
		return err
	}

	err = json.NewDecoder(res.Body).Decode(result)
	if err != nil {
		return err
	}
	return nil
}

func Post(URL string, body []byte) error {
	res, err := http.Post(URL, config.JsonContent, bytes.NewBuffer(body))

	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		err = errors.New("Http post response not ok")
		return err
	}
	return nil
}

func Put(URL string, body []byte) error {
	req, err := http.NewRequest(http.MethodPut, URL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		err = errors.New("Http put response not ok")
		return err
	}
	return nil
}

func Delete(URL string) error {

	req, err := http.NewRequest(http.MethodDelete, URL, nil)

	if err != nil {
		return err
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		err = errors.New("Http delete response not ok")
		return err
	}
	return nil
}