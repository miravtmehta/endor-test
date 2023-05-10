package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

func objectExtraction(key, valueString string) (Object, error) {
	var object Object

	getKindType := SplitObjectStore(key)
	switch getKindType {
	case AnimalKind:
		object = &Animal{}
	case PersonKind:
		object = &Person{}
	default:
		return nil, fmt.Errorf("unknown object type: %s", getKindType)
	}
	err := json.Unmarshal([]byte(valueString), object)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func SplitObjectStore(objectString string) string {

	splitter := strings.Split(objectString, ":")
	// ID : Name : Kind
	return splitter[2]

}
