.PHONY: serve
serve:
	dev_appserver.py app.yaml --host=0.0.0.0 --admin_host=0.0.0.0

.PHONY: lint
lint:
	go vet ./...
	golint -set_exit_status ./...

.PHONY: fmt
fmt:
	go fmt ./...

