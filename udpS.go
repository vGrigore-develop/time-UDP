package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number!")
		return
	}
	PORT := ":" + arguments[1]

	s, err := net.ResolveUDPAddr("udp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}

	connection, err := net.ListenUDP("udp4", s)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer connection.Close()
	buffer := make([]byte, 1024)

	for {
		n, addr, err := connection.ReadFromUDP(buffer)
		fmt.Print("-> ", string(buffer[0:n-2]))

		if strings.TrimSpace(string(buffer[0:n])) == "STOP" {
			fmt.Println("Exiting UDP server!")
			return
		}

		input := string(buffer[0:n])
		s := strings.Split(input, " ")
		timeInput := s[0] + " " + s[1]
		option := s[2]

		then, err := time.Parse(timeFormat, timeInput)
		if err != nil {
			fmt.Println(err)
			return
		}
		t1, _ := time.Parse(timeFormat, "1970-01-01 00:00:00")
		hours := then.Sub(t1).Hours()

		data := []byte("No option")

		switch option {
		case "-s":
			seconds := hours * 3600
			data = []byte(strconv.Itoa(int(seconds)))
			break

		case "-m":
			minutes := hours * 60
			data = []byte(strconv.Itoa(int(minutes)))
			break

		case "-z":
			days := hours / 24
			data = []byte(strconv.Itoa(int(days)))
			break

		case "-a":
			years := hours / 8760
			data = []byte(strconv.Itoa(int(years)))
			break

		}

		fmt.Printf("\ndata: %s\n", string(data))
		_, err = connection.WriteToUDP(data, addr)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
