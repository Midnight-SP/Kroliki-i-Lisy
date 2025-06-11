package main

import (
	"fmt"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	empty  = 0 // Puste pole
	grass  = 1 // Trawa
	rabbit = 2 // Królik
	fox    = 3 // Lis
)

type Cell struct {
	Ground int     // Warstwa ziemi (empty, grass)
	Animal int     // Warstwa zwierząt (empty, rabbit, fox)
	Energy float64 // Energia zwierzęcia (dla królików i lisów)
}

type World struct {
	Grid       [][]Cell // Plansza symulacji
	Width      int      // Szerokość planszy
	Height     int      // Wysokość planszy
	MaxGrass   int      // Maksymalna ilość trawy na polu
	GrowthRate float64  // Współczynnik wzrostu trawy
}

func NewWorld(width, height, maxGrass int, growthRate float64) *World {
	grid := make([][]Cell, height)
	for i := range grid {
		grid[i] = make([]Cell, width)
		for j := range grid[i] {
			grid[i][j] = Cell{Ground: empty, Animal: empty} // Domyślnie puste pole
		}
	}
	return &World{
		Grid:       grid,
		Width:      width,
		Height:     height,
		MaxGrass:   maxGrass,
		GrowthRate: growthRate,
	}
}

// Inicjalizacja planszy z losowym rozmieszczeniem trawy, królików i lisów
func (w *World) Initialize(rabbitCount, foxCount int) {
	rand.Seed(time.Now().UnixNano())

	// Losowe rozmieszczenie trawy
	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			if rand.Float64() < 0.5 { // 50% szans na trawę
				w.Grid[y][x].Ground = grass
			}
		}
	}

	// Losowe rozmieszczenie królików
	for i := 0; i < rabbitCount; i++ {
		x, y := rand.Intn(w.Width), rand.Intn(w.Height)
		w.Grid[y][x].Animal = rabbit
		w.Grid[y][x].Energy = 10.0
	}

	// Losowe rozmieszczenie lisów
	for i := 0; i < foxCount; i++ {
		x, y := rand.Intn(w.Width), rand.Intn(w.Height)
		w.Grid[y][x].Animal = fox
		w.Grid[y][x].Energy = 20.0
	}
}

// Wyświetlanie planszy w konsoli (do debugowania)
func (w *World) Print() {
	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			if w.Grid[y][x].Animal == rabbit {
				fmt.Print("R ")
			} else if w.Grid[y][x].Animal == fox {
				fmt.Print("F ")
			} else if w.Grid[y][x].Ground == grass {
				fmt.Print("G ")
			} else {
				fmt.Print(". ")
			}
		}
		fmt.Println()
	}
}

func (w *World) DrawWorld(cellSize int) {
	// Najpierw rysuj tło
	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			color := rl.LightGray
			if w.Grid[y][x].Ground == grass {
				color = rl.Green
			}
			rl.DrawRectangle(
				int32(x*cellSize), int32(y*cellSize),
				int32(cellSize), int32(cellSize),
				color,
			)
		}
	}

	// Potem rysuj zwierzęta (na wierzchu)
	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			if w.Grid[y][x].Animal == rabbit {
				rl.DrawCircle(
					int32(x*cellSize+cellSize/2),
					int32(y*cellSize+cellSize/2),
					float32(cellSize/3),
					rl.Blue,
				)
			} else if w.Grid[y][x].Animal == fox {
				rl.DrawCircle(
					int32(x*cellSize+cellSize/2),
					int32(y*cellSize+cellSize/2),
					float32(cellSize/3),
					rl.Red,
				)
			}
		}
	}
}

func (w *World) GrowGrass() {
	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			if w.Grid[y][x].Ground == empty && rand.Float64() < w.GrowthRate {
				w.Grid[y][x].Ground = grass
			}
		}
	}
}

func (w *World) MoveRabbits() {
	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			if w.Grid[y][x].Animal == rabbit {
				// Ogranicz ruch do sąsiednich pól (bez wyjścia poza planszę)
				newX := x + rand.Intn(3) - 1 // -1, 0, 1
				newY := y + rand.Intn(3) - 1 // -1, 0, 1

				// Sprawdź, czy nowa pozycja jest w granicach planszy
				if newX >= 0 && newX < w.Width && newY >= 0 && newY < w.Height {
					if w.Grid[newY][newX].Animal == empty {
						// Zjedz trawę, jeśli jest
						if w.Grid[newY][newX].Ground == grass {
							w.Grid[newY][newX].Ground = empty
							w.Grid[y][x].Energy += 5
						}

						// Przenieś królika
						w.Grid[newY][newX].Animal = rabbit
						w.Grid[newY][newX].Energy = w.Grid[y][x].Energy
						w.Grid[y][x].Animal = empty
						w.Grid[y][x].Energy = 0
					}
				}
			}
		}
	}
}

