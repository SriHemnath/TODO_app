package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/SriHemnath/TODO/cmd/database"
	"github.com/jmoiron/sqlx"
	"github.com/manifoldco/promptui"
	"os"
	"os/exec"
	"runtime"
)

var clear map[string]func()
var db *sqlx.DB
var schema = `
CREATE TABLE IF NOT EXISTS tasks (
    title text,
    note text,
    archived_at TIMESTAMPTZ
);
`

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
	db = database.NewDB()
	db.MustExec(schema)
}

func Run() {
	fmt.Println("################ Welcome to TODO app ################")
	//replacing with db
	/*filePath := "todo.txt"
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
	}*/

	for {
		prompt := promptui.Select{
			Label: "Choose option",
			Items: []string{"New Note", "View", "Delete", "View deleted notes", "Exit"},
		}

		_, result, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		switch result {
		case "Exit":
			clearCli()
			db.Close()
			os.Exit(0)
		case "View":
			viewNote()
		case "New Note":
			newNote()
		case "Delete":
			deleteNotes()
		case "View deleted notes":
			viewDeletedNotes()
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

func newNote() {
	clearCli()
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Please enter title for your note:")
	title, _ := reader.ReadString('\n')
	fmt.Println("Enter note:")
	note, _ := reader.ReadString('\n')
	qry := `INSERT INTO tasks (title,note) VALUES ($1,$2)`
	_, err := db.Exec(qry, title, note)
	if err != nil {
		fmt.Println("failed to insert note")
		return
	}
	clearCli()
}

func viewNote() {
	clearCli()
	titles := getAllNotes()
	prompt := promptui.Select{
		Label: "Choose note",
		Items: titles,
		Size:  10,
	}

	_, choice, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	var note string
	qry := `SELECT note FROM tasks WHERE title=$1`
	err = db.Get(&note, qry, choice)
	if err != nil {
		fmt.Println("failed to get notes")
		return
	}
	fmt.Println(note)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	clearCli()
}

func viewDeletedNotes() {
	clearCli()
	titles := getAllDeletedNotes()
	prompt := promptui.Select{
		Label: "Choose note",
		Items: titles,
	}

	_, choice, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	var note string
	qry := `SELECT note FROM tasks WHERE title=$1`
	err = db.Get(&note, qry, choice)
	if err != nil {
		fmt.Println("failed to get notes")
		return
	}
	fmt.Println(note)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	clearCli()
}

func deleteNotes() {
	titles := getAllNotes()
	prompt := promptui.Select{
		Label: "Choose note",
		Items: titles,
	}

	_, choice, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	query := `
	UPDATE tasks 
	SET 
    	archived_at=now()
    WHERE
        title = $1`
	_, err = db.Exec(query, choice)
	if err != nil {
		fmt.Printf("failed to delete note %s", err)
		os.Exit(1)
	}
	clearCli()
}

func getAllNotes() []string {
	clearCli()
	qry := `SELECT title FROM tasks WHERE archived_at IS NULL;`
	var result []string
	err := db.Select(&result, qry)
	if err != nil {
		fmt.Println("failed to get notes")
		return nil
	}
	return result
}

func getAllDeletedNotes() []string {
	clearCli()
	qry := `SELECT title FROM tasks WHERE archived_at IS NOT NULL;`
	var result []string
	err := db.Select(&result, qry)
	if err != nil {
		fmt.Println("failed to get notes")
		return nil
	}
	return result
}
