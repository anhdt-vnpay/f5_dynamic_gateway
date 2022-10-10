START_MISSING?=0
gen_type:
	./schema/scripts/gen_service.sh

gen_example_type:
	./examples/schema/ping/gen_ping.sh

start_example_api:
	go run main.go example_ping_start_api

start_example_gateway:
	go run main.go example_ping_start_gateway