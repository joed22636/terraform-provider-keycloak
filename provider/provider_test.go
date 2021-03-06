package provider

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/meta"
	"github.com/joed22636/terraform-provider-keycloak/keycloak"
)

var testAccProviderFactories map[string]func() (*schema.Provider, error)
var testAccProvider *schema.Provider
var keycloakClient *keycloak.KeycloakClient
var testAccRealm *keycloak.Realm
var testAccRealmTwo *keycloak.Realm
var testAccRealmUserFederation *keycloak.Realm

var requiredEnvironmentVariables = []string{
	"KEYCLOAK_CLIENT_ID",
	"KEYCLOAK_CLIENT_SECRET",
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
	userAgent := fmt.Sprintf("HashiCorp Terraform/%s (+https://www.terraform.io) Terraform Plugin SDK/%s", schema.Provider{}.TerraformVersion, meta.SDKVersionString())

	for _, requiredEnvironmentVariable := range requiredEnvironmentVariables {
		if value := os.Getenv(requiredEnvironmentVariable); value == "" {
			os.Setenv(requiredEnvironmentVariable, requiredEnvironmentVariablesDefaultValues[requiredEnvironmentVariable])
			log.Println("JOED - " + requiredEnvironmentVariable + " set to " + requiredEnvironmentVariablesDefaultValues[requiredEnvironmentVariable])
		}
	}

	os.Setenv("TF_ACC", "1")

	keycloakClient, _ = keycloak.NewKeycloakClient(os.Getenv("KEYCLOAK_URL"), "/auth", os.Getenv("KEYCLOAK_CLIENT_ID"), os.Getenv("KEYCLOAK_CLIENT_SECRET"), os.Getenv("KEYCLOAK_REALM"), "", "", true, 5, "", false, userAgent, map[string]string{
		"foo": "bar",
	})
	testAccProvider = KeycloakProvider(keycloakClient)
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"keycloak": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
}

func TestMain(m *testing.M) {
	testAccRealm = createTestRealm()
	testAccRealmTwo = createTestRealm()
	testAccRealmUserFederation = createTestRealm()

	code := m.Run()

	err := keycloakClient.DeleteRealm(testAccRealm.Realm)
	if err != nil {
		os.Exit(1)
	}

	err = keycloakClient.DeleteRealm(testAccRealmTwo.Realm)
	if err != nil {
		os.Exit(1)
	}

	err = keycloakClient.DeleteRealm(testAccRealmUserFederation.Realm)
	if err != nil {
		os.Exit(1)
	}

	os.Exit(code)
}

func createTestRealm() *keycloak.Realm {
	name := acctest.RandomWithPrefix("tf-acc")
	r := &keycloak.Realm{
		Id:      name,
		Realm:   name,
		Enabled: true,
	}

	err := keycloakClient.NewRealm(r)
	if err != nil {
		os.Exit(1)
	}

	return r
}

func TestProvider(t *testing.T) {
	if err := testAccProvider.InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	for _, requiredEnvironmentVariable := range requiredEnvironmentVariables {
		if value := os.Getenv(requiredEnvironmentVariable); value == "" {
			t.Fatalf("%s must be set before running acceptance tests.", requiredEnvironmentVariable)
		}
	}
}
