# 1. Compile the binary with 'stealth' flags
go build -ldflags="-s -w" -o TimeTracker main.go

# 2. Clear the old version from Applications (if it exists)
rm -rf /Applications/TimeTracker.app

# 3. Create the new Bundle structure
mkdir -p /Applications/TimeTracker.app/Contents/MacOS

# 4. Move the binary into the Bundle
mv TimeTracker /Applications/TimeTracker.app/Contents/MacOS/
