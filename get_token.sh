curl -s -L -X POST http://localhost:3000/realms/vault/protocol/openid-connect/token \
-H 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'client_id=vault' \
--data-urlencode 'grant_type=password' \
--data-urlencode 'client_secret=vault' \
--data-urlencode 'scope=openid' \
--data-urlencode 'username=cosmin' \
--data-urlencode 'password=cosmin'
