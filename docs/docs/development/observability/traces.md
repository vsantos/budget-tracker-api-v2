# Visualizing traces

By default, the application will forward all it's traces to an opentelemetry collector, which forwards traces to a local jaeger installation.

Both containers can be created through `docker-compose up -d`.

To access jaeger UI local container, simply access: [http://localhost:16686](http://localhost:16686){:target="_blank"}.

## Generate traces

Let's explore a new request to see it through Jaeger, according to swagger, we can test any endpoint (even if it fails due to lack of authentication).