package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Making machines as global variables (used in "main" to create and in "people" to use)
var addingMachines []Machine
var multiplayingMachines []Machine

// Set up adding machines
func setAddingMachines() []Machine {

	machines := make([]Machine, MachAddNum)
	//channels := make([]chan Machine, MachAddNum)
	for i := 0; i < MachAddNum; i++ {
		machines[i] = Machine{i, make(chan MachineTask), make(chan bool), true, make(chan bool)}
		go machines[i].addingMachine()
	}
	return machines
}

// Set up multiplaying machines
func setMultiplayingMachines() []Machine {

	machines := make([]Machine, MachMulNum)
	//channels := make([]chan Machine, MachMulNum)
	for i := 0; i < MachMulNum; i++ {
		machines[i] = Machine{i, make(chan MachineTask), make(chan bool), true, make(chan bool)}
		go machines[i].multiplayingMachine()
	}
	return machines
}

// Make those previously set up adding machines running
func (m *Machine) addingMachine() {

	if !Silent {
		fmt.Println("[Mach]   #", m.id, "\tadd machine setten up.")
	}

	for {
		// Backdoor
		select {
		case <-m.fix:
			m.status = true
		default:
		}

		if (rand.Float64() < MachFailure) && m.status {
			m.status = false
		}

		m.ready <- true
		t := <-m.todo

		if !Silent {
			fmt.Println("[Mach]   #", m.id, "\tadd machine started calculating: [", t.taskInMachine.arg1, t.taskInMachine.arg2, "]")
		}

		time.Sleep(time.Duration(MachPerf) * time.Second)

		// Machine works when it sees a task has been given by a worker
		if m.status {
			t.result <- t.taskInMachine.arg1 + t.taskInMachine.arg2
		} else {
			t.result <- ErrorSignal
		}
	}
}

// Make those previously set up multiplaying machines running
func (m *Machine) multiplayingMachine() {

	if !Silent {
		fmt.Println("[Mach]   #", m.id, "\tmul machine setten up.")
	}

	for {
		// Backdoor
		select {
		case <-m.fix:
			m.status = true
		default:
		}

		if (rand.Float64() < MachFailure) && m.status {
			m.status = false
		}

		m.ready <- true

		t := <-m.todo
		if !Silent {
			fmt.Println("[Mach]   #", m.id, "\tmul machine started calculating: [", t.taskInMachine.arg1, t.taskInMachine.arg2, "]")
		}
		time.Sleep(time.Duration(MachPerf) * time.Second)

		// Machine works when it sees a task has been given by a worker
		if m.status {
			t.result <- t.taskInMachine.arg1 * t.taskInMachine.arg2
		} else {
			t.result <- ErrorSignal
		}
	}
}
