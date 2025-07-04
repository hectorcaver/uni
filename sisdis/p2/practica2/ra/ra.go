/*
* AUTOR: Rafael Tolosana Calasanz
* ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
*			Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
* FECHA: septiembre de 2021
* FICHERO: ricart-agrawala.go
* DESCRIPCIÓN: Implementación del algoritmo de Ricart-Agrawala Generalizado en Go
* AUTOR: Héctor Lacueva Sacristán 869637
* AUTOR: Adrián Nasarre Sánchez 869561
 */
package ra

import (
	"practica2/mf"
	"practica2/ms"
	"strconv"
	"sync"

	"github.com/DistributedClocks/GoVector/govec"
	"github.com/DistributedClocks/GoVector/govec/vclock"
)

type Request struct {
	Logger_info []byte
	Pid         int
	Task        Role
}

type Reply struct{}

type Update struct {
	New_text string
}

type Barrier struct{}

type Role struct {
	Fname string
}

type Role_tuple struct {
	F1 Role
	F2 Role
}

type RASharedDB struct {
	OutRepCnt  int
	ReqCS      bool
	RepDefd    []bool
	Ms        *ms.MessageSystem
	Done       chan bool
	Chrep      chan bool
	Mutex      *sync.Mutex
	Task       Role
	Permission chan Reply
	Requests   chan Request
	MapFunct   map[Role_tuple]bool
	VctClock   *govec.GoLog
	Me 		   int
	TotalProcesos int
}

func (ra *RASharedDB)Distribute_received_messages(myFile string, barrier_chan chan Barrier) {

	doneRequest := make(chan bool)

	donePermission := make(chan bool)

	// * Lanzo la gorutina que gestiona la recepción de peticiones
	go handle_received_request(ra, doneRequest)

	// * Lanzo la gorutina que gestiona la recepción de permisos
	go handle_received_permission(ra, donePermission)
	

	for {

		select {
		case <-ra.Done:
			doneRequest <- true
			donePermission <- true
			return
		default:
			received_message := ra.Ms.Receive()
			//fmt.Println("[Distribute_received_messages] Mensaje recibido:", received_message)
			switch message := received_message.(type) {
			case Request:
				//fmt.Println("[Distribute_received_messages] Request recibido de PID:", message.Pid)
				ra.Requests <- message
			case Reply:
				//fmt.Println("[Distribute_received_messages] Reply recibido")
				ra.Permission <- message
			case Update:
				//fmt.Println("[Distribute_received_messages] Update recibido. Escribiendo nuevo texto en archivo.")
				mf.Write_file(myFile, message.New_text)
			case Barrier:
				//fmt.Println("[Distribute_received_messages] Barrier recibido")
				barrier_chan <- Barrier{}
			default:
				//fmt.Println("[Distribute_received_messages] Tipo de mensaje desconocido")
			}
		}	
	}
}

func handle_received_permission(ra *RASharedDB, done chan bool) {
	for {
		select {
		case <-done:
			return
		case <-ra.Permission:
			ra.OutRepCnt--
			//fmt.Println("[handle_received_permission] Permiso recibido. Contador restante:", ra.OutRepCnt)
			if ra.OutRepCnt == 0 {
				//fmt.Println("[handle_received_permission] Todos los permisos recibidos.")
				ra.Chrep <- true
			}
		}
		
	}
}

func happens_before(local_vc vclock.VClock, remote_vc vclock.VClock, local_id int, remote_id int) bool {
	if local_vc.Compare(remote_vc, vclock.Descendant) {
		//fmt.Println("[happens_before] Local VC es descendiente del remoto")
		return true
	} else if local_vc.Compare(remote_vc, vclock.Concurrent) {
		//fmt.Println("[happens_before] Local VC es concurrente con remoto")
		return local_id < remote_id
	} else {
		//fmt.Println("[happens_before] Local VC NO es descendiente ni concurrente con remoto")
		return false
	}
}

func handle_received_request(ra *RASharedDB, done chan bool) {
	for {

		select {
		case <-done:
			return
		case request := <-ra.Requests:

			//fmt.Println("Recibida request del nodo:", request.Pid)

			var messagePayload []byte
			opts := govec.GetDefaultLogOptions()

			ra.Mutex.Lock()
			// * Desempaqueto 
			ra.VctClock.UnpackReceive(
				"Recibe request del nodo " + strconv.Itoa(request.Pid), 
				request.Logger_info, &messagePayload, 
				opts,
			)

			local_vclock := ra.VctClock.GetCurrentVC()
			ra.Mutex.Unlock()

			//fmt.Println("[handle_received_request]: UnpackedPayload: ", messagePayload)

			remote_vclock, _ := vclock.FromBytes(request.Logger_info)

			ra.Mutex.Lock()
			should_defer := (ra.ReqCS) &&
				happens_before(local_vclock, remote_vclock, ra.Me, request.Pid) &&
				ra.MapFunct[Role_tuple{ra.Task, request.Task}]
			ra.Mutex.Unlock()

			if should_defer {
				//fmt.Println("[handle_received_request] Se difiere la respuesta a PID:", request.Pid)
				ra.RepDefd[request.Pid] = true
			} else {
				//fmt.Println("[handle_received_request] Se envía Reply a PID:", request.Pid)
				ra.Ms.Send(request.Pid, Reply{})
			}
		}
	}
}

func New(me int, usersFile string, task Role, barrierChan chan Barrier, 
		myFile string) *RASharedDB {

	//fmt.Println("[New] Inicializando RASharedDB para proceso:", me)
	// * Defino los tipos de mensajes que puede haber
	messageTypes := []ms.Message{Request{}, Reply{}, Update{}, Barrier{}}
	// * Inicializo el serivicio de mensajes
	ms, totalProcesos := ms.New(me, usersFile, messageTypes)
	// * Inicializo el vector de logs
	govector := govec.InitGoVector(strconv.Itoa(me), "log"+strconv.Itoa(me), govec.GetDefaultConfig())
	
	// * Inicializo la matriz de exclusión
	mapFunct := make(map[Role_tuple]bool)
	mapFunct[Role_tuple{Role{"read"}, Role{"read"}}] = false
	mapFunct[Role_tuple{Role{"write"}, Role{"read"}}] = true
	mapFunct[Role_tuple{Role{"read"}, Role{"write"}}] = true
	mapFunct[Role_tuple{Role{"write"}, Role{"write"}}] = true

	// * Creo un objeto de tipo Ra
	ra := RASharedDB{
		OutRepCnt: 0,
		ReqCS: false,
		RepDefd: make([]bool, totalProcesos),
		Ms: &ms,
		Done: make(chan bool),
		Chrep: make(chan bool),
		Mutex: &sync.Mutex{},
		Task: task,
		Permission: make(chan Reply),
		Requests: make(chan Request),
		MapFunct: mapFunct,
		VctClock:  govector,
		Me: me,
		TotalProcesos: totalProcesos, 
	}

	go ra.Distribute_received_messages(myFile, barrierChan)

	return &ra
}

func (ra *RASharedDB) PreProtocol() {
	//fmt.Println("[PreProtocol] Inicio")
	ra.Mutex.Lock()
	ra.ReqCS = true
	ra.OutRepCnt = ra.TotalProcesos - 1
	ra.Mutex.Unlock()
	
	//fmt.Println("[PreProtocol] OutRepCnt seteado a", ra.OutRepCnt)

	payLoad :=  []byte("sample-payload")

	for i := 0; i < ra.TotalProcesos; i++ {
		if i != ra.Me {
			
			opts := govec.GetDefaultLogOptions()

			ra.Mutex.Lock()
			logger_info := ra.VctClock.PrepareSend(
				"Enviar request al nodo: " + strconv.Itoa(i) , 
				payLoad, 
				opts,
			)
			ra.Mutex.Unlock()

			request := Request{logger_info, ra.Me, ra.Task}

			//fmt.Println("[PreProtocol] Enviando request", request, "al nodo:", i)

			ra.Ms.Send(i, request)
		}
	}

	<-ra.Chrep
	//fmt.Println("[PreProtocol] Todos los permisos recibidos, continuando")
}

func (ra *RASharedDB) PostProtocol() {
	//fmt.Println("[PostProtocol] Inicio")
	ra.Mutex.Lock()
	ra.ReqCS = false
	ra.Mutex.Unlock()

	for i := 0; i < ra.TotalProcesos; i++ {
		if i != ra.Me {
			if ra.RepDefd[i] {
				//fmt.Printf("[PostProtocol] Enviando Reply al proceso diferido %d\n", i)
				ra.RepDefd[i] = false
				ra.Ms.Send(i, Reply{})
			}
		}
	}
	//fmt.Println("[PostProtocol] Fin")
}

func (ra *RASharedDB) Stop() {
	//fmt.Println("[Stop] Deteniendo sistema de mensajes")
	ra.Ms.Stop()
	ra.Done <- true
}
