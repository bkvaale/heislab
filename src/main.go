package main

import (
	

	//"fmt"
	."./driver"
	//."./network"
)

	
func main(){
	IO_init()
	for(Elev_get_stop_signal()>0){
		Elev_set_motor_direction(1)
	}
}
