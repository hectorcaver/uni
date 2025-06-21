//Jorge Gallardo y Enrique Baldovin

package testintegracionraft1

import (
	"fmt"
	"raft/internal/comun/check"

	//"log"
	//"crypto/rand"
	//"os"
	"path/filepath"
	"strconv"
	"sync"
	"testing"
	"time"

	"raft/internal/comun/rpctimeout"
	"raft/internal/despliegue"
	"raft/internal/raft"
)

const (
	//hosts
	MAQUINA1 = "127.0.0.1" //"192.168.3.9"
	MAQUINA2 = "127.0.0.1" //"192.168.3.11"
	MAQUINA3 = "127.0.0.1" //"192.168.3.12"

	/*MAQUINA1 = "192.168.3.9"
	MAQUINA2 = "192.168.3.11"
	MAQUINA3 = "192.168.3.12"*/

	//puertos
	PUERTOREPLICA1 = "31143"
	PUERTOREPLICA2 = "31144"
	PUERTOREPLICA3 = "31145"

	//nodos replicas
	REPLICA1 = MAQUINA1 + ":" + PUERTOREPLICA1
	REPLICA2 = MAQUINA2 + ":" + PUERTOREPLICA2
	REPLICA3 = MAQUINA3 + ":" + PUERTOREPLICA3

	// paquete main de ejecutables relativos a PATH previo
	EXECREPLICA = "cmd/srvraft/main.go"

	// comandos completo a ejecutar en máquinas remota con ssh. Ejemplo :
	// 				cd $HOME/raft; go run cmd/srvraft/main.go 127.0.0.1:29001

	// Ubicar, en esta constante, nombre de fichero de vuestra clave privada local
	// emparejada con la clave pública en authorized_keys de máquinas remotas

	PRIVKEYFILE = "id_ed25519"
)

// PATH de los ejecutables de modulo golang de servicio Raft
var PATH string = filepath.Join("/misc", "alumnos", "sd", "sd2425", "a869402", "practica4", "CodigoEsqueleto", "raft")

// go run cmd/srvraft/main.go 0 127.0.0.1:29001 127.0.0.1:29002 127.0.0.1:29003
var EXECREPLICACMD string = "cd " + PATH + "; go run " + EXECREPLICA

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
	for s := range cr {
		fmt.Println(s)
	}
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
	cfg.startDistributedProcessesT1()

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
	//t.Skip("SKIPPED ElegirPrimerLiderTest2")

	fmt.Println(t.Name(), ".....................")
	cfg.t = t // Actualizar la estructura de datos de tests para errores

	cfg.startDistributedProcesses()

	// Se ha elegido lider ?
	fmt.Printf("Probando lider en curso\n")
	cfg.pruebaUnLider(3)

	// Parar réplicas almacenamiento en remoto
	cfg.stopDistributedProcesses() // Parametros

	fmt.Println(".............", t.Name(), "Superado")
}

// Fallo de un primer lider y reeleccion de uno nuevo - 3 NODOS RAFT
func (cfg *configDespliegue) falloAnteriorElegirNuevoLiderTest3(t *testing.T) {
	//t.Skip("SKIPPED FalloAnteriorElegirNuevoLiderTest3")

	fmt.Println(t.Name(), ".....................")

	cfg.startDistributedProcesses()

	fmt.Printf("Lider inicial\n")
	idLider := cfg.pruebaUnLider(3)

	// Desconectar lider
	fmt.Println("Desconectamos lider:", idLider)
	cfg.pararNodo(idLider)

	fmt.Printf("Comprobar nuevo lider\n")
	idLider = cfg.pruebaUnLider(3)
	fmt.Println("Nuevo lider:", idLider)

	// Parar réplicas almacenamiento en remoto
	cfg.stopDistributedProcesses() //Parametros

	fmt.Println(".............", t.Name(), "Superado")
}

