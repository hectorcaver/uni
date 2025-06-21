
// Jorge Gallardo y Enrique Baldovin

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
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
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
type Estado struct {
	ser string
}

// A medida que el nodo Raft conoce las operaciones de las  entradas de registro
// comprometidas, envía un AplicaOperacion, con cada una de ellas, al canal
// "canalAplicar" (funcion NuevoNodo) de la maquina de estados
type AplicaOperacion struct {
	Indice    int // en la entrada de registro
	Operacion TipoOperacion
}

type Entry struct {
	Indice    int
	Mandato   int
	Operacion TipoOperacion
}

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

	// Vuestros datos aqui.
	VotosRecibidos int
	Estado         Estado
	MandatoActual  int
	VotadoA        int

	//Parte Logs
	RegEntry []Entry
	//Parte Comun
	CommitIndex int //Indice de ultima entrada comprometida
	lastApplied int //Indice de la ultima operacion aplicada
	//Parte del LÍDER
	nextIndex  []int
	matchIndex []int

	// Canal para mandar los latidos
	Heartbeat chan bool
	//Canal para notificar que hemos sido votados como lider
	Elegido chan bool

	// Contador de confirmaciones
	HanAplicado int

	// ----- Machine y P4 -----
	//valorADevolver chan string
	aplicado        chan bool
	AplicaOperacion chan AplicaOperacion
	Almacen         map[string]string
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
func NuevoNodo(nodos []rpctimeout.HostPort, yo int, canalAplicarOperacion chan AplicaOperacion) *NodoRaft {
	nr := &NodoRaft{}
	nr.Nodos = nodos
	nr.Yo = yo
	nr.IdLider = -1

	if kEnableDebugLogs {
		nombreNodo := nodos[yo].Host() + "_" + nodos[yo].Port()
		logPrefix := fmt.Sprintf("%s", nombreNodo)

		//fmt.Println("LogPrefix: ", logPrefix)

		if kLogToStdout {
			nr.Logger = log.New(os.Stdout, nombreNodo+" -->> ",
				log.Lmicroseconds|log.Lshortfile)
		} else {
			err := os.MkdirAll(kLogOutputDir, os.ModePerm)
			if err != nil {
				panic(err.Error())
			}
			logOutputFile, err := os.OpenFile(fmt.Sprintf("%s/%s.txt",
				kLogOutputDir, logPrefix), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				panic(err.Error())
			}
			nr.Logger = log.New(logOutputFile,
				logPrefix+" -> ", log.Lmicroseconds|log.Lshortfile)
		}
		nr.Logger.Println("logger initialized")
	} else {
		nr.Logger = log.New(ioutil.Discard, "", 0)
	}

	// Añadir codigo de inicialización
	// código de inicialización
	nr.MandatoActual = 0
	nr.VotadoA = IntNOINICIALIZADO
	nr.Estado.ser = "FOLLOWER"
	nr.VotosRecibidos = 0
	nr.VotadoA = IntNOINICIALIZADO
	nr.Heartbeat = make(chan bool)
	nr.Elegido = make(chan bool)

	//cosas registros
	//Parte Logs
	nr.RegEntry = make([]Entry, 0)
	//Parte Comun
	nr.VotosRecibidos = 0
	nr.CommitIndex = IntNOINICIALIZADO
	nr.lastApplied = IntNOINICIALIZADO
	//Parte del LÍDER
	nr.nextIndex = make([]int, len(nr.Nodos))
	nr.matchIndex = make([]int, len(nr.Nodos))
	for i := range nr.matchIndex {
		nr.matchIndex[i] = -1
	}
	nr.aplicado = make(chan bool)
	nr.AplicaOperacion = make(chan AplicaOperacion)
	nr.Almacen = make(map[string]string)
	go nr.escucha()
	go nr.aplicarOp()
	return nr
}

// aplicarOp es un método de NodoRaft que escucha continuamente las operaciones
// en el canal AplicaOperacion y las aplica a la máquina de estados del nodo.
// Si la operación es una operación de "escribir", actualiza el mapa Almacen con el
// par clave-valor proporcionado. Si el nodo es el líder, señala que la
// operación ha sido aplicada enviando un valor al canal aplicado.
func (nr *NodoRaft) aplicarOp() {
	for {
		operacion := <-nr.AplicaOperacion
		nr.Logger.Printf("Aplicando operación en aplicarOp() : %v", operacion)
		if operacion.Operacion.Operacion == "escribir" {
			nr.Mux.Lock()
			nr.Almacen[operacion.Operacion.Clave] = operacion.Operacion.Valor
			nr.Mux.Unlock()
		}
		if nr.Estado.ser == "LEADER" {
			nr.aplicado <- true
		}
	}
}

