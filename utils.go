package main

import (
	"strconv"
	"strings"
)

func GetStringSliceAtPath(data map[string]interface{}, path string) []string {
	value := GetDataAtPath(data, path)
	if value != nil {
		if val, ok := value.([]string); ok {
			return val
		} else if val, ok := value.([]interface{}); ok {
			result := make([]string, 0, len(val))
			for _, interfaceValue := range val {
				if interfaceValue != nil {
					if interfaceVal, ok := interfaceValue.(string); ok {
						result = append(result, interfaceVal)
					}
				}
			}
			return result
		}
	}
	return nil
}

// data/component:1/abc
func GetDataAtPath(data map[string]interface{}, path string) interface{} {
	pathComponents := strings.Split(path, "/")
	subData := data
	lastIndex := len(pathComponents) - 1
	for index, pathComponent := range pathComponents {
		pathComponentAndIndex := strings.Split(pathComponent, ":")
		var indexPath int
		var finalPathComponent string
		if len(pathComponentAndIndex) == 1 {
			finalPathComponent = pathComponentAndIndex[0]
			if subData[finalPathComponent] != nil {
				if index != lastIndex {
					subData = subData[finalPathComponent].(map[string]interface{})
				} else {
					return subData[finalPathComponent]
				}
			} else {
				return nil
			}
		} else {
			finalPathComponent = pathComponentAndIndex[0]
			indexPath, _ = strconv.Atoi(pathComponentAndIndex[1])

			if val, ok := subData[finalPathComponent].([]interface{}); ok {
				subData = val[indexPath].(map[string]interface{})
			} else if val, ok := subData[finalPathComponent].([]map[string]interface{}); ok {
				subData = val[indexPath]
			}
			if index == lastIndex {
				return subData
			}
		}

	}
	return nil
}

func ContainsByString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
