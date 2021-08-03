package provider

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/joed22636/terraform-provider-keycloak/keycloak"
)

func resourceKeycloakComponent() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeycloakComponentCreate,
		Read:   resourceKeycloakComponentRead,
		Update: resourceKeycloakComponentUpdate,
		Delete: resourceKeycloakComponentDelete,
		// If this resource uses authentication, then this resource must be imported using the syntax {{realm_id}}/{{provider_id}}/{{bind_credential}}
		// Otherwise, this resource can be imported using {{realm}}/{{provider_id}}.
		// The Provider ID is displayed in the GUI when editing this provider
		Importer: &schema.ResourceImporter{
			State: resourceKeycloakComponentImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Display name of the provider when displayed in the console.",
			},
			"realm_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The realm this provider will provide user federation for.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "When false, this provider will not be used when performing queries for users.",
			},
			"priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Priority of this provider when looking up users. Lower values are first.",
			},
			"import_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "When true, LDAP users will be imported into the Keycloak database.",
			},
			"sync_registrations": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When true, newly created users will be synced back to LDAP.",
			},
			"username_ldap_attribute": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the LDAP attribute to use as the Keycloak username.",
			},
			"rdn_ldap_attribute": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the LDAP attribute to use as the relative distinguished name.",
			},
			"uuid_ldap_attribute": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the LDAP attribute to use as a unique object identifier for objects in LDAP.",
			},
			"user_object_classes": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "All values of LDAP objectClass attribute for users in LDAP.",
			},
			"connection_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Connection URL to the LDAP server.",
			},
			"users_dn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Full DN of LDAP tree where your users are.",
			},
			"bind_dn": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "DN of LDAP admin, which will be used by Keycloak to access LDAP server.",
			},
			"bind_credential": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				DiffSuppressFunc: func(_, remoteBindCredential, _ string, _ *schema.ResourceData) bool {
					return remoteBindCredential == "**********"
				},
				Description: "Password of LDAP admin.",
			},
			"custom_user_search_filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Additional LDAP filter for filtering searched users. Must begin with '(' and end with ')'.",
			},
			"validate_password_policy": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When true, Keycloak will validate passwords using the realm policy before updating it.",
			},
			"trust_email": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If enabled, email provided by this provider is not verified even if verification is enabled for the realm.",
			},
			"connection_timeout": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "LDAP connection timeout (duration string)",
				DiffSuppressFunc: suppressDurationStringDiff,
			},
			"read_timeout": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "LDAP read timeout (duration string)",
				DiffSuppressFunc: suppressDurationStringDiff,
			},
			"pagination": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "When true, Keycloak assumes the LDAP server supports pagination.",
			},

			"connection_pooling":                {Type: schema.TypeBool, Required: true},
			"connection_pooling_authentication": {Type: schema.TypeString, Optional: true},
			"connection_pool_debug_level":       {Type: schema.TypeString, Optional: true},
			"connection_pool_initial_size":      {Type: schema.TypeInt, Optional: true},
			"connection_pool_maximum_size":      {Type: schema.TypeInt, Optional: true},
			"connection_pool_preferred_size":    {Type: schema.TypeInt, Optional: true},
			"connection_pool_protocol":          {Type: schema.TypeString, Optional: true},
			"connection_pool_timeout":           {Type: schema.TypeInt, Optional: true},

			"batch_size_for_sync": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1000,
				Description: "The number of users to sync within a single transaction.",
			},
			"full_sync_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      -1,
				ValidateFunc: validateSyncPeriod,
				Description:  "How frequently Keycloak should sync all LDAP users, in seconds. Omit this property to disable periodic full sync.",
			},
			"changed_sync_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      -1,
				ValidateFunc: validateSyncPeriod,
				Description:  "How frequently Keycloak should sync changed LDAP users, in seconds. Omit this property to disable periodic changed users sync.",
			},

			"kerberos": {
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    1,
				Description: "Settings regarding kerberos authentication for this realm.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kerberos_realm": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the kerberos realm, e.g. FOO.LOCAL",
						},
						"server_principal": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The kerberos server principal, e.g. 'HTTP/host.foo.com@FOO.LOCAL'.",
						},
						"key_tab": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Path to the kerberos keytab file on the server with credentials of the service principal.",
						},
						"use_kerberos_for_password_authentication": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Use kerberos login module instead of ldap service api. Defaults to `false`.",
						},
					},
				},
			},
			"cache": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Settings regarding cache policy for this realm.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "DEFAULT",
							ValidateFunc: validation.StringInSlice(keycloakUserFederationCachePolicies, false),
						},
						"max_lifespan": {
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: suppressDurationStringDiff,
							Description:      "Max lifespan of cache entry (duration string).",
						},
						"eviction_day": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      "-1",
							ValidateFunc: validation.All(validation.IntAtLeast(0), validation.IntAtMost(6)),
							Description:  "Day of the week the entry will become invalid on.",
						},
						"eviction_hour": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      "-1",
							ValidateFunc: validation.All(validation.IntAtLeast(0), validation.IntAtMost(23)),
							Description:  "Hour of day the entry will become invalid on.",
						},
						"eviction_minute": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      "-1",
							ValidateFunc: validation.All(validation.IntAtLeast(0), validation.IntAtMost(59)),
							Description:  "Minute of day the entry will become invalid on.",
						},
					},
				},
			},
		},
	}
}

