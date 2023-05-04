package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	backgroundColor = color.RGBA{R: 241, G: 246, B: 249, A: 255}
	heartColor      = color.RGBA{R: 244, G: 80, B: 80, A: 255}
)

var (
	whiteImage    = ebiten.NewImage(3, 3)
	whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	whiteImage.Fill(color.White)
}

type Game struct {
	debug  bool
	width  float64
	height float64
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF3) {
		g.debug = !g.debug
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)

	var path vector.Path

	for t := 0.0; t <= 2*math.Pi; t += 0.001 {
		x := 16 * math.Pow(math.Sin(t), 3)
		y := 13*math.Cos(t) - 5*math.Cos(2*t) - 3*math.Cos(3*t) - math.Cos(4*t)
		if t == 0 {
			path.MoveTo(float32(x), float32(y))
		} else {
			path.LineTo(float32(x), float32(y))
		}
	}
	path.Close()

	var vs []ebiten.Vertex
	var is []uint16
	vs, is = path.AppendVerticesAndIndicesForFilling(nil, nil)

	scale := 20.0 / float32(ebiten.DeviceScaleFactor())
	for i := range vs {
		vs[i].DstX = vs[i].DstX*scale + float32(g.width/2)
		vs[i].DstY = -vs[i].DstY*scale + float32(g.height/2)
		vs[i].SrcX = 1
		vs[i].SrcY = 1
		vs[i].ColorR = float32(heartColor.R) / float32(0xff)
		vs[i].ColorG = float32(heartColor.G) / float32(0xff)
		vs[i].ColorB = float32(heartColor.B) / float32(0xff)
		vs[i].ColorA = 1
	}

	op := &ebiten.DrawTrianglesOptions{}
	op.AntiAlias = true
	screen.DrawTriangles(vs, is, whiteSubImage, op)

	if g.debug {
		ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.2f\nTPS: %.2f\nVertices: %d\nIndices: %d",
			ebiten.ActualFPS(), ebiten.ActualTPS(), len(vs), len(is)))
	}
}

func (g *Game) Layout(_, _ int) (int, int) {
	panic("unreachable")
}

func (g *Game) LayoutF(outsideWidth, outsideHeight float64) (screenWidth, screenHeight float64) {
	g.width = outsideWidth
	g.height = outsideHeight
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowSize(1080, 720)
	ebiten.SetWindowTitle("Equation Of The Heart")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	fmt.Println("Stating...")
	if err := ebiten.RunGame(&Game{}); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Crashed: %s", err)
		os.Exit(1)
	}
	fmt.Println("Bey!")
}
