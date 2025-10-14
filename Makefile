sonar:
	docker run -d -p 9000:9000 sonarqube:community && sleep 60
	docker run --net=host \
		--rm \
		-e SONAR_HOST_URL="http://localhost:9000"  \
		-v "$(pwd):/usr/src" \
		sonarsource/sonar-scanner-cli

test:
	golangci-lint run ./...
	go test ./... -cover

rebuild:
	go test ./... -cover
	docker-compose down; docker-compose up -d --build

rebuild-standalone:
	go test ./... -cover
	docker-compose -f docker-compose-standalone.yml down; docker-compose -f docker-compose-standalone.yml up -d --build

k8s-apply:
	helm template --release-name local-dev ./helm | kubectl apply -n bud -f -

generate-docs:
	./hack/docs/generate_api_docs.sh

serve-docs:
	mkdocs serve -f docs/mkdocs.yml