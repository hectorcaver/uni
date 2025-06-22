/****************************************
 * Autor: Adrián Nasarre Sánchez 869561
 * Autor: Héctor Lacueva Sacristán 869637
 * Implementación del servidor raft
 *****************************************/

package main

import (
	"fmt"
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

	fmt.Println("Iniciando servidor Raft...")

	// Extraer nombre y número de nodo
	identificador := os.Args[1]
	fmt.Printf("Argumento recibido: %s\n", identificador)

	partes := strings.Split(identificador, "-")
	if len(partes) < 2 {
		panic("El identificador debe tener formato <nombre>-<indice>")
	}
	nombreBase := partes[0]
	indiceNodo, err := strconv.Atoi(partes[1])
	check.CheckError(err, "Error al convertir el índice del nodo a entero")
	fmt.Printf("Nombre base: %s\nÍndice del nodo: %d\n", nombreBase, indiceNodo)

	dominio := "raft.default.svc.cluster.local" + raftServicePort
	fmt.Printf("Dominio completo: %s\n", dominio)

	// Construir direcciones DNS de los nodos
	direcciones := make([]string, 0)
	miDireccion := ""
	for i := 0; i < 3; i++ {
		dir := nombreBase + "-" + strconv.Itoa(i) + "." + dominio
		direcciones = append(direcciones, dir)
		if i == indiceNodo {
			miDireccion = dir
		}
		fmt.Printf("Dirección nodo %d: %s\n", i, dir)
	}

	fmt.Printf("Mi dirección para escuchar: %s\n", miDireccion)

	// Convertir a HostPort
	var nodos []rpctimeout.HostPort
	// Resto de argumento son los end points como strings
	// De todas las replicas -> pasarlos a HostPort
	for _, endPoint := range direcciones {
		fmt.Printf("Convertir endpoint a HostPort: %s\n", endPoint)
		nodos = append(nodos, rpctimeout.HostPort(endPoint))
	}

	// Inicializar nodo Raft
	fmt.Println("Inicializando nodo Raft...")
	nodoRaft := raft.NuevoNodo(nodos, indiceNodo, make(chan raft.AplicaOperacion, 1000))

	fmt.Println("Registrando nodo Raft en RPC...")
	err = rpc.Register(nodoRaft)
	check.CheckError(err, "Error al registrar el nodo Raft en RPC")

	fmt.Printf("Abriendo listener TCP en %s...\n", miDireccion)
	listener, err := net.Listen("tcp", miDireccion)
	check.CheckError(err, "Error al abrir el listener TCP para el nodo")

	fmt.Println("Servidor Raft listo para aceptar conexiones RPC.")

	for {
		rpc.Accept(listener)
	}
}