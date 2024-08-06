obu:
	@echo "Building obu binary..."
	@go build -o bin/obu obu/main.go
	@echo "Startig OBU app"
	@./bin/obu

receiver:
	@echo "Building receiver binary..."
	@go build -o bin/receiver ./data_receiver
	@echo "Startig RECEVIER app"
	@./bin/receiver

calculator:
	@echo "Building calculator binary..."
	@go build -o bin/calculator ./distance_calculator
	@echo "Startig CALCULATOR app"
	@./bin/calculator

kafka:
	@echo "Runninig kafka with zookeeper..."
	@docker compose up -d

.PHONY: obu receiver calculator kafka