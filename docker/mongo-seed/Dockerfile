FROM mongo:8.0.15

COPY init.json /init.json

CMD mongoimport --host localhost --db reach-engine --collection users --type json --file /init.json