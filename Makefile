BINARY_NAME=celeritasApp

build:
	@go mod vendor
	@echo "Building Celeritas..."
	@go build -o tmp/$(BINARY_NAME)
	@echo "Celeritas built successfully!"

run: build
	@echo "Starting Celeritas..."
	@./tmp/$(BINARY_NAME) &
	@echo "Celeritas started successfully!"

clean:
	@echo "Cleaning up..."
	@go clean
	@rm tmp/$(BINARY_NAME)
	@echo "Cleaned up successfully!"

start_compose:
	docker-compose up -d

stop_compose:
	docker-compose down

test:
	@echo "Running tests..."
	@go test ./...
	@echo "Done!"

start: start_compose run

stop: stop_compose
	@echo "Stopping Celeritas..."
	@-pkill -SIGTERM -f "./tmp/$(BINARY_NAME)"
	@echo "Celeritas stopped successfully!"

restart: stop start