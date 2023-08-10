package main

import (
	"bufio"
	"fmt"

	// _ "fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	args := os.Args
	handleArgs(args)
}

func handleArgs(args []string) bool {
	if args[1] == "launch" {
		log.Println("Launching...")
		launch()
		return true
	}
	if strings.Contains(args[1], "new") {
		log.Println("Creating new task...")
		newTask(args)
		return true
	}
	log.Fatal("No arguments provided! Use new to create a new task or launch to manually launch all tasks \n Launch should only be used for testing purposes, systemd will handle this automatically.")
	os.Exit(2)
	return false
}

func launch() {
	//read tasks file
	file, err := os.Open("tasks.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		task := scanner.Text()
		// do something with task
		strings.Split(task, "@")
		var output = runCommand(task)
		logOutput(output)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func runCommand(command string) string {
	// run command
	out, err := exec.Command(command).Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}

func logOutput(output string) bool {
	file, err := os.Create("log.txt")
	if err != nil {
		log.Panic(err)
		return false
	}
	defer file.Close()

	_, err = file.WriteString(output)
	if err != nil {
		log.Panic(err)
		return false
	}

	return true
}

func newTask(args []string) bool {
	var taskcommand = args[2]
	log.Println("Task command: " + taskcommand)
	log.Println("Continue? (y/n)")
	var input string
	fmt.Scanln(&input)
	if input != "y" {
		log.Println("Aborting...")
		return false
	}
	log.Println("Creating task...")

	//var taskcommand = args[2]
	if taskcommand == "" {
		log.Println("Task name or command not provided! \n Usage: new <taskcommand>")
		return false
	}
	checkForFile()
	file, err := os.OpenFile("tasks.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Panic(err)
		return false
	}
	defer file.Close()
	_, err = file.WriteString(taskcommand + "\n")
	if err != nil {
		return false
	}
	log.Println("Task created!")
	return true
}

func checkForFile() {
	//check if tasks file exists
	_, err := os.Stat("tasks.txt")
	if err != nil {
		log.Println("Tasks file not found, creating...")
		_, err := os.Create("tasks.txt")
		if err != nil {
			log.Panic(err)
		}
		log.Println("Tasks file created!")
	}
}
