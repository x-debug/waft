OUTPUT=build

.PHONY: clean all run install
all: clean test build

clean:
	@echo -e "\nCLEANING $(OUTPUT) DIRECTORY"
	rm -rf ./$(OUTPUT)

build: clean test
	@echo -e "\nBUILDING $(OUTPUT)/waft BINARY"
	mkdir -p ./$(OUTPUT) && go build -o $(OUTPUT)/waft waft/main.go

test:
	@echo -e "\nTESTING"
	go test -v ./... -coverprofile=coverage.out

run:
	go run waft/main.go

install: build
	sudo cp $(OUTPUT)/waft /usr/local/bin/