package main

import (
	
	"fmt"
	."./driver"
	//."./network"
)
func main(){
	IO_init()
	//var matrix [3][2] int

	button_channel_matrix := [N_FLOORS][N_BUTTONS]int{ 
		{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1}, 
		{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
		{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
		{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
	}

	fmt.Println(button_channel_matrix)
	fmt.Println(BUTTON_UP1)
	IO_set_bit(LIGHT_COMMAND3) 
	IO_clear_bit(LIGHT_COMMAND2)
	fmt.Println("hello")
}
