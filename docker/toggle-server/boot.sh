#!/bin/bash

echo "Installing service dependencies..."
cd /toggle/server

go mod download
fswatch -config /scripts/.fsw.yml
