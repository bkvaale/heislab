package order

import (
	"../network"
	"../driver"
	"../dataTypes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	UP       = 0
	DOWN     = 1
	INTERNAL = 2
)

var (
	//Local Channals
	UpdateLightCh  		= make(chan string, 5)
	CalculateCostInternalCh	= make(chan dataTypes.Message)	
	CostOfOrderInternalCh = make(chan dataTypes.Message)	
	ExternalButtonPressedCh = make(chan dataTypes.Message)

	AddOrderInternalCh = make(chan dataTypes.Message) //Made this one to send who "won" the order by claim LOCALLY
	RemoveOrderLocalCh = make(chan []int)
	Direction      	int
	ExternalQueue   dataTypes.Matrix
	InternalQueue 	dataTypes.Array
	osChan         	chan os.Signal
)

func Run() {

	done := make(chan bool)

	ExternalQueue = dataTypes.NewExternalQueue()
	InternalQueue = dataTypes.NewInternalQueue()
	Direction = UP

	go CheckAndAddInternalOrder()
	ReadBackupFile()
	go UpdateLights()
	go CheckForExternalOrder()
	go DistributeExternalOrder()

	go SendCostOverNetwork()
	go SendCostLocally()	
	//go CheckButtonPressedAndDistributeExternalOrder()
	go AddExternalOrderLocallyToQueue()
	go AddExternalOrderFromNetworkToQueue()
	
	go RemoveExternalOrderDoneOverNetwork()
	go RemoveExternalOrderDoneLocally()
	//go Redistribute()
	
/*
	//go CheckAndAddInternalOrder()
	ReadBackupFile()
	go DistributeExternalOrder(ExternalQueue)
	go AddOrder()
	//go UpdateLights()
	//go CalculateAndSendCost()
	go RemoveExternalOrder()
	//go PrintTables()
	//go PrintOrderDirection()
	go CheckAndAddExternalOrder() //Bad name? Due to DistributeExternalOrder()
	go Redistribute()
*/
	<-done
}


/*		not needed here!!!! CHANGED!!
func CheckError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
*/

func CheckAndAddInternalOrder() {
	var buttonPushed int
	for {
		time.Sleep(10 * time.Millisecond)
		for i := 0; i < dataTypes.N_FLOORS; i++ { 
			if InternalQueue[i] != 1 {
				buttonPushed = driver.ElevGetButtonSignal(INTERNAL, i) 
				if buttonPushed == 1 {
					InternalQueue[i] = buttonPushed
					UpdateLightCh <- "internal"
				}
			}
		}
	}
}

func UpdateLights() {
	var msg string
	for {
		time.Sleep(1 * time.Millisecond)
		msg = <-UpdateLightCh

		switch msg {
		case "internal":
			for i := range InternalQueue {
				time.Sleep(1 * time.Millisecond)
				driver.ElevSetLights(i, 2, InternalQueue[i])
			}
		case "external": 

			for j := 0; j < dataTypes.N_FLOORS; j++ {
				for k := 0; k < 2; k++ {
					time.Sleep(1 * time.Millisecond)
					if ExternalQueue[j][k] != 0 {
						driver.ElevSetLights(j, k, 1)
					} else {
						driver.ElevSetLights(j, k, 0)
					}
				}
			}
		}
	}
}


func CheckForExternalOrder() []int { //////////////////////////////////////////////////////////////////////
	for {
		sendMessage := dataTypes.Message{Head: "externalButtonPanelPressed"} //UNECESSARY WITH HEAD
		sendMessage.WhichExternalPanelPressed = dataTypes.ElevatorID
		sendMessage.ID = dataTypes.ElevatorID
		//fmt.Println("CheckForExternalOrder() elevID: ", sendMessage.ID)
		time.Sleep(1 * time.Millisecond)
		for floor := 0; floor < dataTypes.N_FLOORS; floor++ {
			time.Sleep(1 * time.Millisecond)

			buttonPushed := driver.ElevGetButtonSignal(driver.BUTTON_CALL_UP, floor) 			
			if buttonPushed != 0 {
				sendMessage.Order = []int{floor, UP}
				ExternalButtonPressedCh <- sendMessage
				for driver.ElevGetButtonSignal(driver.BUTTON_CALL_UP, floor) == 1 {
					time.Sleep(1 * time.Millisecond)
				}

			} else{
				buttonPushed = driver.ElevGetButtonSignal(driver.BUTTON_CALL_DOWN, floor)
				if buttonPushed != 0 {
					sendMessage.Order = []int{floor, DOWN}
					ExternalButtonPressedCh <- sendMessage
					for driver.ElevGetButtonSignal(driver.BUTTON_CALL_DOWN, floor) == 1 {
						time.Sleep(1 * time.Millisecond)
					}
				}
			}
		}
	}
}

