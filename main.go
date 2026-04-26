package main

// GoEieruhrGo - Simple Timerapp in Go
// Author: Michael Janssen <m.janssen@lyrah.net>
// License: GPLv3 (See README.md for details)

import (
	"fmt"
	"time"
	"os/exec"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func playAlarm() {
	// 1. Priorität: VLC (kann fast alles)
	if _, err := exec.LookPath("cvlc"); err == nil {
		exec.Command("cvlc", "--play-and-exit", "/usr/share/sounds/freedesktop/stereo/complete.oga").Run()
		return
	}

	// 2. Priorität: Pipewire (moderner Standard)
	if _, err := exec.LookPath("pw-play"); err == nil {
		exec.Command("pw-play", "/usr/share/sounds/freedesktop/stereo/complete.oga").Run()
		return
	}

	// 3. Priorität: ALSA (Das Urgestein - nur für .wav!)
	if _, err := exec.LookPath("aplay"); err == nil {
		exec.Command("aplay", "/usr/share/sounds/alsa/Front_Center.wav").Run()
		return
	}
}

func main() {

	Version := "1.1-3"

	titel := fmt.Sprintf("GoEieruhrGo - V: %s", Version)

	myApp := app.New()
	myWindow := myApp.NewWindow(titel)
	myWindow.Resize(fyne.NewSize(300, 200))

	label := widget.NewLabelWithStyle("00:00", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	input := widget.NewEntry()
	input.SetPlaceHolder("Enter minutes (z.B. 5)")

	stopChan := make(chan bool)
	running := false

	updateTime := func(seconds int) {
		mins := seconds / 60
		secs := seconds % 60
		label.SetText(fmt.Sprintf("%02d:%02d", mins, secs))
		label.Refresh()
	}

	startBtn := widget.NewButton("Start", func() {
		if running { return }
		
		var minutes int
		fmt.Sscanf(input.Text, "%d", &minutes)
		if minutes <= 0 { return }

		timeLeft := minutes * 60
		running = true
		
		go func() {
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()

			for timeLeft > 0 {
				select {
				case <-stopChan:
					running = false
					return
				case <-ticker.C:
					timeLeft--
					updateTime(timeLeft)
				}
			}
			running = false
			//label.SetText("DONE!")

			go func() {
                playAlarm()
                time.Sleep(1 * time.Second)
                playAlarm()
            }()

			go func() {
                for i := 0; i < 4; i++ {
                    label.Hide()
                    time.Sleep(250 * time.Millisecond)
                    label.Show()
                    time.Sleep(250 * time.Millisecond)
                }
                label.SetText("DONE!")
            }()

		}()
	})

	ftenBtn := widget.NewButton("15 Mins", func() {
		input.SetText("15")
		label.SetText("15:00")
	})

	thirtyBtn := widget.NewButton("30 Mins", func() {
		input.SetText("30")
		label.SetText("30:00")
	})

	forthyBtn := widget.NewButton("45 Mins", func() {
		input.SetText("45")
		label.SetText("45:00")
	})

	sixtyBtn := widget.NewButton("60 Mins", func() {
		input.SetText("60")
		label.SetText("60:00")
	})

	stopBtn := widget.NewButton("Stop", func() {
		if running {
			stopChan <- true
		}

	})
	resetBtn := widget.NewButton("Reset", func() {
		label.SetText("00:00")
	})

	presetButtons := container.NewHBox(ftenBtn, thirtyBtn, forthyBtn, sixtyBtn)

	myWindow.SetContent(container.NewVBox(
		label,
		input,
		presetButtons,
		startBtn,
		stopBtn,
		resetBtn,
	))

	myWindow.ShowAndRun()
}

