// go run Exercise4.go
package network

import (
"fmt"
"net"
"time"
"encoding/json"
)

type Message struct {
	ID string
	Word string
	CurrTime time.Time
	LocalIP string
	RemoteIP string
	RawWord string
}




func UDP_receive(readFromPort string, receiveCh chan Message)(err error) {
	var receiveMessage Message


	addr, err := net.ResolveUDPAddr("udp", readFromPort)
	if err != nil {
		return err
	}
	
	lConnection, err := net.ListenUDP("udp", addr)
	
	if err != nil {
		return err
	}
	
	for {
		buffer := make([]byte, 2048)
		messageLength,addr,err := lConnection.ReadFromUDP(buffer[0:])
		if err != nil {fmt.Println(err); return err}

		err = json.Unmarshal(buffer[:messageLength], &receiveMessage)
		receiveMessage.LocalIP = addr.String()
		receiveMessage.CurrTime = time.Now()
		if err != nil {
			//fmt.Println(err)
			return err
		}

		receiveMessage.RawWord = string(buffer)
		receiveCh <- receiveMessage
	}
}

func UDP_broadcast(baddr string, sendCh chan string) (error){
	var transmitMessage Message
	transmitMessage.ID = "1"
	transmitMessage.Word = <- sendCh
	transmitMessage.RemoteIP = baddr


	tempConn, err := net.Dial("udp", baddr)
	if err != nil {
		return err
	}

	buffer, err := json.Marshal(transmitMessage)
	if err != nil {
		return err
	}

	for {
		tempConn.Write([]byte(buffer))
		time.Sleep(100*time.Millisecond)
	}
}



func run() {
	receiveChannel := make(chan Message, 1024)
	sendChannel := make(chan string, 1024)

	//message := Message{}

	go UDP_broadcast("192.168.1.255:24568", sendChannel)
	go UDP_receive("24568", receiveChannel)
	time.Sleep(100*time.Millisecond)

	for {
		sendChannel <-"NOt generic"
		i := <- receiveChannel
		fmt.Println("\n\nMessage received on: ", i.CurrTime)
		fmt.Println("\nMessage ID was: ", i.ID)
		fmt.Println("\nMessage contents: ", i.Word)
		fmt.Println("\nLocal IP was: ", i.LocalIP)
		fmt.Println("\nRemote IP was: ", i.RemoteIP)
		fmt.Println("\nRaw contents: ", i.RawWord)
		fmt.Println("__________________________\n")
		time.Sleep(100*time.Millisecond)
	}
}

