package testintegracionraft1

import (
	"fmt"
	"raft/internal/comun/check"

	//"log"
	//"crypto/rand"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"raft/internal/comun/rpctimeout"
	"raft/internal/despliegue"
	"raft/internal/raft"
)

const (
	//nodos replicas
	REPLICA1 = "127.0.0.1:29001"//"192.168.3.14:29180"
	REPLICA2 = "127.0.0.1:29002"//"192.168.3.15:29180"
	REPLICA3 = "127.0.0.1:29003"//"192.168.3.16:29180"

	// paquete main de ejecutables relativos a PATH previo
	EXECREPLICA = "cmd/srvraft/main.go"

	// comandos completo a ejecutar en máquinas remota con ssh. Ejemplo :
	// 				cd $HOME/raft; go run cmd/srvraft/main.go 127.0.0.1:29001
)

// PATH de los ejecutables de modulo golang de servicio Raft
var PATH string = filepath.Join(os.Getenv("HOME"), "uni", "sisdis","p4", "CodigoEsqueleto", "raft")
// var PATH string = filepath.Join("misc", "alumnos", "sd", "sd2425", "a869637", "p4", "CodigoEsqueleto", "raft")

// go run cmd/srvraft/main.go 0 127.0.0.1:29001 127.0.0.1:29002 127.0.0.1:29003
var EXECREPLICACMD string = "cd " + PATH + " && /usr/local/go/bin/go run " + EXECREPLICA
//var EXECREPLICACMD string = "cd " + PATH + " ; go run " + EXECREPLICA


// TEST primer rango
func TestPrimerasPruebas(t *testing.T) { // (m *testing.M) {
	// <setup code>
	// Crear canal de resultados de ejecuciones ssh en maquinas remotas
	cfg := makeCfgDespliegue(t,
		3,
		[]string{REPLICA1, REPLICA2, REPLICA3},
		[]bool{true, true, true})

	// tear down code
	// eliminar procesos en máquinas remotas
	defer cfg.stop()

	// Run test sequence

	// Test1 : No debería haber ningun primario, si SV no ha recibido aún latidos
	t.Run("T1:soloArranqueYparada",
		func(t *testing.T) { cfg.soloArranqueYparadaTest1(t) })

	// Test2 : No debería haber ningun primario, si SV no ha recibido aún latidos
	t.Run("T2:ElegirPrimerLider",
		func(t *testing.T) { cfg.elegirPrimerLiderTest2(t) })

	// Test3: tenemos el primer primario correcto
	t.Run("T3:FalloAnteriorElegirNuevoLider",
		func(t *testing.T) { cfg.falloAnteriorElegirNuevoLiderTest3(t) })

	// Test4: Tres operaciones comprometidas en configuración estable
	t.Run("T4:tresOperacionesComprometidasEstable",
		func(t *testing.T) { cfg.tresOperacionesComprometidasEstable(t) })
}

// TEST primer rango
func TestAcuerdosConFallos(t *testing.T) { // (m *testing.M) {
	// <setup code>
	// Crear canal de resultados de ejecuciones ssh en maquinas remotas
	cfg := makeCfgDespliegue(t,
		3,
		[]string{REPLICA1, REPLICA2, REPLICA3},
		[]bool{true, true, true})

	// tear down code
	// eliminar procesos en máquinas remotas
	defer cfg.stop()

	// Test5: Se consigue acuerdo a pesar de desconexiones de seguidor
	t.Run("T5:AcuerdoAPesarDeDesconexionesDeSeguidor ",
		func(t *testing.T) { cfg.AcuerdoApesarDeSeguidor(t) })

	t.Run("T5:SinAcuerdoPorFallos ",
		func(t *testing.T) { cfg.SinAcuerdoPorFallos(t) })

	t.Run("T5:SometerConcurrentementeOperaciones ",
		func(t *testing.T) { cfg.SometerConcurrentementeOperaciones(t) })

}

// ---------------------------------------------------------------------
//
// Canal de resultados de ejecución de comandos ssh remotos
type canalResultados chan string

func (cr canalResultados) stop() {
	close(cr)

	// Leer las salidas obtenidos de los comandos ssh ejecutados
	/*for s := range cr {
		fmt.Println(s)
	}*/
}

// ---------------------------------------------------------------------
// Operativa en configuracion de despliegue y pruebas asociadas
type configDespliegue struct {
	t           *testing.T
	conectados  []bool
	numReplicas int
	nodosRaft   []rpctimeout.HostPort
	cr          canalResultados
}

