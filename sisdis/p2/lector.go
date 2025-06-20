package main

import (
	"os"
	"practica2/ra"
	"strconv"
)

func main() {
	me, _ := strconv.Atoi(os.Args[1])
	file_name := "file_" + os.Args[1] + "-txt"
	users_file := "./ms/users.txt"
	mf.create_file(file_name)
	my_ra := ra.New(me, users_file, ra.Role{fname: "read"})
	barrier_chan := make(chan ra.Barrier)
	go ra.distribute_received_messages(my_ra.msg, make(chan ra.Request), make(chan ra.Reply), barrier_chan)
	for i := 0; i < ra.total_processes; i++ {
		if i+1 != me {
			my_ra.msg.Send(i+1, ra.Barrier{})
		}
	}
	for i := 0; i < ra.total_processes-1; i++ {
		<-barrier_chan
	}
	for {
		my_ra.PreProtocol()
		_ = mf.read_file(file_name)
		my_ra.PostProtocol()
	}
}
