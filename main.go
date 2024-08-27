package main

import (
    "server/src/server"
)


func main() {
    server := server.NewServer()
    server.ListenAndServe()
}
