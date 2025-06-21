#!/bin/bash
# Lanzador automatizado de entorno distribuido
# Autor: Adrián Nasarre Sánchez 869561
# Autor: Héctor Lacueva Sacristán 869637
set -euo pipefail

SEPARATOR="--------------------------------------------------"

log() {
  echo -e "\n$SEPARATOR"
  echo "$1"
}

# Añadir Go al PATH si es necesario
env_setup() {
  export PATH=$PATH:/usr/local/go/bin
}

# Borrar clúster y recursos previos
clean_previous() {
  log "Limpiando recursos previos..."
  if [ -f ./delete.sh ]; then
    ./delete.sh
  else
    echo "Falta delete.sh. No se puede continuar."; exit 1
  fi
}

# Crear clúster kind
create_cluster() {
  log "Configurando clúster nuevo..."
  if [ -f ./kind-with-registry.sh ]; then
    ./kind-with-registry.sh
  else
    echo "Falta kind-with-registry.sh. Abortando."; exit 1
  fi
}

# Eliminar binarios antiguos
remove_binaries() {
  log "Eliminando binarios previos..."
  rm -f Dockerfiles/servidor/servidor Dockerfiles/cliente/cliente
}

# Compilar los ejecutables
build_binaries() {
  log "Compilando binarios cliente y servidor..."
  if [ -d raft/cmd/srvraft ]; then
    (cd raft/cmd/srvraft && CGO_ENABLED=0 go build -o ../../../Dockerfiles/servidor/servidor .)
    (cd raft/pkg/cltraft && CGO_ENABLED=0 go build -o ../../../Dockerfiles/cliente/cliente .)
  else
    echo "Directorio raft/cmd/srvraft no encontrado."; exit 1
  fi
}

# Construir y subir imágenes Docker
build_and_push_images() {
  log "Construyendo y subiendo imágenes Docker..."
  (cd Dockerfiles/servidor && docker build . -t localhost:5001/servidor:latest && docker push localhost:5001/servidor:latest)
  (cd Dockerfiles/cliente && docker build . -t localhost:5001/cliente:latest && docker push localhost:5001/cliente:latest)
}

# Desplegar en Kubernetes
deploy_k8s() {
  log "Desplegando statefulset en Kubernetes..."
  if [ -f statefulset_go.yaml ]; then
    kubectl create -f statefulset_go.yaml
  else
    echo "Falta statefulset_go.yaml. Abortando."; exit 1
  fi
}

# Esperar a que el pod cliente esté listo
wait_for_client() {
  log "Esperando a que el pod cliente esté disponible..."
  kubectl wait --for=condition=Ready pod/client --timeout=120s
}

# Ejecutar el cliente dentro del pod
run_client() {
  log "Ejecutando cliente en el contenedor..."
  kubectl exec client -ti -- sh -c "/usr/local/bin/cliente"
}

# Secuencia principal
env_setup
clean_previous
create_cluster
remove_binaries
build_binaries
build_and_push_images
deploy_k8s
wait_for_client
run_client
