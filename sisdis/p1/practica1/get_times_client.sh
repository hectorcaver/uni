#!/bin/bash

# This script is in charge of getting data from the execution of some commands
# and storing it in a file.

# First, as flags, we receive -t for the TCP connection and -u for the UDP
# connection. We also receive the number of times we want to execute the
# command.

# Another flag will be whether we want to execute the client and server in the same host or in a different one. 
# If we want to execute in the same host, we will receive -l. F.e. ./get_times_client.sh -t 10 -l 192.168.3.x:8080
# If we want to execute in a different host, we will receive -r. F.e ./get_times_client.sh -t 10 -r 192.168.3.x:8080

# We will store the data in a file called times_{t || u}_10_{l || r}.txt.


# Usage: ./get_times.sh -t|-u <number_of_executions> -l|-r <host1:port> 

# Check if enough arguments are provided
if [ $# -ne 4 ]; then
    echo "Usage: $0 -t|-u <number_of_executions> -l|-r <host1:port>"
    exit 1
fi

# Parse arguments
protocol=""
exec_count=0
mode=""
host1=""

# path="/misc/practicas/alumnos/sd2425..."
path="~/uni/sisdis/p1/practica1/"

while [[ "$1" != "" ]]; do
    case $1 in
        -t) protocol="t" ;;  # TCP
        -u) protocol="u" ;;  # UDP
        -l) mode="l" ;;      # Local mode
        -r) mode="r" ;;      # Remote mode
        *)
            if [[ $exec_count -eq 0 ]]; then
                exec_count=$1
            elif [[ -z $host1 ]]; then
                host1=$1
            else
                echo "Unexpected argument: $1"
                exit 1
            fi
            ;;
    esac
    shift
done

# Validate required arguments
if [[ -z $protocol || -z $exec_count || -z $mode || -z $host1 ]]; then
    echo "Missing required arguments."
    exit 1
fi

# Define output file
output_file="times_${protocol}_${exec_count}_${mode}.txt"
echo "Storing results in $output_file"

results=""
results2=""
tryouts=""
spaces=""

ip=$(echo "$host1" | cut -d':' -f1)

if [[ $protocol == "u" ]]; then

    echo "Estoy en la mierda"
    
    ssh "$ip" "bash --login -c 'cd $path && ./get_times_server.sh -u$exec_count $host1'" &
    
    sleep 2

fi

# Execute commands
for ((i = 1; i <= exec_count; i++)); do
    tryouts="$tryouts|Prueba $i"  # Concatenar correctamente la variable tryouts
    spaces="$spaces|:-:"

    if [[ $protocol == "u" ]]; then
        output=$(go run network/client_udp/client_udp.go $host1)

        # Extraer el número decimal con grep y sed
        decimal_number=$(echo "$output" | grep -oP '(?<=: )[0-9]+\.[0-9]+[^ 0-9]*')

        # Reemplazar el punto por coma y añadirlo a la variable results
        results="$results|${decimal_number//./,}"  # Concatenar correctamente la variable results
    
    else

        ssh "$ip" "bash --login -c 'cd $path && ./get_times_server.sh -t $host1'" &

        output=$(go run network/client_tcp/tcp_client.go "$host1" | tee /dev/tty)

        # Extraer el número decimal con grep y sed
        decimal_number=($(echo "$output" | grep -oP '(?<=: )[0-9]+\.[0-9]+[^ 0-9]*'))

        first_decimal="${decimal_number[0]}"
        second_decimal="${decimal_number[1]}"

        # Reemplazar el punto por coma y añadirlo a la variable results
        results="$results|${first_decimal//./,}"  # Concatenar correctamente la variable results
        results2="$results2|${second_decimal//./,}"

        sleep 5
    fi

done

tryouts="$tryouts|"  # Concatenar correctamente la variable tryouts
results="$results|"  # Concatenar correctamente la variable results
spaces="$spaces|"   # Concatenar

if [[ $protocol == "t" ]]; then
    results2="$results2|"
fi

if [ ! -d "results" ]; then
    mkdir results
fi

echo "$tryouts" >> "results/$output_file"
echo "$spaces" >> "results/$output_file"
echo "$results" >> "results/$output_file"
echo "$results2" >> "results/$output_file"

echo "Execution completed. Results saved to $output_file"

