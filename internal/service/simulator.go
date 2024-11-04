package service

import (
    "math"
    "math/rand"
    "sync"
    "time"
    "parking/internal/models"
    "parking/internal/domain"
    "parking/internal/ui"
)

type ParkingSimulator struct {
    parking         *domain.ParkingLot
    vehicles        chan *models.Vehicle
    updateUI        chan struct{}
    uiSimulator     *ui.UISimulator
    running         bool
    waitGroup       sync.WaitGroup
    waitingCars     int
    processedCars   int
    rng             *rand.Rand
}

func NewParkingSimulator(uiSimulator *ui.UISimulator) *ParkingSimulator {
    return &ParkingSimulator{
        parking:     &domain.ParkingLot{},
        vehicles:    make(chan *models.Vehicle, 100),
        updateUI:    make(chan struct{}, 1),
        uiSimulator: uiSimulator,
        rng:         rand.New(rand.NewSource(time.Now().UnixNano())),
    }
}

func (ps *ParkingSimulator) poissonArrival(lambda float64) time.Duration {
    L := math.Exp(-lambda)
    k := 0
    p := 1.0

    for p > L {
        k++
        p *= ps.rng.Float64()
    }
    
    duracionBase := 500 * time.Millisecond
    duracionExtra := time.Duration(k * 100) * time.Millisecond
    if duracionExtra > 2500 * time.Millisecond {
        duracionExtra = 2500 * time.Millisecond
    }
    return duracionBase + duracionExtra
}

func (ps *ParkingSimulator) generateVehicles() {
    totalCars := 0
    id := 1
    lambda := 1.0

    for ps.running && totalCars < 100 {
        vehicle := &models.Vehicle{
            ID:     id,
            Status: "esperando",
        }
        
        ps.vehicles <- vehicle
        ps.waitingCars++
        ps.triggerUIUpdate()
        
        id++
        totalCars++
        
        waitTime := ps.poissonArrival(lambda)
        time.Sleep(waitTime)
    }

    if totalCars >= 100 {
        close(ps.vehicles)
    }
}

