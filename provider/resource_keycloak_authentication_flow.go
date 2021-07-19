package provider

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mrparkers/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakAuthenticationFlow() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakAuthenticationFlowCreate,
		Read:   resourceKeycloakAuthenticationFlowRead,
		Delete: resourceKeycloakAuthenticationFlowDelete,
		Update: resourceKeycloakAuthenticationFlowUpdate,
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakAuthenticationFlowImport,
		},
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"alias": {
				Type:     schema.TypeString,
				Required: true,
			},
			"provider_id": {
				Type:         schema.TypeString,
				Default:      "basic-flow",
				ValidateFunc: validation.StringInSlice([]string{"basic-flow", "client-flow"}, false), //it seems toplevel can only one of these and not 'form-flow'
				Optional:     true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func mapFromDataToAuthenticationFlow(data *schema.ResourceData) *keycloak.AuthenticationFlow {
	authenticationFlow := &keycloak.AuthenticationFlow{
		Id:          data.Id(),
		RealmId:     data.Get("realm_id").(string),
		Alias:       data.Get("alias").(string),
		ProviderId:  data.Get("provider_id").(string),
		Description: data.Get("description").(string),
	}

	return authenticationFlow
}

func mapFromAuthenticationFlowToData(data *schema.ResourceData, authenticationFlow *keycloak.AuthenticationFlow) {
	data.SetId(authenticationFlow.Id)
	data.Set("realm_id", authenticationFlow.RealmId)
	data.Set("alias", authenticationFlow.Alias)
	data.Set("provider_id", authenticationFlow.ProviderId)
	data.Set("description", authenticationFlow.Description)
}

func mapFromAuthenticationFlowInfoToData(data *schema.ResourceData, authenticationFlow *keycloak.AuthenticationFlow) {
	data.SetId(authenticationFlow.Id)
	data.Set("realm_id", authenticationFlow.RealmId)
	data.Set("alias", authenticationFlow.Alias)
}

func resourceKeycloakAuthenticationFlowCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	authenticationFlow := mapFromDataToAuthenticationFlow(data)

	err := keycloakClient.NewAuthenticationFlow(authenticationFlow)

	if err != nil {
		var ae *keycloak.ApiError
		if !errors.As(err, &ae) {
			return err
		}

		if ae.Code != 409 {
			return err
		}

		log.Println("flow already exists, may be hardcoded flow, try to update")
		flow, err := keycloakClient.GetAuthenticationFlowFromAlias(authenticationFlow.RealmId, authenticationFlow.Alias)
		if err != nil {
			return err
		}
		data.SetId(flow.Id)
		authenticationFlow.Id = flow.Id
		err = resourceKeycloakAuthenticationFlowUpdate(data, meta)
		if err != nil {
			return err
		}
		err = keycloakClient.DeleteBuiltInFlowExecutors(authenticationFlow)
		if err != nil {
			return err
		}
	}

	mapFromAuthenticationFlowToData(data, authenticationFlow)
	return resourceKeycloakAuthenticationFlowRead(data, meta)
}

func resourceKeycloakAuthenticationFlowRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	authenticationFlow, err := keycloakClient.GetAuthenticationFlow(realmId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	mapFromAuthenticationFlowToData(data, authenticationFlow)
	return nil
}

func resourceKeycloakAuthenticationFlowUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	authenticationFlow := mapFromDataToAuthenticationFlow(data)

	err := keycloakClient.UpdateAuthenticationFlow(authenticationFlow)
	if err != nil {
		return err
	}

	mapFromAuthenticationFlowToData(data, authenticationFlow)
	return nil
}

func resourceKeycloakAuthenticationFlowDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	err := keycloakClient.DeleteAuthenticationFlow(realmId, id)

	if err != nil {
		var ae *keycloak.ApiError
		if !errors.As(err, &ae) {
			return err
		}

		if ae.Code != 400 && ae.Code != 500 {
			return err
		}

		if !strings.Contains(ae.Message, "Can't delete built in flow") &&
			!strings.Contains(ae.Message, "unknown_error") { // unknown_error is common when the flow is still referenced by a client or realm binding
			return err
		}

		log.Println("flow cannot be deleted, ignore error (probably built-in flow or still referenced by the realm or a client application)")
	}

	return nil
}

func resourceKeycloakAuthenticationFlowImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{authenticationFlowId}}")
	}

	d.Set("realm_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
