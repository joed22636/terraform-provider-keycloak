package keycloak

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"golang.org/x/net/publicsuffix"
)

// Some actions, such as creating a realm, require a refresh
// before a GET can be performed on that realm
//
// This test ensures that, after creating a realm and performing
// a GET, the access token and refresh token have changed
//
// Any action that returns a 403 or a 401 could be used for this test
// Creating a realm is just the only one I'm aware of
//
// This appears to have been fixed as of Keycloak 12.x
func TestAccKeycloakApiClientRefresh(t *testing.T) {
	for _, requiredEnvironmentVariable := range requiredEnvironmentVariables {
		if value := os.Getenv(requiredEnvironmentVariable); value == "" {
			t.Fatalf("%s must be set before running acceptance tests.", requiredEnvironmentVariable)
		}
	}

	if v := os.Getenv("KEYCLOAK_CLIENT_SECRET"); v == "" {
		if v := os.Getenv("KEYCLOAK_USER"); v == "" {
			t.Fatal("KEYCLOAK_USER must be set for acceptance tests")
		}
		if v := os.Getenv("KEYCLOAK_PASSWORD"); v == "" {
			t.Fatal("KEYCLOAK_PASSWORD must be set for acceptance tests")
		}
	}

	// Disable [DEBUG] logs which terraform typically handles for you. Re-enable when finished
	if tfLogLevel := os.Getenv("TF_LOG"); tfLogLevel == "" {
		log.SetOutput(ioutil.Discard)
		defer log.SetOutput(os.Stdout)
	}

	// Convert KEYCLOAK_CLIENT_TIMEOUT to int
	clientTimeout, err := strconv.Atoi(os.Getenv("KEYCLOAK_CLIENT_TIMEOUT"))
	if err != nil {
		t.Fatal("KEYCLOAK_CLIENT_TIMEOUT must be an integer")
	}

	keycloakClient, err := NewKeycloakClient(os.Getenv("KEYCLOAK_URL"), "/auth", os.Getenv("KEYCLOAK_CLIENT_ID"), os.Getenv("KEYCLOAK_CLIENT_SECRET"), os.Getenv("KEYCLOAK_REALM"), os.Getenv("KEYCLOAK_USER"), os.Getenv("KEYCLOAK_PASSWORD"), true, clientTimeout, "", false, "", map[string]string{
		"foo": "bar",
	})
	if err != nil {
		t.Fatalf("%s", err)
	}

	// skip test if running 12.x or greater
	if keycloakClient.VersionIsGreaterThanOrEqualTo(Version_12) {
		t.Skip()
	}

	realmName := "terraform-" + acctest.RandString(10)
	realm := &Realm{
		Realm: realmName,
		Id:    realmName,
	}

	err = keycloakClient.NewRealm(realm)
	if err != nil {
		t.Fatalf("%s", err)
	}

	_, err = keycloakClient.GetRealm(realmName) // This should not fail since it will automatically refresh and try again
	if err != nil {
		t.Fatalf("%s", err)
	}

	// Clean up - the realm doesn't need to exist in order for us to assert against the refreshed tokens
	err = keycloakClient.DeleteRealm(realmName)
	if err != nil {
		t.Fatalf("%s", err)
	}

}

func TestKeycloakApiClientRefresh(t *testing.T) {
	t.Skip("just a play around test case - can be removed")
	cookieJar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		t.Fatal("cookiejar issue")
	}
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
		Proxy:           http.ProxyFromEnvironment,
	}
	httpClient := &http.Client{
		Timeout:   time.Second * time.Duration(30),
		Transport: transport,
		Jar:       cookieJar,
	}

	username := "admin"
	password := "admin"
	clientCredentials := &ClientCredentials{
		ClientId:     "admin-cli",
		ClientSecret: "",
	}
	clientCredentials.Username = username
	clientCredentials.Password = password
	clientCredentials.GrantType = "password"

	keycloakClient := KeycloakClient{
		baseUrl:           "http://localhost:8080" + "/auth",
		clientCredentials: clientCredentials,
		httpClient:        httpClient,
		initialLogin:      false,
		realm:             "master",
		userAgent:         "HashiCorp Terraform/%s (+https://www.terraform.io) Terraform Plugin SDK/%s",
		additionalHeaders: nil,
	}

	r, e := keycloakClient.GetIdentityProviderMapper("test", "oidc", "7718c4d0-5c6b-4b1b-9a83-fb40aa2746a6")
	log.Println(r, e)
	keycloakClient.UpdateIdentityProviderMapper(r)
}
