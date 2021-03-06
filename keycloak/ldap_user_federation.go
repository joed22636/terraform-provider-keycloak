package keycloak

import (
	"fmt"
	"strconv"
	"strings"
)

type LdapUserFederation struct {
	Id      string
	Name    string
	RealmId string

	Enabled  bool
	Priority int

	ImportEnabled     bool
	EditMode          string // can be "READ_ONLY", "WRITABLE", or "UNSYNCED"
	SyncRegistrations bool   // I think this field controls whether or not BatchSizeForSync, FullSyncPeriod, and ChangedSyncPeriod are needed

	Vendor                 string // can be "other", "edirectory", "ad", "rhds", or "tivoli". honestly I don't think this field actually does anything
	UsernameLDAPAttribute  string
	RdnLDAPAttribute       string
	UuidLDAPAttribute      string
	UserObjectClasses      []string // api expects comma + space separated for some reason
	ConnectionUrl          string
	UsersDn                string
	BindDn                 string
	BindCredential         string
	CustomUserSearchFilter string // must start with '(' and end with ')'
	SearchScope            string // api expects "1" or "2", but that means "One Level" or "Subtree"

	ValidatePasswordPolicy bool
	TrustEmail             bool
	UseTruststoreSpi       string // can be "ldapsOnly", "always", or "never"
	ConnectionTimeout      string // duration string (ex: 1h30m)
	ReadTimeout            string // duration string (ex: 1h30m)
	Pagination             bool

	ConnectionPooling               bool
	ConnectionPoolingAuthentication *string // "none", "simple"
	ConnectionPoolDebugLevel        *string // "fine", "all"
	ConnectionPoolInitialSize       *int
	ConnectionPoolMaximumSize       *int
	ConnectionPoolPreferredSize     *int
	ConnectionPoolProtocol          *string // "plain", "ssl"
	ConnectionPoolTimeout           *int

	ServerPrincipal                      string
	UseKerberosForPasswordAuthentication bool
	AllowKerberosAuthentication          bool
	KeyTab                               string
	KerberosRealm                        string

	BatchSizeForSync  int
	FullSyncPeriod    int // either a number, in milliseconds, or -1 if full sync is disabled
	ChangedSyncPeriod int // either a number, in milliseconds, or -1 if changed sync is disabled

	CachePolicy    string
	MaxLifespan    string // duration string (ex: 1h30m)
	EvictionDay    *int
	EvictionHour   *int
	EvictionMinute *int
}

