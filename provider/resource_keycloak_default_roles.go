package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/joed22636/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakDefaultRoles() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakDefaultRolesCreate,
		Read:   resourceKeycloakDefaultRolesRead,
		Delete: resourceKeycloakDefaultRolesDelete,
		Update: resourceKeycloakDefaultRolesUpdate,
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakDefaultRolesImport,
		},
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"default_roles": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},
	}
}

func resourceKeycloakDefaultRolesCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	json, err := getDefaultRolesFromData(data)
	if err != nil {
		return err
	}

	err = keycloakClient.UpdateDefaultRoles(json)
	if err != nil {
		return err
	}

	setDefaultRoles(data, json)

	return resourceKeycloakDefaultRolesRead(data, meta)
}

func resourceKeycloakDefaultRolesRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	result, err := keycloakClient.GetDefaultRoles(data.Id())
	if err != nil {
		return handleNotFoundError(err, data)
	}

	setDefaultRoles(data, result)

	return nil
}

func resourceKeycloakDefaultRolesUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	json, err := getDefaultRolesFromData(data)
	if err != nil {
		return err
	}

	err = keycloakClient.UpdateDefaultRoles(json)
	if err != nil {
		return err
	}

	setDefaultRoles(data, json)

	return nil
}

func resourceKeycloakDefaultRolesDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	json, err := getDefaultRolesFromData(data)
	if err != nil {
		return err
	}

	json.DefaultRoles = []string{}

	err = keycloakClient.UpdateDefaultRoles(json)
	if err != nil {
		return err
	}

	return nil
}

func getDefaultRolesFromData(data *schema.ResourceData) (*keycloak.DefaultRoles, error) {
	result := &keycloak.DefaultRoles{
		RealmId: data.Get("realm_id").(string),
		// DefaultRoles: data.Get("default_roles").([]string),
	}
	defaultRoles := make([]string, 0)
	if v, ok := data.GetOk("default_roles"); ok {
		for _, defaultRole := range v.(*schema.Set).List() {
			defaultRoles = append(defaultRoles, defaultRole.(string))
		}
	}
	result.DefaultRoles = defaultRoles
	return result, nil
}

func setDefaultRoles(data *schema.ResourceData, json *keycloak.DefaultRoles) {
	data.SetId(json.RealmId)
	data.Set("default_roles", json.DefaultRoles)
}

func resourceKeycloakDefaultRolesImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	d.Set("realm_id", d.Id())
	return []*schema.ResourceData{d}, nil
}
