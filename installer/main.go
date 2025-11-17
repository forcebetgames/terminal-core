package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	w := myApp.NewWindow("TRM Cassino Installer")

	statusLabel := widget.NewLabel("Welcome to the TRM Cassino installer!")

	installBtn := widget.NewButton("Install", func() {
		statusLabel.SetText("Install button clicked! (Here we will install)")
	})

	cancelBtn := widget.NewButton("Cancel", func() {
		myApp.Quit()
	})

	w.SetContent(
		container.NewVBox(
			statusLabel,
			installBtn,
			cancelBtn,
		),
	)

	w.Resize(fyne.NewSize(400, 150))
	w.ShowAndRun()
}
