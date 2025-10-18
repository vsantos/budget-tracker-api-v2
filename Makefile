sonar:
	docker run -d -p 9000:9000 sonarqube:community && sleep 60
	docker run --net=host \
		--rm \
		-e SONAR_HOST_URL="http://localhost:9000"  \
		-v "$(pwd):/usr/src" \
		sonarsource/sonar-scanner-cli

static-docker-build:
	docker build . -t budget-tracker-api:local

test:
	staticcheck -checks='-S1021' ./...
	go test ./... -cover -coverprofile=coverage.out 
	$(MAKE) helm-test
	shellcheck hack/docs/generate_api_docs.sh

helm-test:
	helm unittest -f helm/budget-tracker/templates/tests/deployment_test.yaml helm/budget-tracker --failfast --color
	helm unittest -f helm/budget-tracker/templates/tests/service_test.yaml helm/budget-tracker --failfast --color
	helm unittest -f helm/budget-tracker/templates/tests/secret_test.yaml helm/budget-tracker --failfast --color
	helm unittest -f helm/budget-tracker/templates/tests/hpa_test.yaml helm/budget-tracker --failfast --color
	helm unittest -f helm/budget-tracker/templates/tests/namespace_test.yaml helm/budget-tracker --failfast --color

helm-docs:
	helm-docs helm/budget-tracker/

rebuild:
	$(MAKE) helm-test
	docker-compose down; docker-compose up -d --build

rebuild-standalone:
	$(MAKE) helm-test
	docker-compose -f docker-compose-standalone.yml down; docker-compose -f docker-compose-standalone.yml up -d --build

k8s-apply:
	if [[ -f helm/budget-tracker/Chat.lock ]]; then rm helm/budget-tracker/Chart.lock; fi
	$(MAKE) helm-test
	$(MAKE) helm-docs
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm repo update
	helm dependency build helm/budget-tracker
	helm template --release-name local-dev ./helm/budget-tracker | kubectl apply -f -
	if [[ -f helm/budget-tracker/Chart.lock ]]; then rm helm/budget-tracker/Chart.lock; fi

k8s-destroy:
	if [[ -f helm/budget-tracker/Chat.lock ]]; then rm helm/budget-tracker/Chart.lock; fi
	helm repo add bitnami https://charts.bitnami.com/bitnami
	helm repo update
	helm dependency build helm/budget-tracker
	helm template --release-name local-dev ./helm/budget-tracker | kubectl delete -f -
	if [[ -f helm/budget-tracker/Chart.lock ]]; then rm helm/budget-tracker/Chart.lock; fi

generate-docs:
	$(MAKE) helm-test
	$(MAKE) helm-docs
	./hack/docs/generate_api_docs.sh

serve-docs:
	$(MAKE) generate-docs
	$(MAKE) helm-docs
	mkdocs serve -f docs/mkdocs.yml