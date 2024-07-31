#!/bin/bash
COUNTS=("100" "1000" "10000")

wait_for_rabbitmq() {
    echo "Waiting for RabbitMQ to be ready..."
    local url="http://localhost:15672/api/overview"
    local retries=30
    local wait_time=5

    for ((i=0; i<$retries; i++)); do
        response=$(curl -s -o /dev/null -w "%{http_code}" -u "guest:guest" "$url")
        if [ "$response" -eq 200 ]; then
            echo "RabbitMQ is ready."
            return
        else
            echo "Retry $((i + 1)): RabbitMQ is not ready yet (HTTP code: $response). Waiting $wait_time seconds..."
            sleep $wait_time
        fi
    done

    echo "Failed to connect to RabbitMQ after $((retries * wait_time)) seconds."
    exit 1
}

for COUNT in "${COUNTS[@]}"; do
    echo "Running tests with count=${COUNT}"

    docker run -d --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.13-management

    wait_for_rabbitmq

    export COUNT

    echo "Starting server..."
    go run ./server/main.go &
    SERVER_PID=$!

    sleep 5

    echo "Starting client..."
    go run ./client/main.go &
    CLIENT_PID=$!

    wait $CLIENT_PID

    docker stop rabbitmq

    echo "Checking for processes using port 5672..."
    PORT_PID=$(lsof -ti :5672)
    if [ -n "$PORT_PID" ]; then
        echo "Killing process $PORT_PID using port 5672"
        kill $PORT_PID
    else
        echo "No process found using port 5672"
    fi

    sleep 10

    echo "Done!"
    echo "--------"
    echo "--------"
done
