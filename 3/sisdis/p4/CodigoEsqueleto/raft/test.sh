#!/bin/bash

PATH="$HOME/uni/sisdis/p4/CodigoEsqueleto/raft"

echo "Inicio de ciclo de pruebas en: $PATH"
sleep 1

while true; do
    echo "===================================="
    echo "🧪 Ejecutando tests con 'go test'..."
    cd "$PATH" && /usr/local/go/bin/go test -v ./...

    echo "🧹 Matando procesos con 'clearProcess.sh'..."
    cd "$PATH" && source clearProcess.sh

    echo "🧼 Limpiando caché de pruebas con 'go clean -testcache'..."
    cd "$PATH" && /usr/local/go/bin/go clean -testcache

    echo "⏳ Esperando 2 segundos antes del siguiente ciclo..."
    sleep 2
done
