package keycloak

import (
	"fmt"
	"strings"
)

type Key struct {
	Algorithm        *string `json:"algorithm,omitempty"`
	Certificate      *string `json:"certificate,omitempty"`
	ProviderId       *string `json:"providerId,omitempty"`
	ProviderPriority *int    `json:"providerPriority,omitempty"`
	PublicKey        *string `json:"publicKey,omitempty"`
	Kid              *string `json:"kid,omitempty"`
	Status           *string `json:"status,omitempty"`
	Type             *string `json:"type,omitempty"`
}

type Keys struct {
	Keys []Key `json:"keys"`
}

type Realm struct {
	Id                string `json:"id,omitempty"`
	Realm             string `json:"realm"`
	Enabled           bool   `json:"enabled"`
	DisplayName       string `json:"displayName"`
	DisplayNameHtml   string `json:"displayNameHtml"`
	UserManagedAccess bool   `json:"userManagedAccessAllowed"`

	// Login Config
	RegistrationAllowed         bool   `json:"registrationAllowed"`
	RegistrationEmailAsUsername bool   `json:"registrationEmailAsUsername"`
	EditUsernameAllowed         bool   `json:"editUsernameAllowed"`
	ResetPasswordAllowed        bool   `json:"resetPasswordAllowed"`
	RememberMe                  bool   `json:"rememberMe"`
	VerifyEmail                 bool   `json:"verifyEmail"`
	LoginWithEmailAllowed       bool   `json:"loginWithEmailAllowed"`
	DuplicateEmailsAllowed      bool   `json:"duplicateEmailsAllowed"`
	SslRequired                 string `json:"sslRequired,omitempty"`

	//SMTP Server
	SmtpServer SmtpServer `json:"smtpServer"`

	// Themes
	LoginTheme   string `json:"loginTheme,omitempty"`
	AccountTheme string `json:"accountTheme,omitempty"`
	AdminTheme   string `json:"adminTheme,omitempty"`
	EmailTheme   string `json:"emailTheme,omitempty"`

	// Tokens
	DefaultSignatureAlgorithm           string `json:"defaultSignatureAlgorithm"`
	RevokeRefreshToken                  bool   `json:"revokeRefreshToken"`
	RefreshTokenMaxReuse                int    `json:"refreshTokenMaxReuse"`
	SsoSessionIdleTimeout               int    `json:"ssoSessionIdleTimeout,omitempty"`
	SsoSessionMaxLifespan               int    `json:"ssoSessionMaxLifespan,omitempty"`
	SsoSessionIdleTimeoutRememberMe     int    `json:"ssoSessionIdleTimeoutRememberMe,omitempty"`
	SsoSessionMaxLifespanRememberMe     int    `json:"ssoSessionMaxLifespanRememberMe,omitempty"`
	OfflineSessionIdleTimeout           int    `json:"offlineSessionIdleTimeout,omitempty"`
	OfflineSessionMaxLifespan           int    `json:"offlineSessionMaxLifespan,omitempty"`
	OfflineSessionMaxLifespanEnabled    bool   `json:"offlineSessionMaxLifespanEnabled,omitempty"`
	AccessTokenLifespan                 int    `json:"accessTokenLifespan,omitempty"`
	AccessTokenLifespanForImplicitFlow  int    `json:"accessTokenLifespanForImplicitFlow,omitempty"`
	AccessCodeLifespan                  int    `json:"accessCodeLifespan,omitempty"`
	AccessCodeLifespanLogin             int    `json:"accessCodeLifespanLogin,omitempty"`
	AccessCodeLifespanUserAction        int    `json:"accessCodeLifespanUserAction,omitempty"`
	ActionTokenGeneratedByUserLifespan  int    `json:"actionTokenGeneratedByUserLifespan,omitempty"`
	ActionTokenGeneratedByAdminLifespan int    `json:"actionTokenGeneratedByAdminLifespan,omitempty"`

	//internationalization
	InternationalizationEnabled bool     `json:"internationalizationEnabled"`
	SupportLocales              []string `json:"supportedLocales"`
	DefaultLocale               string   `json:"defaultLocale"`

	//extra attributes of a realm
	Attributes map[string]interface{} `json:"attributes"`

	// client-scope mapping defaults
	DefaultDefaultClientScopes  []string `json:"defaultDefaultClientScopes,omitempty"`
	DefaultOptionalClientScopes []string `json:"defaultOptionalClientScopes,omitempty"`

	DefaultRoles []string `json:"defaultRoles,omitempty"`

	BrowserSecurityHeaders BrowserSecurityHeaders `json:"browserSecurityHeaders"`

	BruteForceProtected          bool `json:"bruteForceProtected"`
	PermanentLockout             bool `json:"permanentLockout"`
	FailureFactor                int  `json:"failureFactor"` //Max Login Failures
	WaitIncrementSeconds         int  `json:"waitIncrementSeconds"`
	QuickLoginCheckMilliSeconds  int  `json:"quickLoginCheckMilliSeconds"`
	MinimumQuickLoginWaitSeconds int  `json:"minimumQuickLoginWaitSeconds"`
	MaxFailureWaitSeconds        int  `json:"maxFailureWaitSeconds"` //Max Wait
	MaxDeltaTimeSeconds          int  `json:"maxDeltaTimeSeconds"`   //Failure Reset Time

	PasswordPolicy string `json:"passwordPolicy"`

	//flow bindings
	BrowserFlow              string `json:"browserFlow,omitempty"`
	RegistrationFlow         string `json:"registrationFlow,omitempty"`
	DirectGrantFlow          string `json:"directGrantFlow,omitempty"`
	ResetCredentialsFlow     string `json:"resetCredentialsFlow,omitempty"`
	ClientAuthenticationFlow string `json:"clientAuthenticationFlow,omitempty"`
	DockerAuthenticationFlow string `json:"dockerAuthenticationFlow,omitempty"`

	// WebAuthn
	WebAuthnPolicyAcceptableAaguids               []string `json:"webAuthnPolicyAcceptableAaguids"`
	WebAuthnPolicyAttestationConveyancePreference string   `json:"webAuthnPolicyAttestationConveyancePreference"`
	WebAuthnPolicyAuthenticatorAttachment         string   `json:"webAuthnPolicyAuthenticatorAttachment"`
	WebAuthnPolicyAvoidSameAuthenticatorRegister  bool     `json:"webAuthnPolicyAvoidSameAuthenticatorRegister"`
	WebAuthnPolicyCreateTimeout                   int      `json:"webAuthnPolicyCreateTimeout"`
	WebAuthnPolicyRequireResidentKey              string   `json:"webAuthnPolicyRequireResidentKey"`
	WebAuthnPolicyRpEntityName                    string   `json:"webAuthnPolicyRpEntityName"`
	WebAuthnPolicyRpId                            string   `json:"webAuthnPolicyRpId"`
	WebAuthnPolicySignatureAlgorithms             []string `json:"webAuthnPolicySignatureAlgorithms"`
	WebAuthnPolicyUserVerificationRequirement     string   `json:"webAuthnPolicyUserVerificationRequirement"`

	// WebAuthn Passwordless
	WebAuthnPolicyPasswordlessAcceptableAaguids               []string `json:"webAuthnPolicyPasswordlessAcceptableAaguids"`
	WebAuthnPolicyPasswordlessAttestationConveyancePreference string   `json:"webAuthnPolicyPasswordlessAttestationConveyancePreference"`
	WebAuthnPolicyPasswordlessAuthenticatorAttachment         string   `json:"webAuthnPolicyPasswordlessAuthenticatorAttachment"`
	WebAuthnPolicyPasswordlessAvoidSameAuthenticatorRegister  bool     `json:"webAuthnPolicyPasswordlessAvoidSameAuthenticatorRegister"`
	WebAuthnPolicyPasswordlessCreateTimeout                   int      `json:"webAuthnPolicyPasswordlessCreateTimeout"`
	WebAuthnPolicyPasswordlessRequireResidentKey              string   `json:"webAuthnPolicyPasswordlessRequireResidentKey"`
	WebAuthnPolicyPasswordlessRpEntityName                    string   `json:"webAuthnPolicyPasswordlessRpEntityName"`
	WebAuthnPolicyPasswordlessRpId                            string   `json:"webAuthnPolicyPasswordlessRpId"`
	WebAuthnPolicyPasswordlessSignatureAlgorithms             []string `json:"webAuthnPolicyPasswordlessSignatureAlgorithms"`
	WebAuthnPolicyPasswordlessUserVerificationRequirement     string   `json:"webAuthnPolicyPasswordlessUserVerificationRequirement"`
}

