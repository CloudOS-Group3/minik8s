package workflow

import (
	"errors"
	"fmt"
	"minik8s/pkg/api"
	"minik8s/pkg/config"
	"minik8s/pkg/serverless/function/function_util"
	"minik8s/util/httputil"
	"minik8s/util/log"
	"strconv"
	"strings"
)

func RunWorkflow(workflow api.Workflow, args []string) {
	// Run workflow
	log.Info("Run workflow %s", workflow.Metadata.Name)

	succssor := &workflow.Graph
	resWithName := map[string]interface{}{}
	for {
		if succssor == nil {
			break
		}
		node := *succssor
		function, err := GetFunction(node.Function.Name, node.Function.NameSpace)
		if err != nil {
			log.Error("Can't find function: %s. %s", node.Function.Name, err.Error())
			return
		}

		if len(resWithName) != 0 {
			args, _ = MakeParamsFromRet(function.Params, resWithName)
		}
		res, _ := RunFunction(*function, args)
		resWithName, _ = function_util.CheckParams(function.Result, res)
		succssor = CheckRule(node.Rule, resWithName)
	}
	// return res
}

// RunFunction : TODO
func RunFunction(function api.Function, args []string) ([]string, error) {
	// Run function
	params, err := function_util.CheckParams(function.Params, args)
	if err != nil {
		log.Error("Error checking params: %s", err.Error())
		return nil, err
	}
	log.Info("Run function %s, %v", function.Metadata.Name, params)
	// Create job
	return []string{"20", "hello"}, nil
}

// CheckRule : return successor ( or nil )
func CheckRule(rule api.NodeRule, res map[string]interface{}) *api.Graph {
	// Check rule
	if len(rule.Case) == 0 && rule.Default == nil {
		return nil
	}
	for _, case_ := range rule.Case {
		flag := true
		for _, expression := range case_.Expression {
			ifrule, err := CheckExpression(expression, res)
			if err != nil {
				log.Error("Error checking expression: %s", err)
				return nil
			}
			if !ifrule {
				flag = false
				break
			}
		}
		if flag {
			return case_.Successor
		}
	}
	return rule.Default
}

// CheckExpression : return if the expression is true
// Only support and operation now: Exp1 && Exp2 && ...
func CheckExpression(expression api.Expression, res map[string]interface{}) (bool, error) {
	// Check expression
	if expression.Type == api.Int {
		intValue, err := strconv.Atoi(expression.Value)
		if err != nil {
			return false, errors.New("Wrong type of rule " + expression.Variable + ", should be int")
		}
		intArg, ok := res[expression.Variable].(int)
		if !ok {
			return false, errors.New("Wrong type of arg " + expression.Variable + ", should be int")
		}
		switch expression.Opt {
		case api.Equal:
			return intArg == intValue, nil
		case api.NotEqual:
			return intArg != intValue, nil
		case api.Greater:
			return intArg > intValue, nil
		case api.Less:
			return intArg < intValue, nil
		case api.GreaterEqual:
			return intArg >= intValue, nil
		case api.LessEqual:
			return intArg <= intValue, nil
		default:
			return false, errors.New("Wrong operator " + expression.Opt)
		}
	} else if expression.Type == api.Float {
		floatValue, err := strconv.ParseFloat(expression.Value, 64)
		if err != nil {
			return false, errors.New("Wrong type of rule " + expression.Variable + ", should be float")
		}
		floatArg, ok := res[expression.Variable].(float64)
		if !ok {
			return false, errors.New("Wrong type of arg " + expression.Variable + ", should be float")
		}
		switch expression.Opt {
		case api.Equal:
			return floatArg == floatValue, nil
		case api.NotEqual:
			return floatArg != floatValue, nil
		case api.Greater:
			return floatArg > floatValue, nil
		case api.Less:
			return floatArg < floatValue, nil
		case api.GreaterEqual:
			return floatArg >= floatValue, nil
		case api.LessEqual:
			return floatArg <= floatValue, nil
		default:
			return false, errors.New("Wrong operator " + expression.Opt)
		}
	} else if expression.Type == api.String {
		stringArg, ok := res[expression.Variable].(string)
		if !ok {
			return false, errors.New("Wrong type of arg " + expression.Variable + ", should be string")
		}
		switch expression.Opt {
		case api.Equal:
			return stringArg == expression.Value, nil
		case api.NotEqual:
			return stringArg != expression.Value, nil
		case api.Greater:
			return len(stringArg) > len(expression.Value), nil
		case api.Less:
			return len(stringArg) < len(expression.Value), nil
		case api.GreaterEqual:
			return len(stringArg) >= len(expression.Value), nil
		case api.LessEqual:
			return len(stringArg) <= len(expression.Value), nil
		default:
			return false, errors.New("Wrong operator " + expression.Opt)
		}
	} else if expression.Type == api.Bool {
		boolValue, err := strconv.ParseBool(expression.Value)
		if err != nil {
			return false, errors.New("Wrong type of rule " + expression.Variable + ", should be bool")
		}
		boolArg, ok := res[expression.Variable].(bool)
		if !ok {
			return false, errors.New("Wrong type of arg " + expression.Variable + ", should be bool")
		}
		switch expression.Opt {
		case api.Equal:
			return boolArg == boolValue, nil
		case api.NotEqual:
			return boolArg != boolValue, nil
		default:
			return false, errors.New("Wrong operator " + expression.Opt)
		}
	} else {
		return false, errors.New("Wrong type " + expression.Type)
	}
}

// MakeParamsFromRet : get target args from []string
// Return is []string, eg. ["20", "hello"]
// paramTemp is rules to choose args, eg. [{{"name": "a", "type": "int"}, {"name": "b", "type": "string"}}]
func MakeParamsFromRet(paramTemp []api.Template, ret map[string]interface{}) ([]string, error) {
	arg := []string{}
	for i := 0; i < len(paramTemp); i++ {
		value, ok := ret[paramTemp[i].Name]
		if !ok {
			return nil, errors.New("Wrong number of arguments, should be " + strconv.Itoa(len(paramTemp)))
		}
		arg = append(arg, fmt.Sprintf("%v", value))
	}
	return arg, nil
}

func GetFunction(name string, namespace string) (*api.Function, error) {
	// Get function
	function := &api.Function{}
	URL := config.GetUrlPrefix() + config.FunctionURL
	URL = strings.Replace(URL, config.NamespacePlaceholder, namespace, -1)
	URL = strings.Replace(URL, config.NamePlaceholder, name, -1)
	err := httputil.Get(URL, function, "data")
	if err != nil {
		return nil, err
	}
	return function, nil
}