func convertFromLdapUserFederationToComponent(ldap *LdapUserFederation) (*Component, error) {
	componentConfig := map[string][]string{
		"cachePolicy": {
			ldap.CachePolicy,
		},
		"enabled": {
			strconv.FormatBool(ldap.Enabled),
		},
		"priority": {
			strconv.Itoa(ldap.Priority),
		},
		"importEnabled": {
			strconv.FormatBool(ldap.ImportEnabled),
		},
		"editMode": {
			ldap.EditMode,
		},
		"syncRegistrations": {
			strconv.FormatBool(ldap.SyncRegistrations),
		},
		"vendor": {
			strings.ToLower(ldap.Vendor),
		},
		"usernameLDAPAttribute": {
			ldap.UsernameLDAPAttribute,
		},
		"rdnLDAPAttribute": {
			ldap.RdnLDAPAttribute,
		},
		"uuidLDAPAttribute": {
			ldap.UuidLDAPAttribute,
		},
		"userObjectClasses": {
			strings.Join(ldap.UserObjectClasses, ", "),
		},
		"connectionUrl": {
			ldap.ConnectionUrl,
		},
		"usersDn": {
			ldap.UsersDn,
		},
		"searchScope": {
			ldap.SearchScope,
		},
		"validatePasswordPolicy": {
			strconv.FormatBool(ldap.ValidatePasswordPolicy),
		},
		"trustEmail": {
			strconv.FormatBool(ldap.TrustEmail),
		},
		"pagination": {
			strconv.FormatBool(ldap.Pagination),
		},
		"connectionPooling": {
			strconv.FormatBool(ldap.ConnectionPooling),
		},
		"batchSizeForSync": {
			strconv.Itoa(ldap.BatchSizeForSync),
		},
		"fullSyncPeriod": {
			strconv.Itoa(ldap.FullSyncPeriod),
		},
		"changedSyncPeriod": {
			strconv.Itoa(ldap.ChangedSyncPeriod),
		},

		"serverPrincipal": {
			ldap.ServerPrincipal,
		},
		"useKerberosForPasswordAuthentication": {
			strconv.FormatBool(ldap.UseKerberosForPasswordAuthentication),
		},
		"allowKerberosAuthentication": {
			strconv.FormatBool(ldap.AllowKerberosAuthentication),
		},
		"keyTab": {
			ldap.KeyTab,
		},
		"kerberosRealm": {
			ldap.KerberosRealm,
		},
	}

	if ldap.BindDn != "" && ldap.BindCredential != "" {
		componentConfig["bindDn"] = []string{ldap.BindDn}
		componentConfig["bindCredential"] = []string{ldap.BindCredential}

		componentConfig["authType"] = []string{"simple"}
	} else {
		componentConfig["authType"] = []string{"none"}
	}

	if ldap.SearchScope == "ONE_LEVEL" {
		componentConfig["searchScope"] = []string{"1"}
	} else {
		componentConfig["searchScope"] = []string{"2"}
	}

	if ldap.CustomUserSearchFilter != "" {
		componentConfig["customUserSearchFilter"] = []string{ldap.CustomUserSearchFilter}
	}

	if ldap.UseTruststoreSpi == "ONLY_FOR_LDAPS" {
		componentConfig["useTruststoreSpi"] = []string{"ldapsOnly"}
	} else {
		componentConfig["useTruststoreSpi"] = []string{strings.ToLower(ldap.UseTruststoreSpi)}
	}

	if ldap.ConnectionTimeout != "" {
		connectionTimeoutMs, err := getMillisecondsFromDurationString(ldap.ConnectionTimeout)
		if err != nil {
			return nil, err
		}

		componentConfig["connectionTimeout"] = []string{connectionTimeoutMs}
	} else {
		componentConfig["connectionTimeout"] = []string{} // the keycloak API will not unset this unless the config is present with an empty array
	}

	if ldap.ReadTimeout != "" {
		readTimeoutMs, err := getMillisecondsFromDurationString(ldap.ReadTimeout)
		if err != nil {
			return nil, err
		}

		componentConfig["readTimeout"] = []string{readTimeoutMs}
	} else {
		componentConfig["readTimeout"] = []string{} // the keycloak API will not unset this unless the config is present with an empty array
	}

	componentConfig["evictionHour"] = []string{}
	componentConfig["evictionMinute"] = []string{}
	componentConfig["evictionDay"] = []string{}
	componentConfig["maxLifespan"] = []string{}

	if ldap.CachePolicy != "" {
		if ldap.EvictionHour != nil {
			componentConfig["evictionHour"] = []string{strconv.Itoa(*ldap.EvictionHour)}
		}
		if ldap.EvictionMinute != nil {
			componentConfig["evictionMinute"] = []string{strconv.Itoa(*ldap.EvictionMinute)}
		}
		if ldap.EvictionDay != nil {
			componentConfig["evictionDay"] = []string{strconv.Itoa(*ldap.EvictionDay)}
		}

		if ldap.MaxLifespan != "" {
			maxLifespanMs, err := getMillisecondsFromDurationString(ldap.MaxLifespan)
			if err != nil {
				return nil, err
			}
			componentConfig["maxLifespan"] = []string{maxLifespanMs}
		}
	}

	if ldap.ConnectionPoolingAuthentication != nil {
		componentConfig["connectionPoolingAuthentication"] = []string{*ldap.ConnectionPoolingAuthentication}
	}
	if ldap.ConnectionPoolDebugLevel != nil {
		componentConfig["connectionPoolingDebug"] = []string{*ldap.ConnectionPoolDebugLevel}
	}
	if ldap.ConnectionPoolProtocol != nil {
		componentConfig["connectionPoolingProtocol"] = []string{*ldap.ConnectionPoolProtocol}
	}
	if ldap.ConnectionPoolInitialSize != nil {
		componentConfig["connectionPoolingInitSize"] = []string{strconv.Itoa(*ldap.ConnectionPoolInitialSize)}
	}
	if ldap.ConnectionPoolMaximumSize != nil {
		componentConfig["connectionPoolingMaxSize"] = []string{strconv.Itoa(*ldap.ConnectionPoolMaximumSize)}
	}
	if ldap.ConnectionPoolPreferredSize != nil {
		componentConfig["connectionPoolingPrefSize"] = []string{strconv.Itoa(*ldap.ConnectionPoolPreferredSize)}
	}
	if ldap.ConnectionPoolTimeout != nil {
		componentConfig["connectionPoolingTimeout"] = []string{strconv.Itoa(*ldap.ConnectionPoolTimeout)}
	}

	return &Component{
		Id:           ldap.Id,
		Name:         ldap.Name,
		ProviderId:   "ldap",
		ProviderType: userStorageProviderType,
		ParentId:     ldap.RealmId,
		Config:       componentConfig,
	}, nil
}

