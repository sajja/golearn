package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/sajja/golearn/statemachine"
)

// type ExampleTestSuite struct {
// suite.Suite
// VariableThatShouldStartAtFive int
// }

func Test_createStateMachine_should_give_ReadyState(t *testing.T) {
	now := time.Now().Unix()
	sm := statemachine.NewStateMachine("d1", nil)
	assert.Equal(t, sm.CurrentState.Name(), "ReadyState", "Initial state should be ReadyState")
	// assert.GreaterOrEqual(t, sm.Last, now, "")
	logic := sm.Last >= now //Pathetic way of doing GreaterThanOrEqual. testify/assert not seem to work
	assert.True(t, logic, "Time should be valid")
}

func Test_state_transistion(t *testing.T) {
	sm := statemachine.NewStateMachine("d1", nil)
	event := statemachine.CreateVoiceCommand("Computer")
	sm.CalculateState(&event)
	assert.Equal(t, "ActiveState", sm.CurrentState.Name(), "Incorrect state")
}

type DummyState struct{}

func (s DummyState) Eval(sm *statemachine.StateMachine, commands []string) {
}

func (s DummyState) Name() string {
	return "DummyState"
}

func (s DummyState) Print() {
}

func customTransfom(sm *statemachine.StateMachine, commands []string) {
	sm.CurrentState = DummyState{}
}

func Test_engine_support_external_state_transform_function(t *testing.T) {
	sm := statemachine.NewStateMachine("d1", customTransfom)
	assert.Equal(t, sm.CurrentState.Name(), "ReadyState", "Initial state should be ReadyState")
	event := statemachine.CreateVoiceCommand("Computer")
	sm.CalculateState(&event)
	assert.Equal(t, "DummyState", sm.CurrentState.Name(), "Incorrect state")
}
