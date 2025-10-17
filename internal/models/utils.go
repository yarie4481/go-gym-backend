package models

import (
	"encoding/json"
	"log"

	"gorm.io/datatypes"
)

// MapToJSON converts a map[string]interface{} into a datatypes.JSON value
func MapToJSON(data map[string]interface{}) datatypes.JSON {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error converting map to JSON: %v", err)
		return datatypes.JSON("{}")
	}
	return datatypes.JSON(jsonData)
}
