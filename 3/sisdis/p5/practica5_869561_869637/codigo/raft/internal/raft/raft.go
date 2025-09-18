/****************************************
 * Autor: Adrián Nasarre Sánchez 869561
 * Autor: Héctor Lacueva Sacristán 869637
 * Fecha: Curso 24-25
 * Asignatura: Sistemas Distribuidos
 * Archivo: raft.go
 *****************************************/

package raft

// #region 1. API a exportar

//
// API
// ===
// Este es el API que vuestra implementación debe exportar
//
// nodoRaft = NuevoNodo(...)
//   Crear un nuevo servidor del grupo de elección.
//
// nodoRaft.Para()
//   Solicitar la parado de un servidor
//
// nodo.ObtenerEstado() (yo, mandato, esLider)
//   Solicitar a un nodo de elección por "yo", su mandato en curso,
//   y si piensa que es el msmo el lider
//
// nodoRaft.SometerOperacion(operacion interface()) (indice, mandato, esLider)

// type AplicaOperacion

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"

	//"crypto/rand"
	"sync"
	"time"

	//"net/rpc"

	"raft/internal/comun/rpctimeout"
)

const (
	// Constante para fijar valor entero no inicializado
	IntNOINICIALIZADO = -1

	//  false deshabilita por completo los logs de depuracion
	// Aseguraros de poner kEnableDebugLogs a false antes de la entrega
	kEnableDebugLogs = true

	// Poner a true para logear a stdout en lugar de a fichero
	kLogToStdout = false

	// Cambiar esto para salida de logs en un directorio diferente
	kLogOutputDir = "./logs_raft/"
)

type TipoOperacion struct {
	Operacion string // La operaciones posibles son "leer" y "escribir"
	Clave     string
	Valor     string // en el caso de la lectura Valor = ""
}

// A medida que el nodo Raft conoce las operaciones de las  entradas de registro
// comprometidas, envía un AplicaOperacion, con cada una de ellas, al canal
// "canalAplicar" (funcion NuevoNodo) de la maquina de estados
type AplicaOperacion struct {
	Indice    int // en la entrada de registro
	Operacion TipoOperacion
}

// EntradaRegistro es una estructura que representa una entrada de registro
// en el Log de un nodo Raft.
type EntradaRegistro struct {
	Indice  int           // Indice de la entrada en el Log
	Mandato int           // Mandato al que pertenece esta entrada
	Comando TipoOperacion // Comando a ejecutar en la máquina de estados
}

type Estado int

const (
	Seguidor  Estado = iota // iota empieza en 0
	Candidato               // 1
	Lider                   // 2
)

// Tipo de dato Go que representa un solo nodo (réplica) de raft
type NodoRaft struct {
	Mux sync.Mutex // Mutex para proteger acceso a estado compartido

	// Host:Port de todos los nodos (réplicas) Raft, en mismo orden
	Nodos   []rpctimeout.HostPort
	Yo      int // indice de este nodos en campo array "nodos"
	IdLider int
	// Utilización opcional de este logger para depuración
	// Cada nodo Raft tiene su propio registro de trazas (logs)
	Logger *log.Logger

	// ! Vuestros datos aqui.

	// ! ------------- ESTADO PERSISTENTE EN TODOS LOS NODOS RAFT -------------

	// * Mandato actual del nodo Raft
	MandatoActual int
	// * Id del candidato al que ha votado este nodo en este mandato
	MiVoto int
	// * Registro de operaciones de este nodo Raft
	Log []EntradaRegistro
	// * Estado actual de la máquina de estados del nodo Raft
	Estado Estado

	// ! -----------------------------------------------------------------------

	// ! -------------- CANALES Y ALMACEN PARA APLICAR OPERACIONES ------------
	
	// * Canal aplicar operación
	CanalAplicarOperacion chan AplicaOperacion
	// * Map de canales por el que recibes el resultado de aplicar una operación
	canalesValorADevolver map[int]chan string // clave: índice del log

	// * Almacen donde se escribirán o leerán datos
	Almacen map[string]string

	// ! -----------------------------------------------------------------------

	// * Canal para recibir votos realizados
	VotacionRecibida chan bool
	// * Canal para recibir latidos por parte del líder
	Latido chan bool

	// ! ------------- ESTADO VOLATIL EN TODOS LOS NODOS RAFT -----------------

	// * Índice de la última entrada de registro comprometida
	IndiceCommit int
	// * Índice de la entrada más alta del log que ha sido aplicada a la 
	// * máquina de estados
	IndiceUltimoAplicado int

	// ! -----------------------------------------------------------------------
	// ! ------------- ESTADO VOLATIL EN LÍDERES RAFT --------------------------

	// * Indices de la siguiente entrada del log a enviar al servidor
	SiguienteIndice []int
	// * Indices de la última entrada del log replicada en cada servidor
	IndiceEntradaReplicada []int

	// ! -----------------------------------------------------------------------

}

// Creacion de un nuevo nodo de eleccion
//
// Tabla de <Direccion IP:puerto> de cada nodo incluido a si mismo.
//
// <Direccion IP:puerto> de este nodo esta en nodos[yo]
//
// Todos los arrays nodos[] de los nodos tienen el mismo orden

