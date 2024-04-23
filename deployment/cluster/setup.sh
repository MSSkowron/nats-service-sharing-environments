#!/bin/sh

nats context save nats1 --server nats://admin:admin@localhost:4222 --description 'NATS Server localhost:4222'
nats context save nats2 --server nats://admin:admin@localhost:4223 --description 'NATS Server localhost:4223'
nats context save nats3 --server nats://admin:admin@localhost:4224 --description 'NATS Server localhost:4224'

nats context select nats1
nats stream add --cluster natscluster --subjects "lights.>,airconditionersturn,airconditionersset,fridges.>,furnances.>,publishers.*" --retention limits --max-msg-size 1MB --max-msgs 1000 --max-bytes 1GB --storage file --discard old --replicas 3 devices
