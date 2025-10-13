# Running the code locally

In this section, let's explore how you can run the platform locally using container images.

!!! info "Running with or without containers?"
    My personal preference is to always work with containers for everything. My advice is to go through the documentation on “running without containers,” but in the end, stick with using containers.

## With containers

You can simply (at the root of your directory) trigger the [`docker-compose`](https://github.com/vsantos/budget-tracker-api-v2/blob/feat/initial_version/docker-compose.yml) manually or through the `Makefile` command `make rebuild`:

Make sure you edit the file `docker-compose` at the root of your directory to make sure to pass mongodb's credentials:

=== "docker-compose.yml"
```yaml
  budger-tracker-api-v2:
    build: ./
    container_name: budger-tracker-api-v2
    environment:
      MONGODB_HOST: "mongodb+srv://<REPLACE_ME>/"
      MONGODB_USER: "<REPLACE_ME>"
      MONGODB_PASS: "<REPLACE_ME>"
```

!!! warning "don't persist your credentials to git"
    This file is ignored by `.gitconfig` so if you accidentally save those credentails, no changes will be known by your `git` process.

=== "Shell"
    ```shell
    # Manually bootstrapping all containers at once in background:
    docker-compose up -d

    # through Makefile. Prefered over the "manual docker-compose" to allow faster
    ## local development interactions 
    make rebuild
    ```
=== "Shell outcome"
    ```shell
    [+] Running 4/4
    ✔ Network budget-tracker-api-v2_otel-network  Created                                                  0.0s 
    ✔ Container jaeger                            Started                                                  0.2s 
    ✔ Container otel-collector                    Started                                                  0.3s 
    ✔ Container budger-tracker-api-v2             Started                                                  0.3s
    ```

All regular containers - including observability ones - will be up and running except MongoDB's container.

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

If you won't want to deal with an external mongoDB but to test with a local one instead. We will be missing the details step-by-step for now but you can simply create a local mongoDB and point the correct `localhost:27017`:


=== "docker-compose.yml"
```yaml
  budger-tracker-api-v2:
    build: ./
    container_name: budger-tracker-api-v2
    environment:
      MONGODB_HOST: "mongodb+srv://localhost:27017/"
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