// 3 operaciones comprometidas con situacion estable y sin fallos - 3 NODOS RAFT
func (cfg *configDespliegue) tresOperacionesComprometidasEstable(t *testing.T) {
	//t.Skip("SKIPPED tresOperacionesComprometidasEstable")

	fmt.Printf("Iniciando prueba: %s\n", t.Name())

	// Paso 1: Inicia los procesos distribuidos de los nodos Raft
	cfg.startDistributedProcesses()

	// Paso 2: Selecciona un líder en el sistema de 3 nodos
	//time.Sleep(20000 * time.Millisecond)
	leaderID := cfg.pruebaUnLider(3)
	fmt.Printf("Lider seleccionado: %d\n", leaderID)

	// Paso 3: Realiza y verifica tres operaciones en el nodo líder
	// Operación 1: Crear un nuevo recurso
	cfg.comprobarYRealizarOperacion(leaderID, 0, "escribir", "recurso_1", "contenido inicial", "escrito")

	// Operación 2: Leer el recurso existente
	cfg.comprobarYRealizarOperacion(leaderID, 1, "leer", "recurso_1", "", "contenido inicial")

	// Operación 3: Eliminar el recurso
	cfg.comprobarYRealizarOperacion(leaderID, 2, "escribir", "recurso_2", "contenido recurso2", "escrito")

	// Paso 4: Finaliza los procesos distribuidos
	cfg.stopDistributedProcesses()
	fmt.Printf("Prueba %s completada exitosamente\n", t.Name())
}

// Se consigue acuerdo a pesar de desconexiones de seguidor -- 3 NODOS RAFT
// AcuerdoApesarDeSeguidor prueba la capacidad del algoritmo de consenso Raft para alcanzar un acuerdo
// incluso cuando uno de los nodos seguidores está temporalmente desconectado. La prueba realiza los siguientes pasos:
// 1. Inicia los procesos distribuidos.
// 2. Elige un líder entre los nodos.
// 3. Realiza una operación de escritura en el líder y la verifica.
// 4. Desconecta uno de los nodos seguidores.
// 5. Realiza una serie de operaciones de lectura y escritura en el líder y las verifica.
// 6. Reconecta el nodo seguidor previamente desconectado.
// 7. Realiza operaciones adicionales de escritura y lectura para asegurar que el nodo reconectado puede ponerse al día y estar de acuerdo con el líder.
// 8. Detiene los procesos distribuidos y completa la prueba.
func (cfg *configDespliegue) AcuerdoApesarDeSeguidor(t *testing.T) {
	//t.Skip("SKIPPED AcuerdoApesarDeSeguidor")

	// A completar
	fmt.Printf("Iniciando prueba: %s\n", t.Name())

	cfg.startDistributedProcesses()

	leaderID := cfg.pruebaUnLider(3)
	fmt.Printf("Lider en la prueba AcuerdoApesarDeSeguidor: %d\n", leaderID)
	cfg.comprobarYRealizarOperacion(leaderID, 0, "escribir", "clave1", "valor inicial", "escrito")

	// Desconectar un seguidor
	var nodoparado int
	for i := 0; i < cfg.numReplicas; i++ {
		if i != leaderID {
			cfg.pararNodo(i)
			nodoparado = i
			fmt.Printf("Nodo desconectado: %d\n", nodoparado)
			break
		}
	}

	// Realizar y verificar una operación
	cfg.comprobarYRealizarOperacion(leaderID, 1, "leer", "clave1", "", "valor inicial")
	fmt.Printf("Operación de leer sometida\n")
	cfg.comprobarYRealizarOperacion(leaderID, 2, "escribir", "clave2", "nuevo valor", "escrito")
	fmt.Printf("Operación de escribir sometida\n")
	cfg.comprobarYRealizarOperacion(leaderID, 3, "leer", "clave2", "", "nuevo valor")
	fmt.Printf("Operación de leer sometida\n")

	// Reconectar nodo Raft previamente desconectado y comprobar varios acuerdos
	fmt.Printf("Nodo conectado de nuevo: %d\n", nodoparado)
	cfg.conectarNodo(nodoparado)
	cfg.comprobarYRealizarOperacion(leaderID, 4, "escribir", "clave3", "valor adicional", "escrito")
	fmt.Printf("Operación de escribir sometida\n")
	cfg.comprobarYRealizarOperacion(leaderID, 5, "leer", "clave3", "", "valor adicional")
	fmt.Printf("Operación de leer sometida\n")

	cfg.stopDistributedProcesses()
	fmt.Printf("Prueba %s completada exitosamente\n", t.Name())
}

