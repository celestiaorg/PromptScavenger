# Prompt Scavenger

Set the `.env` file with the following values:

```.env
NODE_RPC_IP="http://localhost:26658"
NODE_JWT_TOKEN=""
OPENAI_KEY=""
```

You must run a Light Node using `celestia-node` which would be on the same machine as where the script and game must run.

You must get the Node JWT token from [here](https://docs.celestia.org/developers/gateway-api-tutorial/#curl-guide)

Run the following:

```sh
go run main.go
```
