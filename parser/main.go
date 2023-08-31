package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type NodeType string

const (
	NodeTypeModule     NodeType = "Module"
	NodeTypeImport     NodeType = "Import"
	NodeTypeVariable   NodeType = "Variable"
	NodeTypeType       NodeType = "Type"
	NodeTypeFunction   NodeType = "Function"
	NodeTypeConstant   NodeType = "Constant"
	NodeTypeAssignment NodeType = "Assignment"
	NodeTypeExpression NodeType = "Expression"
	NodeTypeCall       NodeType = "Call"
	NodeTypeReturn     NodeType = "Return"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a file path as an argument.")
		return
	}

	// read the file from the first argument
	filePath := os.Args[1]
	file, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// read the whole file into a map
	var data map[string]interface{}
	err = json.Unmarshal(file, &data)
	if err != nil {
		panic(err)
	}

	// get the map of types and preprocess it
	namedTypeMap := make(map[int]string)
	for key, value := range data["types"].(map[string]interface{}) {
		// treat the value as another map
		valueTypeMap := value.(map[string]interface{})["type"].(map[string]interface{})
		name, ok := valueTypeMap["name"]
		if !ok {
			// type is probably anonymous, ignore
			continue
		}

		entryId, err := strconv.Atoi(key)
		if err != nil {
			fmt.Println("Error parsing type id for entry ", key)
			panic(err)
		}
		namedTypeMap[entryId] = name.(string)
		// if the
	}

	// print the map of types
	for key, value := range namedTypeMap {
		fmt.Println(key, value)
	}

	// go through the modules
	for _, module := range data["modules"].([]interface{}) {
		// ignore modules ending in _stdlib or _test
		if module.(map[string]interface{})["name"].(string) == "_stdlib" || module.(map[string]interface{})["name"].(string) == "_test" {
			continue
		}
	}
}
