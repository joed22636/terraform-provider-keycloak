package keycloak

import (
	"fmt"
	"log"
)

type IdentityProviderMapper struct {
	Realm                  string                 `json:"-"`
	Provider               string                 `json:"-"`
	Id                     string                 `json:"id,omitempty"`
	Name                   string                 `json:"name,omitempty"`
	IdentityProviderAlias  string                 `json:"identityProviderAlias,omitempty"`
	IdentityProviderMapper string                 `json:"identityProviderMapper,omitempty"`
	Config                 map[string]interface{} `json:"config,omitempty"`
}

func (keycloakClient *KeycloakClient) NewIdentityProviderMapper(identityProviderMapper *IdentityProviderMapper) error {
	log.Printf("[WARN] Realm: %s", identityProviderMapper.Realm)
	_, location, err := keycloakClient.post(fmt.Sprintf("/realms/%s/identity-provider/instances/%s/mappers", identityProviderMapper.Realm, identityProviderMapper.IdentityProviderAlias), identityProviderMapper)
	if err != nil {
		return err
	}

	identityProviderMapper.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) GetIdentityProviderMapper(realm, idpAlias, id string) (*IdentityProviderMapper, error) {
	var identityProviderMapper IdentityProviderMapper
	identityProviderMapper.Realm = realm
	identityProviderMapper.IdentityProviderAlias = idpAlias

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/identity-provider/instances/%s/mappers/%s", realm, idpAlias, id), &identityProviderMapper, nil)
	if err != nil {
		return nil, err
	}

	return &identityProviderMapper, nil
}

func (keycloakClient *KeycloakClient) UpdateIdentityProviderMapper(identityProviderMapper *IdentityProviderMapper) error {
	return keycloakClient.put(fmt.Sprintf("/realms/%s/identity-provider/instances/%s/mappers/%s", identityProviderMapper.Realm, identityProviderMapper.IdentityProviderAlias, identityProviderMapper.Id), identityProviderMapper)
}

func (keycloakClient *KeycloakClient) DeleteIdentityProviderMapper(realm, idpAlias, id string) error {
	return keycloakClient.delete(fmt.Sprintf("/realms/%s/identity-provider/instances/%s/mappers/%s", realm, idpAlias, id), nil)
}
