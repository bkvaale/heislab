package order
import "fmt" 
//var GlobalTable []int
//var InternalTable []int

var (
	OrderQueue[4][3] int // f,b
)

/*
func CheckError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}
*/





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
}

func AddToQueue(floor int, button int){
	OrderQueue[floor][button] = 1
}

func DeleteFromQueue(floor int, button int){
	OrderQueue[floor][button] = 0;
}
