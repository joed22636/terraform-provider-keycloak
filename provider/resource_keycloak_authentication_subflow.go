package provider

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/joed22636/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakAuthenticationSubFlow() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakAuthenticationSubFlowCreate,
		Read:   resourceKeycloakAuthenticationSubFlowRead,
		Delete: resourceKeycloakAuthenticationSubFlowDelete,
		Update: resourceKeycloakAuthenticationSubFlowUpdate,
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakAuthenticationSubFlowImport,
		},
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"parent_flow_alias": {
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
				ValidateFunc: validation.StringInSlice([]string{"basic-flow", "form-flow", "client-flow"}, false),
				Optional:     true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			//execution parts of the subflow
			"authenticator": {
				Type:        schema.TypeString,
				Description: "Might be needed to be set with certain custom subflow with specific authenticator, in general this will remain empty",
				Optional:    true,
				ForceNew:    true,
			},
			"requirement": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"REQUIRED", "ALTERNATIVE", "OPTIONAL", "CONDITIONAL", "DISABLED"}, false), //OPTIONAL is removed from 8.0.0 onwards
				Default:      "DISABLED",
			},
		},
	}
}

func mapFromDataToAuthenticationSubFlow(data *schema.ResourceData) *keycloak.AuthenticationSubFlow {
	authenticationSubFlow := &keycloak.AuthenticationSubFlow{
		Id:              data.Id(),
		RealmId:         data.Get("realm_id").(string),
		ParentFlowAlias: data.Get("parent_flow_alias").(string),
		Alias:           data.Get("alias").(string),
		ProviderId:      data.Get("provider_id").(string),
		Description:     data.Get("description").(string),
		Authenticator:   data.Get("authenticator").(string),
		Requirement:     data.Get("requirement").(string),
	}

	return authenticationSubFlow
}

func mapFromAuthenticationSubFlowToData(data *schema.ResourceData, authenticationSubFlow *keycloak.AuthenticationSubFlow) {
	data.SetId(authenticationSubFlow.Id)
	data.Set("realm_id", authenticationSubFlow.RealmId)
	data.Set("parent_flow_alias", authenticationSubFlow.ParentFlowAlias)
	data.Set("alias", authenticationSubFlow.Alias)
	data.Set("provider_id", authenticationSubFlow.ProviderId)
	data.Set("description", authenticationSubFlow.Description)
	data.Set("authenticator", authenticationSubFlow.Authenticator)
	data.Set("requirement", authenticationSubFlow.Requirement)
}

func resourceKeycloakAuthenticationSubFlowCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	authenticationFlow := mapFromDataToAuthenticationSubFlow(data)

	err := keycloakClient.NewAuthenticationSubFlow(authenticationFlow)
	if err != nil {
		var ae *keycloak.ApiError
		if !errors.As(err, &ae) {
			return err
		}

		if ae.Code != 409 {
			return err
		}

		log.Println("subflow already exists, may be hardcoded flow, try to update")
		flow, err := keycloakClient.GetAuthenticationSubFlowByAlias(authenticationFlow.RealmId, authenticationFlow.ParentFlowAlias, authenticationFlow.Alias)
		if err != nil {
			return err
		}
		if flow == nil {
			return fmt.Errorf("'%v' cannot be used as subflow alias. Please choose another one", authenticationFlow.Alias)
		}
		data.SetId(flow.Id)
		authenticationFlow.Id = flow.Id
		err = resourceKeycloakAuthenticationSubFlowUpdate(data, meta)
		if err != nil {
			return err
		}
	}
	mapFromAuthenticationSubFlowToData(data, authenticationFlow)
	return resourceKeycloakAuthenticationSubFlowRead(data, meta)
}

func resourceKeycloakAuthenticationSubFlowRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	parentFlowAlias := data.Get("parent_flow_alias").(string)
	id := data.Id()

	authenticationFlow, err := keycloakClient.GetAuthenticationSubFlow(realmId, parentFlowAlias, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}
	mapFromAuthenticationSubFlowToData(data, authenticationFlow)
	return nil
}

func resourceKeycloakAuthenticationSubFlowUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	authenticationFlow := mapFromDataToAuthenticationSubFlow(data)

	err := keycloakClient.UpdateAuthenticationSubFlow(authenticationFlow)
	if err != nil {
		return err
	}
	mapFromAuthenticationSubFlowToData(data, authenticationFlow)
	return nil
}

func resourceKeycloakAuthenticationSubFlowDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	parentFlowAlias := data.Get("parent_flow_alias").(string)
	id := data.Id()

	err := keycloakClient.DeleteAuthenticationSubFlow(realmId, parentFlowAlias, id)
	if err != nil {
		if err != nil {
			var ae *keycloak.ApiError
			if !errors.As(err, &ae) {
				return err
			}

			if ae.Code != 400 {
				return err
			}

			if !strings.Contains(ae.Message, "It is illegal to remove execution from a built in flow") {
				return err
			}

			log.Println("build-in flows cannot be deleted, ignore error")
		}
	}
	return nil
}

func resourceKeycloakAuthenticationSubFlowImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 3 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{parentFlowAlias}}/{{authenticationSubFlowId}}")
	}

	d.Set("realm_id", parts[0])
	d.Set("parent_flow_alias", parts[1])
	d.SetId(parts[2])

	return []*schema.ResourceData{d}, nil
}
