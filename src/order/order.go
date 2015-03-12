package order

import(
	 "fmt"
	."../driver"
) 

var (
	OrderQueue[4][3] int // f,b
	//network.OrderQueue
)

func addOrdersToQueue(){
	for floor:= 0; floor < 4; floor++{ 
		for button:= 0; button < 3; button++{ //b
			if ( (floor == 0) && (button == BUTTON_CALL_DOWN) || (floor == N_FLOORS-1) && (button == BUTTON_CALL_UP) ){
				//Do nothing
			} else {
				buttonPressed := ElevGetButtonSignal (floor, button)
				if(buttonPressed>0){
					AddToQueue(floor, button)
					ElevSetButtonLamp(floor, button, 1)
					TestOrder()
					//time.Sleep(1000*time.Millisecond)
				}
			}
		}
	}		
}


func initializeQueue(arr[4][3] int){
	for i:= 0; i < 4; i++{ 
		for j:= 0; j < 3; j++{ //b
			arr[i][j] = 0
		}
	}		
}


func printQueue(arr[4][3] int){
	for i:= 0; i < 4; i++{ 
		for j:= 0; j < 3; j++{ //b
			fmt.Print(arr[i][j])
		}
	fmt.Println("\n")
	}
}

func TestOrder(){
	initializeQueue(OrderQueue)
	printQueue(OrderQueue)
	fmt.Println("----")
}

func AddToQueue(floor int, button int){
	OrderQueue[floor][button] = 1
}

func DeleteFromQueue(floor int, button int){
	OrderQueue[floor][button] = 0;
}
