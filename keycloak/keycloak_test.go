package keycloak

import (
	"log"
	"os"
)

var requiredEnvironmentVariables = []string{
	"KEYCLOAK_CLIENT_ID",
	"KEYCLOAK_CLIENT_SECRET",
	"KEYCLOAK_CLIENT_TIMEOUT",
	"KEYCLOAK_REALM",
	"KEYCLOAK_URL",
}

var requiredEnvironmentVariablesDefaultValues = map[string]string{
	"KEYCLOAK_CLIENT_ID":      "terraform",
	"KEYCLOAK_CLIENT_SECRET":  "884e0f95-0f42-4a63-9b1f-94274655669e",
	"KEYCLOAK_CLIENT_TIMEOUT": "5",
	"KEYCLOAK_REALM":          "master",
	"KEYCLOAK_URL":            "http://localhost:8080",
}

func init() {

	for _, requiredEnvironmentVariable := range requiredEnvironmentVariables {
		if value := os.Getenv(requiredEnvironmentVariable); value == "" {
			os.Setenv(requiredEnvironmentVariable, requiredEnvironmentVariablesDefaultValues[requiredEnvironmentVariable])
			log.Println("JOED - " + requiredEnvironmentVariable + " set to " + requiredEnvironmentVariablesDefaultValues[requiredEnvironmentVariable])
		}
	}

	os.Setenv("TF_ACC", "1")
}
