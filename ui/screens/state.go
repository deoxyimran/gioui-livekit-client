package screens

// Add your screen states here
type AppStateRoomScreen struct {
	MicOn          bool
	CameraOn       bool
	ChosenCameraId int
}

type AppStateMeetingScreen struct {
}

type StateManager struct {
	states map[Screen]any
}

func NewStateManager() StateManager {
	sm := StateManager{}
	sm.states = make(map[Screen]any)
	// Add states according to the structs defined
	sm.states[JOIN_MEETING] = &AppStateMeetingScreen{}
	sm.states[JOIN_ROOM] = &AppStateRoomScreen{}
	return sm
}

func (sm *StateManager) GetState(screen Screen) any {
	if state, exists := sm.states[screen]; exists {
		return state
	}
	return nil
}
