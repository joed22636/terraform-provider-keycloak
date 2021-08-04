package provider

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/joed22636/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakComponent() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakComponentCreate,
		Read:   resourceKeycloakComponentRead,
		Update: resourceKeycloakComponentUpdate,
		Delete: resourceKeycloakComponentDelete,
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakComponentImport,
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
			"provider_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"provider_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"config": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func getComponentFromData(data *schema.ResourceData) *keycloak.Component {

	config := map[string][]string{}
	if v, ok := data.GetOk("config"); ok {
		for key, value := range v.(map[string][]string) {
			config[key] = value
		}
	}

	component := &keycloak.Component{
		Id:           data.Id(),
		Name:         data.Get("name").(string),
		ProviderId:   data.Get("provider_id").(string),
		ProviderType: data.Get("provider_type").(string),
		ParentId:     data.Get("parent_id").(string),
		Config:       config,
	}

	return component
}

func setComponentData(data *schema.ResourceData, component *keycloak.Component) {
	data.SetId(component.Id)

	data.Set("name", component.Name)
	data.Set("parent_id", component.ParentId)
	data.Set("provider_type", component.ProviderType)
	data.Set("provider_id", component.ProviderId)
	data.Set("config", component.Config)
}

func resourceKeycloakComponentCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	component := getComponentFromData(data)

	err := keycloakClient.CreateComponent(realmId, *component)
	if err != nil {
		return err
	}

	setComponentData(data, component)

	return resourceKeycloakComponentRead(data, meta)
}

func resourceKeycloakComponentRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	parentId := data.Get("parent_id").(string)
	componentType := data.Get("component_type").(string)
	name := data.Get("name").(string)

	components, err := keycloakClient.GetComponents(realmId, parentId, componentType)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	var comp *keycloak.Component = nil
	for _, c := range components {
		if c.Name == name {
			comp = &c
			break
		}
	}
	if comp == nil {
		return fmt.Errorf("Component could not be found (realm, parent, type, name): %v %v %v %v", realmId, parentId, componentType, name)
	}

	setComponentData(data, comp)

	return nil
}

func resourceKeycloakComponentUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	component := getComponentFromData(data)

	err := keycloakClient.UpdateComponent(realmId, *component)
	if err != nil {
		return err
	}

	setComponentData(data, component)

	return nil
}

func resourceKeycloakComponentDelete(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	return keycloakClient.DeleteComponent(realmId, id)
}

func resourceKeycloakComponentImport(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	var realmId, id string
	switch {
	case len(parts) == 2:
		realmId = parts[0]
		id = parts[1]
	default:
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{componentId}}")
	}

	d.Set("realm_id", realmId)
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