// canalAplicar es un canal donde, en la practica 5, se recogerán las
// operaciones a aplicar a la máquina de estados. Se puede asumir que
// este canal se consumira de forma continúa.
//
// NuevoNodo() debe devolver resultado rápido, por lo que se deberían
// poner en marcha Gorutinas para trabajos de larga duracion
func NuevoNodo(nodos []rpctimeout.HostPort, yo int,
	canalAplicarOperacion chan AplicaOperacion) *NodoRaft {
	nr := &NodoRaft{}
	nr.Nodos = nodos
	nr.Yo = yo
	nr.IdLider = -1

	if kEnableDebugLogs {
		nombreNodo := nodos[yo].Host() + "_" + nodos[yo].Port()
		//fmt.Println("nombreNodo: ", nombreNodo)

		if kLogToStdout {
			nr.Logger = log.New(os.Stdout, nombreNodo+" -->> ",
				log.Lmicroseconds|log.Lshortfile)
		} else {
			err := os.MkdirAll(kLogOutputDir, os.ModePerm)
			if err != nil {
				panic(err.Error())
			}
			logOutputFile, err := os.OpenFile(
				fmt.Sprintf("%s/%s.txt", kLogOutputDir, nombreNodo),
				os.O_RDWR|os.O_CREATE|os.O_TRUNC,
				0755)
			if err != nil {
				panic(err.Error())
			}
			nr.Logger = log.New(logOutputFile,
				nombreNodo+" -> ", log.Lmicroseconds|log.Lshortfile)
		}
		nr.Logger.Println("logger initialized")
	} else {
		nr.Logger = log.New(io.Discard, "", 0)
	}

	// Añadir codigo de inicialización
	// * Inicializar los campos del nodo Raft
	nr.MandatoActual = 0
	nr.MiVoto = IntNOINICIALIZADO
	nr.Log = make([]EntradaRegistro, 0)
	// * Añadimos una entrada en el índice 0 ya que empieza en el 1.
	nr.Log = append(nr.Log, EntradaRegistro{
		Mandato: 0,
		Indice: 0,
		Comando: TipoOperacion{
			Operacion: "start",
			Clave: "",
			Valor: "",
		},
	}) ;
	nr.Estado = Seguidor

	nr.CanalAplicarOperacion = canalAplicarOperacion
	nr.canalesValorADevolver = make(map[int] chan string)
	nr.Almacen = make(map[string]string)

	nr.Latido = make(chan bool)
	nr.VotacionRecibida = make(chan bool)

	nr.IndiceCommit = 0
	nr.IndiceUltimoAplicado = 0
	nr.SiguienteIndice = make([]int, len(nodos))
	nr.IndiceEntradaReplicada = make([]int, len(nodos))

	// * Lanzar la máquina de estados del sistema
	go nr.tratarOperaciones()

	// * Lanzar la máquina de estados
	go nr.tratarNodo()

	return nr
}

// Metodo Para() utilizado cuando no se necesita mas al nodo
//
// Quizas interesante desactivar la salida de depuracion
// de este nodo
func (nr *NodoRaft) para() {
	
	go func() { 
		time.Sleep(10 * time.Millisecond);
		os.Exit(0)
	}()
}

// Devuelve "yo", mandato en curso y si este nodo cree ser lider
//
// Primer valor devuelto es el indice de este  nodo Raft el el conjunto de nodos
// la operacion si consigue comprometerse.
// El segundo valor es el mandato en curso
// El tercer valor es true si el nodo cree ser el lider
// Cuarto valor es el lider, es el indice del líder si no es él
func (nr *NodoRaft) obtenerEstado() (int, int, bool, int) {
	var yo int = nr.Yo
	var mandato int
	var esLider bool
	var idLider int 

	// Vuestro codigo aqui
	nr.Mux.Lock()
	mandato = nr.MandatoActual
	esLider = (nr.IdLider == nr.Yo)
	idLider = nr.IdLider
	nr.Mux.Unlock()

	return yo, mandato, esLider, idLider
}

func (nr *NodoRaft) obtenerEstadoReplicacion() ([]EntradaRegistro, 
	map[string]string) {
		nr.Mux.Lock()
		log := nr.Log
		almacen := nr.Almacen

		nr.Logger.Println("^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^")
		nr.Logger.Println("---------- Estado del log y el almacen ---------")
		printLog(nr.Logger, log)
		printAlmacen(nr.Logger, almacen)
		nr.Logger.Println("vvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvv")
		nr.Mux.Unlock()

		return log, almacen
	}

// El servicio que utilice Raft (base de datos clave/valor, por ejemplo)
// Quiere buscar un acuerdo de posicion en registro para siguiente operacion
// solicitada por cliente.

// Si el nodo no es el lider, devolver falso
// Sino, comenzar la operacion de consenso sobre la operacion y devolver en
// cuanto se consiga
//
// No hay garantía que esta operación consiga comprometerse en una entrada de
// de registro, dado que el lider puede fallar y la entrada ser reemplazada
// en el futuro.
// Resultado de este método :
// - Primer valor devuelto es el indice del registro donde se va a colocar
// - la operacion si consigue comprometerse.
// - El segundo valor es el mandato en curso
// - El tercer valor es true si el nodo cree ser el lider
// - Cuarto valor es el lider, es el indice del líder si no es él
// - Quinto valor es el resultado de aplicar esta operación en máquina de estados
func (nr *NodoRaft) someterOperacion(operacion TipoOperacion) (int, int,
	bool, int, string) {
	indice := -1
	mandato := -1
	EsLider := false
	idLider := -1
	valorADevolver := ""

	nr.Mux.Lock()

	nr.Logger.Println("Operación recibida: ", operacion)

	nr.Logger.Println("Nodo: ", nr.Yo, ", recibida: ", operacion)

	EsLider = (nr.IdLider == nr.Yo)
	idLider = nr.IdLider

	if EsLider {

		// * Índice de la entrada de registro que se va a rellenar
		indice = len(nr.Log)

		// * Añaadimos la entrada de registro al log
		nuevaEntradaRegistro := EntradaRegistro{
			Indice:  indice,
			Mandato: nr.MandatoActual,
			Comando: operacion,
		}

		nr.Log = append(nr.Log, nuevaEntradaRegistro)

		// * Establecemos los valores de retorno
		mandato = nr.MandatoActual

		nr.Logger.Println("Añadida entrada al log")
		nr.Logger.Println(nr.Log)

		ch := make(chan string, 1)
		nr.canalesValorADevolver[indice] = ch

		nr.Logger.Println("Intentando comprometer la operación")
		nr.Mux.Unlock()
		go nr.comprometerOperacion()

		nr.Logger.Println("SometerOperacion: Esperando al resultado de aplicar operación")
		valorADevolver = <-nr.canalesValorADevolver[indice]

	} else {
		nr.Mux.Unlock()
	}

	return indice, mandato, EsLider, idLider, valorADevolver
}

// #endregion







// #region 2. LLAMADAS RPC al API

// -----------------------------------------------------------------------
// LLAMADAS RPC al API
//
// Si no tenemos argumentos o respuesta estructura vacia (tamaño cero)
type Vacio struct{}

func (nr *NodoRaft) ParaNodo(args Vacio, reply *Vacio) error {

	defer nr.para()
	return nil
}

type EstadoParcial struct {
	Mandato int
	EsLider bool
	IdLider int
}

type EstadoRemoto struct {
	IdNodo int
	EstadoParcial
}

type EstadoReplicacionRemoto struct {
	Log []EntradaRegistro
	Almacen map[string]string
}