func DistributeExternalOrder() { //ANOTHER NAME IF POSSIBLE
	for{
		receivedMessage := <- ExternalButtonPressedCh
		sendMessage := dataTypes.Message{Head: "CalculateCostforNewExternalOrder"}
	

		var receiverID int
		var elevatorCostMap = make(map[int]int)	
		lowestCost := 50
		counter := 0

		sendMessage.Order = receivedMessage.Order
		sendMessage.ID = receivedMessage.ID
		sendMessage.WhichExternalPanelPressed = receivedMessage.WhichExternalPanelPressed

		startTimer:=time.Now()

		CalculateCostInternalCh <- sendMessage
		network.UDPCh <- sendMessage
		Tag:
			for{
				select {
					case elevator:=<-CostOfOrderInternalCh:
						counter = counter + 1
						elevatorCostMap[elevator.ID] = elevator.Cost
						fmt.Println("Penalty map cost: ", elevatorCostMap[elevator.ID], " for internal elev: ", elevator.ID)
						if counter == dataTypes.NumberOfElevatorsConnected{
							break Tag
						}
					case elevator:=<-network.CostOfOrderOverNetworkCh:
						fmt.Println("WhichExternalPanelPressed: ", elevator.WhichExternalPanelPressed)  
						fmt.Println("ElevatorID panel: ", dataTypes.ElevatorID)

						if(elevator.WhichExternalPanelPressed == dataTypes.ElevatorID){
							counter = counter + 1							
							elevatorCostMap[elevator.ID] = elevator.Cost
							fmt.Println("Penalty map cost in checkbuttonpresseddistributeexternalorder: ", elevatorCostMap[elevator.ID], " for elev: ", elevator.ID)
				
							if counter == dataTypes.NumberOfElevatorsConnected{
								break Tag
							}
						}else{
							fmt.Println("Elevator ", elevator.ID, " wasn't the elevator panel the order was asked on")
							break
						}
					default:
						if time.Since(startTimer)>3*time.Second {	// changed from 6 seconds
							fmt.Println("DistributeExternalOrder() Have not received message from all connected elevators. Breaking! Counter: ", counter)
							break Tag
						}
						
				}
			}

		if counter == dataTypes.NumberOfElevatorsConnected{
			for elevatorID := 1; elevatorID <= len(elevatorCostMap); elevatorID++ {
				if elevatorCostMap[elevatorID] < lowestCost {
					receiverID = elevatorID
					fmt.Println("distributeExternalOrder elevID: ", elevatorID)
					lowestCost = elevatorCostMap[receiverID]

					fmt.Print("Elevator: ", elevatorID)				
					fmt.Println("has to current lowest cost")
					fmt.Println("The cost is ", lowestCost)		
				}
			}
			fmt.Print("Elevator: ", receiverID)
			fmt.Println(" got the order ", sendMessage.Order)
			Claim(sendMessage.Order, receiverID)
		}else if(counter != dataTypes.NumberOfElevatorsConnected && sendMessage.WhichExternalPanelPressed == dataTypes.ElevatorID){ 
			fmt.Println("Running claim: ") 
			Claim(sendMessage.Order, dataTypes.ElevatorID)	
		}
	}
}

func AddExternalOrderLocallyToQueue() {
	var receivedMessage dataTypes.Message
	var receiverID int
	order := make([]int, 2)
	for {
		receivedMessage = <-AddOrderInternalCh
		order = receivedMessage.Order
		receiverID = receivedMessage.ID
		ExternalQueue[order[0]][order[1]] = receiverID
		//network.UDPCh <- dataTypes.receivedMessage{Head: "table", Order: order, Table: ExternalQueue} //JUST TO SEE WHAT WE GOT
		fmt.Println("new external queue Internally: ", ExternalQueue)
		UpdateLightCh <- "external"
		time.Sleep(25 * time.Millisecond)
	}
}


func AddExternalOrderFromNetworkToQueue() {
	var receivedMessage dataTypes.Message
	var receiverID int
	order := make([]int, 2)
	for {
		receivedMessage = <-network.AddOrderCh
		order = receivedMessage.Order
		receiverID = receivedMessage.ID
		ExternalQueue[order[0]][order[1]] = receiverID
		//network.UDPCh <- dataTypes.Message{Head: "table", Order: order, Table: ExternalQueue} //JUST TO SEE WHAT WE GOT
		fmt.Println("new external queue NETWORK: ", ExternalQueue)
		UpdateLightCh <- "external"
		time.Sleep(25 * time.Millisecond)
	}
}






