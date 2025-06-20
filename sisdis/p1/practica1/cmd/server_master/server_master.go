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
	"bufio"
	"encoding/gob"
	"log"
	"math/rand"
	"net"
	"os"
	"practica1/com"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh"
)

// Función para iniciar un worker en una máquina remota usando SSH
func startWorker(
	ip string, 
	port string, 
	user string, 
	path string, 
	folder_path string,
	config *ssh.ClientConfig ) (error) {

	client, err := ssh.Dial("tcp", ip+":22", config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
		os.Exit(1)
	}

	newSession, err := client.NewSession()

	if err != nil {
		log.Fatal("Failed to create session: ", err)
		os.Exit(1)
	}

	cmd := "cd " + folder_path + " && go run " + path + " " + ip+":"+port

	payload := ssh.Marshal(struct {
    Command string
	}{
		Command: cmd,
	})

	ok, err := newSession.SendRequest("exec", true, payload)
	if err != nil {
		log.Fatal(err)
	}
	if !ok {
		log.Fatal("request rejected by server")
	}

	return err
}

func configssh(user string) (config *ssh.ClientConfig) {
	key, err := os.ReadFile("/home/" + user + "/.ssh/id_rsa")
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
		os.Exit(1)
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
		os.Exit(1)
	}

	return &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
}

// Inicia los workers de manera remota mediante SSH
func startWorkers(
	workers []string, 
	user string, 
	path string, 
	folder_path string) {

	config := configssh(user)

	for _, worker := range workers {
		parts := strings.Split(worker, ":")
		ip := parts[0]
		port := parts[1]
		// Comando para ejecutar el worker en la máquina remota mediante SSH
		err := startWorker(ip, port, user, path, folder_path, config)

		if err != nil {
			log.Fatalf("Error arrancando worker en %s: %v", ip, err)
		}
		log.Printf("Worker iniciado en %s\n", ip)
	}
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
		log.Println("GORUTINE "+strconv.Itoa(id)+": accepted new connection")

		// Recojo una petición nueva
		var request com.Request
		decoder := gob.NewDecoder(conn)
		err := decoder.Decode(&request)
		com.CheckError(err)

		// Selecciono un worker aleatorio
		worker := getRandomWorker(workers)

		// Le envío una petición al worker y un canal para recibir la respuesta
		sendTaskToWorker(request, worker, reply_chan)

		// Recibe la respuesta del worker
		reply := <- reply_chan

		// Envía la respuesta de la petición al cliente
		encoder := gob.NewEncoder(conn)
		encoder.Encode(&reply)
		conn.Close()
	}
}

func getRandomWorker(workers []string)(string){
	if len(workers) == 0 {
		log.Println("No workers available")
		os.Exit(1)
	}
	return workers[rand.Intn(len(workers))]
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

	WORKERS_FILE := "cmd/server_master/workers.txt"
	GORUTINE_POOL_SIZE := 50
	USER := "a869637"
	FOLDER_PATH := "/misc/alumnos/sd/sd2024/" + USER + "/practica1"
	PATH := "cmd/server_worker/server_worker.go"

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

	if err != nil {
		log.Println("Error reading workers file:", err)
		os.Exit(1)
	}

	log.Println("Launching workers")

	startWorkers(workers, USER, PATH, FOLDER_PATH) 

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
