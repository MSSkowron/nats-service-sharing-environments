#!/bin/bash

docker build -t admin .

docker run --net suu -d admin

wait
