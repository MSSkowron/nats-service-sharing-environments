#!/bin/bash

if [ "$#" -lt 1 ]; then
    echo "Usage: $0 <number_of_publishers>"
    exit 1
fi

NETWORK_NAME="suu"

docker build -t publisher .

for ((i=1; i<=$1; i++)); do
    PUBLISHER_NAME="publisher$i"

    docker run --net "$NETWORK_NAME" -e PUBLISHER_NAME="$PUBLISHER_NAME" -d --name "$PUBLISHER_NAME" publisher &
done

wait