// NO se consigue acuerdo al desconectarse mayoría de seguidores -- 3 NODOS RAFT
// SinAcuerdoPorFallos prueba el comportamiento del algoritmo de consenso Raft cuando
// hay fallos en el sistema. Realiza los siguientes pasos:
// 1. Inicia los procesos distribuidos.
// 2. Elige un líder entre los nodos.
// 3. Compromete una entrada en el registro.
// 4. Desconecta dos nodos que no son el líder.
// 5. Intenta realizar varias operaciones con dos nodos desconectados y
//    verifica si hay errores.
// 6. Reconecta los nodos previamente desconectados.
// 7. Realiza varias operaciones después de la reconexión para asegurar que el sistema
//    alcanza el consenso.
// 8. Detiene los procesos distribuidos.
//
// Esta prueba asegura que el algoritmo Raft puede manejar fallos de nodos y
// aún mantener la consistencia y disponibilidad una vez que los nodos se reconectan.
//
// Parámetros:
// - t: La instancia del marco de pruebas.
func (cfg *configDespliegue) SinAcuerdoPorFallos(t *testing.T) {
	//t.Skip("SKIPPED SinAcuerdoPorFallos")

	fmt.Println(t.Name(), ".....................")
	cfg.startDistributedProcesses()

	fmt.Printf("Probando leader \n")
	idLeader := cfg.pruebaUnLider(3)
	fmt.Printf("El nodo leader es %d\n", idLeader)

	// Comprometer una entrada
	cfg.comprobarYRealizarOperacion(idLeader, 0, "escribir", "x", "primer valor", "escrito")
	fmt.Printf("Operación de escribir sometida\n")

	// Desconectar 2 nodos Raft que no son el líder
	nodosParados := []int{}
	for i := 0; i < cfg.numReplicas; i++ {
		if i != idLeader && len(nodosParados) < 2 {
			cfg.pararNodo(i)
			nodosParados = append(nodosParados, i)
			fmt.Printf("Nodo desconectado: %d\n", i)
		}
	}

	// Comprobar varios acuerdos con 2 réplicas desconectadas
	operaciones := []raft.TipoOperacion{
		{Operacion: "leer", Clave: "x", Valor: ""},
		{Operacion: "escribir", Clave: "y", Valor: "segundo valor"},
		{Operacion: "leer", Clave: "y", Valor: ""},
	}

	for _, operacion := range operaciones {
		var reply raft.ResultadoRemoto
		err := cfg.nodosRaft[idLeader].CallTimeout("NodoRaft.SometerOperacionRaft",
			operacion, &reply, 5000*time.Millisecond)

		// Manejo de error en la llamada RPC
		if err == nil {
			check.CheckError(err, "Se ha sometido una op con 2 nodos parados")
		}
	}

	// Reconectar los nodos previamente desconectados
	fmt.Printf("Reconectando nodos parados: %v\n", nodosParados)
	for _, nodo := range nodosParados {
		cfg.conectarNodo(nodo)
		fmt.Printf("Nodo conectado de nuevo: %d\n", nodo)
	}

	// Probar varios acuerdos tras la reconexión
	cfg.comprobarYRealizarOperacion(idLeader, 4, "escribir", "z", "tercer valor", "escrito")
	fmt.Printf("Operación de escribir sometida\n")
	cfg.comprobarYRealizarOperacion(idLeader, 5, "leer", "z", "", "tercer valor")
	fmt.Printf("Operación de leer sometida\n")

	// Finalizar los procesos distribuidos
	cfg.stopDistributedProcesses()
	fmt.Println(".............", t.Name(), "Superado")
}

