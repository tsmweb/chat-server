#!/bin/bash

openssl genrsa -out server.pem 4096
openssl rsa -in server.pem -pubout > server.pub

openssl req -new -x509 -sha256 -key server.pem -out server.crt -days 365
