help:
	- start: Run the application
	- test: Run the tests
	- build: Build the application
	- clean: Clean the application


start:
	- cd docker && docker compose up && go run main.go --config=./config/config.yml 
