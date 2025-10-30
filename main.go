package main

import (
    "flag"
    "fmt"
    "os"
)

func main() {
    // Define command flags
    add := flag.String("add", "", "Add a new todo")
    list := flag.Bool("list", false, "List all todos")
    done := flag.Int("done", 0, "Mark todo as done (by ID)")
    del := flag.Int("delete", 0, "Delete todo (by ID)")
    
    flag.Parse()
    
    // Handle commands based on flags
    if *add != "" {
        // Add new todo
        fmt.Printf("Adding: %s\n", *add)
    } else if *list {
        // List all todos
        fmt.Println("Listing todos...")
    } else if *done != 0 {
        // Mark as done
        fmt.Printf("Marking todo %d as done\n", *done)
    } else if *del != 0 {
        // Delete todo
        fmt.Printf("Deleting todo %d\n", *del)
    } else {
        // Show help if no valid command
        flag.Usage()
        os.Exit(1)
    }
}