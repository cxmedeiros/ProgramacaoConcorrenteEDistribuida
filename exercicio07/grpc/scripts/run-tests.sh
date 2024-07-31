#!/bin/bash
COUNTS=("300" "3000" "30000")

for COUNT in "${COUNTS[@]}"; do
    echo "Running tests with count=${COUNT}"

    export COUNT

    echo "Starting server..."
    go run ./server/main.go &
    SERVER_PID=$!

    sleep 5

    echo "Starting client..."
    go run ./client/main.go &
    CLIENT_PID=$!

    wait $CLIENT_PID

    echo "Checking for processes using port 50051..."
    PORT_PID=$(lsof -ti :50051)
    if [ -n "$PORT_PID" ]; then
        echo "Killing process $PORT_PID using port 50051"
        kill $PORT_PID
    else
        echo "No process found using port 50051"
    fi

    sleep 5

    echo "Done!"
    echo "--------"
    echo "--------"
done
