package network
import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
)



func RunElev1() {
	fmt.Println("Starting server!")
	done := make(chan bool)
	broadcastAddr := "192.168.1.255:12010" // connecting to elevator 2
	//broadcastAddr := "localhost:12000" //localhost
	listenAddr := ":12012"	// listening to elevator 2
	lAddr, err := net.ResolveUDPAddr("udp", listenAddr)
	CheckError(err)
	bAddr, err := net.ResolveUDPAddr("udp", broadcastAddr)
	CheckError(err)
	lConn, err := net.ListenUDP("udp", lAddr)
	CheckError(err)
	bConn, err := net.DialUDP("udp", nil, bAddr)
	CheckError(err)
	go CastDataElev1(bConn)
	go ReceiveDataElev1(lConn)
	
	broadcastAddr = "192.168.1.255:12011" // connecting to elevator 3
	//broadcastAddr := "localhost:12000" //localhost
	listenAddr = ":12014"	// listening to elevator 3
	lAddr, err = net.ResolveUDPAddr("udp", listenAddr)
	CheckError(err)
	bAddr, err = net.ResolveUDPAddr("udp", broadcastAddr)
	CheckError(err)
	lConn, err = net.ListenUDP("udp", lAddr)
	CheckError(err)
	bConn, err = net.DialUDP("udp", nil, bAddr)
	CheckError(err)
	go CastDataElev1(bConn)
	go ReceiveDataElev1(lConn)
	fmt.Println("Sockets created successfully Elev 1")
	<-done
}
func CheckError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

func checkPeerLife(msg Message) bool{
	tdiff := time.Since(msg.CurrTime)
	if(tdiff > 8){
		fmt.Println("Connection to elev", msg.ID, " was lost!")
	}
}

func ReceiveDataElev1(conn *net.UDPConn) {
	var msg Message
	//var err error
	var b = make([]byte, 1024)
	for {
			time.Sleep(100 * time.Millisecond)
			n, _, err := conn.ReadFromUDP(b)
			CheckError(err)
			err = json.Unmarshal(b[:n], &msg)
			CheckError(err)
			fmt.Println("Data Received on Elev1:", msg)
			//fmt.Println("Received case Server:", msg.Word)
		}
}

func CastDataElev1(conn *net.UDPConn) {
	var msg Message
	var err error
	for {
		//msg = <-OutputCh
		//msg.ID = types.CART_ID
		msg.CurrTime = time.Now()
		msg.Word = "Message from Elev1!"
		for i := 0; i < 1; i++ {
			//fmt.Println("Data casted on Server:", msg)
			b := make([]byte, 1024)
			b, err = json.Marshal(msg)
			CheckError(err)
			_, err = conn.Write(b)
			CheckError(err)
			time.Sleep(15 * time.Second)
		}
	}
}

