# Getting started

This guide will help you to: build, run and observe `budget-tracker-api` locally for development purposes.

## Requisites

- [X] Operational System
    * [X] Linux
    * [X] Unix (MacOS)
    * [ ] Windows
- [x] Binaries and requisites needed
    * [X] Docker Desktop / Rancher Desktop
    * [X] `docker` and `docker-compose` binaries
    * [X] git
    * [X] Python 3.x 
    * [X] Golang 1.25+
- [ ] External requisites
    * [ ] [MongoDB Atlas database](https://www.mongodb.com/cloud/atlas/register) - (Optional)
    * [X] Github - for [Github actions](https://github.com/features/actions) usage

### Fetching the code

First things first, ensure you have the latest code-base fetched:
 
```
# Feel free to use SSH instead of HTTPS
git clone https://github.com/vsantos/budget-tracker-api-v2
```

This command will create a local directory called `budget-tracker-api-v2` with the source code within.


## Running the code locally

In this section, let's explore how you can run the platform locally with or without using container images.

### Without containers

If you don't have/want any container engines running the app, you can build and run the platform by following these steps:


```
cd budget-tracker-api-v2

# Ensure you have successfully `cd`ed to the correct directory by running:
pwd
ls -lh
```

The expected outcome for this command is:

```
$ pwd
> /Users/${my-user}/Code/budget-tracker-api-v2

$ ls
>
CHECKLIST.md			Dockerfile			go.sum				LICENSE				otel-collector-config.yaml	swagger
docker				docs				hack				main.go				README.md
docker-compose.yml		go.mod				internal			Makefile			sonar-project.properties
```

#### building and running your app

Since our app was written in `go`, we will use the binary to compile the locally:

```
# Install the needed dependencies:
go mod tidy

# Build the binary
go build . -o budget-tracker
```

This command - if successful - will generate a local binary called `budget-tracker`.

Now, you can simply run it:

``` 
# This command will make your binary "runnable"
chmod +x budget-tracker

# Finally, run the binary
./budget-tracker-api
```

Ops! The app wasn't able to run, according to it's message:

> {"level":"fatal","msg":"empty MONGODB_HOST, MONGODB_USER or MONGODB_PASS env vars for mongodb atlas","time":"2025-10-13T13:24:46-03:00"}

The only external dependency for this application is a MongoDB's database, as we can check through the architecture page. This means we will need to pass a functional mongoDB database to the application, locally (example: through a container) or externally (example: through Atlas).

If you try to pass a "fake" host or a non-existent mongoDB URL, the app will break as well. This is an example of a functional command passing environment variables:

```
# This command will pass environment variables directly to your application's process
MONGODB_HOST="mongodb+srv://localhost:27017" \
MONGODB_USER="user" \
MONGODB_PASS="pass" \
./budget-tracker-api
```

You could also export these variables globally but for simplicity we are not using this option for now:

```
export MONGODB_HOST="mongodb+srv://localhost:27017" \
export MONGODB_USER="user" \
export MONGODB_PASS="pass"
```

```
./budget-tracker-api
```

#### Dealing with MongoDB's dependency

Let's use this opportunity (if you didn't setup your MongoDB Atlas instance yet) to get the credentials from MongoDB Atlas and finally run our app for the first time. Just replace the environment variables values to your command:

```
MONGODB_HOST="mongodb+srv://my-example.fj2qq.mongodb.net/" \
MONGODB_USER="my-user" \
MONGODB_PASS="my-pass" \
./budget-tracker-api
```

You should expect then the following outcome:

> {"level":"info","msg":"Server running on :8080","time":"2025-10-13T13:42:04-03:00"}

Now, we can simply test it locally by running a simple `curl` command:

```
# According to the app logs, the platform is running under port `8080`
## According to the swagger, we can check there is a `/health` endpoint

curl http://localhost:8080/health
```

This is the expected outcome:

> {"database": true}

This means we made a request to our app, which validated if it's connection with MongoDB is working properly. This endpoint is particularly useful when you have High-Availability mechanisms such as [Kubernetes' healthcheck](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/).

### With containers

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