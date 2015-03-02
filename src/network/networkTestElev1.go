package network
import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
	"strings"
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
		if(strings.Contains(err.Error(),"use of closed network connection")){
			fmt.Println("closed network connection error found!")	
		}else{
			fmt.Println("Fatal error ", err.Error())
			os.Exit(1)
		}
	}
}

func CheckErrorUDP(err error, conn *net.UDPConn) bool {
	if err != nil {
		if(strings.Contains(err.Error(),"i/o timeout")){
			fmt.Println("IO Timeout: " ,err.Error())
			switch{
			case strings.Contains(err.Error(),":12010")||strings.Contains(err.Error(),":12011"):
				fmt.Println("problem with elev1: closing connection")
				conn.Close()
				return true
			case strings.Contains(err.Error(),":12012")||strings.Contains(err.Error(),":12013"):
				fmt.Println("problem with elev2: closing connection")
				conn.Close()
				return true
			case strings.Contains(err.Error(),":12014")||strings.Contains(err.Error(),":12015"):
				fmt.Println("problem with elev3: closing connection")
				conn.Close()
				return true
			}
		}else if(strings.Contains(err.Error(),"use of closed network connection")){
				fmt.Println("UDP connection closed error!")
				return true
		}else{
			fmt.Println("Fatal error ", err.Error())
			os.Exit(1)
			return true
		}
	}
	return false
}
/*
func checkPeerLife(msg Message) bool{
	tdiff := time.Since(msg.CurrTime)
	if(tdiff > 10*time.Second){	// 2 second timeout
		fmt.Println("Connection to elev", msg.ID, " was lost!, time diff: " , tdiff)
		return false
	}
	return true
}
*/
func ReceiveDataElev1(conn *net.UDPConn) {
	var msg Message
	//var err error
	var b = make([]byte, 1024)
	time.Sleep(10*time.Second)
	fmt.Println("Starting elev1!")
	for {
			//fmt.Println("msg: " , msg)
			//checkPeerLife(msg)
			time.Sleep(100 * time.Millisecond)
			err := conn.SetDeadline(time.Now().Add(30 * time.Second))
			CheckError(err)
			n, _, err := conn.ReadFromUDP(b)
			if(CheckErrorUDP(err,conn)==false){
				err = json.Unmarshal(b[:n], &msg)
				CheckError(err)
				fmt.Println("Data Received on Elev1:", msg)
			}else if(strings.Contains(err.Error(),"use of closed network connection")){
				fmt.Println("stopping function")
				break	// stop the function
			}
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
		msg.ID = "elev 1"
		for i := 0; i < 1; i++ {
			//fmt.Println("Data casted on Server:", msg)
			b := make([]byte, 1024)
			b, err = json.Marshal(msg)
			CheckError(err)
			_, err = conn.Write(b)
			CheckError(err)
			time.Sleep(5 * time.Second)
		}
	}
}

