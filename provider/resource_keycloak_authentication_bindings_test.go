package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func TestAccKeycloakAuthenticationBindings_basic(t *testing.T) {

	flowAlias := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakAuthenticationBindingsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeycloakAuthenticationBindings(flowAlias),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakAuthenticationBindings("keycloak_authentication_bindings.bindings"),
					resource.TestCheckResourceAttr("keycloak_authentication_bindings.bindings", "realm_id", testAccRealm.Realm),
				),
			},
		},
	})
}

func TestAccKeycloakAuthenticationBindings_import(t *testing.T) {

	flowAlias := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakAuthenticationBindingsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeycloakAuthenticationBindings(flowAlias),
			},
			{
				ResourceName:      "keycloak_authentication_bindings.bindings",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getAuthenticationBindingsImportId("keycloak_authentication_bindings.bindings"),
			},
		},
	})
}

func getAuthenticationBindingsImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource %s not found", resourceName)
		}

		return rs.Primary.ID, nil
	}
}

func testAccCheckKeycloakAuthenticationBindings(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceName)
		}

		bindings, err := keycloakClient.GetAuthenticationBindings(rs.Primary.Attributes["realm_id"])
		if err != nil {
			return fmt.Errorf("error fetching authentication execution config: %v", err)
		}

		if bindings.BrowserFlowAlias == "browser" {
			return fmt.Errorf("Flows were not updated")
		}

		if bindings.ClientAuthenticationFlowAlias == "clients" {
			return fmt.Errorf("Flows were not updated")
		}

		return nil
	}
}

func testAccCheckKeycloakAuthenticationBindingsDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "keycloak_authentication_bindings" {
			continue
		}

		if json, err := keycloakClient.GetAuthenticationBindings(rs.Primary.Attributes["realm_id"]); err == nil {
			if json.BrowserFlowAlias != "browser" {
				return fmt.Errorf("authentication bindings were not cleaned up")
			}
		} else if !keycloak.ErrorIs404(err) {
			return fmt.Errorf("could not fetch realm: %v", err)
		}
	}

	return nil
}

func testAccKeycloakAuthenticationBindings(flowAlias string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_authentication_flow" "flow" {
	realm_id = data.keycloak_realm.realm.id
	alias    = "%s"
}

resource "keycloak_authentication_flow" "client-flow" {
	realm_id = data.keycloak_realm.realm.id
	alias    = "client-%s"
	provider_id = "client-flow"
}

resource "keycloak_authentication_execution" "execution" {
	realm_id          = data.keycloak_realm.realm.id
	parent_flow_alias = keycloak_authentication_flow.flow.alias
	authenticator     = "identity-provider-redirector"
}

resource "keycloak_authentication_bindings" "bindings" {
	realm_id     					= data.keycloak_realm.realm.id
	browser_flow_alias				= keycloak_authentication_flow.flow.alias
	registration_flow_alias			= keycloak_authentication_flow.flow.alias
	direct_grant_flow_alias			= keycloak_authentication_flow.flow.alias
	reset_credentials_flow_alias		= keycloak_authentication_flow.flow.alias
	client_authentication_flow_alias	= keycloak_authentication_flow.client-flow.alias
	docker_authentication_flow_alias	= keycloak_authentication_flow.flow.alias
}`, testAccRealm.Realm, flowAlias, flowAlias)
}
