FROM mongo:8.0.15

COPY init.json /init.json

RUN addgroup -S mongo-seed && adduser -S mongo-seed -G bmongo-seed
USER mongo-seed

SHELL ["/bin/bash", "-c"]
ENTRYPOINT mongoimport --host localhost --db reach-engine --collection users --type json --file /init.json
