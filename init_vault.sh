vault auth enable jwt
vault write auth/jwt/config oidc_discovery_url="http://localhost:8080/realms/vault" oidc_client_id="" oidc_client_secret="" default_role=cosmin
vault write auth/jwt/role/cosmin role_type=jwt bound_audiences=account allowed_redirect_uris=http://localhost:8200/ui/vault/auth/jwt/jwt/callback user_claim=preferred_username
