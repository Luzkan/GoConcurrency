package main

// Guards - Prezentacja str 143 o Adzie, uwaga nt. Golang

import (
	"fmt"
)

// Guards for tasks and solving todo
func taskGuard(cond bool, c <-chan Task) <-chan Task {
	if !cond {
		return nil
	}
	return c
}

func doTodoGuard(cond bool, c <-chan DoTodo) <-chan DoTodo {
	if !cond {
		return nil
	}
	return c
}

// Guards for products and buying answer
func productGuard(cond bool, c <-chan Product) <-chan Product {
	if !cond {
		return nil
	}
	return c
}

func buyGuard(cond bool, c <-chan Buy) <-chan Buy {
	if !cond {
		return nil
	}
	return c
}

// Guard for repair manager
func repairGuard(cond bool, c <-chan RepairNeed) <-chan RepairNeed {
	if !cond {
		return nil
	}
	return c
}

// Task Manager
func taskMagazine(getTodo <-chan DoTodo, storeTask <-chan Task, showMagazine <-chan bool) {
	tasks := make([]Task, 0)
	for {
		select {
		case request := <-doTodoGuard(len(tasks) > 0, getTodo):
			// Take one, cut to start from next one
			t := tasks[0]
			tasks = tasks[1:]
			request.response <- t
		case newTodo := <-taskGuard(len(tasks) < 10000, storeTask):
			tasks = append(tasks, newTodo)
		case <-showMagazine:
			fmt.Println("Task Magazine: ", tasks)
		}
	}
}

// Product Manager
func productMagazine(storeProduct <-chan Product, Buy <-chan Buy, showMagazine <-chan bool) {
	product := make([]Product, 0)
	for {
		select {
		case newProduct := <-productGuard(len(product) < 10000, storeProduct):
			product = append(product, newProduct)
		case buying := <-buyGuard(len(product) >= 1, Buy):
			// Take one, cut to start from next one
			buying.response <- product[0]
			product = product[1:]
		case <-showMagazine:
			fmt.Println("Product Magazine:", product)
		}
	}
}

// Repait Manager
func repairFacility(repairRequest <-chan RepairNeed, doneRepairs <-chan RepairMach, repairTodos <-chan RepairMach, showBrokenMachines <-chan bool) {

	// Getters for first item in list
	toFixList := make([]RepairMach, 0)
	nowBeingFixed := make([]RepairMach, 0)

	for {
		// Actions depending if new errand was given or if repairman reports completed task
		select {
		case toRepair := <-repairTodos:
			// If it isn't yet on the todo fix machine list - add it
			if browseList(toFixList, toRepair) == NotOnList {
				toFixList = append(toFixList, toRepair)
				nowBeingFixed = append(nowBeingFixed, toRepair)
			}

		case repaired := <-doneRepairs:
			// Cross out repaired machine from the list to fix (so it can be added again upon another problem)
			numOnList := browseList(toFixList, repaired)
			// https://golang.org/ref/spec#Passing_arguments_to_..._parameters
			toFixList = append(toFixList[:numOnList], toFixList[numOnList+1:]...)

		case request := <-repairGuard(len(nowBeingFixed) > 0, repairRequest):
			// Secondment to a repairmen, cutting from list
			request.response <- nowBeingFixed[0]
			nowBeingFixed = nowBeingFixed[1:]

		case <-showBrokenMachines:
			fmt.Println("Broken Machines: ", toFixList)
		}
	}
}

func browseList(toFixList []RepairMach, lookFor RepairMach) int {
	for numOnList, mach := range toFixList {
		if mach.id == lookFor.id && mach.model == lookFor.model {
			return numOnList
		}
	}
	return NotOnList
}
