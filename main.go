package main

import (
		// "code.google.com/p/go.net/websocket"
		"github.com/gorilla/websocket"
    "net/http"
    "os/exec"
    "log"
    "os"
    "path"
    "flag"
)

type Request struct {
    Cmd string
}

type Response struct {
    Result string
    Error string
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// func execHandler(ws *websocket.Conn) {
func execHandler(w http.ResponseWriter, r *http.Request) {
    var data Request
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		defer ws.Close()
		for {
			err := ws.ReadJSON(&data)
			if err != nil {
					log.Fatal(err)
			}

			log.Println("Received request", data)

			var out []byte
			if len(data.Cmd) > 2 && data.Cmd[0:3] == "cd " {
					out, err = changeWorkingDirectory(data.Cmd[3:])
			} else {
					var cmd = exec.Command("bash", []string{"-c", data.Cmd}...)
					out, err = cmd.CombinedOutput()
			}

			var error = ""
			if err != nil {
					log.Println("Error ocurred executing command", err)
					error = err.Error()
			}

			err = ws.WriteJSON(Response{Result:string(out), Error:error})
			if err != nil {
					log.Fatal("Could not send response", err)
			}

		}
}

func changeWorkingDirectory(newPath string) (out []byte, err error) {
    wd, wderr := os.Getwd()
    if wderr != nil {
        // No idea in which cases this happens, so not really sure
        // how to recover from it.
        log.Fatal("Could not get current working directory", err)
    }

    newWd := path.Clean(wd + "/" + newPath)

    log.Println("Changing directory to", newWd)
    err = os.Chdir(newWd)
    if err != nil {
        log.Println("Error changing directory", err)
    } else {
        out = []byte("Changed directory to" + newWd)
    }

    return
}

func main() {
    var path string

    flag.StringVar(&path, "path", "", "Working directory, i.e. where commands will be executed.")
    flag.Parse()

    if path == "" {
        log.Fatal("Please specify the 'path' option.")
    }

    err := os.Chdir(path)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Starting with working directory '%s'", path)

    // http.Handle("/exec", websocket.Handler(execHandler))
		http.HandleFunc("/exec", execHandler)
    err = http.ListenAndServe("0.0.0.0:8080", nil)
    if err != nil {
        log.Fatal("ListenAndServe:" + err.Error())
    }
}
