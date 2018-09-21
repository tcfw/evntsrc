#!/bin/bash

openssl req -x509 -newkey rsa:2048 -keyout priv.pem -out pub.pem -days 365 -nodes;