package keycloak

import "fmt"

// https://www.keycloak.org/docs-api/4.2/rest-api/index.html#_component_resource

type Component struct {
	Id           string              `json:"id,omitempty"`
	Name         string              `json:"name"`
	ProviderId   string              `json:"providerId"`
	ProviderType string              `json:"providerType"`
	ParentId     string              `json:"parentId"`
	Config       map[string][]string `json:"config"`
}

func (component *Component) getConfig(val string) string {
	if len(component.Config[val]) == 0 {
		return ""
	}

	return component.Config[val][0]
}

func (component *Component) getConfigOk(val string) (string, bool) {
	if v, ok := component.Config[val]; ok {
		return v[0], true
	}

	return "", false
}

func (keycloakClient *KeycloakClient) GetComponents(realm, parent, providerType string) ([]Component, error) {
	result := []Component{}
	err := keycloakClient.get(fmt.Sprintf("/realms/%s/components?parent=%s&type=%s", realm, parent, providerType), &result, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (keycloakClient *KeycloakClient) CreateComponent(realm string, component Component) error {
	_, _, err := keycloakClient.post(fmt.Sprintf("/realms/%s/components", realm), component)
	if err != nil {
		return err
	}

	return nil
}

func (keycloakClient *KeycloakClient) DeleteComponent(realmId, id string) error {
	return keycloakClient.delete(fmt.Sprintf("/realms/%s/components/%s", realmId, id), nil)
}
