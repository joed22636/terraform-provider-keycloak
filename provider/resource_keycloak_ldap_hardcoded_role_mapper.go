package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/joed22636/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakLdapHardcodedRoleMapper() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakLdapHardcodedRoleMapperCreate,
		Read:   resourceKeycloakLdapHardcodedRoleMapperRead,
		Update: resourceKeycloakLdapHardcodedRoleMapperUpdate,
		Delete: resourceKeycloakLdapHardcodedRoleMapperDelete,
		// This resource can be imported using {{realm}}/{{provider_id}}/{{mapper_id}}. The Provider and Mapper IDs are displayed in the GUI
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakLdapGenericMapperImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name of the mapper when displayed in the console.",
			},
			"realm_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The realm in which the ldap user federation provider exists.",
			},
			"ldap_user_federation_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ldap user federation provider to attach this mapper to.",
			},
			"role": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Role to grant to user.",
			},
		},
	}
}

func getLdapHardcodedRoleMapperFromData(data *schema.ResourceData) *keycloak.LdapHardcodedRoleMapper {
	return &keycloak.LdapHardcodedRoleMapper{
		Id:                   data.Id(),
		Name:                 data.Get("name").(string),
		RealmId:              data.Get("realm_id").(string),
		LdapUserFederationId: data.Get("ldap_user_federation_id").(string),
		Role:                 data.Get("role").(string),
	}
}

func setLdapHardcodedRoleMapperData(data *schema.ResourceData, ldapMapper *keycloak.LdapHardcodedRoleMapper) {
	data.SetId(ldapMapper.Id)
	data.Set("name", ldapMapper.Name)
	data.Set("realm_id", ldapMapper.RealmId)
	data.Set("ldap_user_federation_id", ldapMapper.LdapUserFederationId)
	data.Set("role", ldapMapper.Role)
}

func resourceKeycloakLdapHardcodedRoleMapperCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	ldapMapper := getLdapHardcodedRoleMapperFromData(data)

	err := keycloakClient.ValidateLdapHardcodedRoleMapper(ldapMapper)
	if err != nil {
		return err
	}

	err = keycloakClient.NewLdapHardcodedRoleMapper(ldapMapper)
	if err != nil {
		return err
	}

	setLdapHardcodedRoleMapperData(data, ldapMapper)

	return resourceKeycloakLdapHardcodedRoleMapperRead(data, meta)
}

func resourceKeycloakLdapHardcodedRoleMapperRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	ldapMapper, err := keycloakClient.GetLdapHardcodedRoleMapper(realmId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	setLdapHardcodedRoleMapperData(data, ldapMapper)

	return nil
}

func resourceKeycloakLdapHardcodedRoleMapperUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	ldapMapper := getLdapHardcodedRoleMapperFromData(data)

	err := keycloakClient.ValidateLdapHardcodedRoleMapper(ldapMapper)
	if err != nil {
		return err
	}

	err = keycloakClient.UpdateLdapHardcodedRoleMapper(ldapMapper)
	if err != nil {
		return err
	}

	setLdapHardcodedRoleMapperData(data, ldapMapper)

	return nil
}

func resourceKeycloakLdapHardcodedRoleMapperDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	return keycloakClient.DeleteLdapHardcodedRoleMapper(realmId, id)
}
