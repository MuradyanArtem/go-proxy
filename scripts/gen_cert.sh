#!/bin/bash

openssl req -new -key ./ssl/cert.key -subj "/CN=$1" -sha256 | \
openssl x509 -req -days 3650 -CA ./ssl/ca.crt -CAkey ./ssl/ca.key -set_serial "$2" > ./ssl/"$1".crt
