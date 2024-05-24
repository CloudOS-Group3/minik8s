package httputil

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"minik8s/pkg/config"
	"minik8s/util/log"
	"net/http"
)

func Get(URL string, result interface{}, key string) error {
	res, err := http.Get(URL)

	if err != nil {
		log.Error("Error http get")
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		err = errors.New("Http get response not ok")
		return err
	}
	var resMap map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&resMap)

	if err != nil {
		log.Error("Error json decode: %s", err.Error())
		return err
	}
	value, ok := resMap[key]
	if !ok {
		log.Warn("Empty data with key: %s", key)
		return nil
	}
	dataStr := fmt.Sprint(value)
	err = json.Unmarshal([]byte(dataStr), result)
	if err != nil {
		log.Error("Error json unmarshal: %s", err.Error())
		return err
	}
	return nil
}

func Post(URL string, body []byte) error {
	res, err := http.Post(URL, config.JsonContent, bytes.NewBuffer(body))

	if err != nil {
		log.Error("Error http post")
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Error("res is %v", res)
		err = errors.New("Http post response not ok")
		return err
	}
	return nil
}

func Put(URL string, body []byte) error {
	req, err := http.NewRequest(http.MethodPut, URL, bytes.NewBuffer(body))
	if err != nil {
		log.Error("Error create put request")
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Debug("Error client do put: %s", err.Error())
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
		log.Error("Error create delete request")
		return err
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Error("Error client do delete")
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		err = errors.New("Http delete response not ok")
		return err
	}
	return nil
}
