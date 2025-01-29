.PHONY: gen
gen:
	go generate ./...

.PHONY: build
build:
	$(START_LOG)
	@docker build \
		-t nonodox:latest \
		-f Dockerfile .
	$(END_LOG)