// escucha es un método para NodoRaft que maneja el bucle principal de eventos para un nodo Raft.
// Escucha los latidos del corazón y gestiona las transiciones de estado entre SEGUIDOR, LÍDER y CANDIDATO.
// El método utiliza temporizadores para manejar los tiempos de espera y tareas periódicas como enviar latidos y comenzar elecciones.
func (nr *NodoRaft) escucha() {
	//rand.Seed(time.Now().UnixNano())
	nr.Logger.Println("Estamos escuchando", nr.MandatoActual, nr.Estado.ser)
	for {
		switch nr.Estado.ser {
		case "FOLLOWER":
			for nr.Estado.ser == "FOLLOWER" {
				TimerTimeoutHeartbeat := time.NewTimer(time.Duration(500+rand.Intn(400)) * time.Millisecond)
				//nr.Logger.Println("Nodo en estado FOLLOWER, reiniciando temporizador de heartbeat")
				select {
				case <-nr.Heartbeat:
					//nr.Logger.Println("Latido recibido como seguidor")
					TimerTimeoutHeartbeat.Reset(time.Duration(500 + rand.Intn(400)))
				case <-TimerTimeoutHeartbeat.C:
					//nr.Logger.Println("Temporizador de latidos expirado como seguidor, convirtiéndose en CANDIDATE")
					//TimerTimeoutHeartbeat.Stop()
					nr.VolverseCandidate()
				}
			}

		case "LEADER":
			for nr.Estado.ser == "LEADER" {
				//nr.Logger.Println("Nodo en estado LEADER, enviando latidos a seguidores")
				nr.enviarLatidos()
				Enviar_Heartbeat := time.NewTimer(50 * time.Millisecond)
				select {
				case <-nr.Heartbeat:
					nr.Logger.Println("Latido recibido como líder")
					Enviar_Heartbeat.Stop()
				case <-Enviar_Heartbeat.C:
					//nr.Logger.Println("Temporizador de envío de latidos expirado, verificando entradas")
					if nr.CommitIndex > nr.lastApplied {
						nr.Logger.Printf("%d Hay entradas sin aplicar CommitIndex: %d lastApplied: %d", nr.Yo, nr.CommitIndex, nr.lastApplied)
						nr.lastApplied++

						operacion := AplicaOperacion{
							nr.lastApplied,
							nr.RegEntry[nr.lastApplied].Operacion,
						}
						nr.AplicaOperacion <- operacion
					}
				}
			}

		case "CANDIDATE":
			for nr.Estado.ser == "CANDIDATE" {
				nr.Logger.Printf("%d Soy candidato y mi mandato es %d \n", nr.Yo, nr.MandatoActual)
				TiempoEntreEleccion := time.NewTimer(time.Duration(200+rand.Intn(400)) * time.Millisecond)
				//nr.Logger.Println("Nodo en estado CANDIDATE, iniciando elección")
				nr.iniciarEleccion()
				//nr.Logger.Println("Comienzan nuevas elecciones, reiniciando temporizador")
				select {
				case <-nr.Heartbeat:
					nr.Estado.ser = "FOLLOWER"
					TiempoEntreEleccion.Stop()
					nr.VotosRecibidos = 0
					nr.Logger.Println("Latido recibido como candidato")
				case <-TiempoEntreEleccion.C:
					nr.Logger.Println("Temporizador de elecciones expirado, intentando nuevamente como candidato")
					nr.VolverseCandidate()
				case <-nr.Elegido:
					TiempoEntreEleccion.Stop()
				}
			}

		}
	}
}

// VolverseCandidate transiciona el nodo actual al estado de Candidato.
// Incrementa el mandato actual, vota por sí mismo y reinicia el ID del líder.
// El estado del nodo se actualiza a "CANDIDATE" y registra la transición junto con el nuevo mandato.
func (nr *NodoRaft) VolverseCandidate() {
	nr.Estado.ser = "CANDIDATE"
	nr.IdLider = IntNOINICIALIZADO
	nr.VotadoA = nr.Yo
	nr.MandatoActual++
	nr.VotosRecibidos = 1 // cuenta su propio voto
	nr.Logger.Println("Convirtiéndose en CANDIDATE en mandato", nr.MandatoActual)
}

