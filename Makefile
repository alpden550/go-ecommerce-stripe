GOSTRIPE_PORT=4000
API_PORT=4001
INVOICE_PORT=5001

## build: builds all binaries
build: clean build_front build_back build_invoice
	@printf "All binaries built!\n"

## clean: cleans all binaries and runs go clean
clean:
	@echo "Cleaning..."
	@- rm -f dist/*
	@go clean
	@echo "Cleaned!"

## build_front: builds the front end
build_front:
	@echo "Building front end..."
	@go build -o dist/gostripe ./cmd/web
	@echo "Front end built!"

## build_back: builds the back end
build_back:
	@echo "Building back end..."
	@go build -o dist/gostripe_api ./cmd/api
	@echo "Back end built!"

build_invoice:
	@echo "Building invoice microservice..."
	@go build -o dist/gostripe_invoice ./cmd/micro/invoicer
	@echo "Invoicer built!"

## start: starts front and back end
start: db start_front start_back start_invoice

## start_front: starts the front end
start_front: build_front
	@echo "Starting the front end..."
	./dist/gostripe -port=${GOSTRIPE_PORT} &
	@echo "Front end running!"

## start_back: starts the back end
start_back: build_back
	@echo "Starting the back end..."
	./dist/gostripe_api -port=${API_PORT} &
	@echo "Back end running!"

start_invoice: build_invoice
	@echo "Starting the Invoice microservice..."
	./dist/gostripe_invoice -port=${INVOICE_PORT} &
	@echo "Back end running!"

## stop: stops the front and back end
stop: stop_front stop_back stop_db stop_invoice
	@echo "All applications stopped"

## stop_front: stops the front end
stop_front:
	@echo "Stopping the front end..."
	@-pkill -SIGTERM -f "gostripe -port=${GOSTRIPE_PORT}"
	@echo "Stopped front end"

## stop_back: stops the back end
stop_back:
	@echo "Stopping the back end..."
	@-pkill -SIGTERM -f "gostripe_api -port=${API_PORT}"
	@echo "Stopped back end"

stop_invoice:
	@echo "Stopping the invoicer..."
	@-pkill -SIGTERM -f "gostripe_invoice -port=${INVOICE_PORT}"
	@echo "Stopped invoicer"

db:
	@echo "Starting the postgres..."
	docker-compose up db -d

stop_db:
	@echo "Stopping the postgres..."
	docker-compose stop db

air_api:
	air -c .air-api.toml

air_app:
	air -c .air-app.toml

air_invoice:
	air -c .air-invoice.toml
