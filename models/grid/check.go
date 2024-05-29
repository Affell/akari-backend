package grid

func (startGrid Grid) CheckSolution(playerGrid Grid) bool {

	if len(startGrid) != len(playerGrid) {
		return false
	}

	var ligths Grid
	n := len(startGrid)
	for i := 0; i < n; i++ {
		var l []int
		for j := 0; j < n; j++ {
			l = append(l, startGrid[i][j])
			if playerGrid[i][j] >= 5 {
				playerGrid[i][j] = -3
			}
			if playerGrid[i][j] <= -4 {
				playerGrid[i][j] = -2
			}
			if startGrid[i][j] >= -1 && startGrid[i][j] <= 4 && (playerGrid[i][j] < -1 || playerGrid[i][j] > 4) {
				return false
			}
			if startGrid[i][j] >= -1 && startGrid[i][j] <= 4 && (playerGrid[i][j] < -1 || playerGrid[i][j] > 4) {
				return false
			}
		}
		ligths = append(ligths, l)
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if playerGrid[i][j] == -3 {
				for k := i + 1; k < n && startGrid[k][j] == -2; k++ {
					if playerGrid[k][j] == -3 {
						return false
					}
					ligths[k][j] = -4
				}
				for k := i - 1; k >= 0 && startGrid[k][j] == -2; k-- {
					if playerGrid[k][j] == -3 {
						return false
					}
					ligths[k][j] = -4
				}
				for k := j + 1; k < n && startGrid[i][k] == -2; k++ {
					if playerGrid[i][k] == -3 {
						return false
					}
					ligths[i][k] = -4
				}
				for k := j - 1; k >= 0 && startGrid[k][j] == -2; k-- {
					if playerGrid[i][k] == -3 {
						return false
					}
					ligths[i][k] = -4
				}
			}
		}
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {

			if startGrid[i][j] == -2 && ligths[i][j] != -4 {
				return false
			}

			if startGrid[i][j] >= 0 && startGrid[i][j] <= 4 {
				// Wall with condition
				nb := 0 // Counter to count the number of bulbs around the cell
				// Condition to count bulbs north, south, east, west
				if i-1 >= 0 && playerGrid[i-1][j] == -3 {
					nb++
				}
				if i+1 < n && playerGrid[i+1][j] == -3 {
					nb++
				}
				if j-1 >= 0 && playerGrid[i][j-1] == -3 {
					nb++
				}
				if j+1 < n && playerGrid[i][j+1] == -3 {
					nb++
				}
				if nb != playerGrid[i][j] {
					// If the expected number does not match the found number then return false
					return false
				}
			}
		}
	}
	return true
}