func (w *World) MoveFoxes() {
	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			if w.Grid[y][x].Animal == fox {
				newX := x + rand.Intn(3) - 1
				newY := y + rand.Intn(3) - 1

				if newX >= 0 && newX < w.Width && newY >= 0 && newY < w.Height {
					if w.Grid[newY][newX].Animal == empty {
						// Przenieś lisa
						w.Grid[newY][newX].Animal = fox
						w.Grid[newY][newX].Energy = w.Grid[y][x].Energy
						w.Grid[y][x].Animal = empty
						w.Grid[y][x].Energy = 0
					} else if w.Grid[newY][newX].Animal == rabbit {
						// Lis zjada królika
						w.Grid[newY][newX].Animal = fox
						w.Grid[newY][newX].Energy = w.Grid[y][x].Energy + 10
						w.Grid[y][x].Animal = empty
						w.Grid[y][x].Energy = 0
					}
				}
			}
		}
	}
}

func (w *World) UpdateEnergy() {
	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			if w.Grid[y][x].Animal == rabbit || w.Grid[y][x].Animal == fox {
				w.Grid[y][x].Energy -= 1
				if w.Grid[y][x].Energy <= 0 {
					// Zwierzę umiera z braku energii
					w.Grid[y][x].Animal = empty
					w.Grid[y][x].Energy = 0
				}
			}
		}
	}
}

func (w *World) SimulateWithVisualization(steps, cellSize int) {
	screenWidth := int32(w.Width * cellSize)
	screenHeight := int32(w.Height * cellSize)

	rl.InitWindow(screenWidth, screenHeight, "Symulacja Ekosystemu")
	defer rl.CloseWindow()

	rl.SetTargetFPS(10)

	// Buforowanie stanu
	renderState := w.Copy()
	updateChan := make(chan *World, 1)
	quitChan := make(chan bool)

	// Goroutine do aktualizacji stanu
	go func() {
		for step := 0; step < steps; step++ {
			select {
			case <-quitChan:
				return
			default:
				w.GrowGrass()
				w.MoveRabbits()
				w.MoveFoxes()
				w.UpdateEnergy()

				// Wyślij kopię do renderowania
				updateChan <- w.Copy()

				// Stałe tempo symulacji
				time.Sleep(100 * time.Millisecond)
			}
		}
		close(updateChan)
	}()

	// Główna pętla renderowania
	for !rl.WindowShouldClose() {
		// Odbierz najnowszy stan jeśli jest dostępny
		select {
		case newState := <-updateChan:
			renderState = newState
		default:
		}

		// Renderowanie
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)
		renderState.DrawWorld(cellSize)

		// Dodatkowe informacje
		currentAnimals := countAnimals(renderState)
		rl.DrawText(fmt.Sprintf("Króliki: %d  Lisy: %d",
			currentAnimals[rabbit], currentAnimals[fox]), 10, 10, 20, rl.Black)

		rl.EndDrawing()
	}

	quitChan <- true
}

// Optymalizowana wersja Copy - synchronizacja bez alokacji
func (w *World) SyncCopy(dest *World) {
	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			// Głęboka kopia struktury Cell
			dest.Grid[y][x].Ground = w.Grid[y][x].Ground
			dest.Grid[y][x].Animal = w.Grid[y][x].Animal
			dest.Grid[y][x].Energy = w.Grid[y][x].Energy
		}
	}
}

func (w *World) Copy() *World {
	newGrid := make([][]Cell, w.Height)
	for y := range w.Grid {
		newGrid[y] = make([]Cell, w.Width)
		for x := range w.Grid[y] {
			newGrid[y][x] = Cell{
				Ground: w.Grid[y][x].Ground,
				Animal: w.Grid[y][x].Animal,
				Energy: w.Grid[y][x].Energy,
			}
		}
	}
	return &World{
		Grid:       newGrid,
		Width:      w.Width,
		Height:     w.Height,
		MaxGrass:   w.MaxGrass,
		GrowthRate: w.GrowthRate,
	}
}

func countAnimals(w *World) map[int]int {
	counts := make(map[int]int)
	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			if w.Grid[y][x].Animal != empty {
				counts[w.Grid[y][x].Animal]++
			}
		}
	}
	return counts
}

func main() {
	// Parametry symulacji
	width, height := 20, 10
	maxGrass := 5
	growthRate := 0.1
	rabbitCount := 10
	foxCount := 5
	cellSize := 30 // Rozmiar komórki w pikselach

	// Tworzenie świata
	world := NewWorld(width, height, maxGrass, growthRate)
	world.Initialize(rabbitCount, foxCount)

	// Uruchomienie symulacji z wizualizacją
	world.SimulateWithVisualization(100, cellSize)
}
