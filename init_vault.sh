#!/bin/sh

export VAULT_ADDR=http://127.0.0.1:8200
export VAULT_TOKEN=root

vault auth enable jwt
vault write auth/jwt/config \
    oidc_discovery_url="http://localhost:8080/realms/vault" \
    oidc_client_id="" \
    oidc_client_secret="" \
    default_role=cosmin

tee secret_read.hcl << EOF
path "secret/data/*" {
  capabilities = [ "read" ]
}
EOF
vault policy write secret_read secret_read.hcl

vault write auth/jwt/role/cosmin \
    role_type="jwt" \
    bound_audiences="account" \
    allowed_redirect_uris="http://localhost:8200/ui/vault/auth/jwt/jwt/callback" \
    user_claim="preferred_username" \
    policies="secret_read"

vault kv put -mount=secret \
    creds \
    key=value