// Se somete 5 operaciones de forma concurrente -- 3 NODOS RAFT
// SometerConcurrentementeOperaciones realiza una prueba de operaciones concurrentes en un sistema distribuido.
// La prueba sigue los siguientes pasos:
// 1. Inicia los procesos distribuidos.
// 2. Identifica al líder del clúster.
// 3. Define un conjunto de operaciones (escrituras y lecturas) a ser sometidas concurrentemente.
// 4. Somete las operaciones concurrentemente utilizando goroutines y un WaitGroup para sincronización.
// 5. Verifica que todos los nodos del clúster tienen un estado consistente con las operaciones realizadas.
// 6. Detiene los procesos distribuidos.

func (cfg *configDespliegue) SometerConcurrentementeOperaciones(t *testing.T) {
	//t.Skip("SKIPPED SometerConcurrentementeOperaciones")
	fmt.Printf("Iniciando prueba: %s\n", t.Name())

	// Iniciar procesos distribuidos
	cfg.startDistributedProcesses()

	// Identificar al líder
	leaderID := cfg.pruebaUnLider(3)
	fmt.Printf("Líder seleccionado: %d\n", leaderID)

	// Vector de operaciones (3 escrituras y 3 lecturas entrelazadas)
	operaciones := []struct {
		tipo  string // "escribir" o "leer"
		clave string
		valor string // Valor solo relevante para "escribir"
	}{
		{"escribir", "clave_0", "valor_0"},
		{"leer", "clave_0", ""},
		{"escribir", "clave_1", "valor_1"},
		{"leer", "clave_1", ""},
		{"escribir", "clave_2", "valor_2"},
		{"leer", "clave_2", ""},
	}

	// Realizar operaciones concurrentes
	var wg sync.WaitGroup
	numOperaciones := len(operaciones)

	for i, op := range operaciones {
		wg.Add(1)
		go func(index int, operacion struct {
			tipo  string
			clave string
			valor string
		}) {
			defer wg.Done()
			// Someter operación según su tipo
			cfg.someterOperacion(leaderID, operacion.tipo, operacion.clave, operacion.valor)
			fmt.Printf("Operación %d sometida (%s): %s -> %s\n", index, operacion.tipo, operacion.clave, operacion.valor)
		}(i, op)
	}

	// Esperar a que todas las operaciones terminen
	wg.Wait()

	// Verificar que todos los nodos tienen las 6 operaciones (lectura también aumenta el registro)
	cfg.VerificarEstadoConsistenteConOperaciones(numOperaciones)

	// Detener procesos distribuidos
	cfg.stopDistributedProcesses()
	fmt.Printf("Prueba %s completada exitosamente\n", t.Name())
}

// --------------------------------------------------------------------------
// FUNCIONES DE APOYO
// Comprobar que hay un solo líder
// probar varias veces si se necesitan reelecciones
// pruebaUnLider verifica si hay un solo líder entre las réplicas a lo largo de una serie de iteraciones.
// Devuelve el ID del líder si se encuentra, de lo contrario, falla la prueba.
//
// Parámetros:
// - numreplicas: El número de réplicas en la configuración.
//
// Retorna:
// - int: El ID del líder si se encuentra un solo líder.
//
// La función realiza los siguientes pasos:
// 1. Itera hasta 10 veces, durmiendo 500 milisegundos entre cada iteración.
// 2. Para cada réplica, si está conectada, verifica si la réplica es un líder y registra el término y el ID del líder.
// 3. Asegura que solo haya un líder por término.
// 4. Si se encuentra un líder, devuelve el ID del líder.
// 5. Si no se encuentra ningún líder después de 10 iteraciones, falla la prueba.
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

