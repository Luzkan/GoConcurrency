package Go

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Task generated by boss (ex: 93195423 + 2138217)
type Task struct {
	arg1   int
	action string
	arg2   int
}

// Contains list of todos of type Task that a worker will need to solve
type MagazineTodo struct {
	todos []Task
	mux   sync.Mutex
}

// Contains list of products which were made by workers by solving todos from MagazineTodo
type MagazineReady struct {
	products []int
	mux      sync.Mutex
}

func boss(mgzTodo *MagazineTodo) {

	for {
		// https://golang.org/pkg/math/rand/
		// Generating vars for new task
		arg1 := rand.Intn(100)
		arg2 := rand.Intn(100)
		actions := [3]string{"+", "-", "*"}

		// Sticking them together
		t := Task{arg1, actions[rand.Intn(len(actions))], arg2}

		// Proceed to store new task
		storeTodo(mgzTodo, t)

		// Sleep means "performance" of the boss - defined in config
		time.Sleep(time.Duration(BossPerf) * time.Second)
	}
}

func storeTodo(mgzTodo *MagazineTodo, task Task) {

	// https://tour.golang.org/concurrency/9
	// Adding task to the stack in MagazineTodo
	mgzTodo.mux.Lock()
	mgzTodo.todos = append(mgzTodo.todos, task)
	mgzTodo.mux.Unlock()

	if !Silent {
		fmt.Println("[Added]  New task in the Magazine: ", task)
	}
}

func worker(workerID int, mgzTodo *MagazineTodo, mgzReady *MagazineReady) {

	// https://tour.golang.org/flowcontrol/4
	// Never ending for loop
	for {
		mgzTodo.mux.Lock()

		// Check if something is in the MagazineTodo
		if len(mgzTodo.todos) != 0 {

			// Get first task, then cut it out from the magazine
			task := mgzTodo.todos[0]
			mgzTodo.todos = mgzTodo.todos[1:]

			// Solve the task
			answer := solve(task)
			storeSolved(mgzReady, answer)
			mgzTodo.mux.Unlock()

			if !Silent {
				fmt.Println("[Job] \t Worker ", workerID, " solved a task. Answer: ", answer)
			}

		} else {
			mgzTodo.mux.Unlock()
			//if !Silent {
			//	fmt.Println("[Job] \t Worker ", workerID, " awaits for a better occasion.")
			//}
		}

		// Sleep means "performance" of workers - defined in config
		time.Sleep(time.Duration(WorkPerf) * time.Second)
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

func storeSolved(products *MagazineReady, product int) {

	// https://tour.golang.org/concurrency/9
	// Adding task to the stack in MagazineTodo
	products.mux.Lock()
	products.products = append(products.products, product)
	products.mux.Unlock()

	if !Silent {
		fmt.Println("[Added]  New product in the Magazine: ", product)
	}
}

func customer(customerID int, mgzReady *MagazineReady) {

	for {
		mgzReady.mux.Lock()

		// Check if something is in the MagazineReady
		if len(mgzReady.products) > 0 {

			// Get first product, then cut it out from the magazine
			answer := mgzReady.products[0]
			mgzReady.products = mgzReady.products[1:]
			mgzReady.mux.Unlock()
			if !Silent {
				fmt.Println("[Shop] \t Customer ", customerID, " bought ", answer)
			}
		} else {
			mgzReady.mux.Unlock()
			//if !Silent {
			//	fmt.Println("[Shop] \t Customer ", customerID, " found nothing in the shop.")
			//}
		}
		time.Sleep(time.Duration(CustPerf) * time.Second)
	}
}

func main() {

	fmt.Println("Starting the Business. It will run for one hour before bankruptcy.")
	fmt.Println("Type 'S' for Silent mode or type 'T' for talkactive.")
	var mode string
	_, _ = fmt.Scanf("%s", &mode)

	switch mode {
	case "S", "s", "Silent":
		Silent = true
	case "T", "t", "Talkactive":
		Silent = false
	default:
		fmt.Println("Wrong char. Running with default talkactive mode.")
	}

	mgzTodo := MagazineTodo{todos: make([]Task, 0)}
	mgzReady := MagazineReady{products: make([]int, 0)}

	// Change the values through config.go
	go boss(&mgzTodo)

	for workerID := 0; workerID != WorkNum; workerID++ {
		go worker(workerID, &mgzTodo, &mgzReady)
	}

	for custID := 0; custID != CustNum; custID++ {
		go customer(custID, &mgzReady)
	}

	if Silent {
		go desk(&mgzTodo, &mgzReady)
	}

	// Time before the program will quit
	time.Sleep(time.Hour)
}

func desk(mgzTodo *MagazineTodo, mgzReady *MagazineReady) {
	for {
		fmt.Println(" > Silent mode. \n > Available commands:")
		fmt.Println("tasks - Prints the content of task magazine.")
		fmt.Println("products - Prints the content of products in the shop.")

		var command string
		_, _ = fmt.Scan(&command)

		switch command {
		case "tasks":
			showMagazineTodo(mgzTodo)
		case "products":
			showMagazineReady(mgzReady)
		default:
			fmt.Println("Incorrect input")
		}
	}
}

func showMagazineTodo(mgzTodo *MagazineTodo) {
	mgzTodo.mux.Lock()
	for _, task := range mgzTodo.todos {
		fmt.Println(task)
	}
	mgzTodo.mux.Unlock()
}

func showMagazineReady(mgzReady *MagazineReady) {
	mgzReady.mux.Lock()
	for _, product := range mgzReady.products {
		fmt.Println(product)
	}
	mgzReady.mux.Unlock()
}