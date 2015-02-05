package main
//export GOPATH="/home/bjorkv/Desktop/go"

import "Driver_Test"

func main(){
	C.elev_set_motor_direction(C.DIRN_UP);
}