func (nr *NodoRaft) ObtenerEstadoReplicacionNodo(args Vacio, 
	reply *EstadoReplicacionRemoto) error {
		reply.Log, reply.Almacen = nr.obtenerEstadoReplicacion()
		return nil
	}

func (nr *NodoRaft) ObtenerEstadoNodo(args Vacio, reply *EstadoRemoto) error {
	reply.IdNodo, reply.Mandato, 
	reply.EsLider, reply.IdLider = nr.obtenerEstado()
	return nil
}

type ResultadoRemoto struct {
	ValorADevolver string
	IndiceRegistro int
	EstadoParcial
}

func (nr *NodoRaft) SometerOperacionRaft(operacion TipoOperacion,
	reply *ResultadoRemoto) error {

	reply.IndiceRegistro, reply.Mandato, reply.EsLider,
		reply.IdLider, reply.ValorADevolver = nr.someterOperacion(operacion)

	return nil
}

// #endregion







// #region 3. LLAMADAS RPC protocolo RAFT

// -----------------------------------------------------------------------
// LLAMADAS RPC protocolo RAFT
//
// Structura de ejemplo de argumentos de RPC PedirVoto.
//
// Recordar
// -----------
// Nombres de campos deben comenzar con letra mayuscula !
type ArgsPeticionVoto struct {
	// Vuestros datos aqui

	Mandato          int // * Mandato del nodo candidato
	IdCandidato      int // * Identificador del candidato
	UltimoIndiceLog  int // * Último índice del log del candidato
	UltimoMandatoLog int // * Último mandato del log del candidato

}

// Structura de ejemplo de respuesta de RPC PedirVoto,
//
// Recordar
// -----------
// Nombres de campos deben comenzar con letra mayuscula !
type RespuestaPeticionVoto struct {
	// Vuestros datos aqui

	Mandato  int  // * Mandato para que el candidato se actulice
	VotoDado bool // * true -> recibo voto |||| false -> NO recibo voto

}

func (nr *NodoRaft) concederVoto(peticion *ArgsPeticionVoto,
	reply *RespuestaPeticionVoto) error {

	nr.MiVoto = peticion.IdCandidato
	reply.Mandato = peticion.Mandato
	reply.VotoDado = true

	nr.Logger.Printf(
		"Nodo: %d. Concediendo voto a candidato: %d, mandato: %d\n",
		nr.Yo, peticion.IdCandidato, peticion.Mandato)

	return nil
}

/**
 * * Método de tratamiento de llamadas RPC PedirVoto
 * * Este método se encarga de gestionar las peticiones de voto de otros nodos.
 */
func (nr *NodoRaft) PedirVoto(peticion *ArgsPeticionVoto,
	reply *RespuestaPeticionVoto) error {

	nr.Mux.Lock()
	defer nr.Mux.Unlock()

	if peticion.Mandato >= nr.MandatoActual {

		// * Si la petición tiene un mandato mayor al mio quiere decir que
		// * no he votado en ese mandato y, por tanto, mi voto debe ser nulo
		if(peticion.Mandato > nr.MandatoActual) {
			nr.MiVoto = IntNOINICIALIZADO
		}

		// * Indicamos que se ha recibido una votación para resetear el timeout
		// * de los nodos seguidores.
		nr.VotacionRecibida <- true

		ultimoIndiceLogLocal := len(nr.Log) - 1
		ultimoMandatoLogLocal := nr.Log[ultimoIndiceLogLocal].Mandato

		distintoMandato := peticion.UltimoMandatoLog != ultimoMandatoLogLocal 
		mandatoLocalMayor := peticion.Mandato < ultimoMandatoLogLocal
		indiceLocalMayor := peticion.UltimoIndiceLog < ultimoIndiceLogLocal

		logCandidatoAlDia := !( (distintoMandato && mandatoLocalMayor) || 
						(!distintoMandato && indiceLocalMayor))
		
		noVoto := (nr.MiVoto == IntNOINICIALIZADO)
		votoCandidato := (nr.MiVoto == peticion.IdCandidato)

		// * Si no ha votado o ya me ha votado y el log del candidato 
		// * está, por lo menos, al día con respecto al del receptor
		if ((noVoto || votoCandidato)) && (logCandidatoAlDia) {

			// * Concedemos el voto
			return nr.concederVoto(peticion, reply)
		}
	}

	reply.Mandato = nr.MandatoActual
	reply.VotoDado = false

	return nil
}

type ArgAppendEntries struct {
	// Vuestros datos aqui
	Mandato          int               // * Mandato del nodo líder
	IdLider          int               // * Identificador del líder
	MandatoPrevioLog int               // * Mandato de la entrada de IndiceUltimoLog
	IndicePrevioLog  int               // * Indice de la última entrada del líder
	Entradas         []EntradaRegistro // * Entradas de registro a añadir
	CommitLider      int               // * IndiceCommit del líder

}

type Results struct {
	// Vuestros datos aqui
	Mandato int  // * Mandato para que se actualice el líder
	Exito   bool // * true si el seguidor tenía el log actualizado
}

