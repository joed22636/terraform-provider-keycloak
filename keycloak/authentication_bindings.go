package keycloak

import (
	"fmt"
)

type AuthenticationBindings struct {
	RealmId                       string `json:"-"`
	BrowserFlowAlias              string `json:"browserFlow"`
	RegistrationFlowAlias         string `json:"registrationFlow"`
	DirectGrantFlowAlias          string `json:"directGrantFlow"`
	ResetCredentialsFlowAlias     string `json:"resetCredentialsFlow"`
	ClientAuthenticationFlowAlias string `json:"clientAuthenticationFlow"`
	DockerAuthenticationFlowAlias string `json:"dockerAuthenticationFlow"`
}

func (keycloakClient *KeycloakClient) GetAuthenticationBindings(realmId string) (*AuthenticationBindings, error) {
	var result AuthenticationBindings

	err := keycloakClient.get(fmt.Sprintf("/realms/%s", realmId), &result, nil)
	if err != nil {
		return nil, err
	}
	result.RealmId = realmId

	return &result, nil
}

func (keycloakClient *KeycloakClient) UpdateAuthenticationBindings(authenticationBindings *AuthenticationBindings) error {

	realm, err := keycloakClient.GetRealm(authenticationBindings.RealmId)
	if err != nil {
		return err
	}

	realm.BrowserFlow = authenticationBindings.BrowserFlowAlias
	realm.RegistrationFlow = authenticationBindings.RegistrationFlowAlias
	realm.DirectGrantFlow = authenticationBindings.DirectGrantFlowAlias
	realm.ResetCredentialsFlow = authenticationBindings.ResetCredentialsFlowAlias
	realm.ClientAuthenticationFlow = authenticationBindings.ClientAuthenticationFlowAlias
	realm.DockerAuthenticationFlow = authenticationBindings.DockerAuthenticationFlowAlias

	err = keycloakClient.UpdateRealm(realm)
	if err != nil {
		return err
	}

	return nil
}