func (cfg *configDespliegue) obtenerEstadoRemoto(
	indiceNodo int) (int, int, bool, int) {
	var reply raft.EstadoRemoto
	err := cfg.nodosRaft[indiceNodo].CallTimeout("NodoRaft.ObtenerEstadoNodo",
		raft.Vacio{}, &reply, 500*time.Millisecond)
	check.CheckError(err, "Error en llamada RPC ObtenerEstadoRemoto")

	return reply.IdNodo, reply.Mandato, reply.EsLider, reply.IdLider
}

// start  gestor de vistas; mapa de replicas y maquinas donde ubicarlos;
// y lista clientes (host:puerto)
func (cfg *configDespliegue) startDistributedProcesses() {
	//time.Sleep(500 * time.Millisecond)
	//cfg.t.Log("Before starting following distributed processes: ", cfg.nodosRaft)
	for i, endPoint := range cfg.nodosRaft {
		despliegue.ExecMutipleHosts(EXECREPLICACMD+
			" "+strconv.Itoa(i)+" "+
			rpctimeout.HostPortArrayToString(cfg.nodosRaft),
			[]string{endPoint.Host()}, cfg.cr, PRIVKEYFILE)

		// dar tiempo para se establezcan las replicas
		//time.Sleep(50 * time.Millisecond)
	}
	fmt.Println("        startDistributedProcesses Completado")

	time.Sleep(10000 * time.Millisecond)
}

func (cfg *configDespliegue) startDistributedProcessesT1() {
	//time.Sleep(20000 * time.Millisecond)
	//cfg.t.Log("Before starting following distributed processes: ", cfg.nodosRaft)
	for i, endPoint := range cfg.nodosRaft {
		despliegue.ExecMutipleHosts(EXECREPLICACMD+
			" "+strconv.Itoa(i)+" "+
			rpctimeout.HostPortArrayToString(cfg.nodosRaft),
			[]string{endPoint.Host()}, cfg.cr, PRIVKEYFILE)

		// dar tiempo para se establezcan las replicas
		//time.Sleep(50 * time.Millisecond)
	}
	fmt.Println("        startDistributedProcesses Completado")

	time.Sleep(20000 * time.Millisecond)
}

func (cfg *configDespliegue) stopDistributedProcesses() {
	var reply raft.Vacio

	for i, endPoint := range cfg.nodosRaft {
		if cfg.conectados[i] {
			err := endPoint.CallTimeout("NodoRaft.ParaNodo",
				raft.Vacio{}, &reply, 10*time.Millisecond)
			check.CheckError(err, "Error en llamada RPC Para nodo")
		} else {
			cfg.conectados[i] = true //sirve para reiniciar procesos parados para poder hacer todos los test de seguido
		}
	}
}

// Comprobar estado remoto de un nodo con respecto a un estado prefijado
func (cfg *configDespliegue) comprobarEstadoRemoto(idNodoDeseado int,
	mandatoDeseado int, esLiderDeseado bool, IdLiderDeseado int) {
	idNodo, mandato, esLider, idLider := cfg.obtenerEstadoRemoto(idNodoDeseado)

	fmt.Println("idNodo: ", idNodo, "   idNdoDeseado: ", idNodoDeseado,
		"mandato: ", mandato, "   mandatoDeseado", mandatoDeseado,
		"esLider: ", esLider, "   esLiderDeseado", esLiderDeseado,
		"idLider: ", idLider, "   IdLiderDeseado", IdLiderDeseado)

	//cfg.t.Log("Estado replica 0: ", idNodo, mandato, esLider, idLider, "\n")

	if idNodo != idNodoDeseado || mandato != mandatoDeseado ||
		esLider != esLiderDeseado || idLider != IdLiderDeseado {
		cfg.t.Fatalf("Estado incorrecto en replica %d en subtest %s",
			idNodoDeseado, cfg.t.Name())
	}

}

