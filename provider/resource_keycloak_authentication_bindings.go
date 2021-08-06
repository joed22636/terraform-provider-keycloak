package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/joed22636/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakAuthenticationBindings() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakAuthenticationBindingsCreate,
		Read:   resourceKeycloakAuthenticationBindingsRead,
		Delete: resourceKeycloakAuthenticationBindingsDelete,
		Update: resourceKeycloakAuthenticationBindingsUpdate,
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakAuthenticationBindingsImport,
		},
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"browser_flow_alias": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"registration_flow_alias": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"direct_grant_flow_alias": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"reset_credentials_flow_alias": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"client_authentication_flow_alias": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"docker_authentication_flow_alias": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
		},
	}
}

func resourceKeycloakAuthenticationBindingsCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	json, err := getAuthenticationBindingsFromData(data)
	if err != nil {
		return err
	}

	err = keycloakClient.UpdateAuthenticationBindings(json)
	if err != nil {
		return err
	}

	setAuthenticationBindings(data, json)

	return resourceKeycloakAuthenticationBindingsRead(data, meta)
}

func resourceKeycloakAuthenticationBindingsRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	result, err := keycloakClient.GetAuthenticationBindings(data.Id())
	if err != nil {
		return handleNotFoundError(err, data)
	}

	setAuthenticationBindings(data, result)

	return nil
}

func resourceKeycloakAuthenticationBindingsUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	json, err := getAuthenticationBindingsFromData(data)
	if err != nil {
		return err
	}

	err = keycloakClient.UpdateAuthenticationBindings(json)
	if err != nil {
		return err
	}

	setAuthenticationBindings(data, json)

	return nil
}

func resourceKeycloakAuthenticationBindingsDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)
	json, err := getAuthenticationBindingsFromData(data)
	if err != nil {
		return err
	}

	json.BrowserFlowAlias = "browser"
	json.RegistrationFlowAlias = "registration"
	json.DirectGrantFlowAlias = "direct grant"
	json.ResetCredentialsFlowAlias = "reset credentials"
	json.ClientAuthenticationFlowAlias = "clients"
	json.DockerAuthenticationFlowAlias = "docker auth"

	err = keycloakClient.UpdateAuthenticationBindings(json)
	if err != nil {
		return err
	}

	return nil
}

func getAuthenticationBindingsFromData(data *schema.ResourceData) (*keycloak.AuthenticationBindings, error) {
	return &keycloak.AuthenticationBindings{
		RealmId:                       data.Get("realm_id").(string),
		BrowserFlowAlias:              data.Get("browser_flow_alias").(string),
		RegistrationFlowAlias:         data.Get("registration_flow_alias").(string),
		DirectGrantFlowAlias:          data.Get("direct_grant_flow_alias").(string),
		ResetCredentialsFlowAlias:     data.Get("reset_credentials_flow_alias").(string),
		ClientAuthenticationFlowAlias: data.Get("client_authentication_flow_alias").(string),
		DockerAuthenticationFlowAlias: data.Get("docker_authentication_flow_alias").(string),
	}, nil
}

func setAuthenticationBindings(data *schema.ResourceData, json *keycloak.AuthenticationBindings) {
	data.SetId(json.RealmId)
	data.Set("browser_flow_alias", json.BrowserFlowAlias)
	data.Set("registration_flow_alias", json.RegistrationFlowAlias)
	data.Set("direct_grant_flow_alias", json.DirectGrantFlowAlias)
	data.Set("reset_credentials_flow_alias", json.ResetCredentialsFlowAlias)
	data.Set("client_authentication_flow_alias", json.ClientAuthenticationFlowAlias)
	data.Set("docker_authentication_flow_alias", json.DockerAuthenticationFlowAlias)
}

func resourceKeycloakAuthenticationBindingsImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	d.Set("realm_id", d.Id())
	return []*schema.ResourceData{d}, nil
}
