generate: generate-mocks
	@protoc -I ./proto \
  		--go_out ./backend/gen --go_opt paths=source_relative \
  		--go-grpc_out ./backend/gen --go-grpc_opt paths=source_relative \
  		--grpc-gateway_out ./backend/gen --grpc-gateway_opt paths=source_relative \
  		--js_out=import_style=commonjs:frontend/gen \
 		--grpc-web_out=import_style=commonjs,mode=grpcwebtext:frontend/gen \
		--openapiv2_out ./backend/gen \
  		./proto/*/*.proto
	@$(MAKE) combine-swagger
	

build-backend:
	@cd backend && go build -o bin/backend

run-backend: build-backend
	@./backend/bin/backend

migrate-up:
	@cd backend/migration && go run migrate.go up

migrate-down:
	@cd backend/migration && go run migrate.go down


MOCKGEN=mockgen
OUTPUT_DIR=mock
# Find all directories with Go files
DIRS=$(shell find . -type d -not -path "./$(OUTPUT_DIR)*" -not -path "./vendor*")

# Generate mocks
# this will skip generating mocks for the gen and frontend directories
generate-mocks: clean-mocks
	@for dir in $(DIRS); do \
		if [[ "$$dir" != *"/gen"* && "$$dir" != *"/frontend"* && "$$dir" != *"/$(OUTPUT_DIR)"* ]] && grep -q "interface" "$$dir"/*.go 2>/dev/null; then \
			echo "Generating mocks in $$dir"; \
			for file in "$$dir"/*.go; do \
				if ! grep -q "type.*interface" "$$file"; then \
					echo "Skipping $$file - no interface found"; \
					continue; \
				fi; \
				mkdir -p "$$dir/$(OUTPUT_DIR)"; \
				$(MOCKGEN) -source "$$file" -destination "$$dir/$(OUTPUT_DIR)/$$(basename $$file .go)_mock.go"; \
			done \
		fi \
	done

# Clean generated mocks
clean-mocks:
	find . -type d -name "$(OUTPUT_DIR)" | xargs rm -rf

SWAGGER_FILES := $(wildcard backend/gen/*/*.swagger.json)
COMBINED_SWAGGER := swagger.json

combine-swagger:
	swagger-cli bundle -o ./swagger/$(COMBINED_SWAGGER) $(SWAGGER_FILES) && swagger-cli validate ./swagger/$(COMBINED_SWAGGER)

tests:
	@cd backend && go test ./... --coverprofile=coverage.out && go tool cover -html=coverage.out