// Crear una configuracion de despliegue
func makeCfgDespliegue(t *testing.T, n int, nodosraft []string,
	conectados []bool) *configDespliegue {
	cfg := &configDespliegue{}
	cfg.t = t
	cfg.conectados = conectados
	cfg.numReplicas = n
	cfg.nodosRaft = rpctimeout.StringArrayToHostPortArray(nodosraft)
	cfg.cr = make(canalResultados, 2000)

	return cfg
}

func (cfg *configDespliegue) stop() {
	//cfg.stopDistributedProcesses()

	time.Sleep(50 * time.Millisecond)

	cfg.cr.stop()
}

// --------------------------------------------------------------------------
// FUNCIONES DE SUBTESTS

// Se pone en marcha una replica ?? - 3 NODOS RAFT
func (cfg *configDespliegue) soloArranqueYparadaTest1(t *testing.T) {
	t.Skip("SKIPPED soloArranqueYparadaTest1")

	fmt.Println(t.Name(), ".....................")

	cfg.t = t // Actualizar la estructura de datos de tests para errores

	// Poner en marcha replicas en remoto con un tiempo de espera incluido
	cfg.startDistributedProcesses()

	// Comprobar estado replica 0
	cfg.comprobarEstadoRemoto(0, 0, false, -1)

	// Comprobar estado replica 1
	cfg.comprobarEstadoRemoto(1, 0, false, -1)

	// Comprobar estado replica 2
	cfg.comprobarEstadoRemoto(2, 0, false, -1)

	// Parar réplicas almacenamiento en remoto
	cfg.stopDistributedProcesses()

	fmt.Println(".............", t.Name(), "Superado")
}

// Primer lider en marcha - 3 NODOS RAFT
func (cfg *configDespliegue) elegirPrimerLiderTest2(t *testing.T) {
	t.Skip("SKIPPED ElegirPrimerLiderTest2")

	fmt.Println(t.Name(), ".....................")

	cfg.startDistributedProcesses()

	// Se ha elegido lider ?
	fmt.Printf("Probando lider en curso\n")
	IdLider := cfg.pruebaUnLider(cfg.numReplicas)
	fmt.Println("Nodo", IdLider, "es el lider")

	// Parar réplicas alamcenamiento en remoto
	cfg.stopDistributedProcesses() // Parametros

	fmt.Println(".............", t.Name(), "Superado")
}

// Fallo de un primer lider y reeleccion de uno nuevo - 3 NODOS RAFT
func (cfg *configDespliegue) falloAnteriorElegirNuevoLiderTest3(t *testing.T) {
	t.Skip("SKIPPED FalloAnteriorElegirNuevoLiderTest3")

	fmt.Println(t.Name(), ".....................")

	cfg.startDistributedProcesses()

	fmt.Printf("Lider inicial\n")
	liderId := cfg.pruebaUnLider(3)

	// Desconectar lider
	cfg.pararNodo(liderId, 3)

	// Esperar a que el lider se desconecte
	time.Sleep(2 * time.Second)

	fmt.Printf("Comprobar nuevo lider\n")
	nuevoLider := cfg.pruebaUnLider(3)

	fmt.Printf("Nuevo lider %d", nuevoLider)

	// Parar réplicas almacenamiento en remoto
	cfg.stopDistributedProcesses() //parametros

	fmt.Println(".............", t.Name(), "Superado")
}

// 3 operaciones comprometidas con situacion estable y sin fallos - 3 NODOS RAFT
func (cfg *configDespliegue) tresOperacionesComprometidasEstable(t *testing.T) {
	t.Skip("SKIPPED tresOperacionesComprometidasEstable")

	// A COMPLETAR .....

	fmt.Println(t.Name(), ".....................")

	cfg.startDistributedProcesses()

	fmt.Println("Probando Lider inicial")

	IdLider := cfg.pruebaUnLider(3)

	fmt.Println("Lider en nodo:", IdLider)

	nodoNoLider := (IdLider + 1) % len(cfg.nodosRaft)


	leer1 := raft.TipoOperacion{ 
		Operacion: "leer", 
		Clave: "0x0000",
	}
	escribir1 := raft.TipoOperacion{ 
		Operacion: "escribir", 
		Clave: "0x0000", 
		Valor: "chocolate",
	}

	// * Envío una operación a un no líder 
	// * para comprobar funcionamiento correcto
	_,_,_,lider,resultado := cfg.enviarOperacion(nodoNoLider, escribir1, true)

	fmt.Println(nodoNoLider, "Resultado de operación en nodo no lider --> ", resultado, " y lider: ", lider, ", IdLider == ", IdLider)

	times := 3

	_,_,_,_,resultado = cfg.enviarOperacion(IdLider, escribir1, true)

	fmt.Println("Resultado de operación ", 0, " --> ", resultado)

	for i := 1; i < times; i++ {

		_,_,_,_,resultado := cfg.enviarOperacion(IdLider, leer1, true)

		fmt.Println("Resultado de operación ", i, " --> ", resultado)

	}

	cfg.comprobarEstadoReplicacion(IdLider)

	// Parar réplicas almacenamiento en remoto
	cfg.stopDistributedProcesses() //parametros

	fmt.Println(".............", t.Name(), "Superado")
}


