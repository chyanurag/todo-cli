package main

import (
    "os"
    "encoding/json"
    "fmt"
    "bufio"
    "strings"
    "strconv"
)

type Todo struct {
    Id int
    Text string
    Done bool
}

type TodoHandler struct {
    todos []Todo
    File string
}

func (h TodoHandler) printTodos() {
    fmt.Printf("Id\t\tName\t\t\t\tDone\n")
    for _, todo := range h.todos {
        fmt.Printf("%d\t\t%s\t\t\t%v\n", todo.Id, todo.Text, todo.Done)
    }
}

func (h TodoHandler) listTodos() {
    if len(h.todos) > 0 {
        h.printTodos()
    } else {
        fmt.Println("You currently don't have any todo")
    }
}

func (h *TodoHandler) removeTodo(id string) {
    n, err := strconv.Atoi(id)
    if err != nil {
        fmt.Println("Please provide a valid id")
    }
    if n > len(h.todos) || n < 1 {
        fmt.Println("No such task exists")
        return
    }
    h.todos = append(h.todos[:n], h.todos[n+1:]...)
    writeTodosToFile(*h, h.File)
    fmt.Println("Todo removed successfully")
}

func (h *TodoHandler) addTodo(text string) {
    if strings.Trim(text, " ") == "" {
        fmt.Println("Please provide a valid task")
        return
    }
    todo := Todo { len(h.todos) + 1, text, false }
    h.todos = append(h.todos, todo)
    writeTodosToFile(*h, h.File)
    fmt.Println("Todo added successfully")
}

func (h *TodoHandler) markTodo(id string) {
    s, err := strconv.Atoi(id)
    if err != nil || s > len(h.todos) || s < 1 {
        fmt.Println("Invalid id")
        return
    }
    if h.todos[s-1].Done {
        fmt.Println("Task is already marked done")
    } else {
        h.todos[s-1].Done = true
        fmt.Println("Task marked done successfully")
    }
    writeTodosToFile(*h, h.File)
}

func (h *TodoHandler) unmarkTodo(id string) {
    s, err := strconv.Atoi(id)
    if err != nil || s > len(h.todos) || s < 1 {
        fmt.Println("Invalid id")
        return
    }
    if !h.todos[s-1].Done {
        fmt.Println("Task marked not done")
    } else {
        h.todos[s-1].Done = false 
        fmt.Println("Task marked not done successfully")
    }
    writeTodosToFile(*h, h.File)

}

func (h *TodoHandler) init(file string) {
    todos := getTodosFromFile(file)
    h.File = file
    h.todos = todos
}

func takeInput(prompt string) string {
    fmt.Print(prompt)
    sc := bufio.NewScanner(os.Stdin)
    sc.Scan()
    input := sc.Text()
    for {
        if strings.Trim(input, " ") != "" {
            return input
        }
        fmt.Println("Please input something ")
        continue
    }
}

func getTodosFromFile(file string) []Todo {
    var todos []Todo
    data, err := os.ReadFile(file)
    if err != nil { return todos }
    json.Unmarshal(data, &todos)
    return todos
}

func writeTodosToFile(h TodoHandler, file string) {
    todosS, _ := json.Marshal(h.todos)
    os.WriteFile(file, todosS, 0644)
}

func main() {

    handler := new(TodoHandler)
    handler.init("todos.txt")
    if len(os.Args) < 2 {
        fmt.Println("Not enough arguments")
        return
    }
    command := os.Args[1]
    checkArgument := func (prompt string) {
        if len(os.Args) < 3 {
            fmt.Println(prompt)
            os.Exit(0)
        }
    }
    switch command {
        case "add":
            checkArgument("Pleas provide a task string")
            handler.addTodo(os.Args[2])
        case "list":
            handler.listTodos()
        case "check":
            checkArgument("Pleas provide the task id")
            handler.markTodo(os.Args[2])
        case "uncheck":
            checkArgument("Pleas provide the task id")
            handler.unmarkTodo(os.Args[2])
        case "remove":
            checkArgument("No such task id")
            handler.removeTodo(os.Args[2])

    }


}
