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

func worker(workerID int, doTodo chan<- DoTodo, product chan<- Product, displayWorkerStats <-chan bool, patient bool, restore chan<- Task) {

	// https://tour.golang.org/flowcontrol/4
	// Never ending for loop and start of the jobs done counter
	doneRepairs := 0
	for {
		select {
		case <-displayWorkerStats:
			fmt.Println("Worker #", workerID, "[Patient: ", patient, "] Tasks done: ", doneRepairs)
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

		// Take action accordingly to type of the task
		switch task.action {
		case "+":
			answer = useMachines(task, patient, addingMachines, workerID)
		case "*":
			answer = useMachines(task, patient, multiplayingMachines, workerID)
		}

		// Machine was corrupted and provided wrong answer. Worker puts the task back to Magazine
		if answer == ErrorSignal {
			if !Silent {
				fmt.Println("[Err] \t #", workerID, "\tdetected wrong answer. Restores task [", task.arg1, task.action, task.arg2, "]")
			}
			restore <- task
			continue
		}

		// Solve the task and update stats
		product <- Product{answer: answer}
		doneRepairs = doneRepairs + 1

		if !Silent {
			fmt.Println("[Job] \t #", workerID, "\tWorker solved a task. Answer: ", answer)
		}

		// Sleep means "performance" of workers - defined in config
		time.Sleep(time.Duration(WorkPerf) * time.Second)
	}
}

func useMachines(t Task, behavior bool, machines []Machine, workerID int) int {

	// If [Patient] else [Impatient]
	if behavior {
		// Choose machine
		machineID := rand.Intn(MachAddNum)

		// Wait for machine to be ready
		<-machines[machineID].ready

		// Use Machine
		machineResponse := make(chan int)
		machines[machineID].todo <- MachineTask{t, machineResponse}

		// Wait and check result
		result := <-machineResponse
		if result == ErrorSignal {
			repairTodos <- RepairMach{t.action, machineID}
			return ErrorSignal
		}

		return result
	} else {
		// Move around machines (in mod%numMachines until return)
		for machineID := 0; ; machineID = (machineID + 1) % len(machines) {

			machineResponse := make(chan int)
			select {
			// If Machine is ready to use
			case <-machines[machineID].ready:
				machines[machineID].todo <- MachineTask{t, machineResponse}
				result := <-machineResponse

				if result == ErrorSignal {
					repairTodos <- RepairMach{t.action, machineID}
					return ErrorSignal
				}

				return result
			// If Machine isn't ready after worker inpatience time
			case <-time.After(time.Duration(WorkPat) * time.Second):
				break
			}
		}
	}
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

		// Receive an item
		answer := <-buyProduct.response
		if !Silent {
			fmt.Println("[Shop] \t #", customerID, "\tCustomer bought ", answer)
		}

		time.Sleep(time.Duration(CustPerf) * time.Second)
	}
}

func repairman(repairmanID int, repairRequest chan<- RepairNeed, doneRepair chan<- RepairMach) {

	for {
		// Await for information about corrupted machine
		damagedMachines := make(chan RepairMach)
		repairRequest <- RepairNeed{damagedMachines}
		fixMachine := <-damagedMachines

		// Proceed to fix it after signal with given performance
		time.Sleep(time.Duration(RepairWorkPerf) * time.Second)
		if fixMachine.model == "+" {
			addingMachines[fixMachine.id].fix <- true
		} else {
			multiplayingMachines[fixMachine.id].fix <- true
		}

		// Notice the Facility about finished repair so it can be crossed out of the list to fix
		doneRepair <- RepairMach{fixMachine.model, fixMachine.id}

		if !Silent {
			fmt.Println("[Rep] \t #", repairmanID, "\tRepairman repaired machine #", fixMachine.id, "( of type: ", fixMachine.model, ")")
		}
	}
}
