package main

import (
	//"time"
	//"fmt"
	."./driver"
	//."./network"
	"./order"
)





	
func main(){
	ElevInit()
	for{
		order.addOrdersToQueue()
	}/*
	fmt.Println("How many elevators?")
	var numElev int
	fmt.Scanf("%d",&numElev)
	var toRun int
	if(numElev==2){
		fmt.Println("Do you want to run elevator 1 or 2?")
		fmt.Scanf("%d",&toRun)
	}else if(numElev==3){
		fmt.Println("Do you want to run elevator 1,2 or 3?")
		fmt.Scanf("%d",&toRun)
	}
	if(toRun==1){
		RunElev(toRun,numElev)
	}
	if(toRun==2){
		RunElev(toRun,numElev)
	}
	if(toRun==3){
		RunElev(toRun,numElev)
	}*/
}
/*
	ElevInit()
	
	for {
	addOrdersToQueue()

	}
*/	/*
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
	}*/
