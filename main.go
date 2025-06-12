package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"

	"os/exec"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

type SimParams struct {
	Width, Height  int
	Rabbits, Foxes int
	GrowthRate     float64
}

func ShowMenu() SimParams {
    rl.InitWindow(480, 320, "Ustawienia symulacji")
    defer rl.CloseWindow()
    rl.SetTargetFPS(60)

    params := SimParams{Width: 32, Height: 16, Rabbits: 12, Foxes: 6, GrowthRate: 0.1}
    selected := 0
    options := []string{"Szerokość", "Wysokość", "Króliki", "Lisy", "Wzrost trawy", "Start"}

    // Wczytaj teksturę lisa do menu
    foxMenuTexture := rl.LoadTexture("fox.png")
    defer rl.UnloadTexture(foxMenuTexture)

    for !rl.WindowShouldClose() {
        rl.BeginDrawing()
        rl.ClearBackground(rl.RayWhite)

        // Rysuj lisa w tle menu (np. na środku, lekko przezroczysty)
        foxScale := float32(3.0)
        foxX := float32(240 - int(foxMenuTexture.Width)*int(foxScale)/2)
        foxY := float32(60)
        rl.DrawTextureEx(foxMenuTexture, rl.NewVector2(foxX, foxY), 0, foxScale, rl.Fade(rl.White, 0.18))

        rl.DrawText("Menu symulacji", 120, 20, 28, rl.Black)

        for i, opt := range options {
            color := rl.Black
            if i == selected {
                color = rl.Red
            }
            val := ""
            switch i {
            case 0:
                val = fmt.Sprintf("%d", params.Width)
            case 1:
                val = fmt.Sprintf("%d", params.Height)
            case 2:
                val = fmt.Sprintf("%d", params.Rabbits)
            case 3:
                val = fmt.Sprintf("%d", params.Foxes)
            case 4:
                val = fmt.Sprintf("%.2f", params.GrowthRate)
            }
            rl.DrawText(fmt.Sprintf("%s: %s", opt, val), 80, int32(70+30*i), 24, color)
        }

        rl.DrawText("Strzałki: wybór/opcja, Enter: start", 40, 270, 18, rl.Gray)
        rl.EndDrawing()

        if rl.IsKeyPressed(rl.KeyDown) {
            selected = (selected + 1) % len(options)
        }
        if rl.IsKeyPressed(rl.KeyUp) {
            selected = (selected - 1 + len(options)) % len(options)
        }
        if selected < 5 {
            if rl.IsKeyPressed(rl.KeyRight) {
                switch selected {
                case 0:
                    params.Width += 2
                case 1:
                    params.Height += 2
                case 2:
                    params.Rabbits++
                case 3:
                    params.Foxes++
                case 4:
                    params.GrowthRate += 0.01
                }
            }
            if rl.IsKeyPressed(rl.KeyLeft) {
                switch selected {
                case 0:
                    if params.Width > 4 {
                        params.Width -= 2
                    }
                case 1:
                    if params.Height > 4 {
                        params.Height -= 2
                    }
                case 2:
                    if params.Rabbits > 1 {
                        params.Rabbits--
                    }
                case 3:
                    if params.Foxes > 1 {
                        params.Foxes--
                    }
                case 4:
                    if params.GrowthRate > 0.01 {
                        params.GrowthRate -= 0.01
                    }
                }
            }
        }
        if selected == 5 && rl.IsKeyPressed(rl.KeyEnter) {
            break
        }
    }
    return params
}

type Cell struct {
    Ground            int     // 0=brak trawy, 1=short, 2=medium, 3=tall
    Animal            int     // 0=empty, 4=rabbit, 5=fox
    Energy            float64
    ReproduceCooldown int
    Age               int
}

const (
    empty  = 0
    grassShort = 1
    grassMedium = 2
    grassTall = 3
    rabbit = 4
    fox    = 5

    rabbitReproduceEnergy = 14.0
    foxReproduceEnergy    = 28.0

    rabbitCooldown = 6
    foxCooldown    = 10
)

var (
    texEmpty      rl.Texture2D
    texGrassShort rl.Texture2D
    texGrassMed   rl.Texture2D
    texGrassTall  rl.Texture2D
    texRabbit     rl.Texture2D
    texFox        rl.Texture2D
)

func loadTextures() {
    texEmpty = rl.LoadTexture("empty.png")
    texGrassShort = rl.LoadTexture("grass_short.png")
    texGrassMed = rl.LoadTexture("grass_medium.png")
    texGrassTall = rl.LoadTexture("grass_tall.png")
    texRabbit = rl.LoadTexture("rabbit.png")
    texFox = rl.LoadTexture("fox.png")
}

