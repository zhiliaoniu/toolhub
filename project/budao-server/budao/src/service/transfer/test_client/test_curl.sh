#!/bin/bash
curl --request "POST" \
    --location "http://localhost:8080/twirp/twirp.example.haberdasher.Haberdasher/MakeHat" \
    --header "Content-Type:application/json" \
    --data '{"inches": 10}' \
    --verbose
