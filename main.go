package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

func init() {
}

var login string = `
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <title>Jumanji</title>
</head>

<body>
    <style>
        div {
            margin: 3px;
        }

        label {
            display: inline-block;
            ;
            width: 100px;
        }
    </style>
    <form action="/submit" method="POST">
        <div>
            <label>Email:</label>
            <input type='email' autocomplete="username" id="username" name="username" />
        </div>
        <div>
            <label>Password:</label>
            <input type='password' autocomplete="new-password" id="password" name="password" />
        </div>
        <input type='submit' value="Login">
    </form>
</body>
</html>
`

var home string = `
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <title>Jumanji</title>
</head>

<body>
    <h2> Welcome to Jumanji!</h2>

    <h4> Username - {{ .username }}</h4>
    <h4> Password - {{ .password }}</h4>

</body>

</html>
`

func main() {
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(login)); err != nil {
			http.Error(w, "failed to load login page", http.StatusBadRequest)
			return
		}
	})

	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "failed parse login form", http.StatusBadRequest)
			return
		}

		data := make(map[string]string)
		data["username"] = r.Form.Get("username")
		data["password"] = r.Form.Get("password")

		tmpl, err := template.New("welcome").Parse(home)
		if err != nil {
			http.Error(w, "failed create welcome template", http.StatusBadRequest)
			return
		}

		var out bytes.Buffer
		if err := tmpl.Execute(&out, data); err != nil {
			http.Error(w, "failed to parse welcome template", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(out.Bytes())); err != nil {
			http.Error(w, "failed to load welcome page", http.StatusBadRequest)
			return
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting instance at port %v", port)
	log.Fatalf("%v", http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}
