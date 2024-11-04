package main

import (
	"fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "parking/internal/ui"
    "parking/internal/service"

)

func main() {
    myApp := app.New()
    window := myApp.NewWindow("Simulador de Estacionamiento")
    
    uiSimulator := ui.NewUISimulator(window)
    simulator := service.NewParkingSimulator(uiSimulator)
    
    container := uiSimulator.CreateUI(func() {
        simulator.Stop()
        window.Close()
    })
    
    window.SetContent(container)
    simulator.Start()
    
    window.Resize(fyne.NewSize(1600, 1200))
    window.ShowAndRun()
}