package grid

import (
	"math"

	"golang.org/x/exp/rand"
)

func GenerateGrid(size, difficulty int) (grid Grid) {

	for i := 0; i < size; i++ {
		row := []int{}
		for j := 0; j < size; j++ {
			row = append(row, 0)
		}
		grid = append(grid, row)
	}
	//Calculating the number of walls to be placed
	nbMurs := int(math.Round(ratiosWallsArea[difficulty] * float64(size*size)))

	//Wall Placement
	cpt := 0
	x := rand.Intn(size)
	y := rand.Intn(size)
	for cpt < nbMurs {
		for grid[x][y] != 0 {
			x = rand.Intn(size)
			y = rand.Intn(size)
		}
		grid[x][y] = -1
		cpt++
	}
	//Bulb placement
	nbLitCases := 0
	nbCasesTotalExtinguished := size*size - cpt //Number of squares to light (the whole grid - the walls)
	x = rand.Intn(size)
	y = rand.Intn(size)
	wallOnPath := false
	for nbLitCases < nbCasesTotalExtinguished {
		//Placement of light bulbs
		for grid[x][y] != 0 {
			x = rand.Intn(size)
			y = rand.Intn(size)
		}
		grid[x][y] = 5
		nbCasesTotalExtinguished-- //One less box to light up because a light bulb is placed in it
		//Updating the number of illuminated boxes
		wallOnPath = false
		//Column Down
		for j := x + 1; j < size && !wallOnPath; j++ {
			if grid[j][y] == -1 {
				wallOnPath = true
			}
			if !wallOnPath && grid[j][y] == 0 {
				grid[j][y] = -4
				nbLitCases++
			}
		}
		wallOnPath = false
		//Column up
		for j := x - 1; j >= 0 && !wallOnPath; j-- {
			if grid[j][y] == -1 {
				wallOnPath = true
			}
			if !wallOnPath && grid[j][y] == 0 {
				grid[j][y] = -4
				nbLitCases++
			}
		}
		wallOnPath = false
		//line vers la droite
		for j := y + 1; j < size && !wallOnPath; j++ {
			if grid[x][j] == -1 {
				wallOnPath = true
			}
			if !wallOnPath && grid[x][j] == 0 {
				grid[x][j] = -4
				nbLitCases++
			}
		}
		wallOnPath = false
		//Line to the left
		for j := y - 1; j >= 0 && !wallOnPath; j-- {
			if grid[x][j] == -1 {
				wallOnPath = true
			}
			if !wallOnPath && grid[x][j] == 0 {
				grid[x][j] = -4
				nbLitCases++
			}
		}
	}
	//Temporarily set to 4 squares are set back to -2
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if grid[i][j] == -4 {
				grid[i][j] = -2
			}
		}
	}
	//Setting up wall constraints
	nbConstraintsTotal := int(math.Round(ratiosNumberWalls[difficulty] * float64(cpt)))
	nbConstraints := 0
	//Placing the Right Number of Constraints
	for nbConstraints < nbConstraintsTotal {
		//Browse the grid
		for i := 0; i < size; i++ {
			for j := 0; j < size; j++ {
				if grid[i][j] == -1 {
					// If we have a wall
					if rand.Float64() < ratiosNumberWalls[difficulty] {
						//Random
						grid[i][j] = -5 //Constraint to be placed
						nbConstraints++
					}
				}
			}
		}
	}
	//Update constraints with the number of bulbs around
	for j := 0; j < size; j++ {
		for k := 0; k < size; k++ {
			//If you have to place a constraint
			if grid[j][k] == -5 {
				ampoulesNear := 0
				nbCasesVidesNear := 0
				//We look at the above, below and sides for the number of bulbs and the number of empty boxes
				if IsInGrid(k+1, j, size) {
					if grid[j][k+1] == 5 {
						ampoulesNear++
						nbCasesVidesNear++
					} else if grid[j][k+1] == -2 {
						nbCasesVidesNear++
					}
				}
				if IsInGrid(k-1, j, size) {
					if grid[j][k-1] == 5 {
						ampoulesNear++
						nbCasesVidesNear++
					} else if grid[j][k-1] == -2 {
						nbCasesVidesNear++
					}
				}
				if IsInGrid(k, j+1, size) {
					if grid[j+1][k] == 5 {
						ampoulesNear++
						nbCasesVidesNear++
					} else if grid[j+1][k] == -2 {
						nbCasesVidesNear++
					}
				}
				if IsInGrid(k, j-1, size) {
					if grid[j-1][k] == 5 {
						ampoulesNear++
						nbCasesVidesNear++
					} else if grid[j-1][k] == -2 {
						nbCasesVidesNear++
					}
				}
				//Constraints are removed from walls that are not next to a possible location
				if nbCasesVidesNear != 0 {
					grid[j][k] = ampoulesNear
				} else {
					grid[j][k] = -1
				}
			}
		}
	}
	//The squares temporarily set to 5 (the light bulbs) are reset to -2
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if grid[i][j] == 5 {
				grid[i][j] = -2
			}
		}
	}
	return
}