func convertFromComponentToLdapUserFederation(component *Component) (*LdapUserFederation, error) {
	enabled, err := parseBoolAndTreatEmptyStringAsFalse(component.getConfig("enabled"))
	if err != nil {
		return nil, err
	}

	priority, err := strconv.Atoi(component.getConfig("priority"))
	if err != nil {
		return nil, err
	}

	importEnabled, err := parseBoolAndTreatEmptyStringAsFalse(component.getConfig("importEnabled"))
	if err != nil {
		return nil, err
	}

	syncRegistrations, err := parseBoolAndTreatEmptyStringAsFalse(component.getConfig("syncRegistrations"))
	if err != nil {
		return nil, err
	}

	userObjectClasses := strings.Split(component.getConfig("userObjectClasses"), ", ")

	validatePasswordPolicy, err := parseBoolAndTreatEmptyStringAsFalse(component.getConfig("validatePasswordPolicy"))
	if err != nil {
		return nil, err
	}

	trustEmail, err := parseBoolAndTreatEmptyStringAsFalse(component.getConfig("trustEmail"))
	if err != nil {
		return nil, err
	}

	pagination, err := parseBoolAndTreatEmptyStringAsFalse(component.getConfig("pagination"))
	if err != nil {
		return nil, err
	}

	connectionPooling, err := parseBoolAndTreatEmptyStringAsFalse(component.getConfig("connectionPooling"))
	if err != nil {
		return nil, err
	}

	batchSizeForSync, err := strconv.Atoi(component.getConfig("batchSizeForSync"))
	if err != nil {
		return nil, err
	}

	fullSyncPeriod, err := strconv.Atoi(component.getConfig("fullSyncPeriod"))
	if err != nil {
		return nil, err
	}

	changedSyncPeriod, err := strconv.Atoi(component.getConfig("changedSyncPeriod"))
	if err != nil {
		return nil, err
	}

	useKerberosForPasswordAuthentication, err := parseBoolAndTreatEmptyStringAsFalse(component.getConfig("useKerberosForPasswordAuthentication"))
	if err != nil {
		return nil, err
	}

	allowKerberosAuthentication, err := parseBoolAndTreatEmptyStringAsFalse(component.getConfig("allowKerberosAuthentication"))
	if err != nil {
		return nil, err
	}

	ldap := &LdapUserFederation{
		Id:      component.Id,
		Name:    component.Name,
		RealmId: component.ParentId,

		Enabled:  enabled,
		Priority: priority,

		ImportEnabled:     importEnabled,
		EditMode:          component.getConfig("editMode"),
		SyncRegistrations: syncRegistrations,

		Vendor:                 strings.ToUpper(component.getConfig("vendor")),
		UsernameLDAPAttribute:  component.getConfig("usernameLDAPAttribute"),
		RdnLDAPAttribute:       component.getConfig("rdnLDAPAttribute"),
		UuidLDAPAttribute:      component.getConfig("uuidLDAPAttribute"),
		UserObjectClasses:      userObjectClasses,
		ConnectionUrl:          component.getConfig("connectionUrl"),
		UsersDn:                component.getConfig("usersDn"),
		BindDn:                 component.getConfig("bindDn"),
		BindCredential:         component.getConfig("bindCredential"),
		CustomUserSearchFilter: component.getConfig("customUserSearchFilter"),
		SearchScope:            component.getConfig("searchScope"),

		ValidatePasswordPolicy: validatePasswordPolicy,
		TrustEmail:             trustEmail,
		UseTruststoreSpi:       component.getConfig("useTruststoreSpi"),
		Pagination:             pagination,

		ConnectionPooling: connectionPooling,

		ServerPrincipal:                      component.getConfig("serverPrincipal"),
		UseKerberosForPasswordAuthentication: useKerberosForPasswordAuthentication,
		AllowKerberosAuthentication:          allowKerberosAuthentication,
		KeyTab:                               component.getConfig("keyTab"),
		KerberosRealm:                        component.getConfig("kerberosRealm"),

		BatchSizeForSync:  batchSizeForSync,
		FullSyncPeriod:    fullSyncPeriod,
		ChangedSyncPeriod: changedSyncPeriod,

		CachePolicy: component.getConfig("cachePolicy"),
	}

	if d, ok := component.getConfigOk("connectionPoolingAuthentication"); ok {
		ldap.ConnectionPoolingAuthentication = &d
	}
	if d, ok := component.getConfigOk("connectionPoolingDebug"); ok {
		ldap.ConnectionPoolDebugLevel = &d
	}
	if d, ok := component.getConfigOk("connectionPoolingInitSize"); ok {
		i, err := strconv.Atoi(d)
		if err != nil {
			return nil, err
		}
		ldap.ConnectionPoolInitialSize = &i
	}
	if d, ok := component.getConfigOk("connectionPoolingMaxSize"); ok {
		i, err := strconv.Atoi(d)
		if err != nil {
			return nil, err
		}
		ldap.ConnectionPoolMaximumSize = &i
	}
	if d, ok := component.getConfigOk("connectionPoolingPrefSize"); ok {
		i, err := strconv.Atoi(d)
		if err != nil {
			return nil, err
		}
		ldap.ConnectionPoolPreferredSize = &i
	}
	if d, ok := component.getConfigOk("connectionPoolingProtocol"); ok {
		ldap.ConnectionPoolProtocol = &d
	}
	if d, ok := component.getConfigOk("connectionPoolingTimeout"); ok {
		i, err := strconv.Atoi(d)
		if err != nil {
			return nil, err
		}
		ldap.ConnectionPoolTimeout = &i
	}

	if bindDn := component.getConfig("bindDn"); bindDn != "" {
		ldap.BindDn = bindDn
	}

	if bindCredential := component.getConfig("bindCredential"); bindCredential != "" {
		ldap.BindCredential = bindCredential
	}

	if customUserSearchFilter := component.getConfig("customUserSearchFilter"); customUserSearchFilter != "" {
		ldap.CustomUserSearchFilter = customUserSearchFilter
	}

	if component.getConfig("searchScope") == "1" {
		ldap.SearchScope = "ONE_LEVEL"
	} else {
		ldap.SearchScope = "SUBTREE"
	}

	if useTruststoreSpi := component.getConfig("useTruststoreSpi"); useTruststoreSpi == "ldapsOnly" {
		ldap.UseTruststoreSpi = "ONLY_FOR_LDAPS"
	} else {
		ldap.UseTruststoreSpi = strings.ToUpper(useTruststoreSpi)
	}

	if connectionTimeout, ok := component.getConfigOk("connectionTimeout"); ok {
		connectionTimeoutDurationString, err := GetDurationStringFromMilliseconds(connectionTimeout)
		if err != nil {
			return nil, err
		}

		ldap.ConnectionTimeout = connectionTimeoutDurationString
	}

	if readTimeout, ok := component.getConfigOk("readTimeout"); ok {
		readTimeoutDurationString, err := GetDurationStringFromMilliseconds(readTimeout)
		if err != nil {
			return nil, err
		}

		ldap.ReadTimeout = readTimeoutDurationString
	}

	if maxLifespan, ok := component.getConfigOk("maxLifespan"); ok {
		maxLifespanString, err := GetDurationStringFromMilliseconds(maxLifespan)
		if err != nil {
			return nil, err
		}

		ldap.MaxLifespan = maxLifespanString
	}

	defaultEvictioValue := -1

	if evictionDay, ok := component.getConfigOk("evictionDay"); ok {
		evictionDayInt, err := strconv.Atoi(evictionDay)
		if err != nil {
			return nil, fmt.Errorf("unable to parse `evictionDay`: %w", err)
		}

		ldap.EvictionDay = &evictionDayInt
	} else {
		ldap.EvictionDay = &defaultEvictioValue
	}

	if evictionHour, ok := component.getConfigOk("evictionHour"); ok {
		evictionHourInt, err := strconv.Atoi(evictionHour)
		if err != nil {
			return nil, fmt.Errorf("unable to parse `evictionHour`: %w", err)
		}

		ldap.EvictionHour = &evictionHourInt
	} else {
		ldap.EvictionHour = &defaultEvictioValue
	}
	if evictionMinute, ok := component.getConfigOk("evictionMinute"); ok {
		evictionMinuteInt, err := strconv.Atoi(evictionMinute)
		if err != nil {
			return nil, fmt.Errorf("unable to parse `evictionMinute`: %w", err)
		}

		ldap.EvictionMinute = &evictionMinuteInt
	} else {
		ldap.EvictionMinute = &defaultEvictioValue
	}

	return ldap, nil
}

