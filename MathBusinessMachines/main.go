package main

/* Current work flow:

-> Boss creates a task
-> Task goes to MagazineTodo
-> Worker grabs a task and tries to put it into a machine
-> Either awaits in que or runs around looking for free one
-> Machine solves task
-> Worker takes the task and stores it in MagazineReady
-> Customer buys the product

*/

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {

	// Starting message with a prompt for mode
	fmt.Println("Starting the Business. It will run for one hour before bankruptcy.")
	fmt.Println("Type 'S' for Silent mode or type 'T' for talkactive.")
	var mode string
	_, _ = fmt.Scanf("%s", &mode)

	// At the beginning of the program user decide in which mode program will be running
	switch mode {
	case "S", "s", "Silent":
		Silent = true
	case "T", "t", "Talkactive":
		Silent = false
	default:
		fmt.Println("Wrong char. Running with default talkactive mode.")
	}

	// Setting up people
	bossTasks := make(chan Task)
	workerTodos := make(chan DoTodo)
	workerProducts := make(chan Product)
	customerProducts := make(chan Buy)

	// Setting up magazines and machines
	addingMachines = setAddingMachines()
	multiplayingMachines = setMultiplayingMachines()

	showTaskMagazine := make(chan bool)
	showProductMagazine := make(chan bool)
	go taskMagazine(workerTodos, bossTasks, showTaskMagazine)
	go productMagazine(workerProducts, customerProducts, showProductMagazine)

	// Starting the company
	// Change the values through config.go
	go boss(bossTasks)

	displayWorkerStats := make([]chan bool, WorkNum)
	for i := 0; i < WorkNum; i++ {

		// Display worker stats for silent mode - it can have a delay due to worker being busy
		displayWorkerStats[i] = make(chan bool)

		// Decide behavior on birth
		patient := true
		if rand.Intn(2) == 1 {
			patient = false
		}

		go worker(i, workerTodos, workerProducts, displayWorkerStats[i], patient)
	}

	for i := 0; i < CustNum; i++ {
		go customer(i, customerProducts)
	}

	// Desk for user when mode is set to silent
	if Silent {
		for {
			fmt.Println(" > Silent mode. \n > Available commands:")
			fmt.Println("tasks - Prints the content of task magazine.")
			fmt.Println("products - Prints the content of products in the shop.")
			fmt.Println("workers - Prints stats of all workers in the company.")

			var command string
			_, _ = fmt.Scan(&command)

			switch command {
			case "tasks":
				showTaskMagazine <- true
			case "products":
				showProductMagazine <- true
			case "workers":
				for _, s := range displayWorkerStats {
					s <- true
				}
			default:
				fmt.Println("Incorrect input")
			}
		}
	}

	// Time before the program will quit
	time.Sleep(time.Hour)
}