// Se consigue acuerdo a pesar de desconexiones de seguidor -- 3 NODOS RAFT
func (cfg *configDespliegue) AcuerdoApesarDeSeguidor(t *testing.T) {
	t.Skip("SKIPPED AcuerdoApesarDeSeguidor")

	// A completar ???

	fmt.Println(t.Name(), ".....................")

	cfg.t = t // Actualizar la estructura de datos de tests para errores

	cfg.startDistributedProcesses()

	fmt.Println("Probando Lider inicial")

	IdLider := cfg.pruebaUnLider(3)

	// * Sabes que hay un líder

	fmt.Println("Lider en nodo:", IdLider)

	leer1 := raft.TipoOperacion{ 
		Operacion: "leer", 
		Clave: "0x0000",
	}
	escribir1 := raft.TipoOperacion{ 
		Operacion: "escribir", 
		Clave: "0x0000", 
		Valor: "chocolate",
	}

	leer2 := raft.TipoOperacion{ 
		Operacion: "leer", 
		Clave: "0x0001",
	}
	escribir2 := raft.TipoOperacion{ 
		Operacion: "escribir", 
		Clave: "0x0001", 
		Valor: "con churros",
	}

	// * Comprometer una entrada

	_,_,_,_,resultado := cfg.enviarOperacion(IdLider, escribir1, true)

	fmt.Println("Resultado de operación:", escribir1, " --> ", resultado)

	nodoSeguidor := (IdLider + 1) % len(cfg.nodosRaft) 

	// * Obtener un lider y, a continuación desconectar una de los nodos Raft
	cfg.pararNodo(nodoSeguidor, cfg.numReplicas)

	// * Dejar tiempo para que el nodo se desconecte
	time.Sleep(2 * time.Second)

	// * Comprobar varios acuerdos con una réplica desconectada

	_,_,_,_,resultado = cfg.enviarOperacion(IdLider, escribir2, true)

	fmt.Println("Resultado de operación:", escribir2, " --> ", resultado)

	_,_,_,_,resultado = cfg.enviarOperacion(IdLider, leer1, true)

	fmt.Println("Resultado de operación:", leer1, " --> ", resultado)

	_,_,_,_,resultado = cfg.enviarOperacion(IdLider, leer2, true)

	fmt.Println("Resultado de operación:", leer2, " --> ", resultado)

	// reconectar nodo Raft previamente desconectado y comprobar varios acuerdos

	cfg.reconectarNodo(nodoSeguidor)

	// * Damos tiempo para que se realicen los acuerdos
	time.Sleep(15 * time.Second)

	cfg.comprobarEstadoReplicacion(IdLider)

	cfg.stopDistributedProcesses()

	fmt.Println(".............", t.Name(), "Superado")

}

