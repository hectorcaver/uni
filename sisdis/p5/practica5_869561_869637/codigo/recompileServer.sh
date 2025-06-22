#!/bin/bash

# Salir si hay un error en cualquier parte del script
set -e

# Ruta base del proyecto
PROJECT_ROOT=$(pwd)

# Paso 1: Construir binario sin dependencias C
echo "Compilando servidor Raft..."
cd ./raft/cmd/srvraft
CGO_ENABLED=0 go build -o servidor .

# Paso 2: Mover el binario al Dockerfile
echo "Moviendo binario a carpeta Dockerfile..."
mv servidor "$PROJECT_ROOT/Dockerfiles/servidor"

# Paso 3: Construir la imagen Docker
echo "Construyendo imagen Docker..."
cd "$PROJECT_ROOT/Dockerfiles/servidor"
docker build . -t localhost:5001/servidor:latest

# Paso 4: Subir imagen al registro local
echo "Subiendo imagen al registro local..."
docker push localhost:5001/servidor:latest

echo "Todo listo. Imagen disponible como localhost:5001/servidor:latest"