// Metodo de tratamiento de llamadas RPC AppendEntries
// * Funciona ya que solo hay 1 líder por mandato
func (nr *NodoRaft) AppendEntries(args *ArgAppendEntries,
	results *Results) error {
	// Completar....
	nr.Mux.Lock()
	defer nr.Mux.Unlock()

	if len(args.Entradas) == 0 {
		nr.Logger.Println("LATIDO RECIBIDO -------------------------------------------------------------")

		/*nr.printNodo()

		printArgsAppendEntries(nr.Logger, *args)


		nr.Logger.Println("-----------------------------------------------------------------------------")*/

	} else {
		nr.Logger.Println("APPEND ENTRIES RECIBIDO -------------------------------------------------------------")

		nr.printNodo()

		printArgsAppendEntries(nr.Logger, *args)

		nr.Logger.Println("-----------------------------------------------------------------------------")
	}

	results.Exito = false

	if args.Mandato < nr.MandatoActual {
		results.Mandato = nr.MandatoActual
		return nil
	}

	if len(args.Entradas) != 0 {
		nr.Logger.Println("appendEntries: args.Mandato >= nr.MandatoActual")
	}

	// * Si args.Mandato es mayor o igual al mandato actual, hay un líder vivo
	// * Actualizo el IdLider y el MandatoActual
	nr.IdLider = args.IdLider
	nr.MandatoActual = args.Mandato
	results.Mandato = args.Mandato

	if len(args.Entradas) != 0 {
		nr.Logger.Println("appendEntries: Results.Mandato:", results.Mandato)
		nr.Logger.Println("appendEntries: nr.MandatoActual:", nr.MandatoActual)
		nr.Logger.Println("appendEntries: args.Mandato:", args.Mandato)
	}

	// * y se considera que se ha recibido un latido.
	nr.Latido <- true

	if len(args.Entradas) != 0 {

		lenLog := len(nr.Log)

		// * Si el log previo no es correcto devuelvo !Exito en AppendEntries
		if !nr.logPrevioCorrecto(args, lenLog) {
			nr.Logger.Println("Log previo incorrecto")
			return nil
		}

		nr.Logger.Println("appendEntries: Voy a actualizar el log")
		// * Si hay entradas hay que actualizar el log
		nr.actualizarLog(args)
		
	}

	// * Actualizo el IndiceCommit
	if args.CommitLider > nr.IndiceCommit {
		
		nr.IndiceCommit = min(args.CommitLider, len(nr.Log)-1)
		
		if len(args.Entradas) != 0 {
			nr.Logger.Println("appendEntries: nuevo indiceCommit", nr.IndiceCommit)
		}

		// * Si ha cambiado el indice de commit del líder
		for i:=nr.IndiceUltimoAplicado+1;
		 i <= nr.IndiceCommit ; i++ {
			go nr.aplicarOperacion(i)
		}
	
	}

	if len(args.Entradas) != 0 {
		nr.Logger.Println("appendEntries: exitoso")
	}

	// * RPC AppendEntries exitosa
	results.Exito = true

	return nil

}

// #endregion








// #region 4. MÉTODOS conexión entre nodos

// --------------------------------------------------------------------------
// ----- METODOS/FUNCIONES desde nodo Raft, como cliente, a otro nodo Raft
// --------------------------------------------------------------------------

// Ejemplo de código enviarPeticionVoto
//
// nodo int -- indice del servidor destino en nr.nodos[]
//
// args *RequestVoteArgs -- argumentos par la llamada RPC
//
// reply *RequestVoteReply -- respuesta RPC
//
// Los tipos de argumentos y respuesta pasados a CallTimeout deben ser
// los mismos que los argumentos declarados en el metodo de tratamiento
// de la llamada (incluido si son punteros)
//
// Si en la llamada RPC, la respuesta llega en un intervalo de tiempo,
// la funcion devuelve true, sino devuelve false
//
// la llamada RPC deberia tener un timeout adecuado.
//
// Un resultado falso podria ser causado por una replica caida,
// un servidor vivo que no es alcanzable (por problemas de red ?),
// una petición perdida, o una respuesta perdida
//
// Para problemas con funcionamiento de RPC, comprobar que la primera letra
// del nombre de todo los campos de la estructura (y sus subestructuras)
// pasadas como parametros en las llamadas RPC es una mayuscula,
// Y que la estructura de recuperacion de resultado sea un puntero a estructura
// y no la estructura misma.
func (nr *NodoRaft) enviarPeticionVoto(nodo int, args *ArgsPeticionVoto,
	reply *RespuestaPeticionVoto) bool {

	//fmt.Println("Petición: ", nodo, args, reply)

	Nodo := nr.Nodos[nodo]

	// Pedimos el voto al nodo
	err := Nodo.CallTimeout("NodoRaft.PedirVoto", args,
		reply, 300*time.Millisecond)

	//fmt.Println("Respuesta: ", nodo, reply)

	// Si se produce algún error durante la petición devuelve false
	// Sino, devuelve true
	return (err == nil)

}

/**
 ** Método que envía un latido a un nodo, devuelve la respuesta y si todo va bien
 ** devuelve true
 */
func (nr *NodoRaft) enviarAppendEntries(nodo int, args *ArgAppendEntries,
	reply *Results) bool {

	//fmt.Println("Latido:", nodo, args, reply)

	Nodo := nr.Nodos[nodo]

	if len(args.Entradas) != 0 {
		nr.Mux.Lock()
		nr.Logger.Println("enviarAppendEntries: ENVIANDO APPEND ENTRIES", nodo)
		nr.Mux.Unlock()
	}

	// * Pedimos el voto al nodo
	err := Nodo.CallTimeout("NodoRaft.AppendEntries", args,
		reply, 300*time.Millisecond)

	if len(args.Entradas) != 0 {
		nr.Mux.Lock()
		nr.Logger.Println("enviarAppendEntries: RESPUESTA RECIBIDA DEL NODO", nodo)
		printArgsAppendEntries(nr.Logger, *args)
		printResultados(nr.Logger, reply)
		nr.Logger.Println("------------------------------------------------------")
		nr.Mux.Unlock()
	}
	
	// * Si se produce algún error durante la petición devuelve false
	// * Sino, devuelve true
	return (err == nil)
}

// #endregion








// #region 5. Métodos tratamiento respuestas

/**
 ** Método que se encarga de tratar una respuesta a una petición de voto
 ** Se encarga de llevar el recuento de votos recibidos para que,
 ** cuando haya mayoría el candidato pase a ser el líder
 */
func (nr *NodoRaft) tratarVotoValido(reply *RespuestaPeticionVoto,
	votosRecibidos *int, mayoria *int, votacionFinalizada *bool,
	victoria chan bool) {

	nr.Mux.Lock()
	defer nr.Mux.Unlock()

	// * Si recibo voto
	if reply.VotoDado {

		(*votosRecibidos)++
		//fmt.Println("Voto recibido, ya van: ", *votosRecibidos)

		// * Y no ha acabado la votación
		if !(*votacionFinalizada) {

			(*votacionFinalizada) = (*votosRecibidos >= *mayoria)
			if *votacionFinalizada {
				//fmt.Println("-------- VICTORIA --------")
				victoria <- true
			}
		}

	} else if reply.Mandato > nr.MandatoActual {
		// * Si no me dan el voto porque están en un mandato mayor
		// * Finalizo la votación y vuelvo a ser seguidor
		*votacionFinalizada = true
		nr.hacerseSeguidor(reply.Mandato)

	}
}

/**
 ** Método que se encarga de tratar la respuesta de un nodo ante un latido
 ** enviado por parte del líder. Si el latido no tiene éxito es porque
 ** el nodo que ha recibido el latido tiene mayor mandato.
 ** En ese caso, volvemos a ser seguidores.
 */
