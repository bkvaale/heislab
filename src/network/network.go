package network

import (
	"../dataTypes"
	//"../order"	// should change this...
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
	//"strings"
	"strconv"
)

const (
	startingPort = 12010
)

var (
	PeerMap = NewPeerMap()

	// Global channels
	UDPCh            		= make(chan dataTypes.Message)
	OrderCh            		= make(chan []int)
	TableCh            		= make(chan dataTypes.Matrix)
	DistributeExternalOrderCh       = make(chan dataTypes.Message)
	AddOrderCh          		= make(chan dataTypes.Message)
	RemoveOrderCh       		= make(chan []int)
	UpdateGlobalTableCh 		= make(chan dataTypes.Matrix)
	
	CostOfOrderOverNetworkCh	= make(chan dataTypes.Message)
	CalculateCostNetworkCh		= make(chan dataTypes.Message)

	CalculateCostCh			= make(chan dataTypes.Message) //NY // REGNE cost
	CostToGetOrderCh		= make(chan dataTypes.Message) //NY //BRUKES TIL Å SENDETILBAKE RES

	// Local channels
	peerCh = make(chan int)
)







func Run(){
	fmt.Println("How many elevators?")
	fmt.Scanf("%d",&dataTypes.NumberOfElevatorsConnected)

	if(dataTypes.NumberOfElevatorsConnected==1){
		StartElev(1)

	}else if(dataTypes.NumberOfElevatorsConnected>1){
		fmt.Print("Which elevator do you want to run? Choose whole number between 1-",dataTypes.NumberOfElevatorsConnected, ":")
		fmt.Scanf("%d",&dataTypes.ElevatorID)
		StartElev(dataTypes.ElevatorID)	
	}
}


func StartElev(elevator int){
	if(dataTypes.NumberOfElevatorsConnected > 1){	
		//fmt.Println("Starting server!")
		done := make(chan bool)
		var listenConnection *net.UDPConn 
		//var broadcastConnection *net.UDPConn
		broadcastPort := startingPort+(2*(elevator-1))
		broadcastAddressString := "129.241.187.255:" + strconv.Itoa(broadcastPort)

		broadcastAddress, err := net.ResolveUDPAddr("udp", broadcastAddressString)
		CheckError(err)		
		broadcastConnection, err := net.DialUDP("udp", nil, broadcastAddress)
		CheckError(err)

		for e:=1; e<=dataTypes.NumberOfElevatorsConnected; e++{
			if (e != elevator){
				listenPort := startingPort+(2*(e-1))

				listenAddress, err := net.ResolveUDPAddr("udp", ":" + strconv.Itoa(listenPort))
				CheckError(err)
				listenConnection, err = net.ListenUDP("udp", listenAddress)
				CheckError(err)
				go ReceiveData(listenConnection, elevator)
			}
			
		}
		fmt.Println("Sockets created successfully")
		go SendData(broadcastConnection,elevator)	// CastData
		//go ReceiveData(listenConnection, elevator)
		go UpdatePeerMap(PeerMap) // ???
		go Ping() // ???
		<-done
	}
}

func CheckError(err error) {
	if err != nil{
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

func NewPeerMap() *dataTypes.PeerMap {
	return &dataTypes.PeerMap{M: make(map[int]time.Time)}
}

func CheckPeerLife(p dataTypes.PeerMap, id int) bool {
	_, present := p.M[id]
	if present {
		tdiff := time.Since(p.M[id])
		return tdiff <= dataTypes.TIMEOUT
	}
	return false
}

func UpdatePeerMap(p *dataTypes.PeerMap) {
	var id int
	for {
		time.Sleep(100 * time.Millisecond)
		id = <-peerCh
		p.M[id] = time.Now()
		//fmt.Println("Updated peerMap: ", PeerMap)
	}

}

func ReceiveData(conn *net.UDPConn, elev int) {
	var receivedMessage dataTypes.Message
	var buffer = make([]byte, 1024)
	time.Sleep(3*time.Second) // var 10 sek
	fmt.Println("Starting elev" , elev)
	for {
		time.Sleep(100 * time.Millisecond)
		n, _, err := conn.ReadFromUDP(buffer)
		CheckError(err)
		err = json.Unmarshal(buffer[:n], &receivedMessage)
		CheckError(err)
		//fmt.Println("Data Received:", receivedMessaged
		if receivedMessage.ID > 0 {
			// update peermap
			peerCh <- receivedMessage.ID
		}
		//fmt.Println("Received case:", receivedMessage.Head)
		switch receivedMessage.Head {

		case "order":
			OrderCh <- receivedMessage.Order
			fmt.Println("Order received:", receivedMessage.Order)
			break
		case "table":
			fmt.Println("Table received", receivedMessage.Table)
			break
		case "CalculateCostforNewExternalOrder":
			//fmt.Println("We entered CalculateCostForNewExternalOrder-Case")
			//fmt.Println("Elevator",receivedMessage.ID," With order ",receivedMessage.Order)
			CalculateCostNetworkCh <- receivedMessage
			//fmt.Println("Order sent to CalculateCostNetworkCh")
			break
		case "costCalculatedNetwork":
			CostOfOrderOverNetworkCh <- receivedMessage
			//fmt.Println("Cost calculated: ", receivedMessage.Cost, " and delivered by elevator: ", receivedMessage.ID)
			break
		case "addorder":
			AddOrderCh <- receivedMessage
			fmt.Println("Order added:", receivedMessage.Order, " received by: ", receivedMessage.ID)
			break
		case "removeOrder":
			RemoveOrderCh <- receivedMessage.Order
			fmt.Println("Order removed:", receivedMessage.Order)
			break
		default:
			fmt.Println("Default case entered") 
			break
		}
	}
}

func SendData(conn *net.UDPConn, elev int) {
	var data dataTypes.Message
	var err error
	for {
		data = <-UDPCh
		data.Word = "Message from Elev"+strconv.Itoa(elev)
		data.T = time.Now()
		//fmt.Println("Message sent!")
		for i := 0; i < 1; i++ {
			//fmt.Println("Data casted:", data)
			b := make([]byte, 1024)
			b, err = json.Marshal(data)
			CheckError(err)
			_, err = conn.Write(b)
			CheckError(err)
			time.Sleep(1 * time.Millisecond)
		}
	}
}

func Ping() {
	for {
		UDPCh <- dataTypes.Message{ID: dataTypes.ElevatorID}
		time.Sleep(1 * time.Second)
	}
}
