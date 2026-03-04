DB_URL=postgresql://postgres:postgres@localhost:7000/postgres?sslmode=disable
OPENAPI_GENERATOR := java -jar ~/openapi-generator-cli.jar
CONFIG_FILE := ./config.yaml
API_SRC := ./docs/api.yaml
API_BUNDLED := ./docs/api-bundled.yaml
DOCS_OUTPUT_DIR := ./docs/web
DOCS_INTERNAL_DIR := ./docs/web/docs
RESOURCES_DIR := ./pkg/resources

generate-models:
	test -d $(RESOURCES_DIR) || mkdir -p $(RESOURCES_DIR)
	test -d $(dir $(API_SRC)) || mkdir -p $(dir $(API_SRC))
	test -d $(dir $(API_BUNDLED)) || mkdir -p $(dir $(API_BUNDLED))
	test -d $(DOCS_OUTPUT_DIR) || mkdir -p $(DOCS_OUTPUT_DIR)

	rm -rf $(DOCS_INTERNAL_DIR) && mkdir -p $(DOCS_INTERNAL_DIR)
	rm -rf $(RESOURCES_DIR) && mkdir -p $(RESOURCES_DIR)
	swagger-cli bundle $(API_SRC) --outfile $(API_BUNDLED) --type yaml

	$(OPENAPI_GENERATOR) generate \
		-i $(API_BUNDLED) -g go \
		-o $(DOCS_OUTPUT_DIR) \
		--additional-properties=packageName=resources \
		--import-mappings uuid.UUID=github.com/google/uuid --type-mappings string+uuid=uuid.UUID \
		--import-mappings uint --type-mappings integer+uint=uint \

	mkdir -p $(RESOURCES_DIR)
	find $(DOCS_OUTPUT_DIR) -name '*.go' -exec mv {} $(RESOURCES_DIR)/ \;
	find $(RESOURCES_DIR) -type f -name "*_test.go" -delete

generate-sqlc:
	sqlc generate

build:
	KV_VIPER_FILE=$(CONFIG_FILE) go build -o ./cmd/organizations-svc/main ./cmd/organizations-svc/main.go

migrate-up:
	KV_VIPER_FILE=$(CONFIG_FILE) go build -o ./cmd/organizations-svc/main ./cmd/organizations-svc/main.go
	KV_VIPER_FILE=$(CONFIG_FILE) ./cmd/organizations-svc/main migrate up

migrate-down:
	KV_VIPER_FILE=$(CONFIG_FILE) go build -o ./cmd/organizations-svc/main ./cmd/organizations-svc/main.go
	KV_VIPER_FILE=$(CONFIG_FILE) ./cmd/organizations-svc/main migrate down

run-server:
	KV_VIPER_FILE=$(CONFIG_FILE) go build -o ./cmd/organizations-svc/main ./cmd/organizations-svc/main.go
	KV_VIPER_FILE=$(CONFIG_FILE) ./cmd/organizations-svc/main run service

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-rebuild:
	docker compose up -d --build --force-recreate