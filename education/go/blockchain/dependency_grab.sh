#!/bin/bash

# spew formats structs and slices for the screen.
go get github.com/davecgh/go-spew/spew

# provides a web handler
go get github.com/gorilla/mux

# allows you to access .env files without hard-coding paths
go get github.com/joho/godotenv

# creates the port environment file
touch blockchain.env
echo ADDR=8080 > .env

echo "All Done!"
