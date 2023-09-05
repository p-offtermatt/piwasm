package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func resolveType(typeField map[string]interface{}) Type {
	// match the kind of the type
	switch typeField["kind"].(string) {
	case "rec": // records
		var fields []Field
		fieldsMap := typeField["fields"].(map[string]interface{})["fields"].([]interface{})
		for _, field := range fieldsMap {
			fieldMap := field.(map[string]interface{})
			fieldType := resolveType(fieldMap["fieldType"].(map[string]interface{}))
			fields = append(fields, Field{fieldMap["fieldName"].(string), fieldType})
		}
		return &StructType{Fields: fields}
	case "str":
		return &StrType{}
	case "const":
		// the type is just referenced here by an id and name.
		typeName := typeField["name"].(string)
		return &ConstType{Name: typeName}
	case "list":
		elementType := resolveType(typeField["elem"].(map[string]interface{}))
		return &ListType{ElementType: elementType}
	case "int":
		return &UInt64Type{}
	case "set":
		elementType := resolveType(typeField["elem"].(map[string]interface{}))
		return &SetType{ElementType: elementType}
	case "fun":
		argType := resolveType(typeField["arg"].(map[string]interface{}))
		returnType := resolveType(typeField["res"].(map[string]interface{}))
		return &MapType{Key: argType, Value: returnType}
	case "bool":
		return &BoolType{}
	case "tup":
		var types []Type
		fieldsMap := typeField["fields"].(map[string]interface{})["fields"].([]interface{})
		for _, field := range fieldsMap {
			fieldMap := field.(map[string]interface{})
			fieldType := resolveType(fieldMap["fieldType"].(map[string]interface{}))
			types = append(types, fieldType)
		}
		return &TupleType{Types: types}
	default:
		panic("kind not supported for resolving types: " + typeField["kind"].(string))
	}
}

func resolveDef(defField map[string]interface{}) Decl {
	// match the kind of the type
	switch defField["qualifier"].(string) {
	case "pureval":
		name := defField["name"].(string)
		valType := resolveType(defField["typeAnnotation"].(map[string]interface{}))
		block := resolveBlock(defField["expr"].(map[string]interface{}), valType)
		return &ConstDecl{Name: name, Type: valType, Value: block}
	case "puredef":
		// ====extract parameters====
		var paramNames []string
		var paramTypes []Type
		var returnType Type
		var statements Block
		// if there are no parameters, the shape of the puredef is different.
		if defField["expr"].(map[string]interface{})["params"] == nil {
			// no params
			paramNames = []string{}
			paramTypes = []Type{}

			// return type is the type in typeAnnotation
			returnType = resolveType(defField["typeAnnotation"].(map[string]interface{}))

			// ====extract the expression from expr=====
			statements = resolveBlock(defField["expr"].(map[string]interface{}), returnType)
		} else {
			// parameter names are given in expr.params
			for _, paramField := range defField["expr"].(map[string]interface{})["params"].([]interface{}) {
				// param name is in paramField.name
				paramName := paramField.(map[string]interface{})["name"]
				paramNames = append(paramNames, paramName.(string))
			}

			// types are in typeAnnotation.args
			for _, paramType := range defField["typeAnnotation"].(map[string]interface{})["args"].([]interface{}) {
				paramTypes = append(paramTypes, resolveType(paramType.(map[string]interface{})))
			}

			// ====extract the return type from typeAnnotations.res=====
			returnType = resolveType(defField["typeAnnotation"].(map[string]interface{})["res"].(map[string]interface{}))
			// ====extract the expression from expr.expr - the next layer will always be lambda =====
			statements = resolveBlock(defField["expr"].(map[string]interface{})["expr"].(map[string]interface{}), returnType)
		}

		// construct the params list
		var params []Param
		for i := 0; i < len(paramNames); i++ {
			// Mutable is false because these are puredefs
			params = append(params, Param{Name: paramNames[i], Type: paramTypes[i], Mutable: true})
		}

		return &FunctionDecl{Name: defField["name"].(string), Params: params, ReturnType: returnType, Body: statements.Statements}

	case "val":
		name := defField["name"].(string)
		expr := resolveExpr(defField["expr"].(map[string]interface{}), &ConstType{Name: "Todo"})
		return &ValDecl{Name: name, Value: expr}

	default:
		fmt.Println("qualifier not supported for resolving defs: " + defField["qualifier"].(string))
	}

	return nil
}

