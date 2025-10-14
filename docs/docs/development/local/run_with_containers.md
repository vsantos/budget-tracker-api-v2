# Running the code locally

In this section, let's explore how you can run the platform locally using container images.

!!! info "Running with or without containers?"
    My personal preference is to always work with containers for everything. My advice is to go through the documentation on “running without containers,” but in the end, stick with using containers.

## With containers

You can simply (at the root of your directory) trigger the [`docker-compose`](https://github.com/vsantos/budget-tracker-api-v2/blob/feat/initial_version/docker-compose-standalone.yml) manually or through the `Makefile` command `make rebuild-standalone`. Differently from `make rebuild`, the `make rebuild-standalone` command will also spin up a local mongoDB container with a single user already created.

The created user is a static one, injected to your container automatically to allow fast getting-started.

!!! note "initial credential"
    The best practice even for this scenario is to consult app's `swagger` and create a new user, deleting the `admin` one afterwards. This is needed to ensure you have an user for testing other protected endpoints, such as cards creation.

To use the initial credential, simply pass the following body when requesting a new JWT Token:

=== "'Request new token's body"
```json
{
	"login": "admin",
	"password": "myrandompassword"
}
```

Example:
=== "Shell"
    ```shell
    curl --request POST \
      --url http://localhost:8080/api/v1/jwt/issue \
      --header 'Content-Type: application/json' \
      --header 'User-Agent: my-manual-requester' \
      --data '{
        "login": "admin",
        "password": "myrandompassword"
      }'
    ```
=== "Shell outcome"
    ```json
    {
      "type":"bearer",
      "refresh":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjA0OTYxNjEsInN1YiI6IjUxNDNhZmM2NmQ0NGUxY2ViMzcyMTIxZSJ9.i8NntpiR5w6LiALRxpxvkTtFROTA2EWTYkcuieYXRuQ",
      "token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJleHAiOjE3NjA0MTAwNjEsImlhdCI6MTc2MDQwOTc2MSwibmFtZSI6ImFkbWluIiwic3ViIjoiNTE0M2FmYzY2ZDQ0ZTFjZWIzNzIxMjFlIn0.nX6Tug6FMInA02evdGcHOlr1AHoNe9usi-sr-cOYhJw"
    }
    ```

---

=== "Shell"
    ```shell
    # Manually bootstrapping all containers at once in background:
    docker-compose -f docker-compose-standalone.yml up -d

    # through Makefile. Prefered over the "manual docker-compose" to allow faster
    ## local development interactions 
    make rebuild-standalone
    ```

=== "Shell outcome"
    ```shell
    [+] Running 9/9 
    ✔ otel-collector                                Built    0.0s
    ✔ mongo_seed                                    Built    0.0s
    ✔ Network budget-tracker-api-v2_otel-network    Created  0.0s
    ✔ Container jaeger                              Started  0.2s
    ✔ Container mongodb                             Started  0.2s
    ✔ Container otel-collector                      Started  0.3s
    ✔ Container budger-tracker-api-v2               Started  0.3s
    ✔ Container budget-tracker-api-v2-mongo_seed-1  Started  0.3s
    ```

All regular containers - including observability ones - will be up and running, making all infrastructure-stack transparent.

We can use the same approach as running "without containers" to validate app's health:

=== "Shell"
```shell
# According to the app logs, the platform is running under port `8080`
## According to the swagger, we can check there is a `/health` endpoint

curl http://localhost:8080/health
```

!!! Success "This is the expected outcome"
    ```json
    {
        "message": "healthy",
        "app": true,
        "database": true
    }
    ```

### Dealing with MongoDB's dependency

If you won't want to deal with an external mongoDB instead of a local, simply override the original `docker-compose` file and use `make rebuild` instead of `make rebuild-standalone`.

=== "docker-compose.yml"
```yaml
  budger-tracker-api-v2:
    build: ./
    container_name: budger-tracker-api-v2
    environment:
      MONGODB_HOST: "mongodb+srv://my-mongodb-atlas-host.mongodb.net/"
      MONGODB_USER: "<REPLACE_ME>"
      MONGODB_PASS: "<REPLACE_ME>"
```

<script src="https://giscus.app/client.js"
        data-repo="vsantos/budget-tracker-api-v2-discussions"
        data-repo-id="R_kgDOQApX1g"
        data-category="General"
        data-category-id="DIC_kwDOQApX1s4CwhAe"
        data-mapping="pathname"
        data-strict="0"
        data-reactions-enabled="1"
        data-emit-metadata="0"
        data-input-position="top"
        data-theme="catppuccin_frappe"
        data-lang="en"
        crossorigin="anonymous"
        async>
</script>