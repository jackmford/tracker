# Variables
BINARY_NAME=TimeTracker
APP_NAME=$(BINARY_NAME).app
DEST_APP=/Applications/$(APP_NAME)
WORKING_DIR=/Users/jackfordyce/Lab/git-repos/tracker

.PHONY: all build install clean serve

# Default: build and install to Applications
all: build install

build:
	@echo "🔨 Building binary..."
	go build -ldflags="-s -w" -o $(BINARY_NAME) main.go

install:
	@echo "📦 Creating macOS App Bundle..."
	rm -rf $(DEST_APP)
	mkdir -p $(DEST_APP)/Contents/MacOS
	cp $(BINARY_NAME) $(DEST_APP)/Contents/MacOS/
	@echo "✅ Installed to $(DEST_APP)"

serve:
	@echo "📊 Starting dashboard server from project directory..."
	cd $(WORKING_DIR) && ./$(BINARY_NAME) serve

clean:
	@echo "🧹 Cleaning up..."
	rm -f $(BINARY_NAME)
	rm -rf $(DEST_APP)
