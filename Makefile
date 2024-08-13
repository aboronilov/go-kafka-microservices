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

# invoicer:
# 	@echo "Building invoicer binary..."
# 	@go build -o bin/invoicer ./invoicer
# 	@echo "Startig INVOICER app"
# 	@./bin/invoicer

aggregator:
	@echo "Building aggregatoe binary..."
	@go build -o bin/aggregator ./aggregator
	@echo "Startig AGGREGATOR app"
	@./bin/aggregator

kafka:
	@echo "Runninig kafka with zookeeper..."
	@docker compose up -d

.PHONY: obu receiver calculator kafka aggregator