// NO se consigue acuerdo al desconectarse mayoría de seguidores -- 3 NODOS RAFT
func (cfg *configDespliegue) SinAcuerdoPorFallos(t *testing.T) {
	t.Skip("SKIPPED SinAcuerdoPorFallos")

	// A completar ???
	
	
	fmt.Println(t.Name(), ".....................")

	cfg.startDistributedProcesses()

	fmt.Println("Probando Lider inicial")

	IdLider := cfg.pruebaUnLider(3)

	// * Sabes que hay un líder

	fmt.Println("Lider en nodo:", IdLider)

	leer1 := raft.TipoOperacion{ 
		Operacion: "leer", 
		Clave: "0x0000",
	}
	escribir1 := raft.TipoOperacion{ 
		Operacion: "escribir", 
		Clave: "0x0000", 
		Valor: "chocolate",
	}

	leer2 := raft.TipoOperacion{ 
		Operacion: "leer", 
		Clave: "0x0001",
	}
	escribir2 := raft.TipoOperacion{ 
		Operacion: "escribir", 
		Clave: "0x0001", 
		Valor: "con churros",
	}

	// * Comprometer una entrada

	_,_,_,_,resultado := cfg.enviarOperacion(IdLider, escribir1, true)

	fmt.Println("Resultado de operación:", escribir1, " --> ", resultado)

	// * Obtener un lider y, a continuación desconectar 2 de los nodos Raft
	for i := 0; i < cfg.numReplicas; i++ {
		if i != IdLider {
			cfg.pararNodo(i, cfg.numReplicas)
		}
	}

	time.Sleep(500 * time.Millisecond)

	// Comprobar varios acuerdos con 2 réplicas desconectada

	_,_,_,_,resultado = cfg.enviarOperacion(IdLider, escribir2, false)

	fmt.Println("Resultado de operación:", escribir2, " --> ", resultado)

	_,_,_,_,resultado = cfg.enviarOperacion(IdLider, leer1, false)

	fmt.Println("Resultado de operación:", leer1, " --> ", resultado)

	// reconectar lo2 nodos Raft  desconectados y probar varios acuerdos

	for i := 0; i < cfg.numReplicas; i++ {
		if i != IdLider && !cfg.conectados[i] {
			cfg.reconectarNodo(i)
		}
	}

	// * Damos tiempo a que se establezcan
	time.Sleep(1 * time.Second)


	_,_,_,_,resultado = cfg.enviarOperacion(IdLider, escribir2, false)

	fmt.Println("Resultado de operación:", escribir2, " --> ", resultado)

	_,_,_,_,resultado = cfg.enviarOperacion(IdLider, leer1, false)

	fmt.Println("Resultado de operación:", leer1, " --> ", resultado)

	_,_,_,_,resultado = cfg.enviarOperacion(IdLider, leer2, true)

	fmt.Println("Resultado de operación:", leer2, " --> ", resultado)

	// * Dar tiempo a que se comprometan y apliquen las operaciones en todos
	// * los nodos
	time.Sleep(1 * time.Second)

	cfg.comprobarEstadoReplicacion(IdLider)

	cfg.stopDistributedProcesses()

	fmt.Println(".............", t.Name(), "Superado")
}

// Se somete 5 operaciones de forma concurrente -- 3 NODOS RAFT
func (cfg *configDespliegue) SometerConcurrentementeOperaciones(t *testing.T) {
	//t.Skip("SKIPPED SometerConcurrentementeOperaciones")

	// A completar ???
	
	fmt.Println(t.Name(), ".....................")

	cfg.startDistributedProcesses()

	fmt.Println("Probando Lider inicial")

	IdLider := cfg.pruebaUnLider(3)

	// * Sabes que hay un líder

	fmt.Println("Lider en nodo:", IdLider)

	leer1 := raft.TipoOperacion{ 
		Operacion: "leer", 
		Clave: "0x0000",
	}
	escribir1 := raft.TipoOperacion{ 
		Operacion: "escribir", 
		Clave: "0x0000", 
		Valor: "chocolate",
	}

	leer2 := raft.TipoOperacion{ 
		Operacion: "leer", 
		Clave: "0x0001",
	}
	escribir2 := raft.TipoOperacion{ 
		Operacion: "escribir", 
		Clave: "0x0001", 
		Valor: "con churros",
	}

	escribir3 := raft.TipoOperacion{
		Operacion: "escribir", 
		Clave: "0x0002", 
		Valor: "pechuga",
	}

	escribir4 := raft.TipoOperacion{
		Operacion: "escribir", 
		Clave: "0x0002", 
		Valor: "de pollo",
	}

	// * Comprometer una entrada

	_,_,_,_,resultado := cfg.enviarOperacion(IdLider, escribir1, true)

	fmt.Println("Resultado de operación:", escribir1, " --> ", resultado)


	operaciones := []raft.TipoOperacion{leer1, leer2, escribir2, escribir3, escribir4}

	// Someter 5  operaciones concurrentes

	for _, op := range(operaciones) {
		go func(op raft.TipoOperacion, IdLider int) {
			_,_,_,_,resultado := cfg.enviarOperacion(IdLider, op, true)
			fmt.Println("Resultado de operación:", op, " --> ", resultado)
		}(op, IdLider)	
	}

	// Comprobar estados de nodos Raft, sobre todo
	// el avance del mandato en curso e indice de registro de cada uno
	// que debe ser identico entre ellos
	time.Sleep(10 * time.Second)

	cfg.comprobarEstadoReplicacion(IdLider)

	cfg.stopDistributedProcesses()

	fmt.Println(".............", t.Name(), "Superado")

}

