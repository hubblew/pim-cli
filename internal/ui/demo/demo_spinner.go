package main

import (
	"fmt"
	"time"

	"github.com/hubblew/pim/internal/ui"
)

func demoQuickSpinner() {
	fmt.Println("=== Quick Task Spinner ===")
	fmt.Println()

	err := ui.RunWithSpinner("Loading configuration...", func() error {
		time.Sleep(2 * time.Second)
		return nil
	})
	if err != nil {
		fmt.Printf("\nError: %v\n", err)
		return
	}

	fmt.Println("\n✅ Configuration loaded!")
}

func demoLongSpinner() {
	fmt.Println("=== Long Task Spinner ===")
	fmt.Println()

	err := ui.RunWithSpinner("Detecting CLI agents in your system...", func() error {
		time.Sleep(4 * time.Second)
		return nil
	})
	if err != nil {
		fmt.Printf("\nError: %v\n", err)
		return
	}

	fmt.Println("\n✅ Detection complete!")
}

func demoMultipleSpinners() {
	fmt.Println("=== Multiple Tasks ===")
	fmt.Println()

	tasks := []struct {
		name     string
		duration time.Duration
	}{
		{"Initializing", 1 * time.Second},
		{"Downloading dependencies", 2 * time.Second},
		{"Building project", 3 * time.Second},
		{"Running tests", 2 * time.Second},
	}

	for _, task := range tasks {
		err := ui.RunWithSpinner(task.name+"...", func() error {
			time.Sleep(task.duration)
			return nil
		})
		if err != nil {
			fmt.Printf("\nError: %v\n", err)
			return
		}
		fmt.Printf("✅ %s complete!\n", task.name)
	}

	fmt.Println("\n🎉 All tasks finished!")
}

func runSpinnerDemos() {
	demoQuickSpinner()

	fmt.Println("\n" + "─────────────────────────────────────────────────")
	fmt.Println("")

	demoLongSpinner()

	fmt.Println("\n" + "─────────────────────────────────────────────────")
	fmt.Println("")

	demoMultipleSpinners()
}
