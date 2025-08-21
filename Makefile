BINARY_NAME=lsc
BUILD_DIR=bin

# Default: build for current OS
build:
	@echo "Building $(BINARY_NAME) for local OS..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) main.go

# Cross-platform builds
build-all: clean
	@echo "Building for Linux..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux main.go
	@echo "Building for macOS..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME)-mac main.go
	@echo "Building for Windows..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME).exe main.go

build-windows:
	@echo "Building $(BINARY_NAME) for Windows..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME).exe main.go


# Remove build artifacts
clean:
	@echo "Cleaning build directory..."
	@rm -rf $(BUILD_DIR)
	@mkdir -p $(BUILD_DIR)

clean-windows:
	@echo "Cleaning Windows build..."
	@powershell -Command "Remove-Item -Path '$(BUILD_DIR)' -Force -Recurse"

.PHONY: build build-all build-windows clean clean-windows 