/****************************************
 * Autor: Adrián Nasarre Sánchez 869561
 * Autor: Héctor Lacueva Sacristán 869637
 * Implementación del servidor raft
 *****************************************/

package main

import (
	"net"
	"net/rpc"
	"os"
	"strconv"
	"strings"

	"raft/internal/comun/check"
	"raft/internal/comun/rpctimeout"
	"raft/internal/raft"
)

const raftServicePort = ":6000"

func main() {
	// Comprobar argumentos
	if len(os.Args) < 2 {
		panic("Uso: <nombre-nodo> (ejemplo: nodo-0)")
	}

	// Extraer nombre y número de nodo
	identificador := os.Args[1]
	partes := strings.Split(identificador, "-")
	if len(partes) < 2 {
		panic("El identificador debe tener formato <nombre>-<indice>")
	}
	nombreBase := partes[0]
	indiceNodo, err := strconv.Atoi(partes[1])
	check.CheckError(err, "Error al convertir el índice del nodo a entero")

	dominio := "raft.default.svc.cluster.local" + raftServicePort

	// Construir direcciones DNS de los nodos
	direcciones := make([]string, 0, 3)
	miDireccion := ""
	for i := 0; i < 3; i++ {
		dir := nombreBase + "-" + strconv.Itoa(i) + "." + dominio
		direcciones = append(direcciones, dir)
		if i == indiceNodo {
			miDireccion = dir
		}
	}

	// Convertir a HostPort
	nodos := rpctimeout.StringArrayToHostPortArray(direcciones)

	// Inicializar nodo Raft
	nodoRaft := raft.NuevoNodo(nodos, indiceNodo, make(chan raft.AplicaOperacion, 1000))
	err = rpc.Register(nodoRaft)
	check.CheckError(err, "Error al registrar el nodo Raft en RPC")

	listener, err := net.Listen("tcp", miDireccion)
	check.CheckError(err, "Error al abrir el listener TCP para el nodo")

	for {
		rpc.Accept(listener)
	}
}