func unloadTextures() {
    rl.UnloadTexture(texEmpty)
    rl.UnloadTexture(texGrassShort)
    rl.UnloadTexture(texGrassMed)
    rl.UnloadTexture(texGrassTall)
    rl.UnloadTexture(texRabbit)
    rl.UnloadTexture(texFox)
}

type World struct {
	Grid       [][]Cell
	Width      int
	Height     int
	MaxGrass   int
	GrowthRate float64
}

func NewWorld(width, height, maxGrass int, growthRate float64) *World {
	grid := make([][]Cell, height)
	for i := range grid {
		grid[i] = make([]Cell, width)
		for j := range grid[i] {
			grid[i][j] = Cell{Ground: empty, Animal: empty}
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

    for y := 0; y < w.Height; y++ {
        for x := 0; x < w.Width; x++ {
            r := rand.Float64()
            if r < 0.33 {
                w.Grid[y][x].Ground = grassShort
            } else if r < 0.66 {
                w.Grid[y][x].Ground = grassMedium
            } else {
                w.Grid[y][x].Ground = grassTall
            }
        }
    }

    for i := 0; i < rabbitCount; i++ {
        x, y := rand.Intn(w.Width), rand.Intn(w.Height)
        w.Grid[y][x].Animal = rabbit
        w.Grid[y][x].Energy = 10.0
    }

    for i := 0; i < foxCount; i++ {
        x, y := rand.Intn(w.Width), rand.Intn(w.Height)
        w.Grid[y][x].Animal = fox
        w.Grid[y][x].Energy = 20.0
    }
}

func (w *World) DrawWorld(cellSize int) {
	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			pos := rl.NewVector2(float32(x*cellSize), float32(y*cellSize))
			rl.DrawTextureEx(texEmpty, pos, 0, float32(cellSize)/float32(texEmpty.Width), rl.White)
			switch w.Grid[y][x].Ground {
			case grassShort:
				rl.DrawTextureEx(texGrassShort, pos, 0, float32(cellSize)/float32(texGrassShort.Width), rl.White)
			case grassMedium:
				rl.DrawTextureEx(texGrassMed, pos, 0, float32(cellSize)/float32(texGrassMed.Width), rl.White)
			case grassTall:
				rl.DrawTextureEx(texGrassTall, pos, 0, float32(cellSize)/float32(texGrassTall.Width), rl.White)
			}
			if w.Grid[y][x].Animal == rabbit {
				rl.DrawTextureEx(texRabbit, pos, 0, float32(cellSize)/float32(texRabbit.Width), rl.White)
			} else if w.Grid[y][x].Animal == fox {
				rl.DrawTextureEx(texFox, pos, 0, float32(cellSize)/float32(texFox.Width), rl.White)
			}
		}
	}
}

func neighbors(x, y, width, height int) [][2]int {
	var result [][2]int
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}
			nx, ny := x+dx, y+dy
			if nx >= 0 && nx < width && ny >= 0 && ny < height {
				result = append(result, [2]int{nx, ny})
			}
		}
	}
	return result
}

func (w *World) Copy() *World {
	newGrid := make([][]Cell, w.Height)
	for y := 0; y < w.Height; y++ {
		newGrid[y] = make([]Cell, w.Width)
		copy(newGrid[y], w.Grid[y])
	}
	return &World{
		Grid:       newGrid,
		Width:      w.Width,
		Height:     w.Height,
		MaxGrass:   w.MaxGrass,
		GrowthRate: w.GrowthRate,
	}
}

func (w *World) GrowGrass() {
    for y := 0; y < w.Height; y++ {
        for x := 0; x < w.Width; x++ {
            if w.Grid[y][x].Ground == empty {
                hasGrassNeighbor := false
                for _, n := range neighbors(x, y, w.Width, w.Height) {
                    if w.Grid[n[1]][n[0]].Ground > empty {
                        hasGrassNeighbor = true
                        break
                    }
                }
                if hasGrassNeighbor && rand.Float64() < w.GrowthRate {
                    w.Grid[y][x].Ground = grassShort
                }
            } else if w.Grid[y][x].Ground == grassShort {
                if rand.Float64() < w.GrowthRate {
                    w.Grid[y][x].Ground = grassMedium
                }
            } else if w.Grid[y][x].Ground == grassMedium {
                if rand.Float64() < w.GrowthRate {
                    w.Grid[y][x].Ground = grassTall
                }
            }
        }
    }
}

