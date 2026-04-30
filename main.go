package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    reader := bufio.NewReader(os.Stdin)

    fmt.Print("Enter a host: ")
    host, _ := reader.ReadString('\n')

    fmt.Println("You entered:", host)
}
