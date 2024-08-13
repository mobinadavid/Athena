# Athena App Docs

## Minimum Required Infrastructure (LOM)

- 2 GB Memory
- 1 Core Cpu
- 1 GB Storage
- At least one IPv4 pointing to a domain

## Deployment Instructions

- Clone the project from CVS:
    ```shell
    git clone <repo> && cd <app_dir>
    ```
- Copy <code>.env.example</code> to <code>.env</code> And modify it as you wish.
    ```shell
    cp .env.example .env
    ```
- Download and install dependencies
    ```shell
    go mod download
    ```

- Obtain necessary SSL Certificates with services like certbot or acme
- Finally, build the app:
    ```shell
    go build -o athena main.go
    ```
- Handle the database migrations via:
    ```shell
    ./athena database migrate up
    ```
- Handle the database seeders via:
    ```shell
    ./athena database seed run
    ```
- Bootstrap the application via:
    ```shell
    ./athena app bootstrap
    ```
- To notice new transactions via:
    ```shell
    ./athena wallet-address trace
    ```

### Via Docker:
Docker and docker compose are available as well:
```shell
docker compose up -d
```
<hr>

### Vault Setup:
```shell
docker exec -it <vault_container> sh
vault operator init
vault operator unseal
vault login
vault secrets enable kv-v2
vault kv put kv-v2/<app_name>/<entity> <key>=<value> <key>=<value> ...
```
Setting blockchain-explorer secrets:
```shell
vault kv put kv-v2/athena/blockchain-explorer/api-key bscApiKey="" ethApiKey="" 
```
Setting IPG secrets:
```shell
vault kv put kv-v2/athena/ipg/sep api-key=""
```
Note: After unsealing, you should add <code>Initial Root Token</code> in .env file to let the golang app communicate with vault.
