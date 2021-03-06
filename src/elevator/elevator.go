package elevator

import (
	"../driver"
	"../order"
	"../dataTypes"
	"fmt"
	"os"
	"os/signal"
	"time"
)

const (
	UP   = 0
	DOWN = 1
	IDLE = 2
	OPEN = 3

	ON    = 1
	OFF   = 0
	SPEED = 300
)

var (

	// Global channels sjekke om alle er i bruk
	elevatorDirection int

	// Local channels
	doorTimerStartCh = make(chan bool)
	doorTimerDoneCh  = make(chan bool)
	idleCh           = make(chan bool)
	openCh           = make(chan bool)
	downCh           = make(chan bool)
	upCh             = make(chan bool)
	osChan           chan os.Signal
)

// Starts the go routines and initializing of the elevator
func Run() {

	done := make(chan bool)
	// Initialization
	driver.ElevInit()

	go FloorLights()
	go DoorTimer()
	elevatorDirection = DOWN // ??
	/*for driver.ElevGetFloorSensorSignal() == -1 {
		driver.ElevSetSpeed(-SPEED)
	}
	driver.ElevSetSpeed(0)*/
	go Idle()
	go Open()
	go Down()
	go Up()
	go DoorSafety()
	idleCh <- true
	<-done
}


func Idle() {
	for {
		<-idleCh
		fmt.Println("Idle state entered")
		driver.ElevSetDoorOpenLamp(OFF)
		driver.ElevSetSpeed(0)
		for {
			time.Sleep(10*time.Millisecond)
			if order.CheckCurrentFloorForOrders() {
				openCh <- true
				break
			} else if order.FindDirection() == DOWN {
				downCh <- true
				break
			} else if order.FindDirection() == UP {
				upCh <- true
				break
			}
		}
	}
}

func Open() {
	for {
		<-openCh
		fmt.Println("Open state entered")
		driver.ElevSetSpeed(0)
		order.DeleteOrderDone()
		driver.ElevSetDoorOpenLamp(ON)
		doorTimerStartCh <- true
		<-doorTimerDoneCh	
		driver.ElevSetDoorOpenLamp(OFF)
		idleCh <- true
	}
}

func Down() {
	for {
		<-downCh
		fmt.Println("Down state entered")
		driver.ElevSetSpeed(-SPEED)
		elevatorDirection = DOWN
		for {
			if order.CheckCurrentFloorForOrders() {
				openCh <- true
				break
			} else if Safety() {
				idleCh <- true
				break
			}
		}
	}
}

func Up() {
	for {
		<-upCh
		fmt.Println("Up state entered")
		elevatorDirection = UP
		driver.ElevSetSpeed(SPEED)
		for {
			if order.CheckCurrentFloorForOrders() {
				openCh <- true
				break
			} else if Safety() {
				idleCh <- true
				break
			}
		}
	}
}

func FloorLights() {
	for {
		time.Sleep(1 * time.Millisecond)
		driver.ElevSetFloorIndicator(driver.ElevGetFloorSensorSignal())
	}
}

func Safety() bool {
	if driver.ElevGetFloorSensorSignal() == 0 && !order.CheckCurrentFloorForOrders() && elevatorDirection == DOWN {
		return true
	} else if driver.ElevGetFloorSensorSignal() == (dataTypes.N_FLOORS-1) && !(order.CheckCurrentFloorForOrders()) && elevatorDirection == UP {
		return true
	}
	return false
}

func DoorTimer() {
	for {
		<-doorTimerStartCh
		//driver.ElevSetDoorOpenLamp(ON)
		time.Sleep(3000 * time.Millisecond)
		/*for driver.IoReadBit(driver.OBSTRUCTION) == 1 {
			//doorTimerStartCh <- true
			time.Sleep(2000*time.Millisecond)
		}*/
		doorTimerDoneCh <- true
	}
}

func OsTest() {
	osChan := make(chan os.Signal, 1)
	signal.Notify(osChan, os.Interrupt)
	<-osChan
	order.WriteFile()
	fmt.Println("Programmet er blitt avsluttet")
	driver.ElevSetSpeed(0)
	//stop elevator her...
	os.Exit(1)
}

func DoorSafety() {
	for {
		if driver.ElevGetFloorSensorSignal() == -1 && driver.IoReadBit(driver.DOOR_OPEN) == 1 {
			driver.ElevSetDoorOpenLamp(OFF)
			doorTimerDoneCh <- true
		}
	}
}
