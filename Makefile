gate:
	@echo "Building gateway binary..."
	@go build -o bin/gate gateway/main.go
	@echo "Startig GATEWAY app"
	@./bin/gate

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

prometheus:
	../prometheus/prometheus --config.file=./.config/prometheus.yml

aggregator:
	@echo "Building aggregator binary..."
	@go build -o bin/aggregator ./aggregator
	@echo "Startig AGGREGATOR app"
	@./bin/aggregator

docker:
	@echo "Runninig instrumental services..."
	@docker compose up -d

proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative types/ptypes.proto

.PHONY: obu receiver calculator docker aggregator gate prometheus