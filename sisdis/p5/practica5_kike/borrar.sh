#!/bin/sh
# Jorge Gallardo y Enrique Baldovin
# Detener y eliminar todos los contenedores
echo "Deteniendo y eliminando todos los contenedores..."
docker stop $(docker ps -aq) 2>/dev/null
docker rm $(docker ps -aq) 2>/dev/null

# Eliminar todas las imágenes
echo "Eliminando todas las imágenes de Docker..."
docker rmi -f $(docker images -q) 2>/dev/null

# Eliminar todos los volúmenes
echo "Eliminando todos los volúmenes de Docker..."
docker volume rm $(docker volume ls -q) 2>/dev/null

# Limpiar redes no utilizadas
echo "Eliminando redes no utilizadas..."
docker network prune -f 2>/dev/null

# Eliminar el clúster de kind
echo "Eliminando el clúster de kind..."
kind delete cluster

echo "Operaciones completadas."
