#!/bin/bash

export EXTENSION_OBSERVER_MONGODB_URI="mongodb://localhost:27017"
export EXTENSION_OBSERVER_MONGODB_DATABASE="extension_observer"
export EXTENSION_OBSERVER_MONGODB_COLLECTION="BBBMark"

go run ./cmd/main.go