func getComponentFromData(data *schema.ResourceData) *keycloak.Component {
	var userObjectClasses []string

	for _, userObjectClass := range data.Get("user_object_classes").([]interface{}) {
		userObjectClasses = append(userObjectClasses, userObjectClass.(string))
	}

	component := &keycloak.Component{
		Id:   data.Id(),
		Name: data.Get("name").(string),
		// RealmId: data.Get("realm_id").(string),

		// Enabled:  data.Get("enabled").(bool),
		// Priority: data.Get("priority").(int),

		// ImportEnabled:     data.Get("import_enabled").(bool),
		// EditMode:          data.Get("edit_mode").(string),
		// SyncRegistrations: data.Get("sync_registrations").(bool),

		// Vendor:                 data.Get("vendor").(string),
		// UsernameLDAPAttribute:  data.Get("username_ldap_attribute").(string),
		// RdnLDAPAttribute:       data.Get("rdn_ldap_attribute").(string),
		// UuidLDAPAttribute:      data.Get("uuid_ldap_attribute").(string),
		// UserObjectClasses:      userObjectClasses,
		// ConnectionUrl:          data.Get("connection_url").(string),
		// UsersDn:                data.Get("users_dn").(string),
		// BindDn:                 data.Get("bind_dn").(string),
		// BindCredential:         data.Get("bind_credential").(string),
		// CustomUserSearchFilter: data.Get("custom_user_search_filter").(string),
		// SearchScope:            data.Get("search_scope").(string),

		// ValidatePasswordPolicy: data.Get("validate_password_policy").(bool),
		// TrustEmail:             data.Get("trust_email").(bool),
		// UseTruststoreSpi:       data.Get("use_truststore_spi").(string),
		// ConnectionTimeout:      data.Get("connection_timeout").(string),
		// ReadTimeout:            data.Get("read_timeout").(string),
		// Pagination:             data.Get("pagination").(bool),

		// ConnectionPooling: data.Get("connection_pooling").(bool),

		// BatchSizeForSync:  data.Get("batch_size_for_sync").(int),
		// FullSyncPeriod:    data.Get("full_sync_period").(int),
		// ChangedSyncPeriod: data.Get("changed_sync_period").(int),
	}

	// if d, ok := data.GetOk("connection_pooling_authentication"); ok {
	// 	v := d.(string)
	// 	component.ConnectionPoolingAuthentication = &v
	// }
	// if d, ok := data.GetOk("connection_pool_debug_level"); ok {
	// 	v := d.(string)
	// 	component.ConnectionPoolDebugLevel = &v
	// }
	// if d, ok := data.GetOk("connection_pool_initial_size"); ok {
	// 	v := d.(int)
	// 	component.ConnectionPoolInitialSize = &v
	// }
	// if d, ok := data.GetOk("connection_pool_maximum_size"); ok {
	// 	v := d.(int)
	// 	component.ConnectionPoolMaximumSize = &v
	// }
	// if d, ok := data.GetOk("connection_pool_preferred_size"); ok {
	// 	v := d.(int)
	// 	component.ConnectionPoolPreferredSize = &v
	// }
	// if d, ok := data.GetOk("connection_pool_protocol"); ok {
	// 	v := d.(string)
	// 	component.ConnectionPoolProtocol = &v
	// }
	// if d, ok := data.GetOk("connection_pool_timeout"); ok {
	// 	v := d.(int)
	// 	component.ConnectionPoolTimeout = &v
	// }

	// if cache, ok := data.GetOk("cache"); ok {
	// 	cache := cache.([]interface{})
	// 	cacheData := cache[0].(map[string]interface{})

	// 	evictionDay := cacheData["eviction_day"].(int)
	// 	evictionHour := cacheData["eviction_hour"].(int)
	// 	evictionMinute := cacheData["eviction_minute"].(int)

	// 	component.MaxLifespan = cacheData["max_lifespan"].(string)

	// 	component.EvictionDay = &evictionDay
	// 	component.EvictionHour = &evictionHour
	// 	component.EvictionMinute = &evictionMinute
	// 	component.CachePolicy = cacheData["policy"].(string)
	// }

	// if kerberos, ok := data.GetOk("kerberos"); ok {
	// 	component.AllowKerberosAuthentication = true
	// 	kerberosSettingsData := kerberos.(*schema.Set).List()[0]
	// 	kerberosSettings := kerberosSettingsData.(map[string]interface{})

	// 	component.KerberosRealm = kerberosSettings["kerberos_realm"].(string)
	// 	component.ServerPrincipal = kerberosSettings["server_principal"].(string)
	// 	component.UseKerberosForPasswordAuthentication = kerberosSettings["use_kerberos_for_password_authentication"].(bool)
	// 	component.KeyTab = kerberosSettings["key_tab"].(string)
	// } else {
	// 	component.AllowKerberosAuthentication = false
	// }

	return component
}

