#!/bin/bash
# Autores: Jorge Gallardo y Enrique Baldovin
set -euo pipefail

#---------------------------------------------------
# Configuraciones iniciales
#---------------------------------------------------
# Asegurar que el PATH contenga go
export PATH=$PATH:/usr/local/go/bin

#---------------------------------------------------
# Eliminar posible cluster viejo
#---------------------------------------------------
echo "Eliminando el clúster anterior..."
if [ -f ./borrar.sh ]; then
    ./borrar.sh
else
    echo "El script borrar.sh no existe. Asegúrate de que esté en este directorio."
    exit 1
fi

#---------------------------------------------------
# Crear el clúster
#---------------------------------------------------
echo "Creando clúster..."
if [ -f ./kind-with-registry.sh ]; then
    ./kind-with-registry.sh
else
    echo "El script kind-with-registry.sh no existe. Asegúrate de que esté en este directorio."
    exit 1
fi

#---------------------------------------------------
# Eliminar ejecutables antiguos
#---------------------------------------------------
echo "Eliminando ejecutables antiguos..."
if [ -f Dockerfiles/servidor/servidor ]; then
    rm Dockerfiles/servidor/servidor
fi

if [ -f Dockerfiles/cliente/cliente ]; then
    rm Dockerfiles/cliente/cliente
fi

#---------------------------------------------------
# Compilar cliente y servidor
#---------------------------------------------------
echo "Compilando cliente y servidor..."
if [ -d raft/cmd/srvraft ]; then
    cd raft/cmd/srvraft
    CGO_ENABLED=0 go build -o ../../../Dockerfiles/servidor/servidor .
    cd ../../pkg/cltraft
    CGO_ENABLED=0 go build -o ../../../Dockerfiles/cliente/cliente .
    cd ../../../
else
    echo "No se encontró el directorio raft/cmd/srvraft. Asegúrate de tener la estructura de directorios correcta."
    exit 1
fi

#---------------------------------------------------
# Crear imágenes Docker y pushear al registro local
#---------------------------------------------------
echo "Creando e impulsando imágenes Docker..."
cd Dockerfiles/servidor
docker build . -t localhost:5001/servidor:latest
docker push localhost:5001/servidor:latest

cd ../cliente
docker build . -t localhost:5001/cliente:latest
docker push localhost:5001/cliente:latest

cd ../..

#---------------------------------------------------
# Iniciar Kubernetes con el statefulset
#---------------------------------------------------
echo "Iniciando statefulset..."
if [ -f statefulset_go.yaml ]; then
    kubectl create -f statefulset_go.yaml
else
    echo "No se encontró statefulset_go.yaml. Asegúrate de que esté en este directorio."
    exit 1
fi

#---------------------------------------------------
# Esperar hasta que el pod del cliente esté listo
#---------------------------------------------------
echo "Esperando a que el pod 'client' esté en estado Running..."
# Esperar a que el pod 'client' esté en Running (ajusta el nombre si es necesario).
kubectl wait --for=condition=Ready pod/client --timeout=120s

#---------------------------------------------------
# Ejecutar el cliente dentro del contenedor
#---------------------------------------------------
echo "Ejecutando cliente dentro del contenedor..."
kubectl exec client -ti -- sh -c "/usr/local/bin/cliente"
