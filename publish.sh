#!/usr/bin/env sh

docker build -t ghcr.io/bweston92/http-to-pubsub:latest .

docker push ghcr.io/bweston92/http-to-pubsub:latest