// --------------------------------------------------------------------------
// FUNCIONES DE APOYO
// Comprobar que hay un solo lider
// probar varias veces si se necesitan reelecciones
func (cfg *configDespliegue) pruebaUnLider(numreplicas int) int {
	for iters := 0; iters < 10; iters++ {
		time.Sleep(500 * time.Millisecond)
		mapaLideres := make(map[int][]int)
		for i := 0; i < numreplicas; i++ {
			if cfg.conectados[i] {
				if _, mandato, eslider, _ := cfg.obtenerEstadoRemoto(i); eslider {
					mapaLideres[mandato] = append(mapaLideres[mandato], i)
				}
			}
		}

		ultimoMandatoConLider := -1
		for mandato, lideres := range mapaLideres {
			if len(lideres) > 1 {
				cfg.t.Fatalf("mandato %d tiene %d (>1) lideres",
					mandato, len(lideres))
			}
			if mandato > ultimoMandatoConLider {
				ultimoMandatoConLider = mandato
			}
		}

		if len(mapaLideres) != 0 {

			return mapaLideres[ultimoMandatoConLider][0] // Termina

		}
	}
	cfg.t.Fatalf("un lider esperado, ninguno obtenido")

	return -1 // Termina
}

func (cfg *configDespliegue) pararNodo(nodo int, numReplicas int) {

	if nodo >= 0 && nodo < numReplicas {
		var args, reply raft.Vacio
		err := cfg.nodosRaft[nodo].CallTimeout("NodoRaft.ParaNodo",
			&args, &reply, 50*time.Millisecond)

		if err != nil {
			check.CheckError(err, "Error en llamada RPC Para nodo")
		} else {
			cfg.conectados[nodo] = false
		}
	} else {
		cfg.t.Fatalf("Nodo %d no es un nodo valido", nodo)
	}
}

func (cfg *configDespliegue) obtenerEstadoRemoto(
	indiceNodo int) (int, int, bool, int) {
	var reply raft.EstadoRemoto
	err := cfg.nodosRaft[indiceNodo].CallTimeout("NodoRaft.ObtenerEstadoNodo",
		raft.Vacio{}, &reply, 300*time.Millisecond)
	check.CheckError(err, "Error en llamada RPC ObtenerEstadoRemoto")

	return reply.IdNodo, reply.Mandato, reply.EsLider, reply.IdLider
}

func (cfg *configDespliegue) obtenerEstadoReplicacionRemoto(indiceNodo int) (
	[]raft.EntradaRegistro, map[string]string) {
		var reply raft.EstadoReplicacionRemoto
		err := cfg.nodosRaft[indiceNodo].CallTimeout(
			"NodoRaft.ObtenerEstadoReplicacionNodo", raft.Vacio{}, &reply,
			 5 * time.Second)
		check.CheckError(err, 
			"Error en llamada RPC ObtenerEstadoReplicacionRemoto")

		return reply.Log, reply.Almacen
	}

// start  gestor de vistas; mapa de replicas y maquinas donde ubicarlos;
// y lista clientes (host:puerto)
func (cfg *configDespliegue) startDistributedProcesses() {
	//cfg.t.Log("Before starting following distributed processes: ", cfg.nodosRaft)

	for i, endPoint := range cfg.nodosRaft {
		despliegue.ExecMutipleHosts(EXECREPLICACMD+
			" "+strconv.Itoa(i)+" "+
			rpctimeout.HostPortArrayToString(cfg.nodosRaft),
			[]string{endPoint.Host()}, cfg.cr)
		cfg.conectados[i] = true
		// dar tiempo para se establezcan las replicas
		//time.Sleep(500 * time.Millisecond)
	}

	// aproximadamente 500 ms para cada arranque por ssh en portatil
	time.Sleep(2000 * time.Millisecond)
}

func (cfg *configDespliegue) stopDistributedProcesses() {
	var reply raft.Vacio

	for i, endPoint := range cfg.nodosRaft {
		if cfg.conectados[i] {
			err := endPoint.CallTimeout("NodoRaft.ParaNodo",
			raft.Vacio{}, &reply, 20*time.Millisecond)
			check.CheckError(err, "Error en llamada RPC Para nodo")
			cfg.conectados[i] = false
		}
	}

	time.Sleep(500 * time.Millisecond)
}

