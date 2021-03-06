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
	PeerMap = newPeerMap()

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
	dataTypes.ElevatorID, dataTypes.NumberOfElevatorsConnected = getElevatorsFromTerminal()
	done := make(chan bool)
	
	if(dataTypes.NumberOfElevatorsConnected > 1){	
		var listenConnection *net.UDPConn 
		
		broadcastPort := startingPort+(2*(dataTypes.ElevatorID-1))
		broadcastAddress, err := net.ResolveUDPAddr("udp", "129.241.187.255:" + strconv.Itoa(broadcastPort))
		checkError(err)		
		broadcastConnection, err := net.DialUDP("udp", nil, broadcastAddress)
		checkError(err)

		for e:=1; e<=dataTypes.NumberOfElevatorsConnected; e++{
			if (e != dataTypes.ElevatorID){
				listenPort := startingPort+(2*(e-1))
				listenAddress, err := net.ResolveUDPAddr("udp", ":" + strconv.Itoa(listenPort))
				checkError(err)
				listenConnection, err = net.ListenUDP("udp", listenAddress)
				checkError(err)

				go receivePacket(listenConnection)
			}
			
		}
		fmt.Println("Sockets created successfully")
		go sendPacket(broadcastConnection)
		go updatePeerMap(PeerMap) // ???
		go elevatorIsAlive() // ???
	}
	<-done
}
/*-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/
func checkError(err error) {
	if err != nil{
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
/*-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/
func sendPacket(conn *net.UDPConn) {
	var data dataTypes.Message
	var err error
	for {
		data = <-UDPCh
		data.Word = "Message from Elev"+strconv.Itoa(dataTypes.ElevatorID)
		data.T = time.Now()
		//fmt.Println("Message sent!")
		for i := 0; i < 1; i++ {
			//fmt.Println("Data casted:", data)
			b := make([]byte, 1024)
			b, err = json.Marshal(data)
			checkError(err)
			_, err = conn.Write(b)
			checkError(err)
			time.Sleep(1 * time.Millisecond)
		}
	}
}

func receivePacket(conn *net.UDPConn) {
	var receivedMessage dataTypes.Message
	var buffer = make([]byte, 1024)
	time.Sleep(3*time.Second) // var 10 sek
	fmt.Println("Starting elev" , dataTypes.ElevatorID)
	for {
		time.Sleep(100 * time.Millisecond)
		n, _, err := conn.ReadFromUDP(buffer)
		checkError(err)
		err = json.Unmarshal(buffer[:n], &receivedMessage)
		checkError(err)
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
/*-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/
func newPeerMap() *dataTypes.PeerMap {
	return &dataTypes.PeerMap{M: make(map[int]time.Time)}
}

func updatePeerMap(p *dataTypes.PeerMap) {
	var id int
	for {
		time.Sleep(100 * time.Millisecond)
		id = <-peerCh
		p.M[id] = time.Now()
		//fmt.Println("Updated peerMap: ", PeerMap)
	}

}

func CheckPeerLife(p dataTypes.PeerMap, id int) bool {
	_, present := p.M[id]
	if present {
		tdiff := time.Since(p.M[id])
		return tdiff <= dataTypes.TIMEOUT
	}
	return false
}
/*-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/
func elevatorIsAlive()() {
	for {
		UDPCh <- dataTypes.Message{ID: dataTypes.ElevatorID}
		time.Sleep(1 * time.Second)
	}
}
/*-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------*/
func getElevatorsFromTerminal() (int, int){
	var n int
	var id int	
	fmt.Println("How many elevators?")
	fmt.Scanf("%d",&n)

	if(n==1){
		id = 1
	}else if(n>1){
		fmt.Print("Which elevator do you want to run? Choose whole number between 1-",n, ": ")
		fmt.Scanf("%d",&id)
	}
	return id, n
}
