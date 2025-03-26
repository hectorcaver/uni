// Escribir vuestro código de funcionalidad Raft en este fichero
//

package raft

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

// A medida que el nodo Raft conoce las operaciones de las entradas de registro
// comprometidas, envía un AplicaOperacion, con cada una de ellas, al canal
// "canalAplicar" (funcion NuevoNodo) de la maquina de estados
type AplicaOperacion struct {
	Indice    int // en la entrada de registro
	Operacion TipoOperacion
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

	// Estado persistente en todos los servidores
	MandatoActual int               // valor del mandato actual del nodo
	MiVoto        int               // índice del nodo votado de la lista de nodos
	Log           []EntradaRegistro // Traza de comandos ejecutados en orden
	Estado        Estado            // estado actual del servidor

	// Para notificar cosas
	Latido chan bool // Canal que informa de un latido

}

type EntradaRegistro struct {
	Indice  int
	Mandato int
	Comando TipoOperacion
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
		fmt.Println("nombreNodo: ", nombreNodo)

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
	nr.MandatoActual = 0 // Inicializamos a 0 el término
	nr.MiVoto = -1       // El voto lo inicializamos en -1
	nr.Estado = Seguidor // Al inicializar un nuevo nodo este será seguidor.
	nr.Latido = make(chan bool)

	// Hay que lanzar gorutina para manejo del Nodo.
	go nr.tratarNodo()

	return nr
}

// Metodo Para() utilizado cuando no se necesita mas al nodo
//
// Quizas interesante desactivar la salida de depuracion
// de este nodo
func (nr *NodoRaft) para() {
	go func() { time.Sleep(5 * time.Millisecond); os.Exit(0) }()
}

