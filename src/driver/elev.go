package driver
 
//import "C"

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
/*
button_channel_matrix 	:= [N_FLOORS][N_BUTTONS]int{ 					//KOM HER. FIKSE DISSE OG DE TO NEDERSTE FUNKSJONENE. 
		{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1}, 
		{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
		{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
		{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
}

lamp_channel_matrix 	:= [N_FLOORS][N_BUTTONS]int{
{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
}


func Elev_init() int {
	
	// Initialize the hardware
	if(!Io_init()){
		return 0
	}
	// Reset all floor button lamps
	for i := 0; i<N_FLOORS; i++ {
		if(i!=0){
			Elev_set_button_lamp(BUTTON_CALL_DOWN,i,0)
		}
		if(i != N_FLOORS-1){
			Elev_set_button_lamp(BUTTON_CALL_UP,i,0)
			Elev_set_button_lamp(BUTTON_COMMAND,i,0)
		}
	}
	// Reset stop lamp, door open lamp, set floor indicator to ground floor
	Elev_set_stop_lamp(0)
	Elev_set_door_open_lamp(0)
	Elev_set_floor_indicator(0)
	// Return success
	return 1
}

func Elev_set_motor_direction(motorDirn int){
	if(motorDirn==0){
		IO_write_analog(MOTOR, 0)
	}else if (motorDirn>0) {
		IO_clear_bit(MOTORDIR)
		IO_write_analog(MOTOR,2800)
	}else if(motorDirn<0){
		IO_set_bit(MOTORDIR)
		IO_write_analog(MOTOR,2800)
	}
}

func Elev_set_door_open_lamp(value int){
	if(value){
		IO_set_bit(LIGHT_DOOR_OPEN)
	}else{
		IO_clear_bit(LIGHT_DOOR_OPEN)
	}
}

func Elev_get_obstruction_signal() int{
	return IO_read_bit(OBSTRUCTION)
}

func Elev_get_stop_signal() int{
	return IO_read_bit(STOP)
}

func Elev_set_stop_lamp(value int){
	if(value){
		IO_set_bit(LIGHT_STOP)
	}else{
		IO_clear_bit(LIGHT_STOP)
	}
}

func Elev_get_floor_sensor_signal() int{
	if(IO_read_bit(SENSOR_FLOOR1)){
		return 0
	}else if(IO_read_bit(SENSOR_FLOOR2)){
		return 1
	}else if(IO_read_bit(SENSOR_FLOOR3)){
		return 2
	}else if(IO_read_bit(SENSOR_FLOOR4)){
		return 3
	}else{
		return -1
	}
}

func Elev_set_floor_indicator(floor int){
	if(floor & 0x02){
		IO_set_bit(LIGHT_FLOOR_IND1)
	}else{
		IO_clear_bit(LIGHT_FLOOR_IND1)
	}
	if(floor & 0x01){
		IO_set_bit(LIGHT_FLOOR_IND2)
	}else{
		IO_clear_bit(LIGHT_FLOOR_IND2)
	}
}

func Elev_get_button_signal(button int, floor int) int{
	if(IO_read_bit(???? see elev.c)){
		return 1
	}else{
		return 0
	}
}

func Elev_set_button_lamp(button int, floor int, value int){
	if(value){
		IO_set_bit(??? see elev.c)
	}else{
		IO_clear_bit(??? see elev.c)
	}
}

*/



