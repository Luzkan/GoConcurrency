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
	todo   *Task
	status chan bool
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