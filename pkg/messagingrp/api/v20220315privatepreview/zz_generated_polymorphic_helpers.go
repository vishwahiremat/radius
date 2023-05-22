//go:build go1.18
// +build go1.18

// Licensed under the Apache License, Version 2.0 . See LICENSE in the repository root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package v20220315privatepreview

import "encoding/json"

func unmarshalRabbitMQQueuePropertiesClassification(rawMsg json.RawMessage) (RabbitMQQueuePropertiesClassification, error) {
	if rawMsg == nil {
		return nil, nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(rawMsg, &m); err != nil {
		return nil, err
	}
	var b RabbitMQQueuePropertiesClassification
	switch m["mode"] {
	case "recipe":
		b = &RecipeRabbitMQQueueProperties{}
	case "values":
		b = &ValuesRabbitMQQueueProperties{}
	default:
		b = &RabbitMQQueueProperties{}
	}
	return b, json.Unmarshal(rawMsg, b)
}
