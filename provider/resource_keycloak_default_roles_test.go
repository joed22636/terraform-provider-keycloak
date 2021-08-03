package provider

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/joed22636/terraform-provider-keycloak/keycloak"
)

func TestAccKeycloakDefaultRoles_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakDefaultRolesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeycloakDefaultRoles(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakDefaultRoles("keycloak_default_roles.default_roles", []string{"offline_access", "role_1"}),
					resource.TestCheckResourceAttr("keycloak_default_roles.default_roles", "realm_id", testAccRealm.Realm),
				),
			},
		},
	})
}

func TestAccKeycloakDefaultRoles_import(t *testing.T) {

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakDefaultRolesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeycloakDefaultRoles(),
			},
			{
				ResourceName:      "keycloak_default_roles.default_roles",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getDefaultRolesImportId("keycloak_default_roles.default_roles"),
			},
		},
	})
}

func getDefaultRolesImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource %s not found", resourceName)
		}

		return rs.Primary.ID, nil
	}
}

func testAccCheckKeycloakDefaultRoles(resourceName string, roles []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource %s not found in state", resourceName)
		}

		bindings, err := keycloakClient.GetDefaultRoles(rs.Primary.Attributes["realm_id"])
		if err != nil {
			return fmt.Errorf("error fetching authentication execution config: %v", err)
		}

		sort.Strings(bindings.DefaultRoles)
		if !reflect.DeepEqual(bindings.DefaultRoles, roles) {
			return fmt.Errorf("Expected default roles: %v, Actual: %v", roles, bindings.DefaultRoles)
		}

		return nil
	}
}

func testAccCheckKeycloakDefaultRolesDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "keycloak_default_roles" {
			continue
		}

		if json, err := keycloakClient.GetDefaultRoles(rs.Primary.Attributes["realm_id"]); err == nil {
			if !reflect.DeepEqual(json.DefaultRoles, []string{}) {
				return fmt.Errorf("default roles were not cleaned up")
			}
		} else if !keycloak.ErrorIs404(err) {
			return fmt.Errorf("could not fetch realm: %v", err)
		}
	}

	return nil
}

func testAccKeycloakDefaultRoles() string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_role" "role_1" {
	name     = "role_1"
	realm_id = data.keycloak_realm.realm.id
}

resource "keycloak_default_roles" "default_roles" {
	realm_id = data.keycloak_realm.realm.id
	default_roles    = [ "offline_access", "role_1" ]
}
`, testAccRealm.Realm)
}
