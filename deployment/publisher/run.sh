#!/bin/bash

docker build -t publisher .

docker run --net suu -d publisher

wait
