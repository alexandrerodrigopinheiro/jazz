#!/bin/bash

# Print header
echo "========================="
echo "Running Jazz Project Tests"
echo "========================="

# Run tests for cache and other backend components
go test -v ./backend/pkg/cache/...
go test -v ./backend/pkg/database/...
go test -v ./backend/pkg/logger/...
# go test -v ./backend/configs/...

# Check if tests passed or failed
if [ $? -eq 0 ]; then
    echo "========================="
    echo "All tests passed successfully!"
    echo "========================="
else
    echo "========================="
    echo "Some tests failed. Please check the output above for details."
    echo "========================="
    exit 1
fi
