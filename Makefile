obu:
	@echo "Building binary..."
	@go build -o bin/obu obu/main.go
	@echo "Startig OBU app"
	@./bin/obu

receiver:
	@echo "Building binary..."
	@go build -o bin/receiver data_receiver/main.go
	@echo "Startig RECEVIER app"
	@./bin/receiver

kafka:
	@echo "Runninig kafka with zookeeper..."
	@docker compose up -d

.PHONY: obu receiver