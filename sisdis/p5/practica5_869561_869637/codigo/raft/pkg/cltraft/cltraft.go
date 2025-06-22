/****************************************
 * Autor: Adrián Nasarre Sánchez 869561
 * Autor: Héctor Lacueva Sacristán 869637
 * Implementación del cliente raft
 *****************************************/

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"raft/internal/comun/check"
	"raft/internal/comun/rpctimeout"
	"raft/internal/raft"
)

// Puerto por defecto para los nodos Raft
const raftPort = ":6000"

// Lista de direcciones de los nodos del clúster
var clusterNodes = rpctimeout.StringArrayToHostPortArray([]string{
	"nodo-0.raft.default.svc.cluster.local" + raftPort,
	"nodo-1.raft.default.svc.cluster.local" + raftPort,
	"nodo-2.raft.default.svc.cluster.local" + raftPort,
})

// Muestra el menú de comandos disponibles al usuario
func mostrarMenu() {
	fmt.Println("\nComandos disponibles:")
	fmt.Println("  estado <nodo>      			# Estado de un nodo específico")
	fmt.Println("  commit <nodo>      			# Estado de commit de un nodo")
	fmt.Println("  detener <nodo>     			# Detener un nodo")
	fmt.Println("  get <lider> <clave>        	# Leer valor de una clave")
	fmt.Println("  put <lider> <clave> <valor> 	# Escribir valor en una clave")
	fmt.Println("  exit               			# Salir del cliente")
}

// Imprime el estado general de un nodo
func imprimirEstadoNodo(info raft.EstadoRemoto) {
	fmt.Println("\n[Información del Nodo]")
	fmt.Printf("Nodo: %d | Mandato: %d | EsLíder: %v | Líder: %d\n", info.IdNodo, info.Mandato, info.EsLider, info.IdLider)
}

// Imprime el estado de commit de un nodo
func imprimirEstadoReplicasYAlamacen(reply raft.EstadoReplicacionRemoto) {
	printLog(reply.Log)
	printAlmacen(reply.Almacen)
}

func printLog(log []raft.EntradaRegistro) {
	fmt.Printf("Log (len=%d):", len(log))
	for i, entrada := range log {
		fmt.Printf("  [%d] %v", i, entrada)
	}
}

func printAlmacen(almacen map[string]string) {
	// Almacen (claves limitadas para evitar saturación visual)
	fmt.Println("------ Almacén (claves limitadas) ------")
	maxKeys := 10
	count := 0
	for k, v := range almacen {
		fmt.Printf("  %s: %s", k, v)
		count++
		if count >= maxKeys {
			fmt.Printf("  ...y %d más", len(almacen)-maxKeys)
			break
		}
	}

	fmt.Println("==========================================")
}

// Imprime el resultado de una operación (lectura o escritura)
func imprimirResultadoOp(res raft.ResultadoRemoto) {
	fmt.Println("\n[Resultado de Operación]")
	fmt.Printf("Índice: %d | Mandato: %d | EsLíder: %v | Líder: %d | Valor: %s\n", res.IndiceRegistro, res.Mandato, res.EsLider, res.IdLider, res.ValorADevolver)
}

func comandoEstado(args []string) {
	// Consultar estado de un nodo
	if len(args) != 2 {
		fmt.Println("Uso: estado <nodo>")
		return
	}
	idx, err := strconv.Atoi(args[1])
	if err != nil || idx < 0 || idx >= len(clusterNodes) {
		fmt.Println("Nodo inválido.")
		return
	}
	var res raft.EstadoRemoto
	err = clusterNodes[idx].CallTimeout("NodoRaft.ObtenerEstadoNodo", raft.Vacio{}, &res, 1200*time.Millisecond)
	check.CheckError(err, "Fallo al obtener estado del nodo")
	imprimirEstadoNodo(res)
}

func comandoCommit(args []string) {
	// Consultar commit de un nodo
	if len(args) != 2 {
		fmt.Println("Uso: commit <nodo>")
		return
	}
	idx, err := strconv.Atoi(args[1])
	if err != nil || idx < 0 || idx >= len(clusterNodes) {
		fmt.Println("Nodo inválido.")
		return
	}
	var res raft.EstadoReplicacionRemoto
	err = clusterNodes[idx].CallTimeout("NodoRaft.ObtenerEstadoReplicacionNodo", raft.Vacio{}, &res, 1200*time.Millisecond)
	check.CheckError(err, "Fallo en la operacion RPC ObtenerEstadoReplicacionNodo")
	imprimirEstadoReplicasYAlamacen(res)
}

// Función principal: ciclo interactivo de comandos para el cliente Raft
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Cliente interactivo Raft. Escriba 'exit' para salir.")

	for {
		mostrarMenu() // Mostrar menú de ayuda
		fmt.Print("$ ")
		if !scanner.Scan() {
			fmt.Println("Error de entrada. Intente de nuevo.")
			continue
		}
		entrada := strings.TrimSpace(scanner.Text())
		if entrada == "" {
			continue
		}
		args := strings.Fields(entrada)
		cmd := args[0]

		switch cmd {
		case "exit":
			fmt.Println("Finalizando cliente.")
			return
		case "estado":

			comandoEstado(args)

		case "commit":

			comandoCommit(args)

			
		case "detener":
			// Detener un nodo
			if len(args) != 2 {
				fmt.Println("Uso: detener <nodo>")
				continue
			}
			idx, err := strconv.Atoi(args[1])
			if err != nil || idx < 0 || idx >= len(clusterNodes) {
				fmt.Println("Nodo inválido.")
				continue
			}
			var res raft.Vacio
			err = clusterNodes[idx].CallTimeout("NodoRaft.ParaNodo", raft.Vacio{}, &res, 1200*time.Millisecond)
			check.CheckError(err, "Fallo al detener nodo")
			fmt.Printf("Nodo %d detenido correctamente.\n", idx)
		case "get":
			// Leer valor de una clave
			if len(args) != 3 {
				fmt.Println("Uso: get <lider> <clave>")
				continue
			}
			idx, err := strconv.Atoi(args[1])
			if err != nil || idx < 0 || idx >= len(clusterNodes) {
				fmt.Println("Nodo líder inválido.")
				continue
			}
			clave := args[2]
			var op raft.TipoOperacion
			op.Operacion = "leer"
			op.Clave = clave
			var res raft.ResultadoRemoto
			err = clusterNodes[idx].CallTimeout("NodoRaft.SometerOperacionRaft", op, &res, 5000*time.Millisecond)
			check.CheckError(err, "Fallo en operación de lectura")
			imprimirResultadoOp(res)
		case "put":
			// Escribir valor en una clave
			if len(args) != 4 {
				fmt.Println("Uso: put <lider> <clave> <valor>")
				continue
			}
			idx, err := strconv.Atoi(args[1])
			if err != nil || idx < 0 || idx >= len(clusterNodes) {
				fmt.Println("Nodo líder inválido.")
				continue
			}
			clave := args[2]
			valor := args[3]
			var op raft.TipoOperacion
			op.Operacion = "escribir"
			op.Clave = clave
			op.Valor = valor
			var res raft.ResultadoRemoto
			err = clusterNodes[idx].CallTimeout("NodoRaft.SometerOperacionRaft", op, &res, 10000*time.Millisecond)
			check.CheckError(err, "Fallo en operación de escritura")
			imprimirResultadoOp(res)
		default:
			fmt.Println("Comando no reconocido. Intente de nuevo.")
		}
	}
}