func (keycloakClient *KeycloakClient) ValidateLdapUserFederation(ldap *LdapUserFederation) error {
	if (ldap.BindDn == "" && ldap.BindCredential != "") || (ldap.BindDn != "" && ldap.BindCredential == "") {
		return fmt.Errorf("validation error: authentication requires both BindDN and BindCredential to be set")
	}

	return nil
}

func (keycloakClient *KeycloakClient) NewLdapUserFederation(ldapUserFederation *LdapUserFederation) error {
	component, err := convertFromLdapUserFederationToComponent(ldapUserFederation)
	if err != nil {
		return err
	}

	_, location, err := keycloakClient.post(fmt.Sprintf("/realms/%s/components", ldapUserFederation.RealmId), component)
	if err != nil {
		return err
	}

	ldapUserFederation.Id = getIdFromLocationHeader(location)

	return nil
}

func (keycloakClient *KeycloakClient) GetLdapUserFederation(realmId, id string) (*LdapUserFederation, error) {
	var component *Component

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/components/%s", realmId, id), &component, nil)
	if err != nil {
		return nil, err
	}

	return convertFromComponentToLdapUserFederation(component)
}

func (keycloakClient *KeycloakClient) GetLdapUserFederationMappers(realmId, id string) (*[]interface{}, error) {
	var components []*Component
	var ldapUserFederationMappers []interface{}

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/components?parent=%s&type=org.keycloak.storage.ldap.mappers.LDAPStorageMapper", realmId, id), &components, nil)
	if err != nil {
		return nil, err
	}
	for _, component := range components {
		switch component.ProviderId {
		case "full-name-ldap-mapper":
			mapper, err := convertFromComponentToLdapFullNameMapper(component, realmId)
			if err != nil {
				return nil, err
			}
			ldapUserFederationMappers = append(ldapUserFederationMappers, mapper)
		case "group-ldap-mapper":
			mapper, err := convertFromComponentToLdapGroupMapper(component, realmId)
			if err != nil {
				return nil, err
			}
			ldapUserFederationMappers = append(ldapUserFederationMappers, mapper)
		case "hardcoded-ldap-group-mapper":
			mapper := convertFromComponentToLdapHardcodedGroupMapper(component, realmId)
			ldapUserFederationMappers = append(ldapUserFederationMappers, mapper)
		case "hardcoded-ldap-role-mapper":
			mapper := convertFromComponentToLdapHardcodedRoleMapper(component, realmId)
			ldapUserFederationMappers = append(ldapUserFederationMappers, mapper)
		case "msad-lds-user-account-control-mapper":
			mapper, err := convertFromComponentToLdapMsadLdsUserAccountControlMapper(component, realmId)
			if err != nil {
				return nil, err
			}
			ldapUserFederationMappers = append(ldapUserFederationMappers, mapper)
		case "msad-user-account-control-mapper":
			mapper, err := convertFromComponentToLdapMsadUserAccountControlMapper(component, realmId)
			if err != nil {
				return nil, err
			}
			ldapUserFederationMappers = append(ldapUserFederationMappers, mapper)
		case "user-attribute-ldap-mapper":
			mapper, err := convertFromComponentToLdapUserAttributeMapper(component, realmId)
			if err != nil {
				return nil, err
			}
			ldapUserFederationMappers = append(ldapUserFederationMappers, mapper)
		case "role-ldap-mapper":
			mapper, err := convertFromComponentToLdapRoleMapper(component, realmId)
			if err != nil {
				return nil, err
			}
			ldapUserFederationMappers = append(ldapUserFederationMappers, mapper)
		}
	}

	return &ldapUserFederationMappers, nil
}

