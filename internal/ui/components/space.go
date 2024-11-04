package components

import (
    "fmt"
    "image/color"
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/canvas"
    "fyne.io/fyne/v2/container"
)

func CreateParkingSpace() *fyne.Container {
    outerContainer := container.NewPadded()
    
    rect := canvas.NewRectangle(color.RGBA{0, 0, 0, 255})
    rect.Resize(fyne.NewSize(150, 100))
    
    line1 := canvas.NewLine(color.RGBA{255, 255, 0, 255})
    line1.StrokeWidth = 3
    line1.Position1 = fyne.NewPos(0, 0)
    line1.Position2 = fyne.NewPos(150, 0)
    
    line2 := canvas.NewLine(color.RGBA{255, 255, 0, 255})
    line2.StrokeWidth = 3
    line2.Position1 = fyne.NewPos(150, 0)
    line2.Position2 = fyne.NewPos(150, 100)
    
    line3 := canvas.NewLine(color.RGBA{255, 255, 0, 255})
    line3.StrokeWidth = 3
    line3.Position1 = fyne.NewPos(150, 100)
    line3.Position2 = fyne.NewPos(0, 100)
    
    line4 := canvas.NewLine(color.RGBA{255, 255, 0, 255})
    line4.StrokeWidth = 3
    line4.Position1 = fyne.NewPos(0, 100)
    line4.Position2 = fyne.NewPos(0, 0)

    spaceContainer := container.NewWithoutLayout(rect, line1, line2, line3, line4)
    spaceContainer.Resize(fyne.NewSize(150, 100))
    
    paddedContainer := container.NewPadded(spaceContainer)
    outerContainer.Add(paddedContainer)
    
    return outerContainer
}

func CreateSpaceWithCar(vehicleID int) *fyne.Container {
    container := container.NewWithoutLayout()
    
    parkingSpace := CreateParkingSpace()
    
    carImage := canvas.NewImageFromFile("assets/carroRojo.png")
    carImage.FillMode = canvas.ImageFillOriginal
    carImage.Resize(fyne.NewSize(120, 80))
    carImage.Move(fyne.NewPos(15, 10))
    
    idText := canvas.NewText(fmt.Sprintf("%d", vehicleID), color.White)
    idText.TextSize = 10
    idText.Move(fyne.NewPos(70, 60))
    
    container.Add(parkingSpace)
    container.Add(carImage)
    container.Add(idText)
    container.Resize(fyne.NewSize(150, 100))
    
    return container
}