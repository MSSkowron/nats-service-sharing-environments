#!/bin/bash

if [ "$#" -lt 1 ]; then
    echo "Użycie: $0 <ilość_kontenerów>"
    exit 1
fi

docker build -t  fridge .

for ((i=1; i<=$1; i++)); do
    rooms=("Kitchen" "Bathroom")
    random_room=${rooms[RANDOM % ${#rooms[@]}]}

    FRIDGE_NAME="fridge-$i"
    FRIDGE_ROOM="$random_room"

    docker run --net suu -e FRIDGE_NAME="$FRIDGE_NAME" -e FRIDGE_ROOM="$FRIDGE_ROOM" -d  fridge &
done

wait
