package main

import (
	"fmt"
	"urlshortener/internal/config"
)
func main() {
	fmt.Println("Hello, World!")
	// Additional code can be added here
	fmt.Println("Database Host:", config.LoadDatabaseConfig().DatabaseHost)	
}