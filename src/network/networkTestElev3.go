package network
import (
	"encoding/json"
	"fmt"
	"net"
	//"os"
	"time"
)



func RunElev3() {
	fmt.Println("Starting server!")
	done := make(chan bool)
	broadcastAddr := "192.168.1.255:12014" // connecting to elevator 1
	//broadcastAddr := "localhost:12000" //localhost
	listenAddr := ":12011"	// listening to elevator 1
	lAddr, err := net.ResolveUDPAddr("udp", listenAddr)
	CheckError(err)
	bAddr, err := net.ResolveUDPAddr("udp", broadcastAddr)
	CheckError(err)
	lConn, err := net.ListenUDP("udp", lAddr)
	CheckError(err)
	bConn, err := net.DialUDP("udp", nil, bAddr)
	CheckError(err)
	go CastDataElev3(bConn)
	go ReceiveDataElev3(lConn)
	
	broadcastAddr = "192.168.1.255:12015" // connecting to elevator 2
	//broadcastAddr := "localhost:12000" //localhost
	listenAddr = ":12013"	// listening to elevator 2
	lAddr, err = net.ResolveUDPAddr("udp", listenAddr)
	CheckError(err)
	bAddr, err = net.ResolveUDPAddr("udp", broadcastAddr)
	CheckError(err)
	lConn, err = net.ListenUDP("udp", lAddr)
	CheckError(err)
	bConn, err = net.DialUDP("udp", nil, bAddr)
	CheckError(err)
	go CastDataElev3(bConn)
	go ReceiveDataElev3(lConn)
	fmt.Println("Sockets created successfully Elev 3")
	<-done
}


func ReceiveDataElev3(conn *net.UDPConn) {
	var msg Message
	//var err error
	var b = make([]byte, 1024)
	for {
			time.Sleep(100 * time.Millisecond)
			n, _, err := conn.ReadFromUDP(b)
			CheckError(err)
			err = json.Unmarshal(b[:n], &msg)
			CheckError(err)
			fmt.Println("Data Received on Elev3:", msg)
			//fmt.Println("Received case Client:", msg.Word)
		}
}

func CastDataElev3(conn *net.UDPConn) {
	var msg Message
	var err error
	for {
		//msg = <-OutputCh
		//msg.ID = types.CART_ID
		msg.CurrTime = time.Now()
		msg.Word = "Message from Elev3!"
		for i := 0; i < 1; i++ {
			//fmt.Println("Data casted on Client:", msg)
			b := make([]byte, 1024)
			b, err = json.Marshal(msg)
			CheckError(err)
			_, err = conn.Write(b)
			CheckError(err)
			time.Sleep(15 * time.Second)
		}
	}
}

