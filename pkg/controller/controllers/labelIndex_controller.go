package controllers

import (
	"encoding/json"
	"io/ioutil"
	"minik8s/pkg/api"
	"minik8s/pkg/config"
	"minik8s/pkg/util"
	"minik8s/util/log"
	"net/http"
	"strings"
)

func GetLabelIndex(label map[string]string) (*api.LabelIndex, error) {
	labelString := util.ConvertLabelToString(label)

	URL := config.GetUrlPrefix() + config.LabelIndexURL
	URL = strings.Replace(URL, config.LabelParam, labelString, -1)

	res, err := http.Get(URL)
	if err != nil {
		log.Error("err get label index: %s", labelString)
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	labelIndex := &api.LabelIndex{}
	// deal with not found
	if res.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	err = json.Unmarshal(body, &labelIndex)
	if err != nil {
		log.Error("error unmarshal into label index: %s %v %v", err.Error(), body, labelIndex)
		return nil, err
	}

	return labelIndex, nil

}

func DeleteLabelIndex(label map[string]string) error {
	labelString := util.ConvertLabelToString(label)

	URL := config.GetUrlPrefix() + config.LabelIndexURL
	URL = strings.Replace(URL, config.LabelParam, labelString, -1)

	req, err := http.NewRequest(http.MethodDelete, URL, nil)
	if err != nil {
		log.Error("err delete label index: %s", labelString)
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error("err delete label index: %s", labelString)
		return err
	}

	defer res.Body.Close()

	return nil
}

func UpdateLabelIndex(labelIndex *api.LabelIndex) error {
	log.Info("update label index: %v", labelIndex)
	labelString := util.ConvertLabelToString(labelIndex.Labels)

	URL := config.GetUrlPrefix() + config.LabelIndexURL
	URL = strings.Replace(URL, config.LabelParam, labelString, -1)

	body, err := json.Marshal(labelIndex)
	if err != nil {
		log.Error("err add label index: %s", labelString)
		return err
	}

	res, err := http.Post(URL, "application/json", strings.NewReader(string(body)))
	if err != nil {
		log.Error("err add label index: %s", labelString)
		return err
	}

	defer res.Body.Close()

	return nil
}
