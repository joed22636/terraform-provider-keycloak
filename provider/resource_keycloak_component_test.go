package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/joed22636/terraform-provider-keycloak/keycloak"
)

func TestAccKeycloakComponent_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakComponentDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakComponent_basic(),
				Check:  testAccCheckKeycloakComponentExists("keycloak_component.component"),
			},
			{
				ResourceName:      "keycloak_component.component",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getComponentImportId("keycloak_component.component"),
			},
		},
	})
}

func testAccCheckKeycloakComponentExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getComponentFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakComponentDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_authentication_subflow" {
				continue
			}

			id := rs.Primary.ID
			realm := rs.Primary.Attributes["realm_id"]
			parentId := rs.Primary.Attributes["parent_id"]
			providerType := rs.Primary.Attributes["provider_type"]

			components, _ := keycloakClient.GetComponents(realm, parentId, providerType)
			if components != nil {
				for _, component := range components {
					if component.Id == id {
						return fmt.Errorf("component with id %v still exists", id)
					}
				}
			}
		}

		return nil
	}
}

func getComponentFromState(s *terraform.State, resourceName string) (*keycloak.Component, error) {
	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]
	parentId := rs.Primary.Attributes["parent_id"]
	providerType := rs.Primary.Attributes["provider_type"]

	components, err := keycloakClient.GetComponents(realm, parentId, providerType)
	if err != nil {
		return nil, fmt.Errorf("error getting component with realm/parentId/providerType: %s %s %s - %v", realm, parentId, providerType, err)
	}
	if components != nil {
		for _, component := range components {
			if component.Id == id {
				return &component, nil
			}
		}
	}

	return nil, fmt.Errorf("Component with realm/parentId/providerType/id not found: %s %s %s %s", realm, parentId, providerType, id)
}

func getComponentImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}

		id := rs.Primary.ID
		realmId := rs.Primary.Attributes["realm_id"]

		return fmt.Sprintf("%s/%s", realmId, id), nil
	}
}

func testKeycloakComponent_basic() string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_component" "component" {
	realm_id = data.keycloak_realm.realm.id
	name = "my-ec-provider"
	parent_id = data.keycloak_realm.realm.id
	provider_id = "ecdsa-generated"
	provider_type = "org.keycloak.keys.KeyProvider"
	config = {
		"ecdsaEllipticCurveKey" = "P-521"
		"active" = "true"
		"priority" = "42"
		"enabled" = "true"
	}
}
	`, testAccRealm.Realm)
}
