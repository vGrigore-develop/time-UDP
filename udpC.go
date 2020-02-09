package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
)

func main() {
	regExDate, _ := regexp.Compile("([0-2]\\d|[3][0-1])\\.([0]\\d|[1][0-2])\\.([2][01]|[1][6-9])\\d{2}(\\s([0-1]\\d|[2][0-3])(\\:[0-5]\\d){1,2})?\\s\\-[smza]")
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a host:port string")
		return
	}
	CONNECT := arguments[1]

	s, err := net.ResolveUDPAddr("udp4", CONNECT)
	c, err := net.DialUDP("udp4", nil, s)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("The UDP server is %s \n", c.RemoteAddr().String())
	defer c.Close()

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSuffix(text, "\n")

		if strings.TrimSpace(text) == "STOP" {
			fmt.Println("Exiting UDP client!")
			_, err = c.Write([]byte(text))
			return
		}

		if regExDate.MatchString(text) {

			s := strings.Split(text, " ")
			day := strings.Split(s[0], ".")[0]
			month := strings.Split(s[0], ".")[1]
			year := strings.Split(s[0], ".")[2]

			dateSrvFormat := year + "-" + month + "-" + day + " " + s[1] + " " + s[2]

			data := []byte(dateSrvFormat)
			_, err = c.Write(data)

			if err != nil {
				fmt.Println(err)
				return
			}

			buffer := make([]byte, 1024)
			n, _, err := c.ReadFromUDP(buffer)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("Reply: %s\n", string(buffer[0:n]))
		} else {
			fmt.Printf("Error: Unkown input format. Please use: dd.mm.yyyy hh:mm:ss -[s|m|z|a]\n")
		}
	}
}
