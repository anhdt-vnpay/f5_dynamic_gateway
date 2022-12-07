START_MISSING?=0
gen_example_type:
	./schema/ping/gen_ping.sh

start_example_api:
	go run main.go example_ping_start_api

start_example_gateway:
	go run main.go example_ping_start_gateway
	
start_example_echo:
	go run main.go example_ping_start_echo_server

start_example_echo2:
	go run main.go example_ping_start_echo_server2