func RemoveExternalOrderDoneOverNetwork() {
	order := make([]int, 2)
	for {
		//fmt.Println("her kommer vi i removeExternalOrder()")
		order = <-network.RemoveOrderCh
		fmt.Println("order to be removed is ", order)
		ExternalQueue[order[0]][order[1]] = 0
		//network.UDPCh <- dataTypes.Message{Head: "table", Order: order, Table: ExternalQueue}
		fmt.Println("new external queue removed over network: ", ExternalQueue)
		UpdateLightCh <- "external"
		time.Sleep(25 * time.Millisecond)
	}
}

func RemoveExternalOrderDoneLocally(){
	order := make([]int, 2)
	for {
		//fmt.Println("her kommer vi i removeExternalOrder()")
		order = <-RemoveOrderLocalCh
		fmt.Println("order to be removed is ", order)
		ExternalQueue[order[0]][order[1]] = 0
		//network.UDPCh <- dataTypes.Message{Head: "table", Order: order, Table: ExternalQueue}
		fmt.Println("new external queue removed locally: ", ExternalQueue)
		UpdateLightCh <- "external"
		time.Sleep(25 * time.Millisecond)
	}
}









func DeleteOrderDone() {
	sendOrderLocally := make([]int, 2)
	sendMessageOverNetwork := dataTypes.Message{Head: "removeOrder"}
	floor := GetCurrentFloor()
	dir := GetOrderDirection()
	if(floor != -1){
		sendOrderLocally = []int{floor, dir}
		sendMessageOverNetwork.Order = []int{floor, dir}
	
		fmt.Println("DeleteOrderDone() Internal Queue: ", InternalQueue, "\t floor: ", floor)
		InternalQueue[floor] = 0
		network.UDPCh <- sendMessageOverNetwork
		RemoveOrderLocalCh <- sendOrderLocally
		UpdateLightCh <- "internal"
		//UpdateLightCh <- "external" //skjer ikke dette internt i RemoveExternalOrder?
	}
}

// Checks the floor the elevator is on for orders
func CheckCurrentFloorForOrders() bool {
	currentFloor := driver.ElevGetFloorSensorSignal()
	dir := GetOrderDirection()
	//fmt.Println("CheckCurrentFloorForOrders motor dir: ", dir)
	var oppDir = GetOppositeDirection(dir)

	//ChangeDirectionAtTopOrBottomFloor(currentFloor)

	if currentFloor == -1 {
		return false
	}

	if currentFloor != -1 && currentFloor < dataTypes.N_FLOORS {
		if InternalQueue[currentFloor] == 1 {
			return true
		} else if ExternalQueue[currentFloor][dir] == dataTypes.ElevatorID {
			return true
		} else if CheckForOrdersInSameDirection() == false && ExternalQueue[currentFloor][oppDir] == dataTypes.ElevatorID {
			ChangeOrderDirection(oppDir)
			return true
		}
	}
	return false
}

func GetOppositeDirection(dir int) int{	//help func
	if dir == UP {
		return DOWN
	} else {
		return UP
	}
}

func ChangeDirectionAtTopOrBottomFloor(currentFloor int){ //help func

	if currentFloor == 0 && GetOrderDirection() == DOWN {
		ChangeOrderDirection(UP)
		//fmt.Println("ChangeDirectionAtTopOrBottomFloor bottomFloor Elev ID: " , dataTypes.ElevatorID)
	} else if currentFloor == (dataTypes.N_FLOORS - 1) && GetOrderDirection() == UP {
		ChangeOrderDirection(DOWN)
		//fmt.Println("ChangeDirectionAtTopOrBottomFloor topFloor Elev ID: " , dataTypes.ElevatorID)
	}
}


func CheckForOrdersInSameDirection() bool {
	currentFloor := GetCurrentFloor()
	dir := GetOrderDirection()

	if currentFloor != -1 {
		switch dir {
		case UP:
			for i := currentFloor; i < dataTypes.N_FLOORS; i++ {
				if ExternalQueue[i][UP] == dataTypes.ElevatorID {
					return true
				}
			}
			for j := dataTypes.N_FLOORS - 1; j > currentFloor; j-- {
				if ExternalQueue[j][DOWN] == dataTypes.ElevatorID {
					return true
				}
			}
		case DOWN:
			for k := currentFloor; k >= 0; k-- {
				if ExternalQueue[k][DOWN] == dataTypes.ElevatorID {
					return true
				}
			}
			for l := 0; l < currentFloor; l++ {
				if ExternalQueue[l][UP] == dataTypes.ElevatorID {
					return true
				}
			}
		}
	}
	return false
}



