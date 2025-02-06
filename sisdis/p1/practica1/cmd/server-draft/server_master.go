/*
* AUTOR: Rafael Tolosana Calasanz y Unai Arronategui
* ASIGNATURA: 30221 Sistemas Distribuidos del Grado en Ingeniería Informática
*			Escuela de Ingeniería y Arquitectura - Universidad de Zaragoza
* FECHA: septiembre de 2022
* FICHERO: server-draft.go
* DESCRIPCIÓN: contiene la funcionalidad esencial para realizar los servidores
*				correspondientes a la práctica 1
 */
package main

import (
	"encoding/gob"
	"log"
	"os"
	"net"
	"practica1/com"
	"strconv"
	"bufio"
	//"math/rand"
)

// PRE: verdad = !foundDivisor
// POST: IsPrime devuelve verdad si n es primo y falso en caso contrario
func isPrime(n int) (foundDivisor bool) {
	foundDivisor = false
	for i := 2; (i < n) && !foundDivisor; i++ {
		foundDivisor = (n%i == 0)
	}
	return !foundDivisor
}

// PRE: interval.A < interval.B
// POST: FindPrimes devuelve todos los números primos comprendidos en el
//
//	intervalo [interval.A, interval.B]
func findPrimes(interval com.TPInterval) (primes []int) {
	for i := interval.Min; i <= interval.Max; i++ {
		if isPrime(i) {
			primes = append(primes, i)
		}
	}
	return primes
}

func processRequests(
	id int, 
	conn_chan chan(net.Conn),
	workers []string,
){
	// Creo un canal para la recepción de la respuesta del worker
	reply_chan := make(chan(com.Reply))
	for{
		conn := <- conn_chan
		log.Println("GORUTINE " + strconv.Itoa(id) + ": accepted new connection")

		// Recojo una petición nueva
		var request com.Request
		decoder := gob.NewDecoder(conn)
		err := decoder.Decode(&request)
		com.CheckError(err)

		// Selecciono un worker aleatorio
		worker := getRandomWorker(workers)

		// Le envío una petición al worker y un canal para recibir la respuesta
		go sendTaskToWorker(request, worker, reply_chan)

		// Recibe la respuesta del worker
		reply := <- reply_chan

		// Envía la respuesta de la petición al cliente
		encoder := gob.NewEncoder(conn)
		encoder.Encode(&reply)
		conn.Close()
	}
}

func getRandomWorker(workers []string)(string){
	
	return workers[0]

}

func sendTaskToWorker(
	request com.Request,
	worker string,
	reply_chan chan(com.Reply),
) {

	conn, err := net.Dial("tcp", worker)
	com.CheckError(err)
	encoder := gob.NewEncoder(conn)

	err = encoder.Encode(&request) // send request
	com.CheckError(err)

	go receiveReplyFromWorker(conn, reply_chan)
}

func receiveReplyFromWorker(
	conn net.Conn,
	reply_chan chan(com.Reply),
) {

	var reply com.Reply
	decoder := gob.NewDecoder(conn)
	err := decoder.Decode(&reply) //  receive reply

	com.CheckError(err)

	reply_chan <- reply

	conn.Close()
}

func readWorkers(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var workers []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			workers = append(workers, line)
		}
		log.Println(line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return workers, nil
}

func main() {

	WORKERS_FILE := "cmd/server-draft/workers.txt"

	GORUTINE_POOL_SIZE := 50
	args := os.Args
	if len(args) != 2 {
		log.Println("Error: endpoint missing: go run server.go ip:port")
		os.Exit(1)
	}
	endpoint := args[1]
	listener, err := net.Listen("tcp", endpoint)
	com.CheckError(err)

	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
	
	log.Println("Reading workers endPoints")
	workers, err := readWorkers(WORKERS_FILE)

	conn_chan := make(chan(net.Conn))

	log.Println("Launching gorutine pool")

	for id:=0; id < GORUTINE_POOL_SIZE; id++ {
		go processRequests(id, conn_chan, workers)
	}

	log.Println("***** Listening for new connection in endpoint ", endpoint)
	for {
		conn, err := listener.Accept()
		com.CheckError(err)
		log.Println("New connection: waiting for a gorutine end")
		conn_chan <- conn
	}
}
