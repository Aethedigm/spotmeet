BINARY_NAME=celeritasApp.exe

build:
	@go mod vendor
	@echo "Building Celeritas..."
	@go build -o tmp/celeritasApp.exe
	@echo "Celeritas built!"

run: build
	@echo "Starting Celeritas..."
	@cd tmp
	@./celeritasApp.exe
	@echo "Celeritas started!"

clean:
	@echo "Cleaning..."
	@go clean
	@rm tmp/${BINARY_NAME}
	@echo "Cleaned!"

start_compose:
	docker-compose up -d

stop_compose:
	docker-compose down

test:
	@echo "Testing..."
	@go test ./...
	@echo "Done!"

start: run

stop:
	@echo "Stopping Celeritas..."
	@-pkill -SIGTERM -f "./tmp/${BINARY_NAME}"
	@echo "Stopped Celeritas!"

restart: stop start

## cover: opens coverage in browser
cover:
	@go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

## coverage: displays test coverage
coverage:
	@go test -cover ./...