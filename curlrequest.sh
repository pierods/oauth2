#!/bin/bash

curl -X POST -u "my-client:foobar" -d "grant_type=client_credentials" http://localhost:8080/token