func resolveExpr(exprField map[string]interface{}, exprType Type) Expr {
	switch exprField["kind"].(string) {
	case "str":
		return &StringLiteral{Value: exprField["value"].(string)}

	case "int":
		return &UInt64Literal{Value: uint64(exprField["value"].(float64))}

	case "app":
		// this is an operator application

		// find out the opcode
		opcode := exprField["opcode"].(string)
		switch opcode {
		case "iadd":
			// addition
			args := exprField["args"].([]interface{})
			left := resolveExpr(args[0].(map[string]interface{}), &UInt64Type{})
			right := resolveExpr(args[1].(map[string]interface{}), &UInt64Type{})
			return &Add{Left: left, Right: right}

		case "Rec": // we are building a record
			args := exprField["args"].([]interface{})

			// get the fields of the struct from the args
			fields := make([]FieldValue, len(args)/2)
			// they are in args in the form name, value, name, value, ...

			for i := 0; i < len(args); i += 2 {
				// get the name arg
				nameArg := args[i].(map[string]interface{})
				name := nameArg["value"].(string)

				// get the value arg
				valueArg := args[i+1].(map[string]interface{})

				// TODO: get the type from the type list, since we know the name of the record
				value := resolveExpr(valueArg, nil)

				fields[i/2] = FieldValue{Name: name, Value: value}
			}

			return &StructCons{StructName: exprType.PrettyPrint(0), Fields: fields}

		case "ite":
			// this is an if-then-else expression
			args := exprField["args"].([]interface{})
			cond := resolveExpr(args[0].(map[string]interface{}), &BoolType{})
			then := resolveExpr(args[1].(map[string]interface{}), exprType)
			els := resolveExpr(args[2].(map[string]interface{}), exprType)
			return &IfElse{Condition: cond, Then: then, Else: els}

		case "not":
			// this is a not expression
			args := exprField["args"].([]interface{})
			expr := resolveExpr(args[0].(map[string]interface{}), &BoolType{})
			return &Not{Value: expr}

		case "contains":
			// this maps to `HashSet::contains_key(&self, &value))`
			args := exprField["args"].([]interface{})
			set := resolveExpr(args[0].(map[string]interface{}), &SetType{ElementType: WildcardType})
			value := resolveExpr(args[1].(map[string]interface{}), nil)
			return &MethodCall{
				Value:      set,
				MethodName: "contains_key",
				Arguments:  []Expr{&Borrow{Value: value}},
			}

		case "keys":
			// this maps to `mapExpr.keys().collect::<HashSet<_>>`
			args := exprField["args"].([]interface{})
			mapExpr := resolveExpr(args[0].(map[string]interface{}), &MapType{Key: WildcardType, Value: WildcardType})
			keysExpr := &MethodCall{
				Value:      mapExpr,
				MethodName: "keys",
				TypeArgs:   []Type{},
				Arguments:  []Expr{},
			}
			return &MethodCall{
				Value:      keysExpr,
				MethodName: "collect",
				TypeArgs:   []Type{&SetType{ElementType: WildcardType}},
				Arguments:  []Expr{},
			}

		case "field":
			// this is a field access
			args := exprField["args"].([]interface{})
			value := resolveExpr(args[0].(map[string]interface{}), nil)
			fieldName := args[1].(map[string]interface{})["value"].(string)
			return &FieldAccess{Value: value, Field: fieldName}

		default:
			fmt.Println("app opcode not supported for resolving expr: " + opcode)
		}

	case "name":
		// this is a variable
		return &Variable{VariableName: exprField["name"].(string)}

	case "let":
		// this is a let expression
		opdef := resolveDef(exprField["opdef"].(map[string]interface{})).(*ValDecl)
		body := resolveExpr(exprField["expr"].(map[string]interface{}), exprType)
		return &Let{VariableName: opdef.Name, Value: opdef.Value, Body: body}

	default:
		fmt.Println("kind not supported for resolving expr: " + exprField["kind"].(string))
	}

	return &Todo
}

// resolveBlock resolves an expression block
// the block should return something with the given exprType
// we need the exprType because otherwise it is impossible to tell what type a certain record that will be returned is, and
// rust needs that explicitly
func resolveBlock(exprField map[string]interface{}, exprType Type) Block {
	expr := resolveExpr(exprField, exprType)
	return Block{Statements: []Stmt{&Return{Value: expr}}}
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "  ")
	return string(s)
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <input file path> <output file path>\n", os.Args[0])
		os.Exit(1)
	}

	// read the file from the first argument
	filePath := os.Args[1]
	file, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	outputFilePath := os.Args[2]

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
	// for key, value := range namedTypeMap {
	// 	fmt.Println(key, value)
	// }

	var declarations []Decl

	// go through the modules
	for _, module := range data["modules"].([]interface{}) {
		// ignore modules ending in _stdlib or _test
		moduleMap := module.(map[string]interface{})
		if strings.HasSuffix(moduleMap["name"].(string), "_stdlib") || strings.HasSuffix(moduleMap["name"].(string), "_test") {
			continue
		}

		// collect all declarations
		for _, decl := range moduleMap["declarations"].([]interface{}) {
			declMap := decl.(map[string]interface{})
			switch declMap["kind"] {
			case "typedef":
				var declaration Decl
				name := declMap["name"].(string)
				declType := resolveType(declMap["type"].(map[string]interface{}))

				// if the type is a StructType, this should be a struct decl, otherwise a type decl
				if _, ok := declType.(*StructType); ok {
					structType := declType.(*StructType)

					// this is a struct decl
					declaration = &StructDecl{Name: name, Fields: structType.Fields}
				} else {
					// this is a type decl
					declaration = &TypeDecl{Name: name, Type: declType}
				}
				declarations = append(declarations, declaration)
				// fmt.Println(typeDef)
			case "import":
				// ignore imports
			case "def":
				declarations = append(declarations, resolveDef(declMap))
			default:
				fmt.Println("kind not supported: " + declMap["kind"].(string))
			}
		}
	}

	// hard code some dependencies we might need. rust can just ignore what we do not need
	imports := []Import{
		{Path: "std::collections::HashMap"},
		{Path: "std::collections::HashSet"},
		{Path: "super::neutron_stdlib::*"},
		{Path: "super::wasm_stdlib::*"},
	}

	program := Program{
		Imports: imports,
		Decls:   declarations,
	}

	err = os.WriteFile(outputFilePath, []byte(program.PrettyPrint(0)), 0o644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Println("Wrote output to ", outputFilePath)
}