type BrowserSecurityHeaders struct {
	ContentSecurityPolicy           string `json:"contentSecurityPolicy"`
	ContentSecurityPolicyReportOnly string `json:"contentSecurityPolicyReportOnly"`
	StrictTransportSecurity         string `json:"strictTransportSecurity"`
	XContentTypeOptions             string `json:"xContentTypeOptions"`
	XFrameOptions                   string `json:"xFrameOptions"`
	XRobotsTag                      string `json:"xRobotsTag"`
	XXSSProtection                  string `json:"xXSSProtection"`
}

type SmtpServer struct {
	StartTls           KeycloakBoolQuoted `json:"starttls,omitempty"`
	Auth               KeycloakBoolQuoted `json:"auth,omitempty"`
	Port               string             `json:"port,omitempty"`
	Host               string             `json:"host,omitempty"`
	ReplyTo            string             `json:"replyTo,omitempty"`
	ReplyToDisplayName string             `json:"replyToDisplayName,omitempty"`
	From               string             `json:"from,omitempty"`
	FromDisplayName    string             `json:"fromDisplayName,omitempty"`
	EnvelopeFrom       string             `json:"envelopeFrom,omitempty"`
	Ssl                KeycloakBoolQuoted `json:"ssl,omitempty"`
	User               string             `json:"user,omitempty"`
	Password           string             `json:"password,omitempty"`
}

