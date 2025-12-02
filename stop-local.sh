#!/bin/bash

echo "==================================="
echo "Stopping Microservices"
echo "==================================="
echo ""

if [ ! -f .pids ]; then
    echo "⚠️  No PID file found. Trying to stop services by port..."
    echo ""
    
    # Kill processes by port
    for PORT in 6000 7000 8000; do
        PID=$(lsof -ti:$PORT 2>/dev/null)
        if [ -n "$PID" ]; then
            echo "Stopping service on port $PORT (PID: $PID)..."
            kill -15 $PID 2>/dev/null
        fi
    done
else
    echo "Stopping services using saved PIDs..."
    echo ""
    
    # Read PIDs from file and stop them
    while IFS= read -r PID; do
        if [ -n "$PID" ]; then
            if ps -p $PID > /dev/null 2>&1; then
                echo "Stopping process PID: $PID"
                kill -15 $PID 2>/dev/null
            else
                echo "Process PID: $PID already stopped"
            fi
        fi
    done < .pids
    
    # Remove PID file
    rm -f .pids
fi

# Wait a bit for graceful shutdown
sleep 2

echo ""
echo "Checking if services are stopped..."
echo ""

# Force kill if still running
KILLED=0
for PORT in 6000 7000 8000; do
    PID=$(lsof -ti:$PORT 2>/dev/null)
    if [ -n "$PID" ]; then
        echo "⚠️  Service on port $PORT still running. Force killing..."
        kill -9 $PID 2>/dev/null
        KILLED=1
    fi
done

if [ $KILLED -eq 0 ]; then
    echo "✓ All services stopped successfully"
else
    echo "✓ Services stopped (some required force kill)"
fi

echo ""
