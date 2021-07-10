package provider

import (
	"encoding/json"
	"fmt"
	"log"
)

func logComplexType(data interface{}) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	log.Println("LZA - " + string(b))
}
