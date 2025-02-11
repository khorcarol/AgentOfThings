package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// Create a new Fyne app
	myApp := app.New()

	// Create a new window
	myWindow := myApp.NewWindow("Messages List")

	// List of messages to display
	messages := []string{
		"Hello, World!",
		"Welcome to the Fyne app!",
		"How are you?",
		"Goodbye!",
	}

	// Create a new List widget
	list := widget.NewList(
		func() int { return len(messages) }, // Function to determine the number of items
		func() fyne.CanvasObject {           // Function to create a new list item
			return widget.NewLabel("")
		},
		func(i int, o fyne.CanvasObject) { // Function to update the list item
			o.(*widget.Label).SetText(messages[i])
		},
	)

	// Set the list as the window content inside a scroll container
	scrollContainer := container.NewScroll(list)
	myWindow.SetContent(scrollContainer)

	// Show the window and run the app
	myWindow.ShowAndRun()
}