// Comprobar estado remoto de un nodo con respecto a un estado prefijado
func (cfg *configDespliegue) comprobarEstadoRemoto(idNodoDeseado int,
	mandatoDeseado int, esLiderDeseado bool, IdLiderDeseado int) {
	idNodo, mandato, esLider, idLider := cfg.obtenerEstadoRemoto(idNodoDeseado)

	//cfg.t.Log("Estado replica 0: ", idNodo, mandato, esLider, idLider, "\n")

	if idNodo != idNodoDeseado || mandato != mandatoDeseado ||
		esLider != esLiderDeseado || idLider != IdLiderDeseado {
		cfg.t.Fatalf("Estado incorrecto en replica %d en subtest %s",
			idNodoDeseado, cfg.t.Name())
	}

}


func (cfg *configDespliegue) comprobarEstadoReplicacionRemoto(idNodoDeseado int,
	logDeseado []raft.EntradaRegistro, almacenDeseado map[string]string) {

	logNodo, almacenNodo := cfg.obtenerEstadoReplicacionRemoto(idNodoDeseado)

	// Comparar logs
	if len(logDeseado) != len(logNodo) {
		cfg.t.Fatalf("Diferencia en longitud del log para nodo %d: esperado %d, obtenido %d\n",
			idNodoDeseado, len(logDeseado), len(logNodo))
	} else {
		for i := range logDeseado {
			if logDeseado[i] != logNodo[i] {
				cfg.t.Fatalf("Diferencia en log en índice %d para nodo %d: esperado %+v, obtenido %+v\n",
					i, idNodoDeseado, logDeseado[i], logNodo[i])
			}
		}
	}

	// Comparar almacenes
	for clave, valorEsperado := range almacenDeseado {
		if valorObtenido, ok := almacenNodo[clave]; !ok {
			cfg.t.Fatalf("Clave faltante en nodo %d: %s\n", idNodoDeseado, clave)
		} else if valorObtenido != valorEsperado {
			cfg.t.Fatalf("Valor distinto para clave '%s' en nodo %d: esperado '%s', obtenido '%s'\n",
				clave, idNodoDeseado, valorEsperado, valorObtenido)
		}
	}

	// Claves extras en almacenNodo
	for clave := range almacenNodo {
		if _, ok := almacenDeseado[clave]; !ok {
			cfg.t.Fatalf("Clave extra en nodo %d: %s\n", idNodoDeseado, clave)
		}
	}
}


func (cfg *configDespliegue) enviarOperacion(
		indiceNodo int, args raft.TipoOperacion, checkError bool) (int,
			 int, bool, int, string) {

	var reply raft.ResultadoRemoto

	err := cfg.nodosRaft[indiceNodo].CallTimeout("NodoRaft.SometerOperacionRaft",
		args, &reply, 5 * time.Second)
	
	if checkError {
		check.CheckError(err, "Error en llamada RPC SometerOperacionRaft")
	} else {
		if err != nil {
			fmt.Println("Error en llamada RPC SometerOperacionRaft:", err.Error())
		}
	}
	
	return reply.IndiceRegistro, reply.Mandato, reply.EsLider, reply.IdLider,
		reply.ValorADevolver

}


func (cfg *configDespliegue) reconectarNodo(nodo int) {

	if nodo >= 0 && nodo < cfg.numReplicas && !cfg.conectados[nodo] {
		endPoint := cfg.nodosRaft[nodo]

		despliegue.ExecMutipleHosts(EXECREPLICACMD+
				" "+strconv.Itoa(nodo)+" "+
				rpctimeout.HostPortArrayToString(cfg.nodosRaft),
				[]string{endPoint.Host()}, cfg.cr)

		cfg.conectados[nodo] = true
	} else {
		cfg.t.Fatalf("Nodo %d no es un nodo valido", nodo)
	}
}

func (cfg *configDespliegue) comprobarEstadoReplicacion(IdLider int) {
	fmt.Println(
		"Comprobando estado de los logs y de los almacenes de todos los nodos")

	logLider, almacenLider := cfg.obtenerEstadoReplicacionRemoto(IdLider)

	for i := 0; i < cfg.numReplicas; i++ {
		if i != IdLider {
			cfg.comprobarEstadoReplicacionRemoto(i, logLider, almacenLider)
		}
	}

	fmt.Println("Todo correcto")
}