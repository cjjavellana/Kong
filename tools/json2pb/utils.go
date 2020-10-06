package main

import "strings"

func asCamelCaseClassName(s string) string {
	return toCamelCase(s, false)
}

func asCamelCaseFieldName(s string) string {
	return toCamelCase(s, true)
}

func split(r rune) bool {
	return r == '_' || r == '-'
}

func toCamelCase(s string, isFirstCharLowerCase bool) string {
	var attrNameTokens []string
	for index, token := range strings.FieldsFunc(s, split) {

		if isFirstCharLowerCase && index == 0 {
			token = strings.ToLower(token)
		} else {
			token = strings.Title(token)
		}

		attrNameTokens = append(attrNameTokens, token)
	}

	return strings.Join(attrNameTokens, "")
}
