package domain

import (
    "sync"
    "parking/internal/models"
)
type ParkingLot struct {
    Spaces           [20]int
    Entrance         sync.Mutex
    SpacesInUse      int
    EnteringCount    int
    ExitingCount     int
    Direction        string
    EnteringQueue    []*models.Vehicle
    ExitingQueue     []*models.Vehicle
    sync.Mutex
}