
# Hydra commands 

## Import clients

	docker exec -it hydra hydra clients import "/config/hydra-client.json"

## Include single client with authorization codes

	docker exec -it hydra hydra clients create --endpoint http://localhost:4445 --id piensa --secret piensa --grant-types authorization_code,refresh_token --response-types code,id_token --scope pagovalemia,openid,offline,blahblah --callbacks http://127.0.0.1:5555/callback

## Authorize client.

	docker exec -it hydra hydra token user --client-id piensa --client-secret piensa --endpoint http://localhost:4444 --port 5555 --scope openid,offline,blahblah