func (nr *NodoRaft) VolverseLeader() {
	nr.Logger.Println("Se Vuelve Líder en mandato", nr.MandatoActual)
	nr.Estado.ser = "LEADER"
	nr.IdLider = nr.Yo
	nr.VotosRecibidos = 0
	for i := range nr.nextIndex {
		nr.nextIndex[i] = 0
	}
	for i := range nr.matchIndex {
		nr.matchIndex[i] = -1
	}
	nr.Elegido <- true
}

func (nr *NodoRaft) VolverseFollower(mandato int) {
	//nr.Logger.Println("Se Vuelve Seguidor en mandato", mandato)
	nr.Estado.ser = "FOLLOWER"
	nr.MandatoActual = mandato
	//nr.VotadoA = IntNOINICIALIZADO
	nr.VotosRecibidos = 0
}

// enviarLatidos envía mensajes de latido a todos los nodos en el clúster.
// Itera sobre todos los nodos y envía una RPC AppendEntries a cada uno,
// excepto a sí mismo. La función se ejecuta concurrentemente para cada nodo.
//
// Si la RPC AppendEntries falla, simplemente retorna. Si la RPC tiene éxito,
// verifica el término del resultado. Si el término es mayor que el término actual,
// se convierte en seguidor. Si el término es igual y el nodo sigue siendo el líder,
// actualiza el nextIndex y matchIndex para el nodo.
//
// Si el matchIndex es mayor que el commitIndex y la mayoría de los nodos
// han aplicado la operación, incrementa el commitIndex y reinicia el contador
// de operaciones aplicadas.
func (nr *NodoRaft) enviarLatidos() {
	//nr.Logger.Println("Enviando latidos a todos los nodos")
	for i := range nr.Nodos {
		if i != nr.Yo {
			go func(i int) {
				index := nr.nextIndex[i] - 1
				args := ArgAppendEntries{Mandato: nr.MandatoActual, IdLider: nr.Yo, LeaderCommit: nr.CommitIndex, PrevLogIndex: index, PrevLogTerm: 0}
				if 0 <= index && index < len(nr.RegEntry) {
					args.PrevLogTerm = nr.RegEntry[index].Mandato
				}
				lenEntradas := len(nr.RegEntry) - index - 1
				if lenEntradas > 0 {
					args.Entries = make([]Entry, lenEntradas)
					for j := 0; j < lenEntradas; j++ {
						args.Entries[j] = nr.RegEntry[index+j+1]
					}
				}

				var results Results
				if err := nr.Nodos[i].CallTimeout("NodoRaft.AppendEntries", &args, &results, 500*time.Millisecond); err != nil {
					return
				} else {
					//nr.Mux.Lock()
					//defer nr.Mux.Unlock()
					if results.Mandato > nr.MandatoActual {
						nr.Logger.Printf("Índice: %d Me convierto en follower\n", nr.Yo)
						nr.VolverseFollower(results.Mandato)
						nr.Heartbeat <- true
					} else if results.Mandato == nr.MandatoActual && nr.Estado.ser == "LEADER" {
						if results.Concedido {
							nr.nextIndex[i] = args.PrevLogIndex + len(args.Entries) + 1
							nr.matchIndex[i] = args.PrevLogIndex + len(args.Entries)
							if nr.matchIndex[i] > nr.CommitIndex {
								nr.HanAplicado++
								nr.Logger.Printf("Nodo %d ha aplicado la operación en el índice %d", i, nr.matchIndex[i])

								if nr.HanAplicado >= len(nr.Nodos)/2 {
									nr.CommitIndex++
									nr.Logger.Printf("CommitIndex actualizado a %d", nr.CommitIndex)
									nr.HanAplicado = 0
								}
							}
						} else if nr.nextIndex[i] > 0 {
							nr.nextIndex[i]--
						}
					}
				}
			}(i)
		}
	}
}

// Metodo Para() utilizado cuando no se necesita mas al nodo
//
// Quizas interesante desactivar la salida de depuracion
// de este nodo
func (nr *NodoRaft) para() {
	nr.Logger.Println("Parando nodo")
	go func() { time.Sleep(5 * time.Millisecond); os.Exit(0) }()
}

