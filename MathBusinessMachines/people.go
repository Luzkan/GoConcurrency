package main

import (
	"fmt"
	"math/rand"
	"time"
)

func boss(tasks chan<- Task) {

	for {
		// https://golang.org/pkg/math/rand/
		// Generating vars for new task
		arg1 := rand.Intn(100)
		arg2 := rand.Intn(100)
		actions := [2]string{"+", "*"}

		// Sticking them together
		t := Task{arg1, actions[rand.Intn(len(actions))], arg2, 0}

		// Proceed to store new task
		tasks <- t

		// Sleep means "performance" of the boss - defined in config
		time.Sleep(time.Duration(BossPerf) * time.Second)

		// Message from Boss
		if !Silent {
			fmt.Println("[Boss]  \tNew task in the Magazine: [", t.arg1, t.action, t.arg2, "]")
		}
	}
}

func worker(workerID int, doTodo chan<- DoTodo, product chan<- Product, displayWorkerStats <-chan bool, patient bool) {

	// https://tour.golang.org/flowcontrol/4
	// Never ending for loop and start of the jobs done counter
	jobsDone := 0
	for {
		select {
		case <-displayWorkerStats:
			fmt.Println("Worker #", workerID, "[Patient: ", patient, "] Tasks done: ", jobsDone)
		default:
		}

		// Get first task, then cut it out from the magazine
		request := DoTodo{response: make(chan Task)}
		doTodo <- request
		task := <-request.response

		// Use Machines based on behaviour and task
		var answer int

		if !Silent {
			fmt.Println("[Job] \t #", workerID, "\tis heading to machines with task [", task.arg1, task.action, task.arg2, "]")
		}

		switch task.action {
		case "+":
			answer = useMachines(&task, patient, addingMachines, workerID)
		case "*":
			answer = useMachines(&task, patient, multiplayingMachines, workerID)
		}

		// Solve the task and update stats
		product <- Product{answer: answer}
		jobsDone = jobsDone + 1

		if !Silent {
			fmt.Println("[Job] \t #", workerID, "\tWorker solved a task. Answer: ", answer)
		}

		// Sleep means "performance" of workers - defined in config
		time.Sleep(time.Duration(WorkPerf) * time.Second)
	}
}

func useMachines(t *Task, behavior bool, machines []chan Machine, workerID int) int {

	// If [Patient] else [Impatient]
	if behavior {
		// Choose machine
		machineID := rand.Intn(MachAddNum)
		machine := machines[machineID]
		machineResponse := make(chan bool)

		// Go into que
		machine <- Machine{t, machineResponse}
		for {
			select {
			case <-machineResponse:
				if !Silent {
					fmt.Println("[Que] \t #", workerID, "\tWorker uses machine after being patient. Machine #", machineID)
				}
				return t.answer
			}
		}
	} else {
		// Move around machines
		for _, machine := range machines {
			machineResponse := make(chan bool)
			machine <- Machine{t, machineResponse}

			// Run around after two seconds if no success
			for {
				select {
				case <-machineResponse:
					if !Silent {
						fmt.Println("[Que] \t #", workerID, "\tWorker used machine after being impatient.")
					}
					return t.answer
				case <-time.After(2 * time.Second):
					break
				}
			}
		}
	}
	return t.answer
}

func solve(t Task) int {
	var answer = 0
	switch t.action {
	case "+":
		answer = t.arg1 + t.arg2
	case "-":
		answer = t.arg1 - t.arg2
	case "*":
		answer = t.arg1 * t.arg2
	}
	return answer
}

func customer(customerID int, boughts chan<- Buy) {

	for {
		buyProduct := Buy{response: make(chan Product)}
		boughts <- buyProduct
		// receive an item
		answer := <-buyProduct.response

		if !Silent {
			fmt.Println("[Shop] \t #", customerID, "\tCustomer bought ", answer)
		}

		time.Sleep(time.Duration(CustPerf) * time.Second)
	}
}