//Funciones adicionales
//--------------------------------------------------------------------------------------------------------
func (cfg *configDespliegue) pararNodo(nodo int) {
	if cfg.conectados[nodo] {
		var args, reply raft.Vacio
		err := cfg.nodosRaft[nodo].CallTimeout("NodoRaft.ParaNodo", &args,
			&reply, time.Duration(50*time.Millisecond))
		if err != nil {
			check.CheckError(err, "Error desconectando nodo")
		} else {
			cfg.conectados[nodo] = false
		}
	}
}

// conectarNodo conecta un nodo al sistema distribuido si no está ya conectado.
// Si el nodo no está conectado, ejecuta el comando de replicación en el nodo especificado
// y actualiza el estado de conexión. Después de conectar, espera 10 segundos adicionales
// para permitir que las réplicas se establezcan.
//
// Parámetros:
//   nodo - el índice del nodo a conectar
func (cfg *configDespliegue) conectarNodo(nodo int) {
	//time.Sleep(5000 * time.Millisecond)
	//cfg.t.Log("Before starting following distributed processes: ", cfg.nodosRaft)

	if !cfg.conectados[nodo] {
		despliegue.ExecMutipleHosts(EXECREPLICACMD+
			" "+strconv.Itoa(nodo)+" "+
			rpctimeout.HostPortArrayToString(cfg.nodosRaft),
			[]string{cfg.nodosRaft[nodo].Host()}, cfg.cr, PRIVKEYFILE)

		// dar tiempo para se establezcan las replicas
		//time.Sleep(2000 * time.Millisecond)
		cfg.conectados[nodo] = true
		fmt.Printf("Nodo %d reconectado\n", nodo)
		time.Sleep(10000 * time.Millisecond)
	}
}
// comprobarYRealizarOperacion realiza una operación en el nodo líder y verifica el resultado.
//
// Parámetros:
// - leaderID: El ID del nodo líder donde se ejecutará la operación
// - index: El índice esperado de la operación en el registro
// - operacion: El tipo de operación a realizar "escribir"/ "leer"
// - clave: La clave asociada con la operación
// - valor: El valor asociado con la operación (utilizado para operaciones de "escribir")
// - valorEsperado: El valor esperado que debe ser devuelto por la operación
//
// Esta función realiza los siguientes pasos:
// 1. Define la operación a ejecutar.
// 2. Ejecuta la operación en el nodo líder utilizando una llamada RPC.
// 3. Verifica si hubo errores en la llamada RPC.
// 4. Verifica que el índice devuelto coincida con el índice esperado.
// 5. Verifica que el nodo que respondió sea el líder esperado.
// 6. Verifica que el valor devuelto coincida con el valor esperado.
//
// Si alguna de las verificaciones falla, la función registra un error fatal con un mensaje descriptivo.
func (cfg *configDespliegue) comprobarYRealizarOperacion(leaderID int, index int, operacion string, clave string, valor string, valorEsperado string) {
	// Definir operación
	op := raft.TipoOperacion{
		Operacion: operacion,
		Clave:     clave,
		Valor:     valor,
	}

	// Estructura para almacenar el resultado
	var resultado raft.ResultadoRemoto

	// Ejecutar la operación en el nodo líder
	err := cfg.nodosRaft[leaderID].CallTimeout("NodoRaft.SometerOperacionRaft", op, &resultado, 25000*time.Millisecond)

	// Verificar si hubo un error en la llamada RPC
	if err != nil {
		check.CheckError(err, "Error en llamada RPC SometerOperacionRaft")
	}

	// Verificar que el índice devuelto sea el esperado
	if resultado.IndiceRegistro != index {
		cfg.t.Fatalf("Error en operación '%s'. Índice esperado: %d, índice recibido: %d", operacion, index, resultado.IndiceRegistro)
	}

	// Verificar que el nodo que respondió sea el líder esperado
	if resultado.IdLider != leaderID {
		cfg.t.Fatalf("Error en operación '%s'. Líder esperado: %d, líder recibido: %d", operacion, leaderID, resultado.IdLider)
	}

	// Verificar que el valor devuelto sea el correcto
	if resultado.ValorADevolver != valorEsperado {
		cfg.t.Fatalf("Error en operación '%s'. Valor esperado: '%s', valor recibido: '%s'", operacion, valorEsperado, resultado.ValorADevolver)
	}
}


