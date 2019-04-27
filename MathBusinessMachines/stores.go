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
