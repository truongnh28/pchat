#!/usr/bin/env bash

# This script is used to build docker-image of this project

ver=$1
docker image rm truongnh28/chat-app
docker build --tag truongnh28/chat-app:$ver .
docker push truongnh28/chat-app