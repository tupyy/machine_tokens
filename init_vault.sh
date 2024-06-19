#!/bin/sh

VAULT_ADDR=${VAULT_ADDR:-http://127.0.0.1:8200}
VAULT_TOKEN=${VAULT_TOKEN:-root}

# Setup jwt auth
vault auth enable jwt
vault write auth/jwt/config \
    oidc_discovery_url="http://localhost:8080/realms/vault" \
    oidc_client_id="" \
    oidc_client_secret="" \
    default_role=cosmin


# Creat the approle
vault auth enable approle

tee secret_read.hcl << EOF
path "secret/data/*" {
  capabilities = [ "read" ]
}
EOF
vault policy write secret_read secret_read.hcl

vault write auth/approle/role/awx \
    policies="secret_read" \
    token_ttl="10m" \
    token_max_ttl="30m"


# Policy grating to retreive secret-id
tee pull_secret_id.hcl << EOF
path "auth/approle/role/awx/secret-id" {
  capabilities = [ "update" ]
}
EOF
vault policy write pull_secret_id pull_secret_id.hcl

vault write auth/jwt/role/cosmin \
    role_type="jwt" \
    bound_audiences="account" \
    allowed_redirect_uris="http://localhost:8200/ui/vault/auth/jwt/jwt/callback" \
    user_claim="preferred_username" \
    policies="pull_secret_id" \
    token_policies="pull_secret_id"

# Write some secret
vault kv put -mount=secret \
    creds \
    key=value