func (nr *NodoRaft) tratarLatidoValido(results *Results, nodoID int) {

	nr.Mux.Lock()
	defer nr.Mux.Unlock()

	//printResultados(nr.Logger, results)

	//fmt.Printf("Nodo %d, ha respondido al latido\n", nodoID)
	if !(*results).Exito {
		// * Si (!exito) --> mandato mayor en nodoID
		// * Pasar a follower y actualizar el mandato
		//fmt.Printf("Nodo %d, con mandato superior\n", nodoID)
		nr.hacerseSeguidor(results.Mandato)
	}

}

/**
 ** Método que se encarga de tratar la respuesta de un nodo ante una llamada
 ** a AppendEntries por parte del líder.
 */
func (nr *NodoRaft) tratarAppendEntriesValido(results *Results,
	nodoId int, indiceActual int, numReplicadas *int, primerIndiceEnviado int,
	mayoria int, finChan chan bool) bool {

	nr.Mux.Lock()

	// printTratarAppendEntriesValido(nr.Logger, results,
	//	nodoId, indiceActual, numReplicadas)

	if results.Exito {

		nr.SiguienteIndice[nodoId] = max(nr.SiguienteIndice[nodoId], indiceActual)
		nr.IndiceEntradaReplicada[nodoId] = nr.SiguienteIndice[nodoId] - 1

		// * Marco como replicada
		(*numReplicadas)++

		// * Si se alcanza la mayoría se informará al cliente
		if (*numReplicadas) == mayoria {
			nr.Logger.Println("tratarAppendEntriesValido: Mayoría alcanzada, aplicar operación en máquina de estados")
			// * Si se alcanza la mayoría aplico la operación y 
			// * devuelvo el resultado
			nr.IndiceCommit = max(nr.IndiceCommit, nr.IndiceEntradaReplicada[nodoId])
			indiceaAplicar := nr.IndiceUltimoAplicado + 1
			nr.Mux.Unlock()
			nr.aplicarOperacion(indiceaAplicar)

		} else if (*numReplicadas) == len(nr.Nodos) {
			// * Si todos han replicado correctamente, damos por finalizada la
			// * replicación.
			finChan <- true
			nr.Mux.Unlock()
		}

		nr.Mux.Lock()
		nr.Logger.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		nr.Logger.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		nr.Logger.Println(" SALIENDO DEL BUCLE FOR NODOid:", nodoId, "numReplicadas:", *numReplicadas )
		nr.Logger.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		nr.Logger.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		nr.Mux.Unlock()
		return true

	} else if results.Mandato > nr.MandatoActual {
		
		nr.hacerseSeguidor(results.Mandato)
		nr.Mux.Unlock()
		return true

	} else {
		nr.SiguienteIndice[nodoId] = primerIndiceEnviado - 1
		nr.Logger.Println("=========================================")
		nr.Logger.Println("=========================================")
		nr.printNodo()
		nr.Logger.Println("=========================================")
		nr.Logger.Println("=========================================")
		nr.Mux.Unlock()
		return false
	}
}

// #endregion








// #region 6. métodos de envio a un nodo

/**
 ** Método que se encarga del envío de un latido a cada nodo
 */
func (nr *NodoRaft) enviarLatidos() {

	// * Preparo los argumentos para el latido
	args := nr.argumentosLatido()

	// * Para cada nodo distinto del propio
	for i := 0; i < len(nr.Nodos); i++ {
		if i != nr.Yo {
			go func(nodoID int) {
				var reply Results
				// * Enviamos los latidos y recogemos los resultados
				latidoValido := nr.enviarAppendEntries(nodoID, &args, &reply)

				if latidoValido {

					nr.tratarLatidoValido(&reply, nodoID)

				} else {
					//fmt.Printf("ERROR al enviar latido al Nodo %d\n", nodoID)
				}
			}(i)
		}
	}
}

/**
 ** Método que se encarga del envío de una petición de voto a cada nodo
 */
func (nr *NodoRaft) enviarPeticionesVoto(args ArgsPeticionVoto,
	mayoria int, victoria chan bool) {

	// * Inicializar variables de control
	votacionFinalizada := false
	votosRecibidos := 1

	// * Para cada uno de los otros nodos
	for i := 0; i < len(nr.Nodos); i++ {
		if i != nr.Yo {
			go func(nodoID int) {
				var reply RespuestaPeticionVoto

				votoValido := nr.enviarPeticionVoto(nodoID, &args, &reply)

				if votoValido {
					nr.tratarVotoValido(&reply, &votosRecibidos,
						&mayoria, &votacionFinalizada, victoria)
				} else {
					//fmt.Printf("ERROR al enviar petición de voto al Nodo %d\n",
						//nodoID)
				}

			}(i)
		}
	}
}

// * Esta función se encarga de comprometer la operación
// * Envía AppendEntries a todos los nodos.
// * Si consigue consolidar, aplica a maquina de estados y manda
// * resultado a través de resultadoChan. El envío de AppendEntries
// * se mantiene hasta conseguir replicar en todos los nodos.
func (nr *NodoRaft) comprometerOperacion() {

	// * Si estoy parado o no soy el líder dejo de comprometer
	if nr.Estado != Lider {
		return
	}

	nr.Logger.Println("Comprometiendo operación")

	// * En el nodo líder ya ha sido replicada
	numReplicadas := 1

	finChan := make(chan bool)

	nr.Mux.Lock()
	mayoria := (len(nr.Nodos) / 2) + 1
	nr.Mux.Unlock()

	nr.Logger.Println("Enviando mensajes a los seguidores")
	// * AppendEntries a los nodos hasta que repliquen la operación
	// * en todos los nodos
	for i := 0; i < len(nr.Nodos); i++ {
		if i != nr.Yo {
			go nr.enviarOperacionNodo(i, &numReplicadas,
				finChan, mayoria)
		}
	}

	<-finChan
}

/**
 ** Método que se encarga de enviar una operación a un nodo hasta que este
 ** consiga replicar la operación o nos demos cuenta que ya no somos líderes.
 */
