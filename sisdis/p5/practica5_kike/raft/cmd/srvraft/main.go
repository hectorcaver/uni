// srvraft
// Jorge Gallardo y Enrique Baldovin
package main

import (
	//"errors"
	//"fmt"
	//"log"
	"net"
	"net/rpc"
	"os"
	"raft/internal/comun/check"
	"raft/internal/comun/rpctimeout"
	"raft/internal/raft"
	"strconv"
	"strings"
	//"time"
)

const portDNS = ":6000"

func main() {
	// obtener entero de indice de este nodo
	DNS := "raft.default.svc.cluster.local" + portDNS
	//me, err := strconv.Atoi(os.Args[1])
	//me, err := strconv.Atoi(strings.Split(os.Args[1], "-")[1])

	meStr := os.Args[1]

	// Dividir la cadena en dos partes utilizando el guion "-" como delimitador
	// La primera parte es el nombre del nodo
	name := strings.Split(meStr, "-")[0]

	// La segunda parte es el índice del nodo (ej. "1"), que se convierte a un número entero
	me, err := strconv.Atoi(strings.Split(meStr, "-")[1])
	// Si hay un error al convertir el índice a entero, se maneja el error y se imprime un mensaje
	check.CheckError(err, "Main, mal numero entero de indice de nodo:")

	// Crear una lista de direcciones de nodos
	var direcciones []string
	var MiDir string
	// Generar las direcciones para tres nodos diferentes
	for i := 0; i < 3; i++ {
		// Construir la dirección completa del nodo en el formato "name-i.dns".
		nodo := name + "-" + strconv.Itoa(i) + "." + DNS
		if i == me{
			MiDir=nodo
		}
		// Añadir la dirección generada a la lista de direcciones.
		direcciones = append(direcciones, nodo)
	}

	var nodos []rpctimeout.HostPort
	// Resto de argumento son los end points como strings
	// De todas la replicas-> pasarlos a HostPort
	for _, endPoint := range direcciones {
		nodos = append(nodos, rpctimeout.HostPort(endPoint))
	}

	// Parte Servidor
	nr := raft.NuevoNodo(nodos, me, make(chan raft.AplicaOperacion, 1000))
	rpc.Register(nr) // registrar el objeto nr para que pueda ser llamado remotamente

	//fmt.Println("Replica escucha en :", me, " de ", os.Args[2:])

	l, err := net.Listen("tcp", MiDir)
	check.CheckError(err, "Main listen error:")

	for {
		rpc.Accept(l)
	}
}
