package keycloak

import (
	"fmt"
)

type DefaultRoles struct {
	RealmId      string   `json:"-"`
	DefaultRoles []string `json:"defaultRoles,omitempty"`
}

func (keycloakClient *KeycloakClient) GetDefaultRoles(realmId string) (*DefaultRoles, error) {
	var result DefaultRoles

	err := keycloakClient.get(fmt.Sprintf("/realms/%s", realmId), &result, nil)
	if err != nil {
		return nil, err
	}
	result.RealmId = realmId

	return &result, nil
}

func (keycloakClient *KeycloakClient) UpdateDefaultRoles(defaultRoles *DefaultRoles) error {

	realm, err := keycloakClient.GetRealm(defaultRoles.RealmId)
	if err != nil {
		return err
	}

	realm.DefaultRoles = defaultRoles.DefaultRoles

	err = keycloakClient.UpdateRealm(realm)
	if err != nil {
		return err
	}

	return nil
}