// Devuelve "yo", mandato en curso y si este nodo cree ser lider
//
// Primer valor devuelto es el indice de este  nodo Raft el el conjunto de nodos
// la operacion si consigue comprometerse.
// El segundo valor es el mandato en curso
// El tercer valor es true si el nodo cree ser el lider
// Cuarto valor es el lider, es el indice del líder si no es él
func (nr *NodoRaft) obtenerEstado() (int, int, bool, int) {
	var yo int = nr.Yo //fuera de exclusion porque nuestra propia id no cambia
	var mandato int
	var esLider bool
	var idLider int

	//COMPLETADO
	nr.Mux.Lock()
	//defer nr.Mux.Unlock()
	esLider = (yo == nr.IdLider)
	idLider = nr.IdLider
	mandato = nr.MandatoActual
	nr.Mux.Unlock()

	return yo, mandato, esLider, idLider
}

// El servicio que utilice Raft (base de datos clave/valor, por ejemplo)
// Quiere buscar un acuerdo de posicion en registro para siguiente operacion
// solicitada por cliente.

// Si el nodo no es el lider, devolver falso
// Sino, comenzar la operacion de consenso sobre la operacion y devolver en
// cuanto se consiga
//
// No hay garantia que esta operacion consiga comprometerse en una entrada de
// de registro, dado que el lider puede fallar y la entrada ser reemplazada
// en el futuro.
// Primer valor devuelto es el indice del registro donde se va a colocar
// la operacion si consigue comprometerse.
// El segundo valor es el mandato en curso
// El tercer valor es true si el nodo cree ser el lider
// Cuarto valor es el lider, es el indice del líder si no es él
func (nr *NodoRaft) someterOperacion(operacion TipoOperacion) (int, int,
	bool, int, string) {
	//cosas registros

	nr.Mux.Lock()

	indice := len(nr.RegEntry) //donde vamos a poner la operacion
	mandato := nr.MandatoActual
	EsLider := nr.IdLider == nr.Yo
	idLider := nr.IdLider
	valorADevolver := ""

	if EsLider {
		nr.Logger.Printf("Operación a aplicar: %v", operacion)
		nr.RegEntry = append(nr.RegEntry, Entry{Operacion: operacion, Mandato: mandato})
		nr.Logger.Printf("Esperando a que la entrada se comprometa: %v", operacion)
		// Esperar a que la entrada se comprometa
		nr.Mux.Unlock()
		<-nr.aplicado
		if operacion.Operacion == "leer" {
			valorADevolver = nr.Almacen[operacion.Clave]
		} else {
			nr.Almacen[operacion.Clave] = operacion.Valor
			valorADevolver = "escrito"
		}

	} else {
		nr.Mux.Unlock()
	}
	return indice, mandato, EsLider, idLider, valorADevolver
}

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
	Mandato      int
	IdCandidato  int
	LastLogIndex int
	LastLogTerm  int
}

// Structura de ejemplo de respuesta de RPC PedirVoto,
//
// Recordar
// -----------
// Nombres de campos deben comenzar con letra mayuscula !
type RespuestaPeticionVoto struct {
	Mandato     int
	Garantizado bool
}

// iniciarEleccion inicia el proceso de elección para el nodo actual.
// Incrementa el mandato actual, transiciona el nodo al estado de candidato,
// y envía solicitudes de voto a todos los demás nodos en el clúster.
// Si el nodo recibe la mayoría de los votos, se convierte en líder.
// Si recibe un mandato mayor que su mandato actual, vuelve al estado de seguidor.
func (nr *NodoRaft) iniciarEleccion() {

	nr.Logger.Println("Iniciando elección en mandato", nr.MandatoActual)

	//Registros
	RegIndex := len(nr.RegEntry) - 1
	var RegTerm int
	if RegIndex >= 0 {
		RegTerm = nr.RegEntry[RegIndex].Mandato
	} else {
		RegTerm = 0
	}
	//----------------------------------------------------------------------
	args := ArgsPeticionVoto{Mandato: nr.MandatoActual, IdCandidato: nr.Yo, LastLogIndex: RegIndex, LastLogTerm: RegTerm}

	for i := range nr.Nodos {
		if i != nr.Yo {
			go func(i int) {
				var reply RespuestaPeticionVoto
				if nr.enviarPeticionVoto(i, &args, &reply) {
					nr.Mux.Lock()
					if reply.Garantizado && reply.Mandato == nr.MandatoActual && nr.Estado.ser == "CANDIDATE" { //recibimos un voto
						nr.VotosRecibidos++
						nr.Logger.Printf("Voto recibido de nodo %d, total votos: %d", i, nr.VotosRecibidos)
						if nr.VotosRecibidos > len(nr.Nodos)/2 { //el numero de votos es mayor de la media
							nr.VolverseLeader()
						}
					} else if reply.Mandato > nr.MandatoActual { //el candidato que te solicita esta adelantado
						nr.VolverseFollower(reply.Mandato)
						nr.Heartbeat <- true
					}
					nr.Mux.Unlock()
				}
			}(i)
		}
	}
}