type DefaultClientScope struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type DefaultClientScopes struct {
	Default  []DefaultClientScope
	Optional []DefaultClientScope
}

func (keycloakClient *KeycloakClient) NewRealm(realm *Realm) error {
	_, _, err := keycloakClient.post("/realms", realm)
	if err != nil {
		return err
	}

	return keycloakClient.updateDefaultClientScopes(realm)
}

func (keycloakClient *KeycloakClient) GetRealm(name string) (*Realm, error) {
	var realm Realm

	err := keycloakClient.get(fmt.Sprintf("/realms/%s", name), &realm, nil)
	if err != nil {
		return nil, err
	}

	err = keycloakClient.fillInRealDefaultClientScopes(&realm)
	if err != nil {
		return nil, err
	}

	return &realm, nil
}

func (keycloakClient *KeycloakClient) getCurrentDefaultClientScopes(realm string) (DefaultClientScopes, error) {
	var result DefaultClientScopes

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/default-default-client-scopes", realm), &(result.Default), nil)
	if err != nil {
		return result, err
	}

	err = keycloakClient.get(fmt.Sprintf("/realms/%s/default-optional-client-scopes", realm), &(result.Optional), nil)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (keycloakClient *KeycloakClient) getDefaultClientScopesByDefinition(realm *Realm) (DefaultClientScopes, error) {
	var result DefaultClientScopes
	var allDefaultClientScopes []DefaultClientScope

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/client-scopes", realm.Realm), &allDefaultClientScopes, nil)
	if err != nil {
		return result, err
	}

	for _, scope := range allDefaultClientScopes {
		if contains(realm.DefaultDefaultClientScopes, scope.Name) {
			result.Default = append(result.Default, scope)
		} else if contains(realm.DefaultOptionalClientScopes, scope.Name) {
			result.Optional = append(result.Optional, scope)
		}
	}

	return result, nil
}

func (keycloakClient *KeycloakClient) fillInRealDefaultClientScopes(realm *Realm) error {

	dcs, err := keycloakClient.getCurrentDefaultClientScopes(realm.Realm)
	if err != nil {
		return err
	}

	for _, docs := range dcs.Default {
		realm.DefaultDefaultClientScopes = append(realm.DefaultDefaultClientScopes, docs.Name)
	}

	for _, docs := range dcs.Optional {
		realm.DefaultOptionalClientScopes = append(realm.DefaultOptionalClientScopes, docs.Name)
	}

	return nil
}

func (keycloakClient *KeycloakClient) GetRealms() ([]*Realm, error) {
	var realms []*Realm

	err := keycloakClient.get("/realms", &realms, nil)
	if err != nil {
		return nil, err
	}

	for _, realm := range realms {
		err = keycloakClient.fillInRealDefaultClientScopes(realm)
		if err != nil {
			return nil, err
		}
	}

	return realms, nil
}

func (keycloakClient *KeycloakClient) GetRealmKeys(name string) (*Keys, error) {
	var keys Keys

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/keys", name), &keys, nil)
	if err != nil {
		return nil, err
	}

	return &keys, nil
}

func (keycloakClient *KeycloakClient) UpdateRealm(realm *Realm) error {
	err := keycloakClient.put(fmt.Sprintf("/realms/%s", realm.Realm), realm)
	if err != nil {
		return err
	}
	keycloakClient.updateDefaultClientScopes(realm)
	return nil
}

func (keycloakClient *KeycloakClient) updateDefaultClientScopes(realm *Realm) error {
	dcss, err := keycloakClient.getCurrentDefaultClientScopes(realm.Realm)
	if err != nil {
		return err
	}
	keycloakClient.removeDefaultClientScopesFromRealm(realm, dcss)
	dcss, err = keycloakClient.getDefaultClientScopesByDefinition(realm)
	if err != nil {
		return err
	}
	keycloakClient.assignDefaultClientScopesToRealm(realm, dcss)

	return nil
}

func (keycloakClient *KeycloakClient) removeDefaultClientScopesFromRealm(realm *Realm, dcss DefaultClientScopes) error {
	for _, dcs := range dcss.Default {
		err := keycloakClient.DeleteDefaultDefaultClientScope(realm.Realm, dcs.Id)
		if err != nil {
			return err
		}
	}
	for _, dcs := range dcss.Optional {
		err := keycloakClient.DeleteOptionalDefaultClientScope(realm.Realm, dcs.Id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (keycloakClient *KeycloakClient) assignDefaultClientScopesToRealm(realm *Realm, dcss DefaultClientScopes) error {
	for _, scope := range dcss.Default {
		err := keycloakClient.AddDefaultDefaultClientScope(realm.Realm, scope)
		if err != nil {
			return err
		}
	}
	for _, scope := range dcss.Optional {
		err := keycloakClient.AddOptionalDefaultClientScope(realm.Realm, scope)
		if err != nil {
			return err
		}
	}
	return nil
}

func (keycloakClient *KeycloakClient) DeleteRealm(name string) error {
	err := keycloakClient.delete(fmt.Sprintf("/realms/%s", name), nil)
	if err != nil {
		// For whatever reason, this fails sometimes with a 500 during acceptance tests. try again
		return keycloakClient.delete(fmt.Sprintf("/realms/%s", name), nil)
	}

	return nil
}

func (keycloakClient *KeycloakClient) ValidateRealm(realm *Realm) error {
	if realm.DuplicateEmailsAllowed == true && realm.RegistrationEmailAsUsername == true {
		return fmt.Errorf("validation error: DuplicateEmailsAllowed cannot be true if RegistrationEmailAsUsername is true")
	}

	if realm.DuplicateEmailsAllowed == true && realm.LoginWithEmailAllowed == true {
		return fmt.Errorf("validation error: DuplicateEmailsAllowed cannot be true if LoginWithEmailAllowed is true")
	}

	if realm.SslRequired != "none" && realm.SslRequired != "external" && realm.SslRequired != "all" {
		return fmt.Errorf("validation error: SslRequired should be 'none', 'external' or 'all'")
	}

	// validate if the given theme exists on the server. the keycloak API allows you to use any random string for a theme
	serverInfo, err := keycloakClient.GetServerInfo()
	if err != nil {
		return err
	}

	if realm.LoginTheme != "" && !serverInfo.ThemeIsInstalled("login", realm.LoginTheme) {
		return fmt.Errorf("validation error: theme \"%s\" does not exist on the server", realm.LoginTheme)
	}

	if realm.AccountTheme != "" && !serverInfo.ThemeIsInstalled("account", realm.AccountTheme) {
		return fmt.Errorf("validation error: theme \"%s\" does not exist on the server", realm.AccountTheme)
	}

	if realm.AdminTheme != "" && !serverInfo.ThemeIsInstalled("admin", realm.AdminTheme) {
		return fmt.Errorf("validation error: theme \"%s\" does not exist on the server", realm.AdminTheme)
	}

	if realm.EmailTheme != "" && !serverInfo.ThemeIsInstalled("email", realm.EmailTheme) {
		return fmt.Errorf("validation error: theme \"%s\" does not exist on the server", realm.EmailTheme)
	}

	if realm.InternationalizationEnabled == true && !contains(realm.SupportLocales, realm.DefaultLocale) {
		return fmt.Errorf("validation error: DefaultLocale should be in the SupportLocales")
	}

	if realm.PasswordPolicy != "" {
		policies := strings.Split(realm.PasswordPolicy, " and ")
		for _, policyTypeRepresentation := range policies {
			policy := strings.Split(policyTypeRepresentation, "(")
			if !serverInfo.providerInstalled("password-policy", policy[0]) {
				return fmt.Errorf("validation error: password-policy \"%s\" does not exist on the server, installed providers: %s", policy[0], serverInfo.getInstalledProvidersNames("password-policy"))
			}
		}
	}

	return nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
