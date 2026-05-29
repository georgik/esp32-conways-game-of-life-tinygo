package conway

const (
	Width  = 30
	Height = 30
)

type Grid [Height][Width]int

func (g *Grid) CountCell(x, y int) int {
	var c int
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			nx := x + dx
			ny := y + dy
			if nx >= 0 && nx < Width && ny >= 0 && ny < Height {
				if g[ny][nx] > 0 {
					c++
				}
			}
		}
	}
	return c
}

func (g *Grid) NextStep() Grid {
	var ng Grid
	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			c := g.CountCell(x, y)
			if g[y][x] > 0 {
				if c == 2 || c == 3 {
					ng[y][x] = g[y][x] + 1
				} else {
					ng[y][x] = 0
				}
			} else {
				if c == 3 {
					ng[y][x] = 1
				}
			}
		}
	}
	return ng
}

func (g *Grid) RandomizeWith(randFunc func() uint32) {
	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			if randFunc()%100 < 35 {
				g[y][x] = 1
			} else {
				g[y][x] = 0
			}
		}
	}
}

func (g *Grid) GetCell(x, y int) int {
	if x >= 0 && x < Width && y >= 0 && y < Height {
		return g[y][x]
	}
	return 0
}
