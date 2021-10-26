OUTPUT=build

.PHONY: clean all run install
all: clean test build

clean:
	@echo -e "\nCLEANING $(OUTPUT) DIRECTORY"
	rm -rf ./$(OUTPUT)

build: clean test
	@echo -e "\nBUILDING $(OUTPUT)/waft BINARY"
	mkdir -p ./$(OUTPUT) && go build -ldflags \
	"-X main.buildTime=`date +%Y-%m-%d,%H:%M:%S` -X main.buildVer=1.0.0 -X main.gitCommitId=`git rev-parse HEAD`" \
	-o $(OUTPUT)/waft waft/main.go && go build -o $(OUTPUT)/test_backend test_backend/main.go

test:
	@echo -e "\nTESTING"
	go test -v ./... -coverprofile=coverage.out

run:
	go run waft/main.go

install: build
	sudo cp $(OUTPUT)/waft /usr/local/bin/