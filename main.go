package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	_ "image/png"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var image *ebiten.Image

func init() {
	var err error
	image, _, err = ebitenutil.NewImageFromFile("./image.png")
	if err != nil {
		panic(err)
	}
}

const (
	screenWidth  = 1280
	screenHeight = 720
	gridSize     = 20
	Xgrid        = screenWidth / gridSize
	Ygrid        = screenHeight / gridSize
)
const (
	dirNone = iota
	dirUp
	dirDown
	dirLeft
	dirRight
)

func (g *Game) CollideWithWall() bool {
	if g.SnakeBody[0].x <= 0 ||
		g.SnakeBody[0].y <= 0 ||
		g.SnakeBody[0].x > Xgrid ||
		g.SnakeBody[0].y > Ygrid {
		return true
	}
	return false

}

func (g *Game) CollideWithApple() bool {
	if g.SnakeBody[0].x == g.Apple.x && g.SnakeBody[0].y == g.Apple.y {
		g.SnakeBody[0].x = g.Apple.x
		g.SnakeBody[0].y = g.Apple.y
		return true
	}
	return false
}

func (g *Game) CollideWithSelf() bool {
	for _, body := range g.SnakeBody[1:] {
		if g.SnakeBody[0].x == body.x && g.SnakeBody[0].y == body.y {
			return true
		}
	}
	return false
}

type Position struct {
	x int
	y int
}

type Game struct {
	Direction int        //移动方向
	SnakeBody []Position //蛇的身体
	Apple     Position   //苹果的方位
	Timer     int
	MoveTime  int
	Level     int
}

func (g *Game) NeedMove() bool {
	return g.Timer%g.MoveTime == 0
}

func (g *Game) Reset() {
	g.SnakeBody = []Position{
		{
			x: Xgrid / 2,
			y: Xgrid / 2,
		},
	}
	g.Apple = Position{
		x: rand.Intn(Xgrid - 1),
		y: rand.Intn(Ygrid - 1),
	}
	g.Direction = dirNone
	g.Timer = 1
	g.MoveTime = 3
	g.Level = 1
}

func NewGame() *Game {
	g := Game{
		Apple: Position{
			x: Xgrid / 2,
			y: Ygrid / 2,
		},
		SnakeBody: make([]Position, 1),
		Direction: dirNone,
		Timer:     1,
		MoveTime:  3,
		Level:     1,
	}
	g.SnakeBody[0].x = Xgrid / 2
	g.SnakeBody[0].y = Ygrid / 2
	return &g
}

func (g *Game) Update() error {
	//这是蛇的移动方向的设置，不可以掉头
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		if g.Direction != dirDown {
			g.Direction = dirUp
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		if g.Direction != dirUp {
			g.Direction = dirDown
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		if g.Direction != dirRight {
			g.Direction = dirLeft
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		if g.Direction != dirLeft {
			g.Direction = dirRight
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		panic("game over")
	}
	//处理吃到苹果，撞墙，碰到自己
	if g.NeedMove() {
		if g.CollideWithWall() || g.CollideWithSelf() {
			g.Reset()
		} else if g.CollideWithApple() {
			// 增加蛇的长度
			g.SnakeBody = append(g.SnakeBody, Position{
				x: g.SnakeBody[len(g.SnakeBody)-1].x,
				y: g.SnakeBody[len(g.SnakeBody)-1].y,
			})
			g.Apple = Position{
				x: rand.Intn(Xgrid - 1),
				y: rand.Intn(Ygrid - 1),
			}

		}
		//不断增加难度
		if len(g.SnakeBody) > 0 && len(g.SnakeBody) < 5 {
			g.Level = 1
		} else if len(g.SnakeBody) >= 5 && len(g.SnakeBody) < 10 {
			g.Level = 2
			g.MoveTime = 2
		} else if len(g.SnakeBody) >= 10 {
			g.Level = 3
			g.MoveTime = 1
		}
		for i := len(g.SnakeBody) - 1; i > 0; i-- {
			g.SnakeBody[i].x = g.SnakeBody[i-1].x
			g.SnakeBody[i].y = g.SnakeBody[i-1].y
		}
		switch g.Direction {
		case dirUp:
			g.SnakeBody[0].y--
		case dirDown:
			g.SnakeBody[0].y++
		case dirLeft:
			g.SnakeBody[0].x--
		case dirRight:
			g.SnakeBody[0].x++
		}
	}
	g.Timer++
	return nil

}

func (g *Game) Draw(dst *ebiten.Image) {
	dst.DrawImage(image, nil)
	//先渲染蛇的身体
	for _, body := range g.SnakeBody {
		vector.DrawFilledRect(dst, float32(body.x*gridSize), float32(body.y*gridSize), float32(gridSize), float32(gridSize), color.RGBA{
			//R, G, B, A
			R: 0,
			G: 248,
			B: 255,
			A: 1,
		}, false)
	}
	//接下来渲染苹果
	vector.DrawFilledRect(dst, float32(g.Apple.x*gridSize), float32(g.Apple.y*gridSize), float32(gridSize), float32(gridSize), color.RGBA{
		//R, G, B, A
		R: 165,
		G: 42,
		B: 42,
		A: 1,
	}, false)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowDecorated(true)
	ebiten.SetWindowTitle("snack")
	ebiten.SetWindowSize(screenWidth, screenHeight)
	op := &ebiten.RunGameOptions{}
	op.ScreenTransparent = false
	op.SkipTaskbar = true
	if err := ebiten.RunGameWithOptions(NewGame(), op); err != nil {
		panic(err)
	}

}
