FROM mongo:8.0.15

COPY init.json /init.json
COPY import.sh /import.sh

USER 1000:1000

SHELL ["/bin/bash", "-c"]
ENTRYPOINT  /import.sh
