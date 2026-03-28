.PHONY: build build-windows run clean test

# Default: build for current platform
build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o claude-usage-tracker ./cmd/widget

# Cross-compile for Windows (browser mode only, no CGo)
build-windows:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o claude-usage-tracker.exe ./cmd/widget

# Run in browser mode
run: build
	./claude-usage-tracker --browser

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -f claude-usage-tracker claude-usage-tracker.exe