func setComponentData(data *schema.ResourceData, ldap *keycloak.Component) {
	data.SetId(ldap.Id)

	data.Set("name", ldap.Name)
	// data.Set("realm_id", ldap.RealmId)

	// data.Set("enabled", ldap.Enabled)
	// data.Set("priority", ldap.Priority)

	// data.Set("import_enabled", ldap.ImportEnabled)
	// data.Set("edit_mode", ldap.EditMode)
	// data.Set("sync_registrations", ldap.SyncRegistrations)

	// data.Set("vendor", ldap.Vendor)
	// data.Set("username_ldap_attribute", ldap.UsernameLDAPAttribute)
	// data.Set("rdn_ldap_attribute", ldap.RdnLDAPAttribute)
	// data.Set("uuid_ldap_attribute", ldap.UuidLDAPAttribute)
	// data.Set("user_object_classes", ldap.UserObjectClasses)
	// data.Set("connection_url", ldap.ConnectionUrl)
	// data.Set("users_dn", ldap.UsersDn)
	// data.Set("bind_dn", ldap.BindDn)
	// data.Set("bind_credential", ldap.BindCredential)
	// data.Set("custom_user_search_filter", ldap.CustomUserSearchFilter)
	// data.Set("search_scope", ldap.SearchScope)

	// data.Set("validate_password_policy", ldap.ValidatePasswordPolicy)
	// data.Set("trust_email", ldap.TrustEmail)
	// data.Set("use_truststore_spi", ldap.UseTruststoreSpi)
	// data.Set("connection_timeout", ldap.ConnectionTimeout)
	// data.Set("read_timeout", ldap.ReadTimeout)
	// data.Set("pagination", ldap.Pagination)

	// data.Set("connection_pooling", ldap.ConnectionPooling)
	// data.Set("connection_pooling_authentication", ldap.ConnectionPoolingAuthentication)
	// data.Set("connection_pool_debug_level", ldap.ConnectionPoolDebugLevel)
	// data.Set("connection_pool_initial_size", ldap.ConnectionPoolInitialSize)
	// data.Set("connection_pool_maximum_size", ldap.ConnectionPoolMaximumSize)
	// data.Set("connection_pool_preferred_size", ldap.ConnectionPoolPreferredSize)
	// data.Set("connection_pool_protocol", ldap.ConnectionPoolProtocol)
	// data.Set("connection_pool_timeout", ldap.ConnectionPoolTimeout)

	// if ldap.AllowKerberosAuthentication {
	// 	kerberosSettings := make(map[string]interface{})

	// 	kerberosSettings["server_principal"] = ldap.ServerPrincipal
	// 	kerberosSettings["use_kerberos_for_password_authentication"] = ldap.UseKerberosForPasswordAuthentication
	// 	kerberosSettings["key_tab"] = ldap.KeyTab
	// 	kerberosSettings["kerberos_realm"] = ldap.KerberosRealm

	// 	data.Set("kerberos", []interface{}{kerberosSettings})
	// } else {
	// 	data.Set("kerberos", nil)
	// }

	// data.Set("batch_size_for_sync", ldap.BatchSizeForSync)
	// data.Set("full_sync_period", ldap.FullSyncPeriod)
	// data.Set("changed_sync_period", ldap.ChangedSyncPeriod)

	// if _, ok := data.GetOk("cache"); ok {
	// 	cachePolicySettings := make(map[string]interface{})

	// 	if ldap.MaxLifespan != "" {
	// 		cachePolicySettings["max_lifespan"] = ldap.MaxLifespan
	// 	}

	// 	if ldap.EvictionDay != nil {
	// 		cachePolicySettings["eviction_day"] = *ldap.EvictionDay
	// 	}
	// 	if ldap.EvictionHour != nil {
	// 		cachePolicySettings["eviction_hour"] = *ldap.EvictionHour
	// 	}
	// 	if ldap.EvictionMinute != nil {
	// 		cachePolicySettings["eviction_minute"] = *ldap.EvictionMinute
	// 	}

	// 	cachePolicySettings["policy"] = ldap.CachePolicy

	// 	data.Set("cache", []interface{}{cachePolicySettings})
	// }
}

