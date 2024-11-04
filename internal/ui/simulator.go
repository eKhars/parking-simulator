package ui

import (
    "fmt"
    "image/color"
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/canvas"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/layout"
    "fyne.io/fyne/v2/widget"
    "parking/internal/ui/components"
)

type UISimulator struct {
    spaceContainers [20]*fyne.Container
    statusLabel     *widget.Label
    waitingLabel    *widget.Label
    processedLabel  *widget.Label
    directionLabel  *widget.Label
    entranceRect    *canvas.Rectangle
    window          fyne.Window
}

func NewUISimulator(window fyne.Window) *UISimulator {
    return &UISimulator{
        window: window,
    }
}

func (ui *UISimulator) CreateUI(onStop func()) *fyne.Container {
    mainContainer := container.NewPadded()
    
    parkingGridContainer := container.NewVBox(
        layout.NewSpacer(),
        container.NewGridWithColumns(4),
        layout.NewSpacer(),
    )
    
    parkingGrid := parkingGridContainer.Objects[1].(*fyne.Container)
    
    for i := 0; i < 20; i++ {
        spaceContainer := container.NewPadded(container.NewStack())
        spaceContainer.Add(components.CreateParkingSpace())
        ui.spaceContainers[i] = spaceContainer
        
        paddedSpace := container.NewPadded(
            container.NewPadded(
                container.NewPadded(spaceContainer),
            ),
        )
        parkingGrid.Add(paddedSpace)
    }

    entranceContainer := container.NewPadded()
    ui.entranceRect = canvas.NewRectangle(color.RGBA{128, 128, 128, 255})
    ui.entranceRect.Resize(fyne.NewSize(200, 40))
    entranceText := canvas.NewText("Entrada/Salida", color.White)
    entranceText.TextSize = 16
    entranceText.Move(fyne.NewPos(60, 10))
    entranceStack := container.NewStack(ui.entranceRect, entranceText)
    entranceContainer.Add(entranceStack)

    ui.statusLabel = widget.NewLabel("Simulación en curso...")
    ui.waitingLabel = widget.NewLabel("Carros en espera: 0")
    ui.processedLabel = widget.NewLabel("Carros procesados: 0")
    ui.directionLabel = widget.NewLabel("Puerta libre")
    
    statsContainer := container.NewHBox(
        layout.NewSpacer(),
        ui.statusLabel,
        widget.NewLabel("  |  "),
        ui.waitingLabel,
        widget.NewLabel("  |  "),
        ui.processedLabel,
        widget.NewLabel("  |  "),
        ui.directionLabel,
        layout.NewSpacer(),
    )

    stopButton := widget.NewButton("Detener", onStop)
    
    stopButtonContainer := container.NewHBox(
        layout.NewSpacer(),
        stopButton,
        layout.NewSpacer(),
    )

    mainContainer.Add(container.NewVBox(
        container.NewPadded(widget.NewLabel("\nSimulador de Estacionamiento\n")),
        container.NewPadded(entranceContainer),
        container.NewPadded(parkingGridContainer),
        widget.NewLabel(""),
        container.NewPadded(stopButtonContainer),
        widget.NewLabel(""),
        container.NewPadded(statsContainer),
    ))

    return mainContainer
}

func (ui *UISimulator) UpdateSpaceUI(index int, vehicleID int) {
    if vehicleID == 0 {
        ui.spaceContainers[index].Objects = []fyne.CanvasObject{components.CreateParkingSpace()}
    } else {
        ui.spaceContainers[index].Objects = []fyne.CanvasObject{components.CreateSpaceWithCar(vehicleID)}
    }
    ui.spaceContainers[index].Refresh()
}

func (ui *UISimulator) UpdateLabels(spacesInUse, waitingCars, enteringQueueSize, exitingQueueSize, processedCars, enteringCount, exitingCount int) {
    ui.statusLabel.SetText(fmt.Sprintf("Espacios ocupados: %d/20", spacesInUse))
    ui.waitingLabel.SetText(fmt.Sprintf("Carros en espera: %d (Cola entrada: %d, Cola salida: %d)", 
        waitingCars, enteringQueueSize, exitingQueueSize))
    ui.processedLabel.SetText(fmt.Sprintf("Carros procesados: %d", processedCars))
    
    if enteringCount > 0 {
        ui.directionLabel.SetText(fmt.Sprintf("Entrando: %d carros (Cola: %d)", 
            enteringCount, enteringQueueSize))
    } else if exitingCount > 0 {
        ui.directionLabel.SetText(fmt.Sprintf("Saliendo: %d carros (Cola: %d)", 
            exitingCount, exitingQueueSize))
    } else {
        ui.directionLabel.SetText("Puerta libre")
    }
}

func (ui *UISimulator) UpdateEntranceColor(status string) {
    switch status {
    case "entering":
        ui.entranceRect.FillColor = color.RGBA{0, 255, 0, 255}
    case "exiting":
        ui.entranceRect.FillColor = color.RGBA{255, 0, 0, 255}
    default:
        ui.entranceRect.FillColor = color.RGBA{128, 128, 128, 255}
    }
    ui.entranceRect.Refresh()
}