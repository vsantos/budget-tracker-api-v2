test:
	go test ./... -cover

docs:
	swagger generate spec --scan-models -o ./docs/swagger.yaml