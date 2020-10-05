package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
)

type Args struct {
	jsonFile   string
	outputFile string
}

// Represents an attribute of the message
type MessageAttribute struct {
	Type       interface{}
	Name       string
	Ordinal    int
	JSONName   string
	MessageDef *Message
}

type Message struct {
	MessageName string
	Attribute   []MessageAttribute
}

var (
	// Elements in a array dont have names. We use this counter to give them a generic name
	// e.g. Object1, Object2, Object3, etc
	genericObjectCounter = 1
)

func main() {
	args := getFileToReadFromArgs()
	log.Printf("Generating Protobuf Schema from: %s\n", args.jsonFile)

	byteValue, err := readFile(args.jsonFile)

	var result map[string]interface{}
	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		panic(err)
	}

	rootMessage := Message{MessageName: "RootMessage"}
	createMessageDefinition(&result, &rootMessage)
	prettyPrint(&rootMessage, 0)
}

func prettyPrint(message *Message, indent int) {
	// Top level? Print Known Data Types
	if indent == 0 {
		fmt.Println(keyValuePairMessageDefinition())
		fmt.Println(listenerMessageDefinition())
	}

	tab := strings.Repeat("\t", indent)
	fmt.Println(tab, "message", message.MessageName, "{")

	for _, v := range message.Attribute {
		if v.MessageDef != nil {
			prettyPrint(v.MessageDef, indent+1)
		}

		jsonName := fmt.Sprintf("[json_name = \"%s\"];", v.JSONName)
		fmt.Println(strings.Repeat("\t", indent+1), v.Type, v.Name, "=", v.Ordinal, jsonName)
	}

	fmt.Println(tab, "}")
}

// @param jsonElement may be
// 	simple e.g. {"key": "value"} or
//	complex e.g.
//	{
//		"plugins": {
//			"available_on_server": {
//				"throttle": true
//			}
//		}
//	}
//
//	@param message is an element that we will attach the field definitions to.
func createMessageDefinition(jsonElement *map[string]interface{}, message *Message) {

	for k, elem := range *jsonElement {
		v := reflect.ValueOf(elem)

		switch v.Kind() {
		case reflect.Slice:
			log.Println(k, " :: ", elem)
			genericElementSlice := elem.([]interface{})
			if len(genericElementSlice) == 0 {
				// If it's just an array with no element. Assume string array.
				// Generate a protobuf message definition as follows:
				// repeat string fieldName = <ordinal>;
				addRepeatedField(k, message, "string", nil)
			} else {
				// We assume that an array is composed of a homogeneous data type.
				// Thus we only pass the address of the first element to check
				isAnArrayOfKnownDataTypes, dataType := isKnownDataType(&genericElementSlice[0])
				log.Println(k, " :: ", isAnArrayOfKnownDataTypes, " :: ", dataType)
				if isAnArrayOfKnownDataTypes {
					addRepeatedField(k, message, dataType, nil)
				} else {
					genericEntityName := createGenericEntityName()
					elementAsMap := genericElementSlice[0].(map[string]interface{})

					child := Message{ MessageName: genericEntityName }
					createMessageDefinition(&elementAsMap, &child)

					addRepeatedField(k, message, genericEntityName, &child)
				}
			}
		// Simple types
		case reflect.String, reflect.Float64, reflect.Bool:
			attr := MessageAttribute{
				Type:       kindToProtoBufType(v.Kind()),
				Name:       toCamelCaseWithFirstCharInLowerCase(k),
				Ordinal:    len(message.Attribute) + 1,
				JSONName:   k,
				MessageDef: nil,
			}

			message.Attribute = append(message.Attribute, attr)
		case reflect.Map:
			// We've encountered a map. Check whether we can represent this as
			// map<string, type> in protobuf.
			//
			// We can only represent it as map<string, type> if type is homogenous e.g.
			// of a single type, float, bool etc
			e := elem.(map[string]interface{})

			canBeRepresentedAsMap, protoBufDataType := canBeRepresentedAsProtobufMap(&e)
			if canBeRepresentedAsMap {
				// append the newly created message to its parent
				attr := MessageAttribute{
					Type:       fmt.Sprintf("map<string, %s>", protoBufDataType),
					Name:       toCamelCaseWithFirstCharInLowerCase(k),
					Ordinal:    len(message.Attribute) + 1,
					JSONName:   k,
					MessageDef: nil,
				}

				message.Attribute = append(message.Attribute, attr)

			} else {
				messageName := toCamelCaseWithFirstCharCapitalized(k)
				fieldName := toCamelCaseWithFirstCharInLowerCase(k)

				child := Message{ MessageName: messageName }
				createMessageDefinition(&e, &child)

				attr := MessageAttribute {
					Type:       messageName,
					Name:       fieldName,
					Ordinal:    len(message.Attribute) + 1,
					JSONName:   k,
					MessageDef: &child,
				}

				message.Attribute = append(message.Attribute, attr)
			}
		default:
			log.Printf("No handler for %s, Type %s\n", k, v.Kind())
		}
	}
}

