#!/bin/bash

openssl genrsa -out ssl/ca.key 2048
openssl req -new -x509 -days 3650 -key ssl/ca.key -out ssl/ca.crt -subj "/CN=yngwie proxy CA"
openssl genrsa -out ssl/cert.key 2048
