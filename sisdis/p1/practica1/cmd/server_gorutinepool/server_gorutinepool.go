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
	"net"
	"os"
	"practica1/com"
	"strconv"
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

func processRequests(id int, conn_chan chan(net.Conn)){
	for{
		conn := <- conn_chan
		log.Println("GORUTINE " + strconv.Itoa(id) + ": accepted new connection")
		var request com.Request
		decoder := gob.NewDecoder(conn)
		err := decoder.Decode(&request)
		com.CheckError(err)
		primes := findPrimes(request.Interval)
		reply := com.Reply{Id: request.Id, Primes: primes}
		encoder := gob.NewEncoder(conn)
		encoder.Encode(&reply)
		conn.Close()
	}
}

func main() {

	GORUTINE_POOL_SIZE := 10
	args := os.Args
	if len(args) != 2 {
		log.Println("Error: endpoint missing: go run server.go ip:port")
		os.Exit(1)
	}
	endpoint := args[1]
	listener, err := net.Listen("tcp", endpoint)
	com.CheckError(err)

	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
	
	conn_chan := make(chan(net.Conn))

	log.Println("Launching gorutine pool")

	for id:=0; id < GORUTINE_POOL_SIZE; id++ {
		go processRequests(id, conn_chan)
	}

	log.Println("***** Listening for new connection in endpoint ", endpoint)
	for {
		conn, err := listener.Accept()
		com.CheckError(err)
		log.Println("New connection: waiting for a gorutine end")
		conn_chan <- conn
	}
}
