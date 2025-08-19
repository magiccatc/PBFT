#!/bin/bash

# Ubuntu shell script to kill processes by IP and port

ports=("8000" "8001" "8002" "8003"  "8004" "8005" "8006" "8007" "8008" "8009" "8010" 
"8011" "8012" "8013" "8014" "8015" "8016" "8017" "8018" "8019" "8020" 
"8021" "8022" "8023" "8024" "8025" 
"8888") # List of ports to check

for port in "${ports[@]}"; do
    # Get the process ID (PID) using netstat and awk
    pid=$(netstat -tulnp | grep ":$port " | awk '{print $7}' | cut -d/ -f1)

    # Check if a PID was found
    if [ ! -z "$pid" ]; then
        # Kill the process using the PID
        sudo kill -9 "$pid"
        echo "Process with PID $pid on port $port has been terminated."
    else
        echo "No process found on port $port."
    fi
done