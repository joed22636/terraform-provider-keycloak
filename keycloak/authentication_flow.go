package keycloak

import (
	"errors"
	"fmt"
	"time"
)

type AuthenticationFlow struct {
	Id          string `json:"id,omitempty"`
	RealmId     string `json:"-"`
	Alias       string `json:"alias"`
	Description string `json:"description"`
	ProviderId  string `json:"providerId"` // "basic-flow" or "client-flow"
	TopLevel    bool   `json:"topLevel"`
	BuiltIn     bool   `json:"builtIn"`
}

func (keycloakClient *KeycloakClient) ListAuthenticationFlows(realmId string) ([]*AuthenticationFlow, error) {
	var authenticationFlows []*AuthenticationFlow

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/authentication/flows", realmId), &authenticationFlows, nil)
	if err != nil {
		return nil, err
	}

	for _, authenticationFlow := range authenticationFlows {
		authenticationFlow.RealmId = realmId
	}

	return authenticationFlows, nil
}

func (keycloakClient *KeycloakClient) NewAuthenticationFlow(authenticationFlow *AuthenticationFlow) error {
	authenticationFlow.TopLevel = true
	authenticationFlow.BuiltIn = false

	_, location, err := keycloakClient.post(fmt.Sprintf("/realms/%s/authentication/flows", authenticationFlow.RealmId), authenticationFlow)
	if err != nil {
		return err
	}
	authenticationFlow.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) GetAuthenticationFlow(realmId, id string) (*AuthenticationFlow, error) {
	var authenticationFlow AuthenticationFlow
	err := keycloakClient.get(fmt.Sprintf("/realms/%s/authentication/flows/%s", realmId, id), &authenticationFlow, nil)
	if err != nil {
		return nil, err
	}

	authenticationFlow.RealmId = realmId
	return &authenticationFlow, nil
}

func (keycloakClient *KeycloakClient) GetAuthenticationFlowFromAlias(realmId, alias string) (*AuthenticationFlow, error) {
	var authenticationFlows []*AuthenticationFlow
	var authenticationFlow *AuthenticationFlow = nil

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/authentication/flows", realmId), &authenticationFlows, nil)
	if err != nil {
		return nil, err
	}

	// Retry 3 more times if not found, sometimes it took split milliseconds the Authentication to populate
	if len(authenticationFlows) == 0 {
		for i := 0; i < 3; i++ {
			err := keycloakClient.get(fmt.Sprintf("/realms/%s/authentication/flows", realmId), &authenticationFlows, nil)

			if len(authenticationFlows) > 0 {
				break
			}

			if err != nil {
				return nil, err
			}

			time.Sleep(time.Millisecond * 50)
		}

		if len(authenticationFlows) == 0 {
			return nil, fmt.Errorf("no authentication flow found for alias %s", alias)
		}
	}

	for _, authFlow := range authenticationFlows {
		if authFlow.Alias == alias {
			authenticationFlow = authFlow
		}
	}

	if authenticationFlow == nil {
		return nil, fmt.Errorf("no authentication flow found for alias %s", alias)
	}
	authenticationFlow.RealmId = realmId

	return authenticationFlow, nil
}

func (keycloakClient *KeycloakClient) UpdateAuthenticationFlow(authenticationFlow *AuthenticationFlow) error {
	authenticationFlow.TopLevel = true
	authenticationFlow.BuiltIn = false

	return keycloakClient.put(fmt.Sprintf("/realms/%s/authentication/flows/%s", authenticationFlow.RealmId, authenticationFlow.Id), authenticationFlow)
}

func (keycloakClient *KeycloakClient) DeleteAuthenticationFlow(realmId, id string) error {
	err := keycloakClient.delete(fmt.Sprintf("/realms/%s/authentication/flows/%s", realmId, id), nil)
	if err != nil {
		// For whatever reason, this fails sometimes with a 500 during acceptance tests. try again
		return keycloakClient.delete(fmt.Sprintf("/realms/%s/authentication/flows/%s", realmId, id), nil)
	}
	return nil
}

func (keycloakClient *KeycloakClient) DeleteBuiltInFlowExecutors(flow *AuthenticationFlow) error {
	flow.BuiltIn = false
	err := keycloakClient.UpdateAuthenticationFlow(flow)
	if err != nil {
		return err
	}
	executors, err := keycloakClient.ListAuthenticationExecutions(flow.RealmId, flow.Alias)
	if err != nil {
		return err
	}
	for _, exe := range executors {
		if len(exe.FlowId) > 0 {
			var subFlow AuthenticationSubFlow
			err := keycloakClient.get(fmt.Sprintf("/realms/%s/authentication/flows/%s", flow.RealmId, exe.FlowId), &subFlow, nil)
			if err != nil {
				return err
			}
			subFlow.RealmId = flow.RealmId
			subFlow.BuiltIn = false
			err = keycloakClient.put(fmt.Sprintf("/realms/%s/authentication/flows/%s", subFlow.RealmId, subFlow.Id), subFlow)
			if err != nil {
				return err
			}
			err = keycloakClient.DeleteAuthenticationExecution(flow.RealmId, exe.Id)
			if err != nil && !isHttpError(err, 404) { // resource may be deleted with a subflow already
				return err
			}
		} else {
			err = keycloakClient.DeleteAuthenticationExecution(flow.RealmId, exe.Id)
			if err != nil && !isHttpError(err, 404) { // resource may be deleted with a subflow already
				return err
			}
		}
	}
	return nil
}

func isHttpError(err error, code int) bool {
	var apiError *ApiError
	if !errors.As(err, &apiError) {
		return false
	}
	if apiError.Code != code {
		return false
	}
	return true
}
