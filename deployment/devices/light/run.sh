#!/bin/bash

if [ "$#" -lt 1 ]; then
    echo "Użycie: $0 <ilość_kontenerów>"
    exit 1
fi

docker build -t light .

for ((i=1; i<=$1; i++)); do
    rooms=("Kitchen" "Living Room", "Boiler Room")
    random_room=${rooms[RANDOM % ${#rooms[@]}]}

    LIGHT_NAME="light-$i"
    LIGHT_ROOM="$random_room"

    docker run --net suu -e LIGHT_NAME="$LIGHT_NAME" -e LIGHT_ROOM="$LIGHT_ROOM" -d light &
done

wait
