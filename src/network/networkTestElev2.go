package network
import (
	"encoding/json"
	"fmt"
	"net"
	//"os"
	"time"
	"strings"
)



func RunElev2() {
	fmt.Println("Starting server!")
	done := make(chan bool)
	broadcastAddr := "192.168.1.255:12012" // connecting to elevator 1
	//broadcastAddr := "localhost:12000" //localhost
	listenAddr := ":12010"	// listening to elevator 1
	lAddr, err := net.ResolveUDPAddr("udp", listenAddr)
	CheckError(err)
	bAddr, err := net.ResolveUDPAddr("udp", broadcastAddr)
	CheckError(err,)
	lConn, err := net.ListenUDP("udp", lAddr)
	CheckError(err)
	bConn, err := net.DialUDP("udp", nil, bAddr)
	CheckError(err)
	go CastDataElev2(bConn)
	go ReceiveDataElev2(lConn)
	/*
	broadcastAddr = "192.168.1.255:12013" // connecting to elevator 3
	//broadcastAddr := "localhost:12000" //localhost
	listenAddr = ":12015"	// listening to elevator 2
	lAddr, err = net.ResolveUDPAddr("udp", listenAddr)
	CheckError(err)
	bAddr, err = net.ResolveUDPAddr("udp", broadcastAddr)
	CheckError(err)
	lConn, err = net.ListenUDP("udp", lAddr)
	CheckError(err)
	bConn, err = net.DialUDP("udp", nil, bAddr)
	CheckError(err)
	go CastDataElev2(bConn)
	go ReceiveDataElev2(lConn)
	*/
	fmt.Println("Sockets created successfully Elev 2")
	<-done
}


func ReceiveDataElev2(conn *net.UDPConn) {
	var msg Message
	//var err error
	var b = make([]byte, 1024)
	time.Sleep(10*time.Second)
	fmt.Println("Starting elev2!")
	for {
			time.Sleep(100 * time.Millisecond)
			err := conn.SetDeadline(time.Now().Add(30 * time.Second))
			CheckError(err)
			n, _, err := conn.ReadFromUDP(b)
			if(CheckErrorUDP(err,conn)==false){
				err = json.Unmarshal(b[:n], &msg)
				CheckError(err)
				fmt.Println("Data Received on Elev2:", msg)
			}else if(strings.Contains(err.Error(),"use of closed network connection")){
				fmt.Println("stopping function")
				break	// stop the function
			}
		}
}

func CastDataElev2(conn *net.UDPConn) {
	var msg Message
	var err error
	for {
		//msg = <-OutputCh
		//msg.ID = types.CART_ID
		msg.CurrTime = time.Now()
		msg.Word = "Message from Elev2!"
		msg.ID = "elev 2"
		for i := 0; i < 1; i++ {
			//fmt.Println("Data casted on Client:", msg)
			b := make([]byte, 1024)
			b, err = json.Marshal(msg)
			CheckError(err)
			_, err = conn.Write(b)
			CheckError(err)
			time.Sleep(5 * time.Second)
		}
	}
}
