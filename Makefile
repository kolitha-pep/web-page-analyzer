APP_NAME=web-analyzer
MAIN=main.go
MAIN_DIR=./cmd/server

build:
	@echo "Building the application..."
	go build -o ./bin/$(APP_NAME) $(MAIN_DIR)/$(MAIN)

run: build
	@echo "Running the application..."
	./bin/$(APP_NAME)

clean:
	@echo "Cleaning up..."
	rm -f ./bin/$(APP_NAME)

docker-build:
	@echo "Building Docker image..."
	docker build -t $(APP_NAME) .

docker-run:
	@echo "Running Docker container..."
	docker run --rm -p 8080:8080 $(APP_NAME)