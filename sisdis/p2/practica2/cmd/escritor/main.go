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

	my_ra := ra.New(me, users_file, ra.Role{Fname: "write"}, 
			barrier_chan, file_name)
	fmt.Println("RA inicializado para proceso", me)

	fmt.Println("Esperando a que se establezca la escucha")
	time.Sleep(4 * time.Second)

	fmt.Println("Goroutine Distribute_received_messages lanzada")

	// Cambiar bucles a 1-based para IDs
	for i := 0; i < my_ra.TotalProcesos; i++ {
		if i != me {
			fmt.Println("Enviando Barrier a proceso", i)
			my_ra.Ms.Send(i, ra.Barrier{})
		}
	}

	for i := 0; i < my_ra.TotalProcesos-1; i++ {
		<-barrier_chan
		fmt.Println("Barrier recibido", i)
	}

	reader := bufio.NewReader(os.Stdin)

	numEscritos := 0

	for {
		if shell {
			fmt.Print("> ") // Prompt estilo shell
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			if input == "exit" {
				fmt.Println("Saliendo del programa.")
				//my_ra.Stop()
				break
			}

			if strings.HasPrefix(input, "escribir ") {
				text := strings.TrimPrefix(input, "escribir ")
				
				//fmt.Println("Iniciando PreProtocol")
				my_ra.PreProtocol()

				fmt.Println("Escribiendo en archivo")
				mf.Write_file(file_name, text)

				for i := 0; i < my_ra.TotalProcesos; i++ {
					if i != me {
						//fmt.Println("Enviando Update a proceso", i)
						my_ra.Ms.Send(i, ra.Update{New_text: text})
					}
				}

				//fmt.Println("Ejecutando PostProtocol")
				my_ra.PostProtocol()
			} else {
				fmt.Println("Comando no reconocido. Usa: escribir <texto> o salir.")
			}
		} else {
			text := "Proceso: " + strconv.Itoa(me) + 
					" escribiendo. Escritos: " + strconv.Itoa(numEscritos)
					
			numEscritos++

			//fmt.Println("Iniciando PreProtocol")
			my_ra.PreProtocol()

			fmt.Println("Escribiendo en archivo")
			mf.Write_file(file_name, text)

			for i := 0; i < my_ra.TotalProcesos; i++ {
				if i != me {
					//fmt.Println("Enviando Update a proceso", i)
					my_ra.Ms.Send(i, ra.Update{New_text: text})
				}
			}

			//fmt.Println("Ejecutando PostProtocol")
			my_ra.PostProtocol()
		}
		
	}

}
