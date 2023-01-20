TARGET_OS = "windows" "linux" "darwin"
TARGET_ARCH = "amd64" "arm64"

build:
	for os in $(TARGET_OS); do \
		for arch in $(TARGET_ARCH); do \
			env GOOS=$$os GOARCH=$$arch go build -o bin/$$os/$$arch/go src/main.go; \
		done; \
	done;