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

func main() {
	fileToRead := getFileToReadFromArgs()
	log.Printf("Generating Protobuf Schema from: %s\n", fileToRead)

	byteValue, err := readFile(fileToRead)

	var result map[string]interface{}
	err = json.Unmarshal(byteValue, &result)
	if err != nil {
		panic(err)
	}

	rootMessage := Message{MessageName: "RootMessage"}
	createPBMessageDefinition(&result, &rootMessage)
	prettyPrint(&rootMessage)
}

func prettyPrint(message *Message) {

	fmt.Println("message", message.MessageName, "{")

	for _, v := range message.Attribute {
		if v.MessageDef != nil {
			prettyPrint(v.MessageDef)
		}

		jsonName := fmt.Sprintf("[json_name = \"%s\"];", v.JSONName)
		fmt.Println(" ", v.Type, v.Name, "=", v.Ordinal, jsonName)
	}

	fmt.Println("}")
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
//	@param messageParent is an element that we will attach a new Message element to.
//	For example:
//		(messageParent)
//		  	   |
//   	 (new message)
func createPBMessageDefinition(jsonElement *map[string]interface{}, message *Message) {

	for k, elem := range *jsonElement {
		//log.Printf("%s: Type: %s, %s\n", k, reflect.TypeOf(k), reflect.TypeOf(v))

		v := reflect.ValueOf(elem)
		//log.Println("Is Map: ", v.Kind() == reflect.Map)

		switch v.Kind() {
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

			canBeRepresentedAsMap, protoBufDataType := canBeRepresentedAsProbufMap(&e)
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
				childMessage := Message{
					MessageName: toCamelCaseWithFirstCharCapitalized(k),
				}

				createPBMessageDefinition(&e, &childMessage)

				// append the newly created message to its parent
				attr := MessageAttribute{
					Type:       toCamelCaseWithFirstCharCapitalized(k),
					Name:       toCamelCaseWithFirstCharInLowerCase(k),
					Ordinal:    len(message.Attribute) + 1,
					JSONName:   k,
					MessageDef: &childMessage,
				}

				message.Attribute = append(message.Attribute, attr)
			}
		default:
			log.Printf("No handler for %s, Type %s\n", k, v.Kind())
		}
	}
}

// Determines whether `element`'s value is of the same type
func canBeRepresentedAsProbufMap(element *map[string]interface{}) (bool, string) {
	var expectedKind reflect.Kind
	for _, v := range *element {
		kind := reflect.ValueOf(v).Kind()

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

	return ""
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

func getFileToReadFromArgs() string {
	file := flag.String("file", "", "The JSON file to parse")
	flag.Parse()

	return *file
}
