package main

// Jorge Gallardo y Enrique Baldovin
import (
	"bufio"
	"fmt"

	//"net"
	//"net/rpc"
	"os"
	"raft/internal/comun/check"
	"raft/internal/comun/rpctimeout"
	"raft/internal/raft"
	"strconv"
	"strings"
	"time"
)

// definicion de los nodos-------------------
const port = ":6000"

var NODOS = rpctimeout.StringArrayToHostPortArray([]string{
	"nodo-0.raft.default.svc.cluster.local" + port,
	"nodo-1.raft.default.svc.cluster.local" + port,
	"nodo-2.raft.default.svc.cluster.local" + port,
})

// -----------------------------------------
// Helpers para impresión
func printEstadoRemoto(reply raft.EstadoRemoto) {
	fmt.Println("=== Estado del Nodo ===")
	fmt.Printf("IdNodo:  %d\n", reply.IdNodo)
	fmt.Printf("Mandato: %d\n", reply.Mandato)
	fmt.Printf("EsLider: %v\n", reply.EsLider)
	fmt.Printf("IdLider: %d\n", reply.IdLider)
	fmt.Println("=======================")
}

func printEstadoSometido(reply raft.ResPruebas) {
	fmt.Println("=== Estado Sometido ===")
	fmt.Printf("Índice:  %d\n", reply.Commit)
	fmt.Printf("Mandato: %d\n", reply.Mandatocommit)
	fmt.Printf("Longuitud del log: %d\n", reply.NumOperaciones)
	fmt.Println("=========================")
}

func printSometerOperacion(reply raft.ResultadoRemoto) {
	fmt.Println("=== Resultado de la Operación ===")
	fmt.Printf("IndiceRegistro: %d\n", reply.IndiceRegistro)
	fmt.Printf("Mandato: %d\n", reply.Mandato)
	fmt.Printf("EsLider: %v\n", reply.EsLider)
	fmt.Printf("IdLider: %v\n", reply.IdLider)
	fmt.Printf("ValorADevolver: %s\n", reply.ValorADevolver)
	fmt.Println("=================================")
}

func printUsage() {
	fmt.Println("Comandos disponibles:")
	fmt.Println("  obtenerEstado <idNodo>")
	fmt.Println("  obtenerEstadoSometido <idNodo>")
	fmt.Println("  parar <idNodo>")
	fmt.Println("  leer <idNodoLider> <clave>")
	fmt.Println("  escribir <idNodoLider> <clave> <valor>")
	fmt.Println("  salir")
	fmt.Println()
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Ingrese comandos (escriba 'salir' para terminar).")

	for {
		printUsage()
		fmt.Print(">")
		if !scanner.Scan() {
			fmt.Println("Error al leer la entrada. Intente nuevamente.")
			continue
		}
		linea := scanner.Text()
		if strings.TrimSpace(linea) == "" {
			continue
		}

		partes := strings.Fields(linea)
		comando := partes[0]

		if comando == "salir" {
			fmt.Println("Saliendo...")
			return
		}

		switch comando {
		case "obtenerEstado":
			if len(partes) < 2 {
				fmt.Println("Faltan argumentos. Uso: obtenerEstado <idNodo>")
				continue
			}
			nodo, err := strconv.Atoi(partes[1])
			if err != nil || nodo < 0 || nodo >= len(NODOS) {
				fmt.Println("idNodo inválido.")
				continue
			}
			var reply raft.EstadoRemoto
			err = NODOS[nodo].CallTimeout("NodoRaft.ObtenerEstadoNodo", raft.Vacio{}, &reply, 1000*time.Millisecond)
			check.CheckError(err, "Error en ObtenerEstadoNodo")
			printEstadoRemoto(reply)

		case "obtenerEstadoSometido":
			if len(partes) < 2 {
				fmt.Println("Faltan argumentos. Uso: obtenerEstadoSometido <idNodo>")
				continue
			}
			nodo, err := strconv.Atoi(partes[1])
			if err != nil || nodo < 0 || nodo >= len(NODOS) {
				fmt.Println("idNodo inválido.")
				continue
			}
			var reply raft.ResPruebas
			err = NODOS[nodo].CallTimeout("NodoRaft.EstadoPruebas", raft.Vacio{}, &reply, 1000*time.Millisecond)
			check.CheckError(err, "Error en EstadoPruebas")
			printEstadoSometido(reply)

		case "parar":
			if len(partes) < 2 {
				fmt.Println("Faltan argumentos. Uso: parar <idNodo>")
				continue
			}
			nodo, err := strconv.Atoi(partes[1])
			if err != nil || nodo < 0 || nodo >= len(NODOS) {
				fmt.Println("idNodo inválido.")
				continue
			}
			var reply raft.Vacio
			err = NODOS[nodo].CallTimeout("NodoRaft.ParaNodo", raft.Vacio{}, &reply, 1000*time.Millisecond)
			check.CheckError(err, "Error en ParaNodo")
			fmt.Printf("Nodo %d detenido con éxito.\n", nodo)

		case "leer":
			if len(partes) < 3 {
				fmt.Println("Faltan argumentos. Uso: leer <idNodoLider> <clave>")
				continue
			}
			nodo, err := strconv.Atoi(partes[1])
			if err != nil || nodo < 0 || nodo >= len(NODOS) {
				fmt.Println("idNodo inválido.")
				continue
			}
			clave := partes[2]

			var operacion raft.TipoOperacion
			operacion.Operacion = "leer"
			operacion.Clave = clave

			var reply raft.ResultadoRemoto
			err = NODOS[nodo].CallTimeout("NodoRaft.SometerOperacionRaft", operacion, &reply, 5000*time.Millisecond)
			check.CheckError(err, "Error en SometerOperacionRaft (leer)")
			printSometerOperacion(reply)

		case "escribir":
			if len(partes) < 4 {
				fmt.Println("Faltan argumentos. Uso: escribir <idNodoLider> <clave> <valor>")
				continue
			}
			nodo, err := strconv.Atoi(partes[1])
			if err != nil || nodo < 0 || nodo >= len(NODOS) {
				fmt.Println("idNodo inválido.")
				continue
			}
			clave := partes[2]
			valor := partes[3]

			var operacion raft.TipoOperacion
			operacion.Operacion = "escribir"
			operacion.Clave = clave
			operacion.Valor = valor

			var reply raft.ResultadoRemoto
			err = NODOS[nodo].CallTimeout("NodoRaft.SometerOperacionRaft", operacion, &reply, 10000*time.Millisecond)
			check.CheckError(err, "Error en SometerOperacionRaft (escribir)")
			printSometerOperacion(reply)
		}
	}
}