func (nr *NodoRaft) enviarOperacionNodo(nodoId int,
	numReplicadas *int, finChan chan bool, mayoria int) {

	var results Results

	fin := false

	for !fin {

		nr.Logger.Println("enviarOperacionNodo: Preparo argumentos para el nodo", nodoId)
		args := nr.argumentosAppendEntries(nodoId)

		// * Si al preparar los argumentos no hay entradas por añadir salgo
		// * ya que entonces el nodo réplica ya está actualizado.
		if len(args.Entradas) == 0{
			break
		}

		if nr.enviarAppendEntries(nodoId, &args, &results) {
			fin = nr.tratarAppendEntriesValido(&results, nodoId,
				args.Entradas[len(args.Entradas)-1].Indice + 1,
				numReplicadas, args.Entradas[0].Indice, mayoria,
				finChan)
			
		} else {
			//fmt.Printf("ERROR al enviar operacióno al Nodo %d\n", nodoId)
		}
		if (!fin) {
			time.Sleep(10 * time.Millisecond)
		}
	}
}

// #endregion








// #region 7. métodos auxiliares

func (nr *NodoRaft) actualizarLog(args *ArgAppendEntries) {
	// Buscar el primer conflicto entre el log local y las nuevas entradas

	indiceConflicto := -1
	for _, entrada := range args.Entradas {
		indiceLog := entrada.Indice
		if indiceLog < len(nr.Log) {
			if nr.Log[indiceLog].Mandato != entrada.Mandato {
				indiceConflicto = indiceLog
				break
			}
		}
	}

	if indiceConflicto != -1 {
		// Hay conflicto: truncar el log en el primer conflicto
		nr.Log = nr.Log[:indiceConflicto]
		// Añadir las entradas nuevas a partir del conflicto
		nr.Log = append(nr.Log, args.Entradas[indiceConflicto-(args.IndicePrevioLog+1):]...)
	} else {
		// * Si no hay conflicto y además todas las entradas de los argumentos
		// * ya están en el Log acabo
		if (len(args.Entradas) + args.IndicePrevioLog) < len(nr.Log) {
			return
		}
		// No hay conflicto, pero puede que falten entradas al final
		comienzo := len(nr.Log) - (args.IndicePrevioLog + 1)
		if comienzo < 0 {
			comienzo = 0
		}
		nr.Log = append(nr.Log, args.Entradas[comienzo:]...)
	}
}

/**
 ** Método para comprobar si la entrada del log previo es correcta.
 */
func (nr *NodoRaft) logPrevioCorrecto(args *ArgAppendEntries, lenLog int) bool {

	//fmt.Println("Indice previo log: ", args.IndicePrevioLog)
	//fmt.Println("LenLog: ", lenLog)

	// * Si no hay entradas en el log local con dicho índice, devolver falso
	if lenLog <= args.IndicePrevioLog {
		return false
	}

	// * Si no coinciden los mandatos, devolveer falso
	if nr.Log[args.IndicePrevioLog].Mandato != args.MandatoPrevioLog {
		return false
	}

	nr.Logger.Println("Log previo correcto")
	// * Sino, la entrada previa es correcta
	return true
}

/**
 ** Método que se encarga de llevar a cabo una elección
 */
func (nr *NodoRaft) iniciarEleccion(victoria chan bool) {

	nr.Mux.Lock()

	// * Aumenta el mandato y se vota a sí mismo
	nr.MandatoActual++
	nr.MiVoto = nr.Yo

	// * Preparamos los argumentos de la petición
	args := nr.argumentosPeticionVoto()

	nr.Mux.Unlock()

	// * Establecemos la mayoría
	mayoria := (len(nr.Nodos) / 2) + 1

	//fmt.Printf("Votos para mayoría: %d --- TENGO: %d --- FALTAN: %d\n",
	//	mayoria, 1, mayoria-1)

	// * Enviar las peticiones a los nodos
	nr.enviarPeticionesVoto(args, mayoria, victoria)

}

/**
 ** Método que devuelve al estado de seguidor al nodo y actualiza el mandato
 ** El mutex debe estar bloqueado para un funcionamiento correcto
 */
func (nr *NodoRaft) hacerseSeguidor(mandato int) {
	//fmt.Printf("Vuelvo a ser SEGUIDOR\n")
	nr.MandatoActual = mandato
	nr.Estado = Seguidor

}

// Funcion auxiliar para mínimo
func min(a int, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
}

