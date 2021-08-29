package provider

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/joed22636/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakOpenidClientScope() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakOpenidClientScopeCreate,
		Read:   resourceKeycloakOpenidClientScopeRead,
		Delete: resourceKeycloakOpenidClientScopeDelete,
		Update: resourceKeycloakOpenidClientScopeUpdate,
		// This resource can be imported using {{realm}}/{{client_scope_id}}. The Client Scope ID is displayed in the GUI
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakOpenidClientScopeImport,
		},
		Schema: map[string]*schema.Schema{
			"realm_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"consent_screen_text": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"include_in_token_scope": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"gui_order": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func getClientScopeFromData(data *schema.ResourceData) *keycloak.OpenidClientScope {
	clientScope := &keycloak.OpenidClientScope{
		Id:          data.Id(),
		RealmId:     data.Get("realm_id").(string),
		Name:        data.Get("name").(string),
		Description: data.Get("description").(string),
	}

	if consentScreenText, ok := data.GetOk("consent_screen_text"); ok {
		clientScope.Attributes.ConsentScreenText = consentScreenText.(string)
		clientScope.Attributes.DisplayOnConsentScreen = true
	} else {
		clientScope.Attributes.DisplayOnConsentScreen = false
	}

	clientScope.Attributes.IncludeInTokenScope = strconv.FormatBool(data.Get("include_in_token_scope").(bool))

	// Treat 0 as an empty string for the purpose of omitting the attribute to reset the order
	if guiOrder := data.Get("gui_order").(int); guiOrder != 0 {
		clientScope.Attributes.GuiOrder = strconv.Itoa(guiOrder)
	}

	return clientScope
}

func setClientScopeData(data *schema.ResourceData, clientScope *keycloak.OpenidClientScope) {
	data.SetId(clientScope.Id)

	data.Set("realm_id", clientScope.RealmId)
	data.Set("name", clientScope.Name)
	data.Set("description", clientScope.Description)

	if clientScope.Attributes.DisplayOnConsentScreen {
		data.Set("consent_screen_text", clientScope.Attributes.ConsentScreenText)
	}

	its, err := strconv.ParseBool(clientScope.Attributes.IncludeInTokenScope)
	if err != nil {
		data.Set("include_in_token_scope", true) // default value
	} else {
		data.Set("include_in_token_scope", its)
	}

	if guiOrder, err := strconv.Atoi(clientScope.Attributes.GuiOrder); err == nil {
		data.Set("gui_order", guiOrder)
	}
}

func resourceKeycloakOpenidClientScopeCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	clientScope := getClientScopeFromData(data)

	err := keycloakClient.NewOpenidClientScope(clientScope)
	if err != nil {
		var ae *keycloak.ApiError
		if !errors.As(err, &ae) {
			return err
		}

		if ae.Code != 409 {
			return err
		}

		log.Println("resource already exists, may be hardcoded, try to update")
		sc, err1 := keycloakClient.GetOpenIdClientScopeByName(clientScope.RealmId, clientScope.Name)
		if err1 != nil {
			return err1
		}
		data.SetId(sc.Id)
		clientScope.Id = sc.Id
		err = resourceKeycloakOpenidClientScopeUpdate(data, meta)
		if err != nil {
			return err
		}
	}

	setClientScopeData(data, clientScope)

	return resourceKeycloakOpenidClientScopeRead(data, meta)
}

func resourceKeycloakOpenidClientScopeRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	clientScope, err := keycloakClient.GetOpenidClientScope(realmId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	setClientScopeData(data, clientScope)

	return nil
}

func resourceKeycloakOpenidClientScopeUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	clientScope := getClientScopeFromData(data)

	err := keycloakClient.UpdateOpenidClientScope(clientScope)
	if err != nil {
		return err
	}

	setClientScopeData(data, clientScope)

	return nil
}

func resourceKeycloakOpenidClientScopeDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	err := keycloakClient.DeleteOpenidClientScope(realmId, id)
	if err != nil {
		var ae *keycloak.ApiError
		if !errors.As(err, &ae) {
			return err
		}

		if ae.Code != 400 {
			return err
		}

		if !strings.Contains(ae.Message, "Cannot remove client scope, it is currently in use") {
			return err
		}

		log.Println("resource cannot be deleted, ignore error (probably built-in resource and/or still referenced by another resource)")
	}
	return nil
}

func resourceKeycloakOpenidClientScopeImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{openidClientScopeId}}")
	}

	d.Set("realm_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