// Metodo para RPC PedirVoto

// PedirVoto maneja una solicitud de voto de un candidato en el algoritmo de consenso Raft.
// Verifica si el mandato y el registro del candidato están actualizados y concede o niega el voto en consecuencia.
//
// Parámetros:
// - peticion: Un puntero a ArgsPeticionVoto que contiene el mandato del candidato, el índice del último registro y el término del último registro.
// - reply: Un puntero a RespuestaPeticionVoto donde se almacenará la respuesta.
//
// Retorna:
// - error: Un error si la operación falla.
//
// La función bloquea el mutex del nodo para asegurar la seguridad de los hilos. Luego evalúa si el registro del candidato
// está al menos tan actualizado como el registro del nodo actual. Si el mandato del candidato es menor que el mandato actual
// y el registro no es mejor, se niega el voto. Si el mandato del candidato es mayor o igual y el nodo no ha votado por otro
// candidato en el mandato actual, se concede el voto. El nodo también actualiza su estado a seguidor si el mandato del candidato es mayor.
func (nr *NodoRaft) PedirVoto(peticion *ArgsPeticionVoto, reply *RespuestaPeticionVoto) error {
	nr.Mux.Lock()
	reply.Garantizado = false
	EsMejor := (peticion.LastLogTerm == nr.MandatoActual && peticion.LastLogIndex >= len(nr.RegEntry)-1 || peticion.LastLogTerm > nr.MandatoActual) //CONDICION MEJOR LIDER
	if peticion.Mandato < nr.MandatoActual && !EsMejor {
		// Denegar el voto si la petición es de un mandato menor
		reply.Mandato = nr.MandatoActual
		reply.Garantizado = false
		nr.Logger.Printf("No doy voto a %d", peticion.IdCandidato)
	} else if peticion.Mandato > nr.MandatoActual || (peticion.Mandato == nr.MandatoActual && (nr.VotadoA == IntNOINICIALIZADO || nr.VotadoA == peticion.IdCandidato)) {

		nr.VolverseFollower(peticion.Mandato)
		//nr.IdLider = peticion.IdCandidato //nuevo lider
		reply.Garantizado = true
		reply.Mandato = peticion.Mandato
		nr.VotadoA = peticion.IdCandidato
		nr.Heartbeat <- true
	}

	//reply.Mandato = nr.MandatoActual
	nr.Mux.Unlock()
	nr.Logger.Printf("Me pide %d, voto %d", peticion.IdCandidato, nr.VotadoA)
	return nil
}

type ArgAppendEntries struct {
	Mandato      int
	IdLider      int
	PrevLogIndex int
	PrevLogTerm  int
	Entries      []Entry
	LeaderCommit int
}

type Results struct {
	Mandato   int
	Concedido bool
}

