package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Println("migration files verification")

	dirs := []string{
		"internal/db/migrations",
		"cmd/verify",
	}

	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			log.Fatalf("Directory %s does not exist", dir)
		}
		fmt.Printf("✓ Directory %s exists\n", dir)
	}

	migrations, err := os.ReadDir("internal/db/migrations")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nFound %d migration files:\n", len(migrations))
	for _, m := range migrations {
		fmt.Printf("  - %s\n", m.Name())
	}
}