func max(a int, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

// #endregion







// #region 8. métodos para argumentos ----------------

/**
 ** Método que prepara los argumentos para una petición de voto
 */
func (nr *NodoRaft) argumentosPeticionVoto() ArgsPeticionVoto {

	var args ArgsPeticionVoto

	// * Mandatos e Id propios (mandatoActual ya ha sido actualizado)
	args.Mandato = nr.MandatoActual
	args.IdCandidato = nr.Yo

	indicePrevio := len(nr.Log) - 1

	// * Siempre hay al menos una entrada en el Log por tanto, esto nunca falla
	args.UltimoIndiceLog = nr.Log[indicePrevio].Indice
	args.UltimoMandatoLog = nr.Log[indicePrevio].Mandato

	return args 
}

/**
 ** Método que prepara los argumentos para AppendEntries
 */
func (nr *NodoRaft) argumentosAppendEntries(nodoId int) ArgAppendEntries {

	nr.Mux.Lock()
	defer nr.Mux.Unlock()

	var indicePrevio int
	var mandatoPrevio int

	ultimaEntradaReplicada := nr.SiguienteIndice[nodoId] - 1

	var entradas []EntradaRegistro

	indicePrevio = nr.Log[ultimaEntradaReplicada].Indice
	mandatoPrevio = nr.Log[ultimaEntradaReplicada].Mandato
	entradas = nr.Log[nr.SiguienteIndice[nodoId]:]


	return ArgAppendEntries{
		Mandato:          nr.MandatoActual,
		IdLider:          nr.Yo,
		IndicePrevioLog:  indicePrevio,
		MandatoPrevioLog: mandatoPrevio,
		Entradas:         entradas,
		CommitLider:      nr.IndiceCommit,
	}

}

/**
 ** Método que prepara los argumentos para el envío de un latido
 */
func (nr *NodoRaft) argumentosLatido() ArgAppendEntries {

	nr.Mux.Lock()
	defer nr.Mux.Unlock()

	// Creamos los argumentos del latido
	args := ArgAppendEntries{
		Mandato:          nr.MandatoActual,
		IdLider:          nr.Yo,
		IndicePrevioLog:  IntNOINICIALIZADO,
		MandatoPrevioLog: IntNOINICIALIZADO,
		Entradas:         []EntradaRegistro{},
		CommitLider:      nr.IndiceCommit,
	}

	return args
}

// #endregion








// #region 9. autómata del nodo

// ! Esto seguramente se puede hacer más sencillo
func getRandomMilliseconds(minTime int, maxTime int) time.Duration {
	// Validar que minTime sea menor que maxTime
	if minTime >= maxTime {
		auxTime := minTime
		minTime = maxTime
		maxTime = auxTime
	}

	// Calcular el rango
	rangeMillis := maxTime - minTime

	// Generar un número aleatorio usando crypto/rand
	var randomInt uint64
	err := binary.Read(rand.Reader, binary.BigEndian, &randomInt)
	if err != nil {
		return 0
	}

	// Normalizar el número aleatorio al rango
	randomMillis := minTime + int(randomInt%uint64(rangeMillis))
	return time.Duration(randomMillis) * time.Millisecond
}

/**
* * Esta función se encarga de tratar el estado del nodo Raft.
* * Dependiendo del estado actual del nodo (Seguidor, Candidato o Líder),
* * se llama a la función correspondiente para manejar ese estado.
 */
func (nr *NodoRaft) tratarNodo() {
	for {

		switch nr.Estado {
		case Seguidor:
			nr.Logger.Printf("Nodo: %d. Estado: SEGUIDOR\n", nr.Yo)
			nr.tratarSeguidor()
		case Candidato:
			nr.Logger.Printf("Nodo: %d. Estado: CANDIDATO\n", nr.Yo)
			nr.tratarCandidato()
		case Lider:
			nr.Logger.Printf("Nodo: %d. Estado: LIDER\n", nr.Yo)
			nr.tratarLider()
		}
	}
}

/**
* * Esta función se encarga de tratar el estado de Seguidor del nodo Raft.
* * En este estado, el nodo espera recibir latidos del líder.
* * Se establece un timeout y se eespera a recibir uno de los
* * siguientes eventos:
* *  - Un latido recibido del líder, lo que indica que el nodo sigue siendo
* *    un seguidor y se resetea el voto.
* *  - Un timeout que indica que no se ha recibido un latido en el tiempo
* *    esperado, lo que provoca que el nodo cambie su estado a Candidato.
* *  - Un voto dado, lo que indica que hay un nodo que ya es candidato y no
* *    hace falta que este lo sea.
 */
func (nr *NodoRaft) tratarSeguidor() {

	timeOut := time.NewTimer(getRandomMilliseconds(300, 500))

	select {
	case <-timeOut.C:
		nr.Logger.Printf("Nodo: %d. Vence timeout, ELECCIÓN\n", nr.Yo)
		nr.Estado = Candidato
	case <-nr.Latido:
		nr.Logger.Printf("Nodo: %d. Latido RECIBIDO\n", nr.Yo)
	case <-nr.VotacionRecibida:
		nr.Logger.Printf("Nodo: %d. Voto DADO\n", nr.Yo)
	}

}

/**
* * Esta función se encarga de tratar el estado de Candidato del nodo Raft.
* * En este estado, se inicia un proceso de elección para convertirse en líder.
* * Se establece un timeout y se espera a recibir uno de los
* * siguientes eventos:
* *  - Un timeout que indica que la elección ha fallado.
* *  - Un latido recibido de un líder, lo que indica que el nodo debe volver
* *    al estado de Seguidor.
* *  - Una victoria en la elección, lo que indica que el nodo se ha convertido
* *    en líder.
 */
func (nr *NodoRaft) tratarCandidato() {

	timeOut := time.NewTimer(getRandomMilliseconds(150, 300))
	victoria := make(chan bool)
	go nr.iniciarEleccion(victoria)

	select {
	case <-timeOut.C:
		nr.Logger.Printf("Nodo: %d. Vence timeout, ELECCIÓN\n", nr.Yo)

	case <-nr.Latido:
		nr.Logger.Printf("Nodo: %d. Latido RECIBIDO\n", nr.Yo)
		nr.Mux.Lock()
		nr.Estado = Seguidor
		nr.Mux.Unlock()

	case <-victoria:
		nr.Logger.Printf("Nodo: %d. Elección GANADA\n", nr.Yo)
		nr.Mux.Lock()
		nr.Estado = Lider
		nr.IdLider = nr.Yo
		// * Reinicializar al ganar elección
		for i := 0; i < len(nr.Nodos); i++ {
			nr.SiguienteIndice[i] = len(nr.Log)
			nr.IndiceEntradaReplicada[i] = 0
		}
		nr.Mux.Unlock()
	}
}

/**
* * Esta función se encarga de tratar el estado de Líder del nodo Raft.
* * En este estado, el nodo envía latidos a los demás nodos para mantener
* * su liderazgo y asegurar que los demás nodos están al tanto de su estado.
* * Se establece un timeout y se espera a recibir uno de los
* * siguientes eventos:
* *  - Cada 50 milisegundos, se envía un latido a los demás nodos.
* *  - Un latido recibido de otro nodo, lo que indica que el nodo debe volver
* *    al estado de Seguidor.
 */
func (nr *NodoRaft) tratarLider() {

	timeOut := time.NewTimer(50 * time.Millisecond)

	select {
	case <-timeOut.C:
		nr.Logger.Printf("Nodo: %d. Vence timeout, Mandar Latido\n", nr.Yo)
		nr.enviarLatidos()
	case <-nr.Latido:
		nr.Logger.Printf("Nodo: %d. Latido RECIBIDO\n", nr.Yo)
		nr.Estado = Seguidor
	}
}

// #endregion








// #region 10. Método que se encarga de aplicar las operaciones


/**
 ** Método para aplicar operaciones y devolver el resultado
 */
func (nr *NodoRaft) tratarOperaciones() {

	var valorADevolver string

	for {
		comando := <-nr.CanalAplicarOperacion

		nr.Mux.Lock()

		nr.Logger.Println("Aplicando operacion:", comando)
		
		switch comando.Operacion.Operacion {
		case "leer":
			nr.Logger.Println("Intentando leer de la clave:", comando.Operacion.Clave)
			// * Si la clave no existe devuelve una cadena vacía
			valorADevolver = nr.Almacen[comando.Operacion.Clave]
			nr.Logger.Println("Valor leído:", valorADevolver)
		case "escribir":
			nr.Logger.Println("Intentando escribir en la clave:", comando.Operacion.Clave)
			nr.Almacen[comando.Operacion.Clave] = comando.Operacion.Valor
			valorADevolver = "Escrito correctamente"
		default:
			nr.Logger.Println("Intentando aplicar operación desconocida")
			valorADevolver = "Operación desconocida"
		}
		
		// * Si soy líder envío el resultado de la operación sino solo hago 
		// * los cambios
		if (nr.IdLider == nr.Yo) {
			if ch, ok := nr.canalesValorADevolver[comando.Indice]; ok {
				select {
				case ch <- valorADevolver:
				default:
					// si el cliente ya se fue o hubo timeout
				}
				delete(nr.canalesValorADevolver, comando.Indice) // limpieza
			}
		}
		nr.Mux.Unlock()

	}
}


// * Cada vez que esta operación es realizada, toma el valor el último aplicado
func (nr *NodoRaft) aplicarOperacion(IndiceSiguienteParaAplicar int) {

	for {
		nr.Mux.Lock()

		if nr.IndiceUltimoAplicado + 1 == IndiceSiguienteParaAplicar {
				// * Aplicamos en máquina de estados
			nr.CanalAplicarOperacion <- AplicaOperacion{
				Indice: IndiceSiguienteParaAplicar,
				Operacion: nr.Log[IndiceSiguienteParaAplicar].Comando,
			}

			nr.IndiceUltimoAplicado++

			nr.Mux.Unlock()
			break
		}
		nr.Mux.Unlock()
	}
	
}

// #endregion








// #region 11. Métodos print

func (nr *NodoRaft) printNodo() {
	nr.Logger.Println("========== Estado del NodoRaft ==========")
	nr.Logger.Printf("ID Nodo (Yo): %d", nr.Yo)
	nr.Logger.Printf("ID Líder: %d", nr.IdLider)
	nr.Logger.Printf("Nodos: %v", nr.Nodos)
	nr.Logger.Printf("Estado: %v", nr.Estado)

	// Estado persistente
	nr.Logger.Println("------ Estado Persistente ------")
	nr.Logger.Printf("MandatoActual: %d", nr.MandatoActual)
	nr.Logger.Printf("MiVoto: %d", nr.MiVoto)
	nr.Logger.Printf("Log (len=%d):", len(nr.Log))
	for i, entrada := range nr.Log {
		nr.Logger.Printf("  [%d] %v", i, entrada)
	}

	// Estado volátil común
	nr.Logger.Println("------ Estado Volátil (Todos) ------")
	nr.Logger.Printf("IndiceCommit: %d", nr.IndiceCommit)
	nr.Logger.Printf("IndiceUltimoAplicado: %d", nr.IndiceUltimoAplicado)

	// Estado volátil en líderes
	nr.Logger.Println("------ Estado Volátil (Líderes) ------")
	nr.Logger.Printf("SiguienteIndice: %v", nr.SiguienteIndice)
	nr.Logger.Printf("IndiceEntradaReplicada: %v", nr.IndiceEntradaReplicada)

	// Almacen (claves limitadas para evitar saturación visual)
	nr.Logger.Println("------ Almacén (claves limitadas) ------")
	maxKeys := 5
	count := 0
	for k, v := range nr.Almacen {
		nr.Logger.Printf("  %s: %s", k, v)
		count++
		if count >= maxKeys {
			nr.Logger.Printf("  ...y %d más", len(nr.Almacen)-maxKeys)
			break
		}
	}

	nr.Logger.Println("==========================================")
}

func printArgsAppendEntries(logger *log.Logger, args ArgAppendEntries) {
	logger.Println("----- ARGS ------")
	logger.Println("Mandato: ", args.Mandato)
	logger.Println("IdLider: ", args.IdLider)
	logger.Println("MandatoPrevioLog: ", args.MandatoPrevioLog)
	logger.Println("IndicePrevioLog: ", args.IndicePrevioLog)
	logger.Println("CommitLider: ", args.CommitLider)
	logger.Printf("Entradas (len=%d):", len(args.Entradas))
	for i, entrada := range args.Entradas {
		logger.Printf("  [%d] %v\n", i, entrada)
	}
	logger.Println("----- FIN ARGS ------")
}

func printReply(logger *log.Logger, reply *ResultadoRemoto) {
	logger.Println("----- RESPUESTA REMOTA -----")
	logger.Println("ValorADevolver:", reply.ValorADevolver)
	logger.Println("IndiceRegistro:", reply.IndiceRegistro)
	logger.Println("Mandato:", reply.Mandato)
	logger.Println("EsLider:", reply.EsLider)
	logger.Println("IdLider:", reply.IdLider)
	logger.Println("----- FIN RESPUESTA REMOTA -----")
}

func printResultados(logger *log.Logger, results *Results) {
	logger.Println("----- RESULTADOS ------")
	logger.Println("Exito:", (*results).Exito)
	logger.Println("Mandato:", (*results).Mandato)
	logger.Println("----- FIN RESULTADOS --")
}

func printTratarAppendEntriesValido(logger *log.Logger, results *Results, nodoId int,
	indiceActual int, numReplicadas *int) {
	logger.Println("______________________________________________________")
	printResultados(logger, results)
	logger.Println("Nodo ID:", nodoId)
	logger.Println("IndiceActual:", indiceActual)
	logger.Println("NumReplicadas", *numReplicadas)
	logger.Println("------------------------------------------------------")
}

func printPeticionVoto(logger *log.Logger, peticion *ArgsPeticionVoto) {
	logger.Println("_______________ PETICIÓN DE VOTO ________________")
	logger.Println("Mandato:", peticion.Mandato)
	logger.Println("IdCandidato:", peticion.IdCandidato)
	logger.Println("UltimoIndiceLog:", peticion.UltimoIndiceLog)
	logger.Println("UltimoMandatoLog:", peticion.UltimoMandatoLog)
	logger.Println("-------------------------------------------------")
}

func printLog(logger *log.Logger, log []EntradaRegistro) {
	logger.Printf("Log (len=%d):", len(log))
	for i, entrada := range log {
		logger.Printf("  [%d] %v", i, entrada)
	}
}

func printAlmacen(logger *log.Logger, almacen map[string]string) {
	// Almacen (claves limitadas para evitar saturación visual)
	logger.Println("------ Almacén (claves limitadas) ------")
	maxKeys := 5
	count := 0
	for k, v := range almacen {
		logger.Printf("  %s: %s", k, v)
		count++
		if count >= maxKeys {
			logger.Printf("  ...y %d más", len(almacen)-maxKeys)
			break
		}
	}

	logger.Println("==========================================")
}

// #endregion