func resourceKeycloakComponentCreate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	ldap := getComponentFromData(data)

	err := keycloakClient.ValidateComponent(ldap)
	if err != nil {
		return err
	}

	err = keycloakClient.NewComponent(ldap)
	if err != nil {
		return err
	}

	err = keycloakClient.DeleteComponentMappers(ldap.RealmId, ldap.Id)
	if err != nil {
		return err
	}

	setComponentData(data, ldap)

	return resourceKeycloakComponentRead(data, meta)
}

func resourceKeycloakComponentRead(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	realmId := data.Get("realm_id").(string)
	id := data.Id()

	ldap, err := keycloakClient.GetComponent(realmId, id)
	if err != nil {
		return handleNotFoundError(err, data)
	}

	ldap.BindCredential = data.Get("bind_credential").(string) // we can't trust the API to set this field correctly since it just responds with "**********"
	setComponentData(data, ldap)

	return nil
}

func resourceKeycloakComponentUpdate(data *schema.ResourceData, meta interface{}) error {
	keycloakClient := meta.(*keycloak.KeycloakClient)

	ldap := getComponentFromData(data)

	err := keycloakClient.ValidateComponent(ldap)
	if err != nil {
		return err
	}

	err = keycloakClient.UpdateComponent(ldap)
	if err != nil {
		return err
	}

	setComponentData(data, ldap)

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
	case len(parts) == 3:
		realmId = parts[0]
		id = parts[1]
		d.Set("bind_credential", parts[2])
	default:
		return nil, fmt.Errorf("Invalid import. Supported import formats: {{realmId}}/{{userFederationId}}, {{realmId}}/{{userFederationId}}/{{bindCredentials}}")
	}

	d.Set("realm_id", realmId)
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}