func (keycloakClient *KeycloakClient) DeleteLdapUserFederationMappers(realmId, ldapUserFederationId string) error {
	var components []*Component

	err := keycloakClient.get(fmt.Sprintf("/realms/%s/components?parent=%s&type=org.keycloak.storage.ldap.mappers.LDAPStorageMapper", realmId, ldapUserFederationId), &components, nil)
	if err != nil {
		return err
	}
	for _, component := range components {
		switch component.ProviderId {
		case "full-name-ldap-mapper",
			"group-ldap-mapper",
			"hardcoded-ldap-group-mapper",
			"hardcoded-ldap-role-mapper",
			"msad-lds-user-account-control-mapper",
			"msad-user-account-control-mapper",
			"user-attribute-ldap-mapper",
			"role-ldap-mapper":
			err := keycloakClient.DeleteComponent(realmId, component.Id)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (keycloakClient *KeycloakClient) UpdateLdapUserFederation(ldapUserFederation *LdapUserFederation) error {
	component, err := convertFromLdapUserFederationToComponent(ldapUserFederation)
	if err != nil {
		return err
	}

	return keycloakClient.put(fmt.Sprintf("/realms/%s/components/%s", ldapUserFederation.RealmId, ldapUserFederation.Id), component)
}

func (keycloakClient *KeycloakClient) DeleteLdapUserFederation(realmId, id string) error {
	return keycloakClient.delete(fmt.Sprintf("/realms/%s/components/%s", realmId, id), nil)
}
