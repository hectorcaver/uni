#!/bin/bash
# Autor: Adrián Nasarre Sánchez 869561
# Autor: Héctor Lacueva Sacristán 869637
# Script para iniciar múltiples nodos de srvraft de forma ordenada

# Definir variables
BINARY=./bin/srvraft
HOST=127.0.0.1
PORTS=(29001 29002 29003)

# Detener procesos previos
pkill srvraft

# Iniciar nodos
for i in 0 1 2; do
    echo "Iniciando nodo $i en ${HOST}:${PORTS[$i]}..."
    $BINARY $i "${HOST}:${PORTS[0]}" "${HOST}:${PORTS[1]}" "${HOST}:${PORTS[2]}" &
done

echo "Todos los nodos han sido iniciados."