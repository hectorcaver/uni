#!/bin/sh
# Limpieza completa de recursos Docker y clúster kind
# Adrián Nasarre Sánchez 869561
# Héctor Lacueva Sacristán 869637

SEPARATOR="=================================================="

echo "$SEPARATOR"
echo "Parando y borrando todos los contenedores activos..."
CONTAINERS=$(docker ps -aq)
if [ -n "$CONTAINERS" ]; then
  docker stop $CONTAINERS 2>/dev/null
  docker rm $CONTAINERS 2>/dev/null
else
  echo "No existen contenedores para borrar."
fi

echo "$SEPARATOR"
echo "Eliminando todas las imágenes almacenadas en Docker..."
IMAGES=$(docker images -q)
if [ -n "$IMAGES" ]; then
  docker rmi -f $IMAGES 2>/dev/null
else
  echo "No se encontraron imágenes para eliminar."
fi

echo "$SEPARATOR"
echo "Borrando todos los volúmenes de Docker..."
VOLUMES=$(docker volume ls -q)
if [ -n "$VOLUMES" ]; then
  docker volume rm $VOLUMES 2>/dev/null
else
  echo "No hay volúmenes disponibles para borrar."
fi

echo "$SEPARATOR"
echo "Limpiando redes Docker no utilizadas..."
docker network prune -f 2>/dev/null

echo "$SEPARATOR"
echo "Destruyendo el clúster kind si existe..."
kind delete cluster 2>/dev/null

echo "$SEPARATOR"
echo "Limpieza finalizada."
