
.PHONE: build
build:
	@echo "Building..."
	@go build -o bin/ ./cmd/...

.PHONE: install
install:
	@echo "Installing..."
	@go install ./cmd/...

.PHONE: uninstall
uninstall:
	@echo "Uninstalling..."
	@go clean -i ./cmd/...

.PHONE: clean
clean:
	@echo "Cleaning..."
	@rm -rf bin/