func (w *World) MoveRabbits() {
    newGrid := w.Copy().Grid

    coords := make([][2]int, 0, w.Width*w.Height)
    for y := 0; y < w.Height; y++ {
        for x := 0; x < w.Width; x++ {
            coords = append(coords, [2]int{x, y})
        }
    }
    rand.Shuffle(len(coords), func(i, j int) { coords[i], coords[j] = coords[j], coords[i] })

    for _, pos := range coords {
        x, y := pos[0], pos[1]
        cell := w.Grid[y][x]
        if cell.Animal == rabbit && cell.ReproduceCooldown == 0 {
            ns := neighbors(x, y, w.Width, w.Height)
            done := false

            // 1. Ucieczka przed lisem
            foxes := [][2]int{}
            for _, n := range ns {
                if w.Grid[n[1]][n[0]].Animal == fox {
                    foxes = append(foxes, n)
                }
            }
            if len(foxes) > 0 && !done {
                maxDist := -1.0
                var best [2]int
                for _, n := range ns {
                    if w.Grid[n[1]][n[0]].Animal == empty {
                        minDist := 1000.0
                        for _, f := range foxes {
                            dx := float64(n[0] - f[0])
                            dy := float64(n[1] - f[1])
                            dist := dx*dx + dy*dy
                            if dist < minDist {
                                minDist = dist
                            }
                        }
                        if minDist > maxDist {
                            maxDist = minDist
                            best = n
                        }
                    }
                }
                if maxDist >= 0 {
                    newGrid[best[1]][best[0]] = cell
                    newGrid[y][x] = Cell{Ground: w.Grid[y][x].Ground}
                    cell.Age++
                    done = true
                }
            }

            // 2. Szukanie trawy gdy głodny
            if !done && cell.Energy < rabbitReproduceEnergy {
                for _, n := range ns {
                    ng := w.Grid[n[1]][n[0]]
                    if ng.Animal == empty && ng.Ground > empty {
                        // Jeśli bardzo głodny, zjada całą trawę
                        if cell.Energy < rabbitReproduceEnergy/2 {
                            cell.Energy += float64(ng.Ground) * 8
                            cell.Ground = w.Grid[y][x].Ground
                            newGrid[n[1]][n[0]] = cell
                            newGrid[y][x] = Cell{Ground: w.Grid[y][x].Ground}
                            newGrid[n[1]][n[0]].Ground = empty
                        } else {
                            cell.Energy += 6
                            cell.Ground = w.Grid[y][x].Ground
                            newGrid[n[1]][n[0]] = cell
                            newGrid[y][x] = Cell{Ground: w.Grid[y][x].Ground}
                            // Zmniejsz stadium trawy o 1
                            if ng.Ground == grassTall {
                                newGrid[n[1]][n[0]].Ground = grassMedium
                            } else if ng.Ground == grassMedium {
                                newGrid[n[1]][n[0]].Ground = grassShort
                            } else {
                                newGrid[n[1]][n[0]].Ground = empty
                            }
                        }
                        cell.Age++
                        done = true
                        break
                    }
                }
            }

            // 3. Szukanie królika do rozmnożenia
            if !done && cell.Energy >= rabbitReproduceEnergy {
                for _, n := range ns {
                    other := w.Grid[n[1]][n[0]]
                    if other.Animal == rabbit && other.ReproduceCooldown == 0 {
                        if y < n[1] || (y == n[1] && x < n[0]) {
                            for _, emptyN := range ns {
                                if w.Grid[emptyN[1]][emptyN[0]].Animal == empty {
                                    newGrid[emptyN[1]][emptyN[0]] = Cell{
                                        Ground:            w.Grid[emptyN[1]][emptyN[0]].Ground,
                                        Animal:            rabbit,
                                        Energy:            cell.Energy / 2,
                                        ReproduceCooldown: rabbitCooldown,
                                        Age:               0,
                                    }
                                    cell.Energy = cell.Energy / 2
                                    cell.ReproduceCooldown = rabbitCooldown
                                    newGrid[y][x] = cell
                                    cell.Age++
                                    done = true
                                    break
                                }
                            }
                        }
                        if done {
                            break
                        }
                    }
                }
            }

            // 4. Ruch losowy jeśli nic innego nie zadziałało
            if !done {
                rand.Shuffle(len(ns), func(i, j int) { ns[i], ns[j] = ns[j], ns[i] })
                for _, n := range ns {
                    if w.Grid[n[1]][n[0]].Animal == empty {
                        newGrid[n[1]][n[0]] = cell
                        newGrid[y][x] = Cell{Ground: w.Grid[y][x].Ground}
                        cell.Age++
                        break
                    }
                }
            }
        }
    }
    w.Grid = newGrid
}

func (w *World) MoveFoxes() {
	newGrid := w.Copy().Grid

	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			cell := w.Grid[y][x]
			if cell.Animal == fox && cell.ReproduceCooldown == 0 {
				ns := neighbors(x, y, w.Width, w.Height)
				done := false

				// 1. Szukanie królika gdy bardzo głodny
				if cell.Energy < foxReproduceEnergy/2 && !done {
					for attempt := 0; attempt < 2; attempt++ {
						for _, n := range ns {
							if w.Grid[n[1]][n[0]].Animal == rabbit {
								cell.Energy += 20 // zwiększ energię po zjedzeniu królika
								newGrid[n[1]][n[0]] = cell
								newGrid[y][x] = Cell{Ground: w.Grid[y][x].Ground}
								continue
							}
						}
					}
				}

				// 2. Szukanie królika gdy głodny (standardowo)
				if !done && cell.Energy < foxReproduceEnergy {
					for _, n := range ns {
						if w.Grid[n[1]][n[0]].Animal == rabbit {
							cell.Energy += 20 // zwiększ energię po zjedzeniu królika
							newGrid[n[1]][n[0]] = cell
							newGrid[y][x] = Cell{Ground: w.Grid[y][x].Ground}
							cell.Age++
							done = true
							break
						}
					}
				}

				// 3. Szukanie lisa do rozmnożenia
				if !done && cell.Energy >= foxReproduceEnergy {
					for _, n := range ns {
						other := w.Grid[n[1]][n[0]]
						if other.Animal == fox && other.ReproduceCooldown == 0 {
							if y < n[1] || (y == n[1] && x < n[0]) {
								for _, emptyN := range ns {
									if w.Grid[emptyN[1]][emptyN[0]].Animal == empty {
										newGrid[emptyN[1]][emptyN[0]] = Cell{
											Ground:            w.Grid[emptyN[1]][emptyN[0]].Ground,
											Animal:            fox,
											Energy:            cell.Energy / 2,
											ReproduceCooldown: foxCooldown,
											Age:               0,
										}
										cell.Energy = cell.Energy / 2
										cell.ReproduceCooldown = foxCooldown
										newGrid[y][x] = cell
										cell.Age++
										done = true
										break
									}
								}
							}
							if done {
								break
							}
						}
					}
				}

				// 4. Ruch losowy jeśli nic innego nie zadziałało
				if !done {
					rand.Shuffle(len(ns), func(i, j int) { ns[i], ns[j] = ns[j], ns[i] })
					for _, n := range ns {
						if w.Grid[n[1]][n[0]].Animal == empty {
							newGrid[n[1]][n[0]] = cell
							newGrid[y][x] = Cell{Ground: w.Grid[y][x].Ground}
							cell.Age++
							break
						}
					}
				}
			}
		}
	}
	w.Grid = newGrid
}

func (w *World) UpdateEnergy() {
	for y := 0; y < w.Height; y++ {
		for x := 0; x < w.Width; x++ {
			if w.Grid[y][x].Animal == rabbit || w.Grid[y][x].Animal == fox {
				// Zużycie energii rośnie z wiekiem
				energyLoss := 1.0 + float64(w.Grid[y][x].Age)/10.0
				w.Grid[y][x].Energy -= energyLoss
				if w.Grid[y][x].ReproduceCooldown > 0 {
					w.Grid[y][x].ReproduceCooldown--
				}
				if w.Grid[y][x].Energy <= 0 {
					w.Grid[y][x].Animal = empty
					w.Grid[y][x].Energy = 0
					w.Grid[y][x].ReproduceCooldown = 0
					w.Grid[y][x].Age = 0
				}
			}
		}
	}
}

var popHistory []struct{ Rabbits, Foxes int }
var paused bool

func (w *World) SimulateWithVisualization(cellSize int, worldHeight int, plotPreviewHeight int) {
	rl.SetTargetFPS(10)
	renderState := w.Copy()
	updateChan := make(chan *World, 1)
	quitChan := make(chan struct{})
	popHistory = nil

	go func() {
		for {
			select {
			case <-quitChan:
				close(updateChan)
				return
			default:
				w.GrowGrass()
				w.MoveRabbits()
				w.MoveFoxes()
				w.UpdateEnergy()

				animals := countAnimals(w)
				popHistory = append(popHistory, struct{ Rabbits, Foxes int }{
					Rabbits: animals[rabbit], Foxes: animals[fox],
				})

				select {
				case updateChan <- w.Copy():
				case <-quitChan:
					close(updateChan)
					return
				}

				if animals[rabbit]+animals[fox] == 0 {
					close(updateChan)
					return
				}
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
	var plotTexture rl.Texture2D
    var plotImageLoaded bool
    plotUpdateCounter := 0

loop:
    for !rl.WindowShouldClose() {
        if rl.IsKeyPressed(rl.KeySpace) {
            paused = !paused
        }

        if !paused {
            select {
            case newState, ok := <-updateChan:
                if ok {
                    renderState = newState
                } else {
                    break loop
                }
            default:
            }
        }

        rl.BeginDrawing()
        rl.ClearBackground(rl.RayWhite)
        renderState.DrawWorld(cellSize)

        currentAnimals := countAnimals(renderState)
        rl.DrawText(fmt.Sprintf("Króliki: %d  Lisy: %d",
            currentAnimals[rabbit], currentAnimals[fox]), 10, 10, 20, rl.Black)

        if paused {
            rl.DrawText("PAUZA (spacja)", 10, 40, 20, rl.Red)
        }

        plotUpdateCounter++
        if plotUpdateCounter >= 10 {
            ShowPlot()
            if plotImageLoaded {
                rl.UnloadTexture(plotTexture)
            }
            img := rl.LoadImage("populacje_preview.png")
			plotTexture = rl.LoadTextureFromImage(img)
			rl.UnloadImage(img)
            plotImageLoaded = true
            plotUpdateCounter = 0
        }

        if plotImageLoaded {
			plotW := int32(plotTexture.Width)
			plotH := int32(plotTexture.Height)
			winW := rl.GetScreenWidth()
			rl.DrawTexturePro(
				plotTexture,
				rl.NewRectangle(0, 0, float32(plotW), float32(plotH)),
				rl.NewRectangle(0, float32(cellSize*worldHeight), float32(winW), float32(plotPreviewHeight)),
				rl.NewVector2(0, 0),
				0,
				rl.White,
			)
		}

        rl.EndDrawing()
    }

    if plotImageLoaded {
        rl.UnloadTexture(plotTexture)
    }

    close(quitChan)
    for range updateChan {
    }

    ShowPlot()
    openImage("populacje.png")
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

func ShowPlot() {
	p := plot.New()
	p.Title.Text = "Populacje w czasie"
	p.X.Label.Text = "Tura"
	p.Y.Label.Text = "Liczebność"

	rabbits := make(plotter.XYs, len(popHistory))
	foxes := make(plotter.XYs, len(popHistory))
	for i, v := range popHistory {
		rabbits[i].X = float64(i)
		rabbits[i].Y = float64(v.Rabbits)
		foxes[i].X = float64(i)
		foxes[i].Y = float64(v.Foxes)
	}
	l1, _ := plotter.NewLine(rabbits)
	l2, _ := plotter.NewLine(foxes)
	l1.Color = plotter.DefaultLineStyle.Color
	l2.Color = plotter.DefaultLineStyle.Color
	l2.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
	p.Add(l1, l2)
	p.Legend.Add("Króliki", l1)
	p.Legend.Add("Lisy", l2)
	p.Legend.Top = true

	p.Save(8*vg.Inch, 4*vg.Inch, "populacje.png")
	p.Save(8*vg.Inch, 1.5*vg.Inch, "populacje_preview.png")
}

func openImage(filename string) {
	switch runtime.GOOS {
	case "windows":
		exec.Command("rundll32", "url.dll,FileProtocolHandler", filename).Start()
	case "darwin":
		exec.Command("open", filename).Start()
	default:
		exec.Command("xdg-open", filename).Start()
	}
}

func main() {
	params := ShowMenu()

	world := NewWorld(params.Width, params.Height, 8, params.GrowthRate)
	world.Initialize(params.Rabbits, params.Foxes)

	cellSize := 32
	plotPreviewHeight := int(float32(params.Width*cellSize) * 1.5 / 8.0)
	rl.InitWindow(int32(params.Width*cellSize), int32(params.Height*cellSize+plotPreviewHeight), "Symulacja Ekosystemu")
	loadTextures()
	defer unloadTextures()
	defer rl.CloseWindow()

	world.SimulateWithVisualization(32, params.Height, plotPreviewHeight)
}