func (ps *ParkingSimulator) handleVehicle(v *models.Vehicle) {
    defer ps.waitGroup.Done()

    for ps.running {
        ps.parking.Lock()
        if ps.parking.SpacesInUse < 20 {
            if ps.parking.Direction == "exiting" {
                ps.parking.EnteringQueue = append(ps.parking.EnteringQueue, v)
                ps.parking.Unlock()
                
                for ps.running {
                    time.Sleep(100 * time.Millisecond)
                    ps.parking.Lock()
                    if len(ps.parking.EnteringQueue) > 0 && ps.parking.EnteringQueue[0] == v && 
                       ps.parking.Direction != "exiting" {
                        ps.parking.Direction = "entering"
                        ps.parking.EnteringCount++
                        ps.parking.EnteringQueue = ps.parking.EnteringQueue[1:]
                        ps.parking.Unlock()
                        break
                    }
                    ps.parking.Unlock()
                }
                break
            } else if ps.parking.Direction == "" || ps.parking.Direction == "entering" {
                ps.parking.Direction = "entering"
                ps.parking.EnteringCount++
                ps.parking.Unlock()
                break
            }
        }
        ps.parking.Unlock()
        time.Sleep(100 * time.Millisecond)
    }

    v.Status = "entrando"
    ps.uiSimulator.UpdateEntranceColor("entering")
    ps.waitingCars--
    ps.triggerUIUpdate()
    time.Sleep(500 * time.Millisecond)
    
    ps.parking.Lock()
    spaceFound := -1
    for i := 0; i < 20; i++ {
        if ps.parking.Spaces[i] == 0 {
            spaceFound = i
            ps.parking.Spaces[i] = v.ID
            ps.parking.SpacesInUse++
            break
        }
    }
    ps.parking.Unlock()
    
    ps.parking.Lock()
    ps.parking.EnteringCount--
    if ps.parking.EnteringCount == 0 {
        ps.parking.Direction = ""
        if len(ps.parking.ExitingQueue) > 0 {
            ps.parking.Direction = "exiting"
        }
    }
    ps.parking.Unlock()
    
    ps.uiSimulator.UpdateEntranceColor("")
    
    if spaceFound != -1 {
        v.Status = "estacionado"
        v.Position = spaceFound
        ps.triggerUIUpdate()
        
        stayTime := 3 + ps.rng.Intn(3)
        time.Sleep(time.Duration(stayTime) * time.Second)
        
        for ps.running {
            ps.parking.Lock()
            if ps.parking.Direction == "entering" {
                ps.parking.ExitingQueue = append(ps.parking.ExitingQueue, v)
                ps.parking.Unlock()
                
                for ps.running {
                    time.Sleep(100 * time.Millisecond)
                    ps.parking.Lock()
                    if len(ps.parking.ExitingQueue) > 0 && ps.parking.ExitingQueue[0] == v && 
                       ps.parking.Direction != "entering" {
                        ps.parking.Direction = "exiting"
                        ps.parking.ExitingCount++
                        ps.parking.ExitingQueue = ps.parking.ExitingQueue[1:]
                        ps.parking.Unlock()
                        break
                    }
                    ps.parking.Unlock()
                }
                break
            } else {
                ps.parking.Direction = "exiting"
                ps.parking.ExitingCount++
                ps.parking.Unlock()
                break
            }
        }

        v.Status = "saliendo"
        ps.uiSimulator.UpdateEntranceColor("exiting")
        ps.triggerUIUpdate()
        
        ps.parking.Lock()
        ps.parking.Spaces[v.Position] = 0
        ps.parking.SpacesInUse--
        ps.processedCars++
        ps.parking.Unlock()
        
        time.Sleep(500 * time.Millisecond)
        
        ps.parking.Lock()
        ps.parking.ExitingCount--
        if ps.parking.ExitingCount == 0 {
            ps.parking.Direction = ""
            if len(ps.parking.EnteringQueue) > 0 {
                ps.parking.Direction = "entering"
            }
        }
        ps.parking.Unlock()
        
        ps.uiSimulator.UpdateEntranceColor("")
    }
}

func (ps *ParkingSimulator) processVehicles() {
    for vehicle := range ps.vehicles {
        if !ps.running {
            break
        }
        ps.waitGroup.Add(1)
        go ps.handleVehicle(vehicle)
    }
}

func (ps *ParkingSimulator) updateUILoop() {
    for range ps.updateUI {
        if !ps.running {
            break
        }
        
        ps.parking.Lock()
        for i := 0; i < 20; i++ {
            ps.uiSimulator.UpdateSpaceUI(i, ps.parking.Spaces[i])
        }
        
        enteringQueueSize := len(ps.parking.EnteringQueue)
        exitingQueueSize := len(ps.parking.ExitingQueue)
        
        ps.uiSimulator.UpdateLabels(
            ps.parking.SpacesInUse,
            ps.waitingCars,
            enteringQueueSize,
            exitingQueueSize,
            ps.processedCars,
            ps.parking.EnteringCount,
            ps.parking.ExitingCount,
        )
        ps.parking.Unlock()
    }
}

func (ps *ParkingSimulator) triggerUIUpdate() {
    select {
    case ps.updateUI <- struct{}{}:
    default:
    }
}

func (ps *ParkingSimulator) Start() {
    ps.running = true

    go ps.generateVehicles()
    go ps.processVehicles()
    go ps.updateUILoop()
    
    go func() {
        for ps.running {
            if ps.processedCars >= 100 {
                time.Sleep(2 * time.Second)
                ps.running = false
                break
            }
            time.Sleep(time.Second)
        }
    }()
}

func (ps *ParkingSimulator) Stop() {
    ps.running = false
}