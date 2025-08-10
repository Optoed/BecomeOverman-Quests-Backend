.PHONY: init
init:
	mkdir -p cmd internal/models internal/repositories internal/services internal/handlers configs migrations docs
	touch main.go configs/config.yaml README.md

.PHONY: swag-init
swag-init:
	swag init -g cmd/app/main.go