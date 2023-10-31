package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"os"
	"os/exec"
	"runtime"
)

var clear map[string]func()

func init() {
	clear = make(map[string]func())
	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func Run() {
	fmt.Println("################ Welcome to TODO app ################")
	filePath := "todo.txt"
	isFileExist := checkFileExists(filePath)

	if isFileExist {
		fmt.Println("file exist")
	} else {
		f, err := os.Create(filePath)
		if err != nil {
			fmt.Println("error while creating file")
			return
		}
		f.Close()
	}

	for {
		prompt := promptui.Select{
			Label: "Choose option",
			Items: []string{"New Note", "View", "Delete", "Exit"},
		}

		_, result, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		switch result {
		case "Exit":
			clearCli()
			os.Exit(0)
		case "View":
			clearCli()
			data, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Println(err)
				return
			}
			if len(data) == 0 {
				fmt.Println("nothing to view")
				continue
			} else {
				fmt.Println(string(data))
				scanner := bufio.NewScanner(os.Stdin)
				scanner.Scan()
				clearCli()
			}
		case "New Note":
			clearCli()
			reader := bufio.NewReader(os.Stdin)
			line, _, _ := reader.ReadLine()
			os.WriteFile(filePath, line, 0644)
			clearCli()
		case "Delete":
			os.Truncate(filePath, 0)
			clearCli()
		}
	}
}

func clearCli() {
	v, _ := clear[runtime.GOOS]
	v()
}

func checkFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	//return !os.IsNotExist(err)
	return !errors.Is(error, os.ErrNotExist)
}
