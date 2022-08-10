package main

import (
    "fmt"
    "net/http"
    "html/template"
    "path"
    "strings"
	"strconv"
    "time"
)

type Response struct {
    Success string
    Message string
}

type Task struct {
    Id          string
    TaskName    string
    Assignee  string
    Date string
    Status string
}

var tasks []Task;

func main() {
	fmt.Println("Test golang progate")
	fmt.Println("Hello World")
	http.HandleFunc("/", handlerIndex)
    http.HandleFunc("/create", createTask)
    http.HandleFunc("/save", saveTask)
    http.HandleFunc("/update/", updateTask)
    http.HandleFunc("/delete/", deleteTask)

	http.Handle("/static/",
        http.StripPrefix("/static/",
            http.FileServer(http.Dir("assets"))))

    var address = "localhost:9090"
    fmt.Printf("server started at %s\n", address)
    err := http.ListenAndServe(address, nil)
    if err != nil {
        fmt.Println(err.Error())
    }
}

func handlerIndex(w http.ResponseWriter, r *http.Request) {
    currentTime := time.Now()
    fmt.Println("currentTime", currentTime.Format("2006-01-02"))
    
    for i := range tasks {
        // tasks[i].Status = Datet, err := time.Parse(layout, str)
        s := "N/A"
        if tasks[i].Date != "" {

            if currentTime.Format("2006-01-02") == tasks[i].Date {
                s = "Last Day"
            } else {
                t, err := time.Parse("2006-01-02", tasks[i].Date)
                if err != nil {
                    fmt.Println(err)
                }
                g1 := currentTime.Before(t)
                fmt.Println("today before tomorrow:", g1)

                s = "Expired"
                if g1 == true {
                    s = "Active"
                }
            }
        }

        tasks[i].Status = s
    }

	var filepath = path.Join("views", "index.html")
	var tmpl = template.Must(template.ParseFiles(filepath))
	if err := tmpl.Execute(w, map[string]interface{}{"tasks":tasks}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func createTask(w http.ResponseWriter, r *http.Request) {
    var filepath = path.Join("views", "form.html")
	var tmpl = template.Must(template.New("form").ParseFiles(filepath))
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func updateTask(w http.ResponseWriter, r *http.Request) {
    id := strings.TrimPrefix(r.URL.Path, "/update/")
    var selectedTask Task
    for i := range tasks {
        if tasks[i].Id == id {
            selectedTask = tasks[i]
        }
    }

    var filepath = path.Join("views", "form.html")
	var tmpl = template.Must(template.New("form").ParseFiles(filepath))
	if err := tmpl.Execute(w, selectedTask); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func saveTask(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {

        if err := r.ParseForm(); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

		var Id = r.FormValue("Id")
        var TaskName = r.FormValue("TaskName")
        var Assignee = r.FormValue("Assignee")
        var Date = r.FormValue("Date")

        action := "Create"

        if Id != "" {
            action = "Update"
            for i := range tasks {
                if tasks[i].Id == Id {
                    tasks[i].Id = Id
                    tasks[i].TaskName = TaskName
                    tasks[i].Assignee = Assignee
                    tasks[i].Date = Date
                }
            }
        } else {
            newId := strconv.Itoa(len(tasks))
            tasks = append(tasks, Task{Id: newId, TaskName: TaskName, Assignee:Assignee, Date:Date})
        }

        var response = Response {
            Message : action + " Task berhasil!",
            Success : "true",
        }


        var filepath = path.Join("views", "info.html")
		var tmpl = template.Must(template.ParseFiles(filepath))
		if err := tmpl.Execute(w, response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
        return
    }

    http.Error(w, "", http.StatusBadRequest)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
    id := strings.TrimPrefix(r.URL.Path, "/delete/")

    for i := range tasks {
        if tasks[i].Id == id {
            tasks[i] = tasks[len(tasks)-1]
            tasks = tasks[:len(tasks)-1]
        }
    }

    var response = Response {
        Message : "Data berhasil dihapus!",
        Success : "true",
    }

    var filepath = path.Join("views", "info.html")
	var tmpl = template.Must(template.ParseFiles(filepath))
	if err := tmpl.Execute(w, response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}