// Devuelve "yo", mandato en curso y si este nodo cree ser lider
//
// Primer valor devuelto es el indice de este nodo Raft en el conjunto de nodos
// la operacion si consigue comprometerse.
// El segundo valor es el mandato en curso
// El tercer valor es true si el nodo cree ser el lider
// Cuarto valor es el lider, es el indice del líder si no es él
func (nr *NodoRaft) obtenerEstado() (int, int, bool, int) {
	var yo int = nr.Yo
	var mandato = nr.MandatoActual
	var esLider = nr.IdLider == nr.Yo
	var idLider int = nr.IdLider

	return yo, mandato, esLider, idLider
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
//   - Primer valor devuelto es el indice del registro donde se va a colocar
//   - la operacion si consigue comprometerse.
//   - El segundo valor es el mandato en curso
//   - El tercer valor es true si el nodo cree ser el lider
//   - Cuarto valor es el lider, es el indice del líder si no es él
//   - Quinto valor es el resultado de aplicar esta operación en máquina de
//     estados
func (nr *NodoRaft) someterOperacion(operacion TipoOperacion) (int, int,
	bool, int, string) {
	indice := -1
	mandato := -1
	EsLider := false
	idLider := -1
	valorADevolver := ""

	fmt.Println(operacion)

	// VUESTRO CODIGO AQUI
	_, mandato, EsLider, idLider = nr.obtenerEstado()

	if EsLider {
		// Iniciar consenso de la operación
		//go consensuarOperacion(canal de fin)
		//<- fin
		//valorADevolver = "nuevoValor" // ¿?
	}

	return indice, mandato, EsLider, idLider, valorADevolver
}

// -----------------------------------------------------------------------
// LLAMADAS RPC al API
//
// Si no tenemos argumentos o respuesta estructura vacia (tamaño cero)
type Vacio struct{}

func (nr *NodoRaft) ParaNodo(args Vacio, reply *Vacio) error {
	nr.Logger.Printf("Nodo: %d. PARADO\n", nr.Yo)
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

func (nr *NodoRaft) ObtenerEstadoNodo(args Vacio, reply *EstadoRemoto) error {
	reply.IdNodo, reply.Mandato, reply.EsLider, reply.IdLider = nr.obtenerEstado()
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
	Mandato          int // Mandato del nodo candidato
	IdCandidato      int // Identificador del candidato
	UltimoIndiceLog  int // Último índice del log del candidato
	UltimoMandatoLog int // Último mandato del log del candidato
}

// Structura de ejemplo de respuesta de RPC PedirVoto,
//
// Recordar
// -----------
// Nombres de campos deben comenzar con letra mayuscula !
type RespuestaPeticionVoto struct {
	// Vuestros datos aqui
	Mandato  int  // Mandato actual para que el candidato se actualice
	VotoDado bool // true -> recibo voto |||| false -> NO recibo voto
}

// Metodo para RPC PedirVoto
func (nr *NodoRaft) PedirVoto(peticion *ArgsPeticionVoto,
	reply *RespuestaPeticionVoto) error {
	// Vuestro codigo aqui
	nr.Mux.Lock()
	defer nr.Mux.Unlock()

	if peticion.Mandato < nr.MandatoActual {
		reply.VotoDado = false
		reply.Mandato = nr.MandatoActual
		return nil
	}

	cond1 := peticion.Mandato > nr.MandatoActual // ¿ Esto aquí ?

	cond2 := (peticion.UltimoIndiceLog >= len(nr.Log))
	cond3 := (nr.MiVoto == -1 || nr.MiVoto == peticion.IdCandidato)
	if (cond3 && cond2) || cond1 {
		nr.MandatoActual = peticion.Mandato
		nr.MiVoto = peticion.IdCandidato
		reply.VotoDado = true
		reply.Mandato = nr.MandatoActual
		return nil
	}

	reply.VotoDado = false
	reply.Mandato = nr.MandatoActual
	return nil
}

type ArgAppendEntries struct {
	// Vuestros datos aqui
	Mandato          int               // Mandato del lider
	IdLider          int               // Para que los seguidores redirijan a los clientes
	IndiceLogPrevio  int               // Índice del log que precede a los nuevos
	MandatoLogPrevio int               // Mandato del IndiceLogPrevio
	Entradas         []EntradaRegistro // Entradas del log para guardar
	CommitLider      int               // Índice de operaciones "commiteadas"
}

type Results struct {
	// Vuestros datos aqui
	Mandato int  // Término actual, para que el lider se actualice
	Exito   bool // True -> Si follower al día, False -> en caso contrario
}

// Metodo de tratamiento de llamadas RPC AppendEntries
func (nr *NodoRaft) AppendEntries(args *ArgAppendEntries,
	results *Results) error {
	// Completar....
	nr.Mux.Lock()
	defer nr.Mux.Unlock()

	results.Mandato = nr.MandatoActual

	// Si se detecta algún fallo se devuelve false
	cond1 := args.Mandato < nr.MandatoActual
	// Si aún no se han comprometido operaciones el índice es -1 (problemas rango)
	if args.IndiceLogPrevio > 0 {
		cond2 := nr.Log[args.IndiceLogPrevio].Mandato != args.MandatoLogPrevio
		if cond1 || cond2 {
			results.Exito = false
			return nil
		}
	}

	var notInLog []EntradaRegistro

	// Busco inconsistencias en el Log, si hay devuelve false y acaba
	for indice, entrada := range args.Entradas {
		if entrada.Indice == indice {
			if entrada.Mandato != nr.Log[indice].Mandato {
				results.Exito = false
				return nil
			}
		} else {
			notInLog = append(notInLog, entrada)
		}
	}

	// Se añaden las nuevas entradas al Log
	nr.Logger.Printf("Nodo: %d. Antes de copy\n", nr.Yo)
	copy(nr.Log, notInLog)
	nr.Logger.Printf("Nodo: %d. Después de copy\n", nr.Yo)

	// Si no se han aplicado todas las operaciones, se hacen las que se pueda
	/*if CommitLider > CommitIndice {
		if CommitLider > len(Log) {
			CommitIndice = len(Log)
		} else {
			CommitIndice = CommitLider
		}
	}*/

	nr.Logger.Printf("Nodo: %d. Recibido latido\n", nr.Yo)
	nr.Latido <- true
	results.Exito = true
	return nil
}

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

	fmt.Println(nodo, args, reply)

	Nodo := nr.Nodos[nodo]

	// Pedimos el voto al nodo
	err := Nodo.CallTimeout("NodoRaft.PedirVoto", args,
		reply, 300*time.Millisecond)

	// Habrá que hacer algo con el reply
	if reply.Mandato > nr.MandatoActual {
		nr.Estado = Seguidor
		nr.MandatoActual = reply.Mandato
	}

	// Si se produce algún error durante la petición devuelve false
	// Sino, devuelve true
	return (err == nil)
}

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
func (nr *NodoRaft) enviarLatido(nodo int, args *ArgAppendEntries,
	reply *Results) bool {

	fmt.Println(nodo, args, reply)

	Nodo := nr.Nodos[nodo]

	// Pedimos el voto al nodo
	err := Nodo.CallTimeout("NodoRaft.AppendEntries", args,
		reply, 300*time.Millisecond)

	// Habrá que hacer algo con el reply
	/*if reply.Mandato > nr.MandatoActual {
		nr.Estado = Seguidor
		nr.MandatoActual = reply.Mandato
	}*/

	// Si se produce algún error durante la petición devuelve false
	// Sino, devuelve true
	return (err == nil)
}

func (nr *NodoRaft) iniciarEleccion(victoria chan bool) {
	// Aumentamos el mandato y nos votamos a nosotros mismos
	nr.MandatoActual++
	nr.MiVoto = nr.Yo
	votosRecibidos := 1

	// Preparamos los argumentos de la petición
	args := ArgsPeticionVoto{nr.MandatoActual, nr.Yo, 0, 0}

	// Establecemos la mayoría
	mayoria := len(nr.Nodos)/2 + 1

	fmt.Printf("Votos para mayoría: %d\n", mayoria)

	var obtenida = false

	// Enviamos petición a todos los nodos
	for i := 0; i < len(nr.Nodos); i++ {
		if i == nr.Yo {
			continue
		}
		go func(nodoID int) {
			var reply RespuestaPeticionVoto
			// Si nos votan
			if nr.enviarPeticionVoto(nodoID, &args, &reply) && reply.VotoDado {
				nr.Mux.Lock()
				// Aumentamos los votos recibidos
				votosRecibidos++
				// Si tenemos mayoría vencemos (solo la primera vez)
				if votosRecibidos >= mayoria && !obtenida {
					victoria <- true
					obtenida = true
				}
				nr.Mux.Unlock()
			}
		}(i)
	}
}

func (nr *NodoRaft) enviarLatidos() {

	// Hay que tener cuidado con el acceso a los vectores
	// Posible acceso con -1 como indice del vector
	var indiceActual = len(nr.Log)
	var mandatoPrevio = 0
	if indiceActual != 0 {
		mandatoPrevio = nr.Log[indiceActual-1].Mandato
	}

	// Creamos los argumentos del latido
	args := ArgAppendEntries{
		Mandato:          nr.MandatoActual,
		IdLider:          nr.Yo,
		IndiceLogPrevio:  indiceActual - 1,
		MandatoLogPrevio: mandatoPrevio,
		Entradas:         []EntradaRegistro{},
		CommitLider:      0,
	}

	// Para cada nodo distinto del propio
	for i := 0; i < len(nr.Nodos); i++ {
		if i == nr.Yo {
			continue
		}
		go func(nodoID int) {
			var results Results
			// Enviamos los latidos y recogemos los resultados
			if nr.enviarLatido(nodoID, &args, &results) {
				fmt.Printf("Nodo %d, ha respondido al latido\n"+
					" - Número de mandato: %d\n"+
					" - Éxito: %t\n",
					nodoID, results.Mandato, results.Exito)
			} else {
				fmt.Printf("ERROR al enviar latido al Nodo %d\n",
					nodoID)
			}
			/*if nr.enviarLatido(nodoID, &args, &results) && !results.Exito {
				nr.Mux.Lock()

				nr.Mux.Unlock()
			}*/
		}(i)
	}
}

func (nr *NodoRaft) tratarNodo() {
	for {
		switch nr.Estado {
		case Seguidor:
			nr.Logger.Printf("Nodo: %d. Estado: SEGUIDOR\n", nr.Yo)
			// establecer un tiempo de time-out aleatorio
			nr.tratarSeguidor()

		case Candidato:
			nr.Logger.Printf("Nodo: %d. Estado: CANDIDATO\n", nr.Yo)
			// Gestiona el proceso de elección
			nr.tratarCandidato()

		case Lider:
			nr.Logger.Printf("Nodo: %d. Estado: LIDER\n", nr.Yo)
			// Manda latidos a los procesos
			nr.tratarLider()

		}
	}
}

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

func (nr *NodoRaft) tratarSeguidor() {
	// Establece el tiempo de time-out (si vence crea elección)
	///// ALERTA 1
	timeOut := time.NewTimer(getRandomMilliseconds(1000, 5000))

	select {
	case <-timeOut.C:
		nr.Logger.Printf("Nodo: %d. Vence timeout, ELECCIÓN\n", nr.Yo)
		nr.Estado = Candidato
	case <-nr.Latido:
		nr.Logger.Printf("Nodo: %d. Latido RECIBIDO\n", nr.Yo)
		nr.MiVoto = -1
	}
}

func (nr *NodoRaft) tratarCandidato() {
	// Inicio el timeOut de elección
	///// ALERTA 2
	timeOut := time.NewTimer(getRandomMilliseconds(1000, 3000))
	victoria := make(chan bool)
	go nr.iniciarEleccion(victoria)

	select {
	case <-timeOut.C:
		nr.Logger.Printf("Nodo: %d. Vence timeout, ELECCIÓN\n", nr.Yo)

	case <-nr.Latido:
		nr.Logger.Printf("Nodo: %d. Latido RECIBIDO\n", nr.Yo)
		nr.Estado = Seguidor
		nr.MiVoto = -1

	case <-victoria:
		nr.Logger.Printf("Nodo: %d. Elección GANADA\n", nr.Yo)
		nr.Estado = Lider
		nr.IdLider = nr.Yo
	}
}

func (nr *NodoRaft) tratarLider() {
	// Establece el tiempo para mandar un latido
	///// ALERTA 3
	timeOut := time.NewTimer(50 * time.Millisecond)
	// Envía los latidos
	nr.enviarLatidos()
	select {
	case <-timeOut.C:
		nr.Logger.Printf("Nodo: %d. Vence timeout, Mandar Latido\n", nr.Yo)
	case <-nr.Latido:
		nr.Logger.Printf("Nodo: %d. Latido RECIBIDO\n", nr.Yo)
		nr.Estado = Seguidor
		nr.MiVoto = -1
	}
}
