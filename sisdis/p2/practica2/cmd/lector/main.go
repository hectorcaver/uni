package main

import (
	"bufio"
	"fmt"
	"os"
	"practica2/mf"
	"practica2/ra"
	"strconv"
	"strings"
	"time"
)

func main() {

	if len(os.Args) < 2 && len(os.Args) > 3 {
		fmt.Println("Número de argumentos inválido: uso main.go <me> [<\"shell\"||\"loop\">], default \"loop\"")
		return
	}
	me, _ := strconv.Atoi(os.Args[1])
	fmt.Println("Proceso ID:", me)

	var shell bool

	if len(os.Args) == 3 {
		switch os.Args[2] {
		case "shell": shell = true
		case "loop": shell = false
		default: shell = false
		}
	}

	
	file_name := "file_" + os.Args[1] + ".txt"
	users_file := "./ms/users.txt"

	fmt.Println("Creando archivo:", file_name)
	mf.Create_file(file_name)

	barrier_chan := make(chan ra.Barrier)

	my_ra := ra.New(me, users_file, ra.Role{Fname: "read"}, 
				barrier_chan, file_name)
	
	fmt.Println("RA inicializado para proceso", me)

	fmt.Println("Esperando a que se establezca la escucha")
	time.Sleep(4 * time.Second)

	fmt.Println("Goroutine Distribute_received_messages lanzada")

	for i := 0; i < my_ra.TotalProcesos; i++ {
		if i != me {
			fmt.Println("Enviando Barrier a proceso", i)
			my_ra.Ms.Send(i, ra.Barrier{})
		}
	}

	for i := 0; i < my_ra.TotalProcesos-1; i++ {
		<-barrier_chan
		fmt.Println("Barrier recibido", i+1)
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		if shell {
			fmt.Print("> ") // Shell prompt
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			switch input {
			case "exit":
				fmt.Println("Saliendo del programa.")
				//my_ra.Stop()
				return

			case "leer":
				//fmt.Println("Iniciando PreProtocol")
				my_ra.PreProtocol()

				//fmt.Println("Leyendo archivo:", file_name)
				content := mf.Read_file(file_name)
				
				fmt.Println("======== CONTENIDO ========")
				fmt.Println("")
				fmt.Println(content)
				fmt.Println("=========== EOF ===========")


				//fmt.Println("Ejecutando PostProtocol")
				my_ra.PostProtocol()

			default:
				fmt.Println("Comando no reconocido. Usa: leer o salir.")
			}
		} else {
			//fmt.Println("Iniciando PreProtocol")
			my_ra.PreProtocol()

			//fmt.Println("Leyendo archivo:", file_name)
			content := mf.Read_file(file_name)
			
			fmt.Println("======== CONTENIDO ========")
			fmt.Println("")
			fmt.Println(content)
			fmt.Println("=========== EOF ===========")


			//fmt.Println("Ejecutando PostProtocol")
			my_ra.PostProtocol()

			time.Sleep(300 * time.Millisecond)
		}
		
	}

}