func addRepeatedField(jsonName string, message *Message, dataType string, messageDef *Message) {
	attr := MessageAttribute{
		Type:       fmt.Sprintf("repeated %s", dataType),
		Name:       toCamelCaseWithFirstCharInLowerCase(jsonName),
		Ordinal:    len(message.Attribute) + 1,
		JSONName:   jsonName,
		MessageDef: messageDef,
	}

	message.Attribute = append(message.Attribute, attr)
}

func createGenericEntityName() string {
	genericMessageName := fmt.Sprint("Object", genericObjectCounter)
	genericObjectCounter++
	return genericMessageName
}

// Returns true if all of the array elements are of the same datatype
// and is one of: string, float
//
// Under normal circumstances, an array is composed of a single datatype.
func isKnownDataType(v *interface{}) (bool, string) {
	var expectedKind reflect.Kind
	kind := reflect.ValueOf(*v).Kind()

	switch kind {
	case reflect.Map:
		// cast the value pointed to by `v` to map[string]interface{}
		m := (*v).(map[string]interface{})

		if isKeyValuePairMap(m) {
			return true, "KeyValuePair"
		}

		if isListener(m) {
			return true, "Listener"
		}

	case reflect.String, reflect.Float64, reflect.Float32, reflect.Uint8:
		return true, kindToProtoBufType(expectedKind)
	default:
		// kind is not one of datatype defined above
		return false, ""
	}

	return false, ""
}

// Determines whether `element`'s value is of the same type
func canBeRepresentedAsProtobufMap(element *map[string]interface{}) (bool, string) {
	var expectedKind reflect.Kind
	for _, v := range *element {
		kind := reflect.ValueOf(v).Kind()

		// Nested map?
		if kind == reflect.Map {
			return false, ""
		}

		// set expected kind to the first element
		// that we see.
		if expectedKind == 0 {
			expectedKind = kind
			continue
		}

		if expectedKind != kind {
			return false, ""
		}
	}

	return true, kindToProtoBufType(expectedKind)
}

func kindToProtoBufType(kind reflect.Kind) string {
	switch kind {
	case reflect.String:
		return "string"
	case reflect.Float32, reflect.Float64:
		return "float"
	case reflect.Bool:
		return "bool"
	}

	// We got not idea what kind it is. Assume string.
	return "string"
}

func toCamelCaseWithFirstCharCapitalized(s string) string {
	return toCamelCase(s, false)
}

func toCamelCaseWithFirstCharInLowerCase(s string) string {
	return toCamelCase(s, true)
}

func Split(r rune) bool {
	return r == '_' || r == '-'
}

func toCamelCase(s string, isFirstCharLowerCase bool) string {
	var attrNameTokens []string
	for index, token := range strings.FieldsFunc(s, Split) {

		if isFirstCharLowerCase && index == 0 {
			token = strings.ToLower(token)
		} else {
			token = strings.Title(token)
		}

		attrNameTokens = append(attrNameTokens, token)
	}

	return strings.Join(attrNameTokens, "")
}

func readFile(fileToRead string) ([]byte, error) {
	jsonFile, err := os.Open(fileToRead)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	return ioutil.ReadAll(jsonFile)
}

func getFileToReadFromArgs() Args {
	jsonFile := flag.String("file", "", "The JSON file to parse")
	outputFile := flag.String("out", "", "The file where the output will be written to")
	flag.Parse()

	return Args{
		jsonFile:   *jsonFile,
		outputFile: *outputFile,
	}
}
