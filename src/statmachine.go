package elevator 

type State_t int

const { 
	UNDEFINED 		State_t = 0
	IDLE			State_t = 1
	MOVING			State_t = 2
	AT_FLOOR		State_t = 3
	STOP_AT_FLOOR		state_t = 4
	STOP_BETWEEN_FLOORS	state_t = 5
}


func RunElevator(state state_t) state_t {
	switch (state){
		case UNDEFINED:
			ElevInit()
			return IDLE
		case IDLE:
			if( stopPushed() ) {
				return STOP_AT_FLOOR			
			}else if( QueueNotEmpty() ){
				ElevSetMotorDirection(200);
				return Moving
			} 		
		case MOVING:
			return Moving 

/*
		case AT_FLOOR

		case STOP_AT_FLOOR


		case STOP_BETWEEN_FLOORS		*/
		default:
	}
}
