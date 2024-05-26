package grid

func IsInGrid(x, y, size int) bool {
	return x >= 0 && y >= 0 && x < size && y < size
}