// AppendEntries maneja la llamada RPC AppendEntries del líder para replicar entradas de registro y gestionar latidos.
// Actualiza el término y las entradas de registro del nodo actual en función del término del líder y la consistencia del registro.
//
// Parámetros:
// - args: Un puntero a ArgAppendEntries que contiene el término del líder, el índice y término del registro anterior, las entradas a añadir y el índice de compromiso del líder.
// - results: Un puntero a Results para almacenar el resultado de la llamada AppendEntries.
//
// Retorna:
// - error: Un error si la operación falla.
//
// La función realiza los siguientes pasos:
// 1. Actualiza el término actual si el término del líder es mayor.
// 2. Verifica la consistencia de las entradas del registro.
// 3. Trunca el registro si hay un conflicto.
// 4. Añade nuevas entradas del líder.
// 5. Actualiza el índice de compromiso si el índice de compromiso del líder es mayor.
// 6. Aplica cualquier nueva entrada comprometida a la máquina de estados.
func (nr *NodoRaft) AppendEntries(args *ArgAppendEntries, results *Results) error {

	results.Concedido = false
	results.Mandato = nr.MandatoActual

	if nr.MandatoActual < args.Mandato {
		nr.MandatoActual = args.Mandato
		nr.VotadoA = -1
		nr.VolverseFollower(args.Mandato)
		nr.Heartbeat <- true
	}

	if nr.MandatoActual == args.Mandato {

		// Verificamos la consistencia del log
		if args.PrevLogIndex < len(nr.RegEntry) {
			// Validamos que el término en `PrevLogIndex` sea consistente
			if args.PrevLogIndex == -1 || nr.RegEntry[args.PrevLogIndex].Mandato == args.PrevLogTerm {
				nr.VolverseFollower(args.Mandato)
				nr.IdLider = args.IdLider
				results.Mandato = nr.MandatoActual

				// Notificamos que se ha recibido un latido del líder
				//nr.Mux.Unlock()
				nr.Heartbeat <- true
				//nr.Mux.Lock()
				// Si el término es consistente o el índice es -1 (log vacío), añadimos las nuevas entradas

				// Truncamos el log en caso de conflicto
				if args.PrevLogIndex >= 0 {
					nr.RegEntry = nr.RegEntry[:args.PrevLogIndex+1] //Guardamos el log hasta el indice previo
				}
				// Añadimos las nuevas entradas del líder, a partir del índice previo
				nr.RegEntry = append(nr.RegEntry, args.Entries...)
				results.Concedido = true

				// Actualizamos el índice comprometido si el líder ha comprometido más
				if args.LeaderCommit > nr.CommitIndex {
					nr.CommitIndex = min(args.LeaderCommit, len(nr.RegEntry)-1) //para evitar desbordamiento
					if nr.CommitIndex > nr.lastApplied {
						nr.Logger.Printf("%d  Hay entradas sin aplicar CommitIndex: %d lastApplied: %d", nr.Yo, nr.CommitIndex, nr.lastApplied)
						nr.lastApplied++

						operacion := AplicaOperacion{
							nr.lastApplied,
							nr.RegEntry[nr.lastApplied].Operacion,
						}
						nr.AplicaOperacion <- operacion
					}

				}
			}
		}
	}

	return nil
}

// Funcion auxiliar para mínimo
func min(a int, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
}

// ----- Metodos/Funciones a utilizar como clientes
//
//

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
// de la llamada (incluido si son punteros
//
// Si en la llamada RPC, la respuesta llega en un intervalo de tiempo,
// la funcion devuelve true, sino devuelve false
//
// la llamada RPC deberia tener un timout adecuado.
//
// Un resultado falso podria ser causado por una replica caida,
// un servidor vivo que no es alcanzable (por problemas de red ?),
// una petición perdida, o una respuesta perdida
//
// Para problemas con funcionamiento de RPC, comprobar que la primera letra
// del nombre  todo los campos de la estructura (y sus subestructuras)
// pasadas como parametros en las llamadas RPC es una mayuscula,
// Y que la estructura de recuperacion de resultado sea un puntero a estructura
// y no la estructura misma.

func (nr *NodoRaft) enviarPeticionVoto(nodo int, args *ArgsPeticionVoto, reply *RespuestaPeticionVoto) bool {
	nr.Logger.Printf("Petición de voto enviada a nodo %d, con los args: %+v", nodo, args)
	success := nr.Nodos[nodo].CallTimeout("NodoRaft.PedirVoto", args, reply, time.Duration(500*time.Millisecond)) == nil
	nr.Logger.Printf("Resultado de la petición de voto a nodo %d: %v", nodo, success)
	return success

}

// ResultadoEstadoNodo contiene información sobre el estado del nodo Raft.
type ResultadoEstadoNodo struct {
	NumOperaciones int // Número de operaciones sometidas
	Mandato        int // Mandato actual del nodo
	IdLider        int // ID del líder
}

// ObtenerEstadoYOperaciones devuelve información sobre el número de operaciones, líder y mandato.
func (nr *NodoRaft) ObtenerEstadoYOperaciones(args Vacio, reply *ResultadoEstadoNodo) error {
	nr.Mux.Lock() // Aseguramos exclusión mutua
	defer nr.Mux.Unlock()

	reply.NumOperaciones = len(nr.RegEntry)
	reply.Mandato = nr.MandatoActual
	reply.IdLider = nr.IdLider
	return nil
}

// para que la use el cliente
type ResPruebas struct {
	NumOperaciones int // Número de operaciones dentro del log
	Commit         int // commitindex
	Mandatocommit  int //mandato del ultimo log commited
}

func (nr *NodoRaft) EstadoPruebas(args Vacio, reply *ResPruebas) error {
	nr.Mux.Lock() // Aseguramos exclusión mutua
	defer nr.Mux.Unlock()

	reply.NumOperaciones = len(nr.RegEntry)
	reply.Commit = nr.CommitIndex
	reply.Mandatocommit = nr.RegEntry[nr.CommitIndex].Mandato
	return nil
}
