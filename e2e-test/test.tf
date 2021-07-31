terraform {
  required_providers {
    keycloak = {
      source = "joed22636/keycloak"
      version = "1.0.3"
    }
  }
}

provider "keycloak" {
  client_id     = "admin-cli"
  username      = "admin"
  password      = "admin"
  url           = "http://192.168.1.172:8080"
}

resource "keycloak_realm" "_master" {
  realm = "master-test"
  enabled = true
  display_name = "Keycloak"
  display_name_html = "<div class=\"kc-logo-text\"><span>Keycloak</span></div>"
  user_managed_access = true
  attributes = {
    "actionTokenGeneratedByUserLifespan.execute-actions" = "1080"
    "actionTokenGeneratedByUserLifespan.reset-credentials" = "1080"
    "clientOfflineSessionIdleTimeout" = "420"
    "clientOfflineSessionMaxLifespan" = "480"
    "clientSessionIdleTimeout" = "540"
    "clientSessionMaxLifespan" = "600"
  }
  registration_allowed = true
  registration_email_as_username = true
  edit_username_allowed = true
  reset_password_allowed = true
  remember_me = true
  verify_email = true
  login_with_email_allowed = true
  duplicate_emails_allowed = false
  ssl_required = "external"
  login_theme = "keycloak"
  account_theme = "keycloak"
  admin_theme = "keycloak"
  email_theme = "keycloak"
  default_signature_algorithm = "ES512"
  revoke_refresh_token = true
  refresh_token_max_reuse = 5
  sso_session_idle_timeout = "60s"
  sso_session_max_lifespan = "7200s"
  offline_session_idle_timeout = "432000s"
  offline_session_max_lifespan = "518400s"
  offline_session_max_lifespan_enabled = true
  access_token_lifespan = "660s"
  access_token_lifespan_for_implicit_flow = "720s"
  access_code_lifespan = "780s"
  access_code_lifespan_login = "840s"
  access_code_lifespan_user_action = "900s"
  action_token_generated_by_user_lifespan = "960s"
  action_token_generated_by_admin_lifespan = "61200s"
  smtp_server {
    host = "smpthost"
    port = "545"
    from = "fromme@me.me"
    from_display_name = "dskplkay"
    reply_to = "fromme@me.me"
    reply_to_display_name = "replytothis"
    envelope_from = "fromme@me.me"
    starttls = "true"
    ssl = "true"
    auth {
      username = "lgin"
      password = "**********"
    }
  }

  internationalization {
    default_locale = "en"
    supported_locales =  [ "de", "no", "ru", "sv", "pt-BR", "lt", "en", "it", "fr", "es", "cs", "ja", "sk", "pl", "ca", "nl", "tr" ]
  }
  security_defenses  {
    brute_force_detection  {
      failure_reset_time_seconds =  12
      max_failure_wait_seconds =  900
      max_login_failures =  31
      minimum_quick_login_wait_seconds =  120
      permanent_lockout =  false
      quick_login_check_milli_seconds =  1001
      wait_increment_seconds =  60
    }

    headers  {
      content_security_policy =  "frame-src 'self'; frame-ancestors 'self'; object-src 'none';2"
      content_security_policy_report_only =  "2"
      strict_transport_security =  "max-age=31536000; includeSubDomains2"
      x_content_type_options =  "nosniff2"
      x_frame_options =  "SAMEORIGIN2"
      x_robots_tag =  "none2"
      x_xss_protection =  "1; mode=block2"
    }

  }

  password_policy = ""
  browser_flow = "browser"
  registration_flow = "registration"
  direct_grant_flow = "direct grant"
  reset_credentials_flow = "reset credentials"
  client_authentication_flow = "clients"
  docker_authentication_flow = "docker auth"
  web_authn_policy {
    acceptable_aaguids =  [ "asdf", "uiddd" ]
    attestation_conveyance_preference = "indirect"
    authenticator_attachment = "cross-platform"
    avoid_same_authenticator_register = "true"
    create_timeout = "667"
    relying_party_entity_name = "keycloak"
    relying_party_id = "rpid"
    require_resident_key = "Yes"
    signature_algorithms =  [ "RS512" ]
    user_verification_requirement = "preferred"
  }
  web_authn_passwordless_policy {
    acceptable_aaguids =  [ "errrr", "reeee" ]
    attestation_conveyance_preference = "direct"
    authenticator_attachment = "platform"
    avoid_same_authenticator_register = "true"
    create_timeout = "7776"
    relying_party_entity_name = "keycloak"
    relying_party_id = "diprrr2"
    require_resident_key = "No"
    signature_algorithms =  [ "ES384" ]
    user_verification_requirement = "discouraged"
  }
  default_default_client_scopes = [
    "role_list",
    "profile",
    "email",
    "roles",
    "web-origins",
  ]
  default_optional_client_scopes = [
    "offline_access",
    "phone",
    "microprofile-jwt",
    "scope-test",
  ]
}
resource "keycloak_openid_client_scope" "_master_REALM_scope_test" {
  realm_id = keycloak_realm._master.id
  name = "scope-test"
  description = "desc"
  consent_screen_text = "Consent!"
  include_in_token_scope = true
  gui_order = 2
}
resource "keycloak_generic_client_protocol_mapper" "_master_REALM_scope_test_SCOPE_nnnn" {
  realm_id = keycloak_realm._master.id
  client_scope_id = keycloak_openid_client_scope._master_REALM_scope_test.id
  name = "nnnn"
  protocol = "openid-connect"
  protocol_mapper = "oidc-role-name-mapper"
  config = {
    "new.role.name" = "rerer"
    "role" = "admin"
  }
}
resource "keycloak_generic_client_protocol_mapper" "_master_REALM_scope_test_SCOPE_note" {
  realm_id = keycloak_realm._master.id
  client_scope_id = keycloak_openid_client_scope._master_REALM_scope_test.id
  name = "note"
  protocol = "openid-connect"
  protocol_mapper = "oidc-usersessionmodel-note-mapper"
  config = {
    "access.token.claim" = "true"
    "claim.name" = "noteennnn"
    "id.token.claim" = "false"
    "jsonType.label" = "boolean"
    "user.session.note" = "notennn"
  }
}
resource "keycloak_generic_client_protocol_mapper" "_master_REALM_scope_test_SCOPE_gmp" {
  realm_id = keycloak_realm._master.id
  client_scope_id = keycloak_openid_client_scope._master_REALM_scope_test.id
  name = "gmp"
  protocol = "openid-connect"
  protocol_mapper = "oidc-group-membership-mapper"
  config = {
    "access.token.claim" = "true"
    "claim.name" = "ttt"
    "full.path" = "true"
    "id.token.claim" = "true"
    "userinfo.token.claim" = "true"
  }
}
resource "keycloak_generic_client_protocol_mapper" "_master_REALM_scope_test_SCOPE_nnn" {
  realm_id = keycloak_realm._master.id
  client_scope_id = keycloak_openid_client_scope._master_REALM_scope_test.id
  name = "nnn"
  protocol = "openid-connect"
  protocol_mapper = "oidc-allowed-origins-mapper"
  config = { }
}
resource "keycloak_generic_client_protocol_mapper" "_master_REALM_scope_test_SCOPE_fullnamepm" {
  realm_id = keycloak_realm._master.id
  client_scope_id = keycloak_openid_client_scope._master_REALM_scope_test.id
  name = "fullnamepm"
  protocol = "openid-connect"
  protocol_mapper = "oidc-full-name-mapper"
  config = {
    "access.token.claim" = "true"
    "id.token.claim" = "true"
    "userinfo.token.claim" = "true"
  }
}
resource "keycloak_generic_client_protocol_mapper" "_master_REALM_scope_test_SCOPE_client_rrooool" {
  realm_id = keycloak_realm._master.id
  client_scope_id = keycloak_openid_client_scope._master_REALM_scope_test.id
  name = "client rrooool"
  protocol = "openid-connect"
  protocol_mapper = "oidc-usermodel-client-role-mapper"
  config = {
    "access.token.claim" = "true"
    "claim.name" = "er"
    "id.token.claim" = "true"
    "jsonType.label" = "String"
    "multivalued" = "true"
    "userinfo.token.claim" = "true"
    "usermodel.clientRoleMapping.clientId" = "my-client"
    "usermodel.clientRoleMapping.rolePrefix" = "mycrrrr"
  }
}
resource "keycloak_generic_client_protocol_mapper" "_master_REALM_scope_test_SCOPE_rrr" {
  realm_id = keycloak_realm._master.id
  client_scope_id = keycloak_openid_client_scope._master_REALM_scope_test.id
  name = "rrr"
  protocol = "openid-connect"
  protocol_mapper = "oidc-usermodel-realm-role-mapper"
  config = {
    "access.token.claim" = "true"
    "claim.name" = "rererererrerrr"
    "id.token.claim" = "false"
    "jsonType.label" = "String"
    "multivalued" = "true"
    "userinfo.token.claim" = "true"
    "usermodel.realmRoleMapping.rolePrefix" = "rrr"
  }
}
resource "keycloak_generic_client_protocol_mapper" "_master_REALM_scope_test_SCOPE_hardcrp" {
  realm_id = keycloak_realm._master.id
  client_scope_id = keycloak_openid_client_scope._master_REALM_scope_test.id
  name = "hardcrp"
  protocol = "openid-connect"
  protocol_mapper = "oidc-hardcoded-role-mapper"
  config = {
    "role" = "offline_access"
  }
}
resource "keycloak_generic_client_protocol_mapper" "_master_REALM_scope_test_SCOPE_userpropr" {
  realm_id = keycloak_realm._master.id
  client_scope_id = keycloak_openid_client_scope._master_REALM_scope_test.id
  name = "userpropr"
  protocol = "openid-connect"
  protocol_mapper = "oidc-usermodel-property-mapper"
  config = {
    "access.token.claim" = "true"
    "claim.name" = "proopy"
    "id.token.claim" = "true"
    "jsonType.label" = "String"
    "user.attribute" = "propropr"
    "userinfo.token.claim" = "false"
  }
}
resource "keycloak_generic_client_protocol_mapper" "_master_REALM_scope_test_SCOPE_uattribute" {
  realm_id = keycloak_realm._master.id
  client_scope_id = keycloak_openid_client_scope._master_REALM_scope_test.id
  name = "uattribute"
  protocol = "openid-connect"
  protocol_mapper = "oidc-usermodel-attribute-mapper"
  config = {
    "access.token.claim" = "false"
    "aggregate.attrs" = "true"
    "claim.name" = "atrtrrrr"
    "id.token.claim" = "true"
    "jsonType.label" = "String"
    "multivalued" = "true"
    "user.attribute" = "attrr"
    "userinfo.token.claim" = "true"
  }
}
resource "keycloak_generic_client_protocol_mapper" "_master_REALM_scope_test_SCOPE_hardccc" {
  realm_id = keycloak_realm._master.id
  client_scope_id = keycloak_openid_client_scope._master_REALM_scope_test.id
  name = "hardccc"
  protocol = "openid-connect"
  protocol_mapper = "oidc-hardcoded-claim-mapper"
  config = {
    "access.token.claim" = "false"
    "claim.name" = "hrdrercecccc"
    "claim.value" = "cccc"
    "id.token.claim" = "true"
    "jsonType.label" = "String"
    "userinfo.token.claim" = "true"
  }
}
resource "keycloak_openid_client" "_master_REALM_my_client" {
  realm_id = keycloak_realm._master.id
  client_id = "my-client"
  name = ""
  enabled = true
  description = ""
  access_type = "PUBLIC"
  standard_flow_enabled = true
  implicit_flow_enabled = false
  direct_access_grants_enabled = true
  service_accounts_enabled = false
  valid_redirect_uris = [
    "https://test",
  ]
  web_origins = [
    "https://test",
  ]
  root_url = ""
  admin_url = "https://test"
  base_url = "https://test"
  pkce_code_challenge_method = ""
  full_scope_allowed = true
  consent_required = false
  access_token_lifespan = ""
  client_offline_session_idle_timeout = ""
  client_offline_session_max_lifespan = ""
  client_session_idle_timeout = ""
  client_session_max_lifespan = ""
  login_theme = ""
  exclude_session_state_from_auth_response = false
}
resource "keycloak_role" "_master_REALM_client_role" {
  realm_id = keycloak_realm._master.id
  name = "client role"
  client_id = keycloak_openid_client._master_REALM_my_client.id
  description = "cli rol"
  composite_roles =   [ keycloak_role._master_REALM_testrole.id ]
  attributes =  {
    clik1 = "cliv1"

  }

}
resource "keycloak_role" "_master_REALM_testrole" {
  realm_id = keycloak_realm._master.id
  name = "testrole"
  description = "tr"
  composite_roles =   [ keycloak_role._master_REALM_simp_role.id ]
  attributes =  {
    k1 = "v2"

  }

}
resource "keycloak_role" "_master_REALM_simp_role" {
  realm_id = keycloak_realm._master.id
  name = "simp role"
  client_id = keycloak_openid_client._master_REALM_my_client.id
  description = ""
  attributes =  {

  }

}
resource "keycloak_openid_client_default_scopes" "_master_REALM_my_client_CLIENT_default_scopes" {
  realm_id = keycloak_realm._master.id
  client_id = keycloak_openid_client._master_REALM_my_client.id
  default_scopes = [
    "web-origins",
    "profile",
    "roles",
    "email",
  ]
}
resource "keycloak_openid_client_optional_scopes" "_master_REALM_my_client_CLIENT_optional_scopes" {
  realm_id = keycloak_realm._master.id
  client_id = keycloak_openid_client._master_REALM_my_client.id
  optional_scopes = [
    "address",
    "phone",
    "offline_access",
    "microprofile-jwt",
  ]
}
resource "keycloak_realm_events" "_master_REALM_event_config" {
  realm_id = keycloak_realm._master.id
  admin_events_enabled = true
  admin_events_details_enabled = true
  events_enabled = true
  events_expiration = 198000
  enabled_event_types =   [ "SEND_RESET_PASSWORD", "UPDATE_CONSENT_ERROR", "GRANT_CONSENT", "REMOVE_TOTP", "REVOKE_GRANT", "UPDATE_TOTP", "LOGIN_ERROR", "CLIENT_LOGIN", "RESET_PASSWORD_ERROR", "IMPERSONATE_ERROR", "CODE_TO_TOKEN_ERROR", "CUSTOM_REQUIRED_ACTION", "RESTART_AUTHENTICATION", "IMPERSONATE", "UPDATE_PROFILE_ERROR", "LOGIN", "UPDATE_PASSWORD_ERROR", "FEDERATED_IDENTITY_LINK", "SEND_IDENTITY_PROVIDER_LINK", "SEND_VERIFY_EMAIL_ERROR", "RESET_PASSWORD", "CLIENT_INITIATED_ACCOUNT_LINKING_ERROR", "UPDATE_CONSENT", "REMOVE_TOTP_ERROR", "VERIFY_EMAIL_ERROR", "SEND_RESET_PASSWORD_ERROR", "CLIENT_UPDATE", "CUSTOM_REQUIRED_ACTION_ERROR", "IDENTITY_PROVIDER_POST_LOGIN_ERROR", "UPDATE_TOTP_ERROR", "CODE_TO_TOKEN", "GRANT_CONSENT_ERROR", "IDENTITY_PROVIDER_FIRST_LOGIN_ERROR" ]
  events_listeners =   [ "email" ]
}
