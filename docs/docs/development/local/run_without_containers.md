# Running the code locally

In this section, let's explore how you can run the platform locally without using container images.

!!! info "Running with or without containers?"
    My personal preference is to always work with containers for everything. My advice is to go through the documentation on “running without containers,” but in the end, stick with using containers.

## Without containers

If you don't have/want any container engines running the app, you can build and run the platform by following these steps:

=== "Shell"
    ```shell
    cd budget-tracker-api-v2

    # Ensure you have successfully `cd`ed to the correct directory by running:
    pwd
    ```
=== "Shell outcome"
    ```shell
    /Users/${my-user}/Code/budget-tracker-api-v2
    ```

---

=== "Shell"
    ```shell
    ls
    ```
=== "Shell outocome"
    ```shell
    CHECKLIST.md			Dockerfile			go.sum				LICENSE				otel-collector-config.yaml	swagger
    docker				docs				hack				main.go				README.md
    docker-compose.yml		go.mod				internal			Makefile			sonar-project.properties
    ```
### building and running your app

Since our app was written in `go`, we will use the binary to compile the locally:

=== "Shell"
```shell
# Install the needed dependencies:
go mod tidy

# Build the binary
go build . -o budget-tracker
```

This command - if successful - will generate a local binary called `budget-tracker`.

Now, you can simply run it:

=== "Shell"
```shell
# This command will make your binary "runnable"
chmod +x budget-tracker

# Finally, run the binary
./budget-tracker-api
```

!!! Failure "Ops! The app wasn't able to run, according to it's message"
    ```json
    {
        "level":"fatal",
        "msg":"empty MONGODB_HOST, MONGODB_USER or MONGODB_PASS env vars for mongodb",
        "time":"2025-10-13T13:24:46-03:00"
    }
    ```

The only external dependency for this application is a MongoDB database, as shown on the architecture page. This means we need to provide a functional MongoDB instance to the application—either locally (e.g., through a container) or externally (e.g., via Atlas).

If you try to pass a fake host or a non-existent MongoDB URL, the app will fail to start.

Here’s an example of a functional command that passes environment variables:

```shell
# This command will pass environment variables directly
## to your application's process

MONGODB_HOST="mongodb+srv://my-mongodb-atlas-url.mongodb.net" \
MONGODB_USER="user" \
MONGODB_PASS="pass" \
./budget-tracker-api
```

You could also export these variables globally but for simplicity we are not using this option for now:

=== "Shell"
```shell
export MONGODB_HOST="mongodb+srv://my-mongodb-atlas-url.mongodb.net" \
export MONGODB_USER="user" \
export MONGODB_PASS="pass"
```

=== "Shell"
```
./budget-tracker-api
```

### Dealing with MongoDB's dependency

Let’s take this chance (if you haven’t set up your MongoDB Atlas instance yet) to grab your credentials from Atlas and run the app for the first time. Just replace the environment variable values in your command:

=== "Shell"
```shell
MONGODB_HOST="mongodb+srv://my-example.fj2qq.mongodb.net/" \
MONGODB_USER="my-user" \
MONGODB_PASS="my-pass" \
./budget-tracker-api
```

!!! Success "You should expect then the following outcome"
    ```json
    {"level":"info","msg":"Server running on :8080","time":"2025-10-13T13:42:04-03:00"}
    ```

Now, we can simply test it locally by running a simple `curl` command:

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

This means we made a request to our app, which validated if it's connection with MongoDB is working properly. This endpoint is particularly useful when you have High-Availability mechanisms such as [Kubernetes' healthcheck](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/).


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