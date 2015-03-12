package network
import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
	"strings"
	"strconv"
	//."../order"
)
/*
var 
(
	Global_stopCastDataElev1 bool = false
	Global_stopCastDataElev2 bool = false
	Global_stopCastDataElev3 bool = false
)*/
// data that is sent between elevators
type Data struct {
	OrderQueue[4][3] int
	ID string
	Word string
	CurrTime time.Time
	LocalIP string
	RemoteIP string
	RawWord string
}

// elev 1 uses ports :12010 and :12011, elev 2 :12012 and :12013, elev 3 :12014 and :12015
func RunElev(elev int,numElev int) {
	fmt.Println("Starting server!")
	done := make(chan bool)
	port := 0
	switch elev{
	case 1:
		port = 12010
		break
	case 2:
		port = 12012
		break
	case 3:
		port = 12014
		break
	}
	for i:=0; i<numElev-1; i++{
		port+=i
		portString := strconv.Itoa(port)
		broadcastAddr := "129.241.187.255:"+portString // connecting to elevator 2
		//broadcastAddr := "localhost:12000" //localhost
		listenAddr:=""
		if(i==0 && elev==1){
			listenAddr = ":12012"
		}else if(i==1 && elev==1){
			listenAddr = ":12014"
		}else if(i==0 && elev==2){
			listenAddr = ":12010"
		}else if(i==1 && elev==2){
			listenAddr = ":12015"
		}else if(i==0 && elev==3){
			listenAddr = ":12011"
		}else if(i==1 && elev==3){
			listenAddr = ":12013"
		}
		lAddr, err := net.ResolveUDPAddr("udp", listenAddr)
		CheckError(err)
		bAddr, err := net.ResolveUDPAddr("udp", broadcastAddr)
		CheckError(err)
		lConn, err := net.ListenUDP("udp", lAddr)
		CheckError(err)
		bConn, err := net.DialUDP("udp", nil, bAddr)
		CheckError(err)
		go SendData(bConn,elev)
		go ReceiveData(lConn,elev)
	}
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

func ReceiveData(conn *net.UDPConn, elev int) {
	var data Data
	//var err error
	var b = make([]byte, 1024)
	time.Sleep(10*time.Second)
	fmt.Println("Starting elev" , elev)
	for {
			//fmt.Println("msg: " , msg)
			//checkPeerLife(msg)
			time.Sleep(100 * time.Millisecond)
			err := conn.SetDeadline(time.Now().Add(30 * time.Second))
			CheckError(err)
			n, _, err := conn.ReadFromUDP(b)
			if(CheckErrorUDP(err,conn)==false){
				err = json.Unmarshal(b[:n], &data)
				CheckError(err)
				fmt.Println("Data Received on Elev", elev, data)
			}else if(strings.Contains(err.Error(),"use of closed network connection")){
				fmt.Println("stopping receiveData function")
				/*switch elev{
					case 1:
						Global_stopCastDataElev1 = true
						break
					case 2:
						Global_stopCastDataElev2 = true
						break
					case 3:
						Global_stopCastDataElev3 = true
						break
				}*/
				break	// stop the function
			}
		}
}

func SendData(conn *net.UDPConn, elev int) {
	var data Data
	var err error
	for {
/*
		if((elev==1 && Global_stopCastDataElev1==true)||(elev==2 && Global_stopCastDataElev2==true)||(elev==3 && Global_stopCastDataElev3==true)){
			fmt.Println("Stopping CastData function")
			break		
		}
*/
		//msg = <-OutputCh
		//msg.ID = types.CART_ID
		data.CurrTime = time.Now()
		data.Word = "Message from Elev"+strconv.Itoa(elev)
		data.ID = "elev"
		//data.OrderQueue
		for i := 0; i < 1; i++ {
			//fmt.Println("Data casted on Server:", msg)
			b := make([]byte, 1024)
			b, err = json.Marshal(data)
			CheckError(err)
			_, err = conn.Write(b)
			CheckError(err)
			time.Sleep(5 * time.Second)
		}
	}
}

