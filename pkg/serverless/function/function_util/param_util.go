package function_util

import (
	"errors"
	"minik8s/pkg/api"
	"strconv"
)

func CheckParams(paramTemp []api.Template, args []string) (map[string]interface{}, error) {
	// DO NOT SUPPORT DEFAULT VALUE !!
	if len(paramTemp) != len(args) {
		return nil, errors.New("Wrong number of arguments, should be " + strconv.Itoa(len(paramTemp)))
	}
	params := make(map[string]interface{})
	for i := 0; i < len(args); i++ {
		value := args[i]
		if paramTemp[i].Type == "int" {
			// convert to int
			intValue, err := strconv.Atoi(value)
			if err != nil {
				return nil, errors.New("Wrong type of arg " + value + ", should be int")
			}
			params[paramTemp[i-1].Name] = intValue
		} else {
			params[paramTemp[i-1].Name] = value
		}
	}
	return params, nil
}
