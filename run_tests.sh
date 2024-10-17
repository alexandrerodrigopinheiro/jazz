#!/bin/bash

# Print header
echo "========================="
echo "Running Jazz Project Tests"
echo "========================="

# Navigate to backend directory where tests are located
cd backend

# Run tests for cache and other backend components
go test -v ./cache/...

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
