#!/bin/bash

# Salir si hay un error en cualquier parte del script
set -e

# Ruta base del proyecto
PROJECT_ROOT=$(pwd)

# Paso 1: Construir binario sin dependencias C
echo "Compilando cliente Raft..."
cd ./raft/pkg/cltraft
CGO_ENABLED=0 go build -o cliente .

# Paso 2: Mover el binario al Dockerfile
echo "Moviendo binario a carpeta Dockerfile..."
mv cliente "$PROJECT_ROOT/Dockerfiles/cliente"

# Paso 3: Construir la imagen Docker
echo "Construyendo imagen Docker..."
cd "$PROJECT_ROOT/Dockerfiles/cliente"
docker build . -t localhost:5001/cliente:latest

# Paso 4: Subir imagen al registro local
echo "Subiendo imagen al registro local..."
docker push localhost:5001/cliente:latest

echo "Todo listo. Imagen disponible como localhost:5001/cliente:latest"