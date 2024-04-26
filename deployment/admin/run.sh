#!/bin/bash

docker build -t admin .

docker run --net suu -d --name admin admin

wait
