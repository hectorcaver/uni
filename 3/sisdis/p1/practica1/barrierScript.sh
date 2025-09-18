#!/bin/bash

# Validar argumentos
if [ "$#" -ne 1 ]; then
    echo "Sintaxis: ./barrierScript.sh <num_endpoints>"
    exit 1
fi

# Número de endpoints especificado
NUM_ENDPOINTS=$1

# Ejecutar el comando para cada endpoint
for ((i=1; i<=NUM_ENDPOINTS-1; i++)); do
  # Comando a ejecutar con el número de máquina
  COMMAND="go run ./barrier/barrier.go ./barrier/endpoints.txt $i"

  echo "Ejecutando: $COMMAND"

  # Ejecutar el comando localmente
  $COMMAND &
done

COMMAND="go run ./barrier/barrier.go ./barrier/endpoints.txt $NUM_ENDPOINTS"

echo "Ejecutando: $COMMAND"

# Ejecutar el comando localmente
$COMMAND