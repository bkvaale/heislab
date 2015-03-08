package driver

const N_FLOORS 		= 4 
const N_BUTTONS 	= 3 

//elev_button_type_t
const BUTTON_CALL_UP 	= 0
const BUTTON_CALL_DOWN 	= 1
const BUTTON_COMMAND 	= 2


//elev_motor_direction
const DIRN_DOWN 	= -1
const DIRN_STOP 	= 0
const DIRN_UP	 	= 2


var (	
	lastMotorDirn int 
	lampChannelMatrix = [N_FLOORS][N_BUTTONS]int{
		{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
		{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
		{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
		{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
	}
	buttonChannelMatrix = [N_FLOORS][N_BUTTONS]int{
		{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
		{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
		{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
		{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
	}
)



func ElevInit() int {
	
	// Initialize the hardware
	if(IoInit() == 0){
		return 0
	}
	// Reset all floor button lamps
	for i := 0; i < N_FLOORS; i++ {
		if( i !=0 ){
			ElevSetButtonLamp(i, BUTTON_CALL_DOWN,0)
		}
		if( i != N_FLOORS-1 ){
			ElevSetButtonLamp(i,BUTTON_CALL_UP,0)
		}
		ElevSetButtonLamp(i,BUTTON_COMMAND,0)
	}
	// Reset stop lamp, door open lamp, set floor indicator to ground floor
	ElevSetStopLamp(0)
	ElevSetDoorOpenLamp(0)
	ElevSetFloorIndicator(0)

	lastMotorDirn = 0
	tempCurrentFloor := ElevGetFloorSensorSignal()
	for tempCurrentFloor == -1 {
		ElevSetMotorDirection(-300)
		tempCurrentFloor = ElevGetFloorSensorSignal()
	}
	ElevSetFloorIndicator(tempCurrentFloor)
	ElevSetMotorDirection(0)
	ElevInitLights()
	// Return success
	return 1
}



func ElevSetMotorDirection(speed int){
	if( speed==0 ){
		IoWriteAnalog(MOTOR, 0)
	}else if(speed >0 ) {
		IoClearBit(MOTORDIR)
		IoWriteAnalog(MOTOR,2800)
	}else if( speed<0 ){
		IoSetBit(MOTORDIR)
		IoWriteAnalog(MOTOR,2800)     //they had IOWriteAnalog(MOTOR, 2048+4*int(math.Abs(float64(speed))))
	}
}



func ElevSetDoorOpenLamp(value int){
	if(value > 0){
		IoSetBit(LIGHT_DOOR_OPEN)
	}else{
		IoClearBit(LIGHT_DOOR_OPEN)
	}
}



/*
func ElevGetObstructionSignal() int{
	return IoReadBit(OBSTRUCTION)
}
*/



func ElevGetStopSignal() int{
	return IoReadBit(STOP)
}



func ElevSetStopLamp(value int){
	if( value == 1 ){
		IoSetBit(LIGHT_STOP)
	}else{
		IoClearBit(LIGHT_STOP)
	}
}




func ElevGetFloorSensorSignal() int{
	if(IoReadBit(SENSOR_FLOOR1) == 1){ 
		return 0
	}else if(IoReadBit(SENSOR_FLOOR2) == 1){
		return 1
	}else if(IoReadBit(SENSOR_FLOOR3) == 1){
		return 2
	}else if(IoReadBit(SENSOR_FLOOR4) == 1){
		return 3
	}else{
		return -1
	}
}




func ElevSetFloorIndicator(floor int){
	switch floor {
	case 0:
		IoClearBit(LIGHT_FLOOR_IND1)
		IoClearBit(LIGHT_FLOOR_IND2)
	case 1:
		IoClearBit(LIGHT_FLOOR_IND1)
		IoSetBit(LIGHT_FLOOR_IND2)
	case 2:
		IoSetBit(LIGHT_FLOOR_IND1)
		IoClearBit(LIGHT_FLOOR_IND2)
	case 3:
		IoSetBit(LIGHT_FLOOR_IND1)
		IoSetBit(LIGHT_FLOOR_IND2)
	default:
	}
}




func ElevSetButtonLamp(floor int, button int, value int){
	if value == 1 {
		IoSetBit(lampChannelMatrix[floor][button])
	} else {
		IoClearBit(lampChannelMatrix[floor][button])
	}
}




func ElevGetButtonSignal(floor int, button int) int {
	if (IoReadBit(buttonChannelMatrix[floor][button]) != 0){
		return 1
	}else{
		return 0
	}
}




func ElevInitLights() {
	for i := 0; i < N_FLOORS; i++ {
		for j := 0; j < 3; j++ {
			IoClearBit(lampChannelMatrix[i][j])
		}
	}
}
