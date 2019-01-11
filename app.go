package main

import (
  "html/template"
  "log"
  "net"
  "net/http"
  "os"
  "path/filepath"
  "fmt"
  "io"
)

func upload(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
         fmt.Fprint(w, "This is the upload uri, and you ran a GET. Hmmm.")
    } else if r.Method == "POST" {
        file, handler, err := r.FormFile("file")
        if err != nil {
            fmt.Println(err)
            return
        }
        defer file.Close()

        fmt.Fprintf(w, "%v", handler.Header)
        f, err := os.OpenFile("repo/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
        if err != nil {
            fmt.Println(err)
            return
        }
        defer f.Close()

        io.Copy(f, file)

    } else {
          fmt.Println("Unknown HTTP "+ r.Method +" Method")
    }
}

func GetLocalIP() string {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return ""
    }
    for _, address := range addrs {
        if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() { 
            if ipnet.IP.To4() !=nil {
                return ipnet.IP.String()
            }
        }
    }
    return ""
}

func main() {
  ip := GetLocalIP()
  port := "3000"

  fs_static := http.FileServer(http.Dir("static"))
  fs_repo := http.FileServer(http.Dir("repo"))
  http.Handle("/static/", http.StripPrefix("/static/", fs_static))
  http.Handle("/repo/", http.StripPrefix("/repo/", fs_repo))
  http.HandleFunc("/", serveTemplate)
  http.HandleFunc("/upload", upload)

  log.Println("Go Repo listening @", ip, ":", port)
  http.ListenAndServe(":3000", nil)
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
  log.Println("Get /")
  lp := filepath.Join("templates", "layout.html")
  fp := filepath.Join("templates", filepath.Clean(r.URL.Path))

  // Return a 404 if the template doesn't exist
  info, err := os.Stat(fp)
  if err != nil {
    if os.IsNotExist(err) {
      http.NotFound(w, r)
      return
    }
  }

  // Return a 404 if the request is for a directory
  if info.IsDir() {
    http.NotFound(w, r)
    return
  }

  tmpl, err := template.ParseFiles(lp, fp)
  if err != nil {
    // Log the detailed error
    log.Println(err.Error())
    // Return a generic "Internal Server Error" message
    http.Error(w, http.StatusText(500), 500)
    return
  }

  if err := tmpl.ExecuteTemplate(w, "layout", nil); err != nil {
    log.Println(err.Error())
    http.Error(w, http.StatusText(500), 500)
  }
  log.Println(fp) 
}
