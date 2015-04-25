package order

import (
	"../network"
	"../driver"
	"../dataTypes"
	"fmt"
	"math"
	//"time"
)


// Calculates the cost with concern on direction and how far the order is from the elevator (HAVE TO CHANGE PARAMETERS)
func SendCostOverNetwork(){
	for {
		sendMessage := dataTypes.Message{Head: "costCalculatedNetwork"}


		//fmt.Println("Vi kommer her i SendCostOverNetwork, and still waiting for the channel")
	
		receivedMessage :=<- network.CalculateCostNetworkCh
	

		//fmt.Println("Vi kommer her i SendCostOverNetwork, and GOT information from channel")
		sendMessage.Cost = CalculateCost(receivedMessage.Order)
		sendMessage.Order = receivedMessage.Order
		sendMessage.ID = dataTypes.ElevatorID
		sendMessage.WhichExternalPanelPressed =  receivedMessage.WhichExternalPanelPressed
		fmt.Println("Elevator: ", sendMessage.ID, " Cost calculated and sent. cost:", sendMessage.Cost)
		network.UDPCh <- sendMessage
	}
}

func SendCostLocally(){
	for{
		var sendMessage dataTypes.Message

		//fmt.Println("Vi kommer her i SendCostInternal, and still waiting for the channel")
		receivedMessage :=<- CalculateCostInternalCh
		//fmt.Println("Vi kommer her i CalculateAndSendCost, and GOT information from channel")
		sendMessage.Cost = CalculateCost(receivedMessage.Order)
		sendMessage.Order = receivedMessage.Order
		sendMessage.ID = dataTypes.ElevatorID
		sendMessage.WhichExternalPanelPressed =  receivedMessage.WhichExternalPanelPressed
		fmt.Println("Elevator: ", sendMessage.ID, " Cost calculated and sent INTERNAL. cost:", sendMessage.Cost)
		CostOfOrderInternalCh <- sendMessage
	}
}

func CalculateCost(order []int) int{
	cost := 0
	var state int
	var motorDir int

	punishWrongDirection:= 2 
	punishFloorDifferenceMultiplier := 4 

	//direction := GetOrderDirection()	can be included in switch statement
	currentFloor := GetCurrentFloor()
	
	

	if driver.ElevGetFloorSensorSignal() == -1 {
		motorDir = driver.IoReadBit(driver.MOTORDIR)
		if motorDir == 1 {
			state = UP
		} else if motorDir == 0 {
			state = DOWN
		}
	}	

	switch state {
		case UP:
			if order[1] == DOWN {
				cost = cost + punishWrongDirection
			}
			break
		case DOWN:
			if order[1] == UP {
				cost = cost + punishWrongDirection
			}
			break
		default:
			break		
	}
	cost = cost + punishFloorDifferenceMultiplier*(int(math.Abs(float64(currentFloor - order[0]))))
	return cost
}


func Claim(order []int, receiverID int) {
	var sendMessageLocally dataTypes.Message
	sendMessageOverNetwork := dataTypes.Message{Head: "addorder"}
		
	

	sendMessageLocally.Order = order
	sendMessageOverNetwork.Order = order
	
	sendMessageLocally.ID = receiverID
	sendMessageOverNetwork.ID = receiverID

	fmt.Println("Claimed(): ", order, " by elev: ", sendMessageLocally.ID)

	AddOrderInternalCh <- sendMessageLocally
	network.UDPCh <- sendMessageOverNetwork	
}

func ContainsAll(carts []int) bool {
	for _, t := range carts {
		if t == 0 {
			return false
		}
	}
	return true
}
