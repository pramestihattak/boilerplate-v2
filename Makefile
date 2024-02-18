generate:
	@protoc -I ./proto \
  		--go_out ./backend/gen --go_opt paths=source_relative \
  		--go-grpc_out ./backend/gen --go-grpc_opt paths=source_relative \
  		--grpc-gateway_out ./backend/gen --grpc-gateway_opt paths=source_relative \
  		--js_out=import_style=commonjs:frontend/gen \
 		--grpc-web_out=import_style=commonjs,mode=grpcwebtext:frontend/gen \
		--openapiv2_out ./backend/gen \
  		./proto/*/*.proto

build-backend:
	@cd backend && go build -o bin/backend

run-backend: build-backend
	@./backend/bin/backend

migrate-up:
	@cd backend/migration && go run migrate.go up

migrate-down:
	@cd backend/migration && go run migrate.go down