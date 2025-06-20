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
	"ms"
	"strconv"
	"sync"

	"github.com/DistributedClocks/GoVector/govec"
	"github.com/DistributedClocks/GoVector/govec/vclock"
)

type Request struct {
	logger_info []byte
	Pid         int
	task        Role
}

type Reply struct{}

type Update struct {
	new_text string
}

type Barrier struct{}

type Role struct {
	fname string
}

type Role_tuple struct {
	f1 Role
	f2 Role
}

const (
	total_processes = 4
)

// RASharedDB is a structure used in a distributed system, likely to manage
// mutual exclusion or synchronization across nodes. It implements message
// passing and uses concurrency mechanisms such as channels and mutexes to
// ensure thread-safe access to shared variables.
type RASharedDB struct {
	OutRepCnt  int               // Number of outstanding replies not yet received
	ReqCS      bool              // Indicates if the node is requesting access to the critical section (CS)
	RepDefd    []bool            // Array to track which replies have been deferred for later response
	msg        *ms.MessageSystem // Pointer to the system managing message passing between nodes
	done       chan bool         // Channel used for signaling when operations are done/completed
	chrep      chan bool         // Channel used for signaling when all replys are received
	Mutex      *sync.Mutex       // Mutex to protect concurrent access to shared variables (OurSeqNum, HigSeqNum, etc.)
	task       Role
	permission chan Reply          // Channel used for receiving permission messages
	requests   chan Request        // Channel used for receiving request messages
	mapFunct   map[Role_tuple]bool // Boolean matrix that asigns a bool value for a tuple of functions
	vctClock   *govec.GoLog        // Contains the value of the local vectorial clock
}

func distribute_received_messages(msgs *ms.MessageSystem, myFile string, request_channel chan Request, reply_channel chan Reply, barrier_chan chan Barrier) {
	for {
		received_message := msgs.Receive()
		switch message := received_message.(type) {
		case Request:
			request_channel <- message
		case Reply:
			reply_channel <- message
		case Update:
			mf.write_file(myFile, message.new_text)
		case Barrier:
			barrier_chan <- Barrier{}
		}
	}
}

func handle_received_permission(ra *RASharedDB) {
	for {
		<-ra.permission
		ra.OutRepCnt--
		if ra.OutRepCnt == 0 {
			ra.chrep <- true
		}
	}
}

func happens_before(local_vc vclock.VClock, remote_vc vclock.VClock, local_id int, remote_id int) bool {
	if local_vc.Compare(remote_vc, vclock.Descendant) {
		return true
	} else if local_vc.Compare(remote_vc, vclock.Concurrent) {
		return local_id < remote_id
	} else {
		return false
	}
}

func handle_received_request(ra *RASharedDB) {
	for {
		request := <-ra.requests
		messagePayload := []byte("shample-payload")
		ra.vctClock.UnpackReceive("Recibe request", request.logger_info, &messagePayload, govec.GetDefaultConfig())
		remote_vclock := vclock.FromBytes(request.logger_info)
		local_vclock := ra.vctClock.GetCurrentVC()
		ra.Mutex.Lock()
		should_defer := (ra.ReqCS != false) && happens_before(local_vclock, remote_vclock, ra.msg.me, request.Pid) && ra.mapFunct[Role_tuple{ra.task, request.task}]
		ra.Mutex.Unlock()
		if should_defer {
			ra.RepDefd[request.Pid-1] = true
		} else {
			ra.msg.Send(request.Pid, Reply{})
		}
	}
}

func New(me int, usersFile string, task Role) *RASharedDB {
	messageTypes := []ms.Message{Request{}, Reply{}, Update{}, Barrier{}}
	msgs := ms.New(me, usersFile, messageTypes)
	govector := govec.InitGoVector(strconv.Itoa(me), "log"+strconv.Itoa(me), govec.GetDefaultConfig())
	mapFunct := make(map[Role_tuple]bool)
	mapFunct[Role_tuple{Role{"read"}, Role{"read"}}] = false
	mapFunct[Role_tuple{Role{"write"}, Role{"read"}}] = true
	mapFunct[Role_tuple{Role{"read"}, Role{"write"}}] = true
	mapFunct[Role_tuple{Role{"write"}, Role{"write"}}] = true

	ra := RASharedDB{0, false, []bool{}, &msgs, make(chan bool), make(chan bool),
		&sync.Mutex{}, task, make(chan Reply), make(chan Request), mapFunct, govector}

	go handle_received_request(&ra)
	go handle_received_permission(&ra)
	return &ra
}

// Pre: Verdad
// Post: Realiza  el  PreProtocol  para el  algoritmo de
//
//	Ricart-Agrawala Generalizado
func (ra *RASharedDB) PreProtocol() {
	ra.Mutex.Lock()                    // Cogemos el mutex
	ra.ReqCS = true                    // Marcamos estado como intentando
	ra.Mutex.Unlock()                  // Devolvemos el mutex
	ra.OutRepCnt = total_processes - 1 // Indicamos cuantos permisos hay que recibir
	for i := 0; i < total_processes; i++ {
		if (i + 1) != ra.msg.me { // Para todos los procesos que no son el local
			logger_info := ra.vctClock.PrepareSend("Received request", govec.GetDefaultConfig()) // Actualizamos el reloj vectorial
			ra.msg.Send(i+1, Request{logger_info, ra.msg.me, ra.task})                           // Y enviamos el mensaje
		}
	}
	<-ra.chrep // Esperas la llegada de todos los permisos necesarios
}

// Pre: Verdad
// Post: Realiza  el  PostProtocol  para el  algoritmo de
//
//	Ricart-Agrawala Generalizado
func (ra *RASharedDB) PostProtocol() {
	ra.Mutex.Lock()
	ra.ReqCS = false
	ra.Mutex.Unlock()
	for i := 0; i < total_processes; i++ {
		if ra.RepDefd[i] {
			ra.RepDefd[i] = false
			ra.msg.Send(i+1, Reply{})
		}
	}
}

func (ra *RASharedDB) Stop() {
	ra.msg.Stop()
	ra.done <- true
}
