#!/bin/bash

if [ "$#" -lt 1 ]; then
    echo "Użycie: $0 <ilość_kontenerów>"
    exit 1
fi

docker build -t furnance .

for ((i=1; i<=$1; i++)); do
    rooms=("Kitchen" "Bathroom")
    random_room=${rooms[RANDOM % ${#rooms[@]}]}

    FURNANCE_NAME="furnance-$i"
    FURNANCE_ROOM="$random_room"

    docker run --net suu -e FURNANCE_NAME="$FURNANCE_NAME" -e FURNANCE_ROOM="$FURNANCE_ROOM" -d furnance &
done

wait
