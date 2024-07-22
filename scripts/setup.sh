#!/bin/bash

# Initialize the Go module
go mod init github.com/scilcon2/revolut_homework || true

# Set the desired Go version
go mod edit -go=1.19

# Add dependencies
go get github.com/gorilla/mux@v1.8.0
go get github.com/lib/pq@v1.10.2

# Tidy up the dependencies
go mod tidy

echo "Go module initialized and dependencies installed."
