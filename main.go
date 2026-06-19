package main

import "myapp/database"

func main() {
    database.Connect()
    database.Migrate()
}
