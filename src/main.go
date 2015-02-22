package main

import (
	

	//"fmt"
	."./driver"
	//."./network"
)

	
func main(){
	IO_init()
	for{
		currFloor := Elev_get_floor_sensor_signal()
		if(currFloor==3){
			Elev_set_motor_direction(-1)
		}else if(currFloor==0){
			Elev_set_motor_direction(1)
		}
		if(currFloor!=-1){
			Elev_set_floor_indicator(currFloor)
		}
	}
}
