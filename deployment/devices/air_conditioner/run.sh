#!/bin/bash

if [ "$#" -lt 1 ]; then
    echo "Użycie: $0 <ilość_kontenerów>"
    exit 1
fi

docker build -t air_conditioner .

for ((i=1; i<=$1; i++)); do
    AIR_CONDITIONER_NAME="air-conditioner-$i"

    docker run --net suu -e AIR_CONDITIONER_NAME="$AIR_CONDITIONER_NAME" -d air_conditioner &
done

wait
