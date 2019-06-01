package main

// Task generated by boss (ex: 93195423 + 2138217)
type Task struct {
	arg1   int
	action string
	arg2   int
	answer int
}

// Worker can check Machine status and place a todo to calculate
type Machine struct {
	id     int
	todo   chan MachineTask
	ready  chan bool
	status bool
	fix    chan bool
}

// Task in a machine with printed result
type MachineTask struct {
	taskInMachine Task
	result        chan int
}

// Customer shop request
type Buy struct {
	response chan Product
}

// Worker store
type Product struct {
	answer int
}

// Worker todo request
type DoTodo struct {
	response chan Task
}

// What Machine needs to be repaired
type RepairMach struct {
	model string
	id    int
}

// Machine fix request
type RepairNeed struct {
	response chan RepairMach
}

// Global variables
var repairTodos chan RepairMach
var doneRepairs chan RepairMach