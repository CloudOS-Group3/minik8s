package controllers

import (
	"encoding/json"
	"minik8s/pkg/api"
	"minik8s/pkg/config"
	"minik8s/pkg/util"
	"minik8s/util/httputil"
	"minik8s/util/log"
	"strings"
)

func GetLabelIndex(label map[string]string) (*api.LabelIndex, error) {
	labelString := util.ConvertLabelToString(label)
	log.Info("get label index: %v", labelString)

	URL := config.GetUrlPrefix() + config.LabelIndexURL
	URL = strings.Replace(URL, config.LabelPlaceholder, labelString, -1)

	labelIndex := &api.LabelIndex{}
	err := httputil.Get(URL, &labelIndex, "data")
	if err != nil {
		log.Error("err get label index: %s", labelString)
		return nil, err
	}

	return labelIndex, nil

}

func DeleteLabelIndex(label map[string]string) error {
	labelString := util.ConvertLabelToString(label)

	URL := config.GetUrlPrefix() + config.LabelIndexURL
	URL = strings.Replace(URL, config.LabelPlaceholder, labelString, -1)

	err := httputil.Delete(URL)

	if err != nil {
		log.Error("error deleting pod")
		return err
	}
	return nil
}

func UpdateLabelIndex(labelIndex *api.LabelIndex) error {
	log.Info("update label index: %v", labelIndex)
	labelString := util.ConvertLabelToString(labelIndex.Labels)

	URL := config.GetUrlPrefix() + config.LabelIndexURL
	URL = strings.Replace(URL, config.LabelPlaceholder, labelString, -1)

	body, err := json.Marshal(labelIndex)
	if err != nil {
		log.Error("err add label index: %s", labelString)
		return err
	}

	err = httputil.Post(URL, body)
	if err != nil {
		log.Error("err add label index: %s", labelString)
		return err
	}

	return nil
}
