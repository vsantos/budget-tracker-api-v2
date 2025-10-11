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

generate-docs:
	./hack/docs/generate_api_docs.sh
	mkdocs build -f docs/mkdocs.yml

serve-docs:
	mkdocs serve -f docs/mkdocs.yml