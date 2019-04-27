package main

import (
	"fmt"
	"time"
)

// Making machines as global variables (used in "main" to create and in "people" to use)
var addingMachines []chan Machine
var multiplayingMachines []chan Machine

// Set up adding machines
func setAddingMachines() []chan Machine {

	channels := make([]chan Machine, MachAddNum)
	for i := 0; i < MachAddNum; i++ {
		channels[i] = make(chan Machine)
		go addingMachine(i, channels[i])
	}

	return channels
}

// Set up multiplaying machines
func setMultiplayingMachines() []chan Machine {

	channels := make([]chan Machine, MachMulNum)
	for i := 0; i < MachMulNum; i++ {
		channels[i] = make(chan Machine)
		go multiplayingMachine(i, channels[i])
	}

	return channels
}

// Make those previously set up adding machines running
func addingMachine(machineID int, addTasks <-chan Machine) {

	if !Silent {
		fmt.Println("[Mach]   #", machineID, "\tadd machine setten up.")
	}

	// Machine works when it sees a task has been given by a worker
	for t := range addTasks {

		if !Silent {
			fmt.Println("[Mach]   #", machineID, "\tadd machine started calculating: [", t.todo.arg1, t.todo.arg2, "]")
		}

		time.Sleep(time.Duration(MachPerf) * time.Second)
		t.status <- true
		t.todo.answer = t.todo.arg1 + t.todo.arg2
	}
}

// Make those previously set up multiplaying machines running
func multiplayingMachine(machineID int, mulTasks <-chan Machine) {

	if !Silent {
		fmt.Println("[Mach]   #", machineID, "\tmul machine setten up.")
	}

	// Machine works when it sees a task has been given by a worker
	for t := range mulTasks {

		if !Silent {
			fmt.Println("[Mach]   #", machineID, "\tmul machine started calculating.")
		}

		time.Sleep(time.Duration(MachPerf) * time.Second)
		t.status <- true
		t.todo.answer = t.todo.arg1 * t.todo.arg2
	}
}
