#!/bin/bash

if [ "$#" -lt 4 ]; then
    echo "Usage: $0 <number_of_air_conditioners> <number_of_fridges> <number_of_furnaces> <number_of_lights>"
    exit 1
fi

NETWORK_NAME="suu"

# Air Conditioner
docker build -t air_conditioner ./air_conditioner

for ((i=1; i<=$1; i++)); do
    AIR_CONDITIONER_NAME="air-conditioner$i"

    docker run --net "$NETWORK_NAME" -e AIR_CONDITIONER_NAME="$AIR_CONDITIONER_NAME" -e AIR_CONDITIONER_ROOM="$AIR_CONDITIONER_ROOM" -d --name "$AIR_CONDITIONER_NAME" air_conditioner &
done

# Fridge
docker build -t  fridge ./fridge

for ((i=1; i<=$2; i++)); do
    rooms=("Kitchen" "Bathroom")
    random_room=${rooms[RANDOM % ${#rooms[@]}]}

    FRIDGE_NAME="fridge$i"
    FRIDGE_ROOM="$random_room"

    docker run --net "$NETWORK_NAME" -e FRIDGE_NAME="$FRIDGE_NAME" -e FRIDGE_ROOM="$FRIDGE_ROOM" -d --name "$FRIDGE_NAME" fridge &
done

# Furnance
docker build -t furnance ./furnance

for ((i=1; i<=$3; i++)); do
    rooms=("Kitchen" "Bathroom")
    random_room=${rooms[RANDOM % ${#rooms[@]}]}

    FURNANCE_NAME="furnance$i"
    FURNANCE_ROOM="$random_room"

    docker run --net "$NETWORK_NAME" -e FURNANCE_NAME="$FURNANCE_NAME" -e FURNANCE_ROOM="$FURNANCE_ROOM" -d --name "$FURNANCE_NAME" furnance &
done

# Light
docker build -t light ./light

for ((i=1; i<=$4; i++)); do
    rooms=("Kitchen" "LivingRoom", "BoilerRoom")
    random_room=${rooms[RANDOM % ${#rooms[@]}]}

    LIGHT_NAME="light$i"
    LIGHT_ROOM="$random_room"

    docker run --net "$NETWORK_NAME" -e LIGHT_NAME="$LIGHT_NAME" -e LIGHT_ROOM="$LIGHT_ROOM" -d --name "$LIGHT_NAME" light &
done

wait