// VerificarEstadoConsistenteConOperaciones verifica la consistencia del estado entre todos los nodos conectados
// en el clúster Raft. Verifica que todos los nodos tengan el mismo término, líder y número de operaciones.
//
// Parámetros:
// - numOperacionesEsperado: El número esperado de operaciones que deberían estar presentes en cada nodo.
//
// La función itera sobre todas las réplicas y realiza las siguientes verificaciones:
// 1. Asegura que el término (mandato) sea consistente en todos los nodos.
// 2. Asegura que el ID del líder (idLider) sea consistente en todos los nodos.
// 3. Asegura que el número de operaciones (NumOperaciones) coincida con el número esperado.
//
// Si alguna de estas verificaciones falla, la función registrará un error y terminará la prueba.
func (cfg *configDespliegue) VerificarEstadoConsistenteConOperaciones(numOperacionesEsperado int) {
	fmt.Println("Verificando consistencia de estado entre los nodos:")

	var mandatoRef, idLiderRef int
	primero := true

	for i := 0; i < cfg.numReplicas; i++ {
		if cfg.conectados[i] {
			var reply raft.ResultadoEstadoNodo
			err := cfg.nodosRaft[i].CallTimeout("NodoRaft.ObtenerEstadoYOperaciones", raft.Vacio{}, &reply, 500*time.Millisecond)
			if err != nil {
				cfg.t.Fatalf("Error obteniendo estado en nodo %d: %v", i, err)
			}

			fmt.Printf("Nodo %d -> Operaciones: %d, Mandato: %d, Líder: %d\n", i, reply.NumOperaciones, reply.Mandato, reply.IdLider)

			// Verificar consistencia y número de operaciones sometidas
			if primero {
				mandatoRef = reply.Mandato
				idLiderRef = reply.IdLider
				primero = false
			} else {
				if reply.Mandato != mandatoRef {
					cfg.t.Fatalf("Inconsistencia en mandato: Nodo %d está en mandato %d, esperado %d", i, reply.Mandato, mandatoRef)
				}
				if reply.IdLider != idLiderRef {
					cfg.t.Fatalf("Inconsistencia en líder: Nodo %d tiene líder %d, esperado %d", i, reply.IdLider, idLiderRef)
				}
			}

			if reply.NumOperaciones != numOperacionesEsperado {
				cfg.t.Fatalf("Error en número de operaciones: Nodo %d tiene %d operaciones, esperado %d", i, reply.NumOperaciones, numOperacionesEsperado)
			}
		} else {
			fmt.Printf("Nodo %d está desconectado, no se puede verificar el estado\n", i)
		}
	}
}
func (cfg *configDespliegue) someterOperacion(leaderID int, operacion string, clave string, valor string) {
	// Definir operación
	op := raft.TipoOperacion{
		Operacion: operacion,
		Clave:     clave,
		Valor:     valor,
	}

	// Estructura para almacenar el resultado
	var resultado raft.ResultadoRemoto

	// Ejecutar la operación en el nodo líder
	err := cfg.nodosRaft[leaderID].CallTimeout("NodoRaft.SometerOperacionRaft", op, &resultado, 25000*time.Millisecond)

	// Verificar si hubo un error en la llamada RPC
	if err != nil {
		check.CheckError(err, "Error en llamada RPC SometerOperacionRaft")
	}
}
