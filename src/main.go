package main

import (
	
	//"time"
	"fmt"
	//."./driver"
	."./network"
)

	
func main(){
	fmt.Println("Do you want to run elevator 1,2 or 3?")
	var toRun string
	fmt.Scanf("%s",&toRun)
	if(toRun=="1"){
		RunElev1()
	}
	if(toRun=="2"){
		RunElev2()
	}
	if(toRun=="3"){
		RunElev3()
	}
}