func GetCurrentFloor() int {
	return driver.ElevGetFloorSensorSignal() 
}

func CheckOtherFloors() int {
	currentFloor := GetCurrentFloor()
	dir := GetOrderDirection()


	ChangeDirectionAtTopOrBottomFloor(currentFloor)


	switch dir {
		case UP:
			for floor := currentFloor; floor < dataTypes.N_FLOORS; floor++ {
				if floor != currentFloor {
					if InternalQueue[floor] == 1 || ExternalQueue[floor][UP] == dataTypes.ElevatorID {
						return floor
					} else if ExternalQueue[floor][DOWN] == dataTypes.ElevatorID {
						return floor
					}
				}
			}
			for floor := currentFloor; floor >= 0; floor-- {
				if floor != currentFloor {
					if InternalQueue[floor] == 1 {
						ChangeOrderDirection(DOWN)
						return floor
					}
				}
			}
			ChangeOrderDirection(DOWN)

		case DOWN:
			for floor := currentFloor; floor >= 0; floor-- {
				if floor != currentFloor {
					if InternalQueue[floor] == 1 || ExternalQueue[floor][DOWN] == dataTypes.ElevatorID {
						return floor
					} else if ExternalQueue[floor][UP] == dataTypes.ElevatorID {
						return floor
					}
				}
			}
			for floor := currentFloor; floor < dataTypes.N_FLOORS; floor++ {
				if floor != currentFloor {
					if InternalQueue[floor] == 1 {
						ChangeOrderDirection(UP)
						return floor
					}
				}
			}
			ChangeOrderDirection(UP)
	}
	return -1
}

func ChangeOrderDirection(dir int) {
	Direction = dir
}

func GetOrderDirection() int {
	return Direction
}

func PrintOrderDirection() {
	for {
		time.Sleep(3000 * time.Millisecond)
		dir := GetOrderDirection()
		switch dir {
		case UP:
			fmt.Println("Order direction: UP")
		case DOWN:
			fmt.Println("Order direction: DOWN")
		}
	}
}

func FindDirection() int {
	var diff int
	currentFloor := driver.ElevGetFloorSensorSignal()
	if currentFloor != -1 && currentFloor < dataTypes.N_FLOORS && CheckOtherFloors() != -1 {
		diff = currentFloor - CheckOtherFloors()
	}
	if diff > 0 {
		return DOWN
	} else if diff < 0 {
		return UP
	} else {
		return -1
	}
}

func PrintTables() {
	for {
		time.Sleep(2*3000 * time.Millisecond)
		fmt.Println("Internal:", InternalQueue)
		fmt.Println("External:", ExternalQueue)
	}
}

func Redistribute() {
	var message dataTypes.Message
	var alive bool
	for {
		time.Sleep(1000 * time.Millisecond)
		for i:=0; i < dataTypes.N_FLOORS; i++ {
			for j := 0; j < 2; j++ {
				if ExternalQueue[i][j] != 0 { //Some Elevator is handling a order
					_, alive = network.PeerMap.M[ExternalQueue[i][j]]
					if !alive {
						fmt.Print(ExternalQueue[i][j])
						fmt.Println(" is probably dead") 
						
						message = dataTypes.Message{Head: "removeOrder", Order: []int{i, j}}
						network.UDPCh <- message

						time.Sleep(20 * time.Millisecond)

						message = dataTypes.Message{Head: "order", Order: []int{i, j}}
						network.UDPCh <- message
					}
				}
			}
		}
	}
}


// Reads from a backup file and sees if there are some orders left after the last elevator disconnected from the server.
func ReadBackupFile() {
	b, err := ioutil.ReadFile("backup.txt")
	if err != nil {
		panic(err)
	}
	internal := strings.Split(string(b), "")
	for i := 0; i < dataTypes.N_FLOORS; i++ {
		InternalQueue[i], _ = strconv.Atoi(internal[i])
	}
	time.Sleep(25 * time.Millisecond)
	UpdateLightCh <- "internal"
}

func WriteFile() {
	fmt.Println("Backup skrevet til fil")
	msg := strconv.Itoa(InternalQueue[0]) + strconv.Itoa(InternalQueue[1]) + strconv.Itoa(InternalQueue[2]) + strconv.Itoa(InternalQueue[3])
	buf := []byte(msg)
	_ = ioutil.WriteFile("backup.txt", buf, 0644)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
