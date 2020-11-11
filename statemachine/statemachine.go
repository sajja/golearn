package statemachine

import (
	"fmt"
	"strings"
	"time"
)

var stateMachines = make(map[string]*StateMachine)

func GetAllStateMachines() map[string]*StateMachine {
	//todo: Deep copy
	return stateMachines
}

const RESET_SECONDS = 3

type State interface {
	Eval(sm *StateMachine, commands []string)
	Print()
	Name() string
}

type ReadyState struct {
}

type ActiveState struct {
}

type LightManipulationPartialState struct {
}

type transform func(stateMachine *StateMachine, commands []string)

func defaultTransform(sm *StateMachine, commands []string) {
	println("Do nothing default transformer")
}

func evalCommand(sm *StateMachine, commands []string) {
	if len(commands) != 0 {
		tail := commands[1:]
		switch commands[0] {
		case "COMPUTER":
			sm.CurrentState = ActiveState{}
			sm.CurrentState.Print()
			evalCommand(sm, tail)
		case "LIGHT":
			if sm.CurrentState.Name() == "ActiveState" {
				sm.CurrentState = LightManipulationPartialState{}
				sm.CurrentState.Print()
				sm.CurrentState.Eval(sm, tail)
			}
		case "AC":
			if sm.CurrentState.Name() == "ActiveState" {
				sm.CurrentState = LightManipulationPartialState{}
				sm.CurrentState.Print()
				sm.CurrentState.Eval(sm, tail)
			}
		default:
			evalCommand(sm, tail)
		}
	}
}

func (s LightManipulationPartialState) Eval(sm *StateMachine, commands []string) {
	if len(commands) != 0 {
		switch commands[0] {
		case "ON":
			fmt.Println("Turning light ON")
			//mqtt publish
			sm.CurrentState = ReadyState{}
		case "OFF":
			fmt.Println("Turning light OFF")
			//mqtt publish
			sm.CurrentState = ReadyState{}
		default:
			fmt.Printf("Unknown command %s. Still listeninig", commands[0])
		}
	}

}

func (s LightManipulationPartialState) Name() string {
	return "LightManipulationPartialState"
}

func (s LightManipulationPartialState) Print() {
	fmt.Println("State: LightManipulationPartialState")
}

func (s ReadyState) Eval(sm *StateMachine, commands []string) {
	if sm.Transform != nil {
		sm.Transform(sm, commands)
	} else {
		evalCommand(sm, commands)
	}
}

func (s ReadyState) Name() string {
	return "ReadyState"
}

func (s ReadyState) Print() {
	fmt.Println("State: Ready")
}

func (s ActiveState) Eval(sm *StateMachine, commands []string) {

}

func (s ActiveState) Name() string {
	return "ActiveState"
}

func (s ActiveState) Print() {
	fmt.Println("State: Active, Listening to commands")
}

type StateMachine struct {
	id           string
	CurrentState State
	Transform    transform
	Last         int64
}

func NewStateMachine(id string, transform transform) *StateMachine {
	sm := StateMachine{id: id, CurrentState: ReadyState{}, Last: time.Now().Unix(), Transform: transform}
	return &sm
}

func GetStateMachine(deviceId string) *StateMachine {
	sm := stateMachines[deviceId]
	if sm == nil {
		sm = NewStateMachine(deviceId, nil)
		stateMachines[deviceId] = sm
	} else if time.Now().Unix()-sm.Last > int64(RESET_SECONDS) {
		sm.CurrentState = ReadyState{}
	}
	return sm
}

func PrintAllStateMachines() {
	for k, v := range stateMachines {
		fmt.Printf("Device: %s\t State: %s\n", k, v.CurrentState.Name())
	}
}

func (s *StateMachine) CalculateState(command Command) {
	s.CurrentState.Eval(s, strings.Split(command.GetCommand(), " "))
}

type Command interface {
	GetCommand() string
	SetCommand(command string)
}

type VoiceCommand struct {
	command string
}

func (vc *VoiceCommand) SetCommand(command string) {
	vc.command = strings.ToUpper(command)
}

func (vc VoiceCommand) GetCommand() string {
	return vc.command
}

func CreateVoiceCommand(cmd string) VoiceCommand {
	command := VoiceCommand{}
	command.SetCommand(cmd)
	return command
}

func TestStateMachine() {
	// reader := bufio.NewReader(os.Stdin)
	// fmt.Print("Enter your city: ")
	// city, _ := reader.ReadString('\n')
	// fmt.Print("You live in " + city)

	sm := NewStateMachine("d1", nil)
	sm.CurrentState.Print()
	fmt.Println("Sending a command : Computer")
	computer := VoiceCommand{}
	computer.SetCommand("Computer yo Light OFF")
	sm.CalculateState(&computer)
	// sm.CurrentState.Print()
}

func TestMultipleStateMachines() {
	GetStateMachine("d1")
	GetStateMachine("d2")
	GetStateMachine("d3")

	computer := VoiceCommand{}
	computer.SetCommand("computer")
	d4 := GetStateMachine("d4")
	d4.CalculateState(&computer)
	PrintAllStateMachines()

	// time.Sleep(4 * time.Second)
	// GetStateMachine("d4")
	// PrintAllStateMachines()
}
