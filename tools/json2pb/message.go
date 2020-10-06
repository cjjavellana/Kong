package main

import "strings"

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

func (m *Message) addAttribute(
	jsonName string,
	dataType string,
	messageDef *Message,
	repeated bool,
) {

	var t string
	if repeated {
		t = strings.Join([]string{"repeated", dataType}, " ")
	} else {
		t = dataType
	}

	attr := MessageAttribute{
		Type:       t,
		Name:       asCamelCaseFieldName(jsonName),
		Ordinal:    len(m.Attribute) + 1,
		JSONName:   jsonName,
		MessageDef: messageDef,
	}

	m.Attribute = append(m.Attribute, attr)
}
