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
	$(MAKE) helm-test

helm-test:
	helm unittest -f helm/templates/tests/deployment_test.yaml helm
	helm unittest -f helm/templates/tests/service_test.yaml helm
	helm unittest -f helm/templates/tests/deployment.yaml helm

rebuild:
	$(MAKE) helm-test
	docker-compose down; docker-compose up -d --build

rebuild-standalone:
	$(MAKE) helm-test
	docker-compose -f docker-compose-standalone.yml down; docker-compose -f docker-compose-standalone.yml up -d --build

k8s-apply:
	$(MAKE) helm-test
	kubectl create namespace demo
	helm template --release-name local-dev ./helm | kubectl apply -n demo -f -

generate-docs:
	$(MAKE) helm-test
	./hack/docs/generate_api_docs.sh

serve-docs:
	$(MAKE) generate-docs
	mkdocs serve -f docs/mkdocs.yml