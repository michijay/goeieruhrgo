package main

import (
	"time"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Go-Eieruhr")

	label := widget.NewLabel("00:00")
	input := widget.NewEntry()
	input.SetPlaceHolder("Minuten...")

	var stopChan chan bool

	startBtn := widget.NewButton("Start", func() {
		// Logik für den Timer
		stopChan = make(chan bool)
		go func() {
			// Hier würde die Zeitschleife laufen
			label.SetText("Läuft...")
		}()
	})

	stopBtn := widget.NewButton("Stop", func() {
		if stopChan != nil {
			stopChan <- true
		}
	})

	myWindow.SetContent(container.NewVBox(
		label,
		input,
		startBtn,
		stopBtn,
	))

	myWindow.ShowAndRun()
}
