package main

import (
	"fmt"
	"os"

	"github.com/hubblew/pim/internal/ui"
)

type Demo struct {
	Name string
	Run  func()
}

func main() {
	demos := []Demo{
		{Name: "Choice Component Demos", Run: runChoiceDemos},
		{Name: "Press Any Key Demo", Run: runPressAnyKeyDemos},
		{Name: "Spinner Dialog Demos", Run: runSpinnerDemos},
	}

	choices := make([]ui.Choice, len(demos))
	for i, demo := range demos {
		choices[i] = ui.Choice{
			Label: demo.Name,
			Value: demo,
		}
	}

	fmt.Println()
	fmt.Println("═════════════════════════════════════════════════")
	fmt.Println(" PIM UI Component Demos")
	fmt.Println("═════════════════════════════════════════════════")
	fmt.Println()

	model := ui.NewChoiceDialog("Select a demo to run:", choices).Vertical()
	choice, err := model.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if choice == nil {
		fmt.Println("\nBye.")
		os.Exit(0)
	}

	selectedDemo := choice.Value.(Demo)
	fmt.Println("")
	selectedDemo.Run()
	fmt.Println("")
	fmt.Println("✅ Demo completed!")

	main()
}
