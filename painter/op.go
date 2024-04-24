package painter

import (
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/exp/shiny/screen"
)


type Operation interface {

	Do(t screen.Texture) (ready bool)
}

type StateTweaker interface {

	SetState(sol *StatefulOperationList)
}

type StatefulOperationList struct {
	BgOperation      Operation
	BgRectOperation  Operation
	FigureOperations []*OperationFigure
}


func (sol StatefulOperationList) Do(t screen.Texture) (ready bool) {
	if sol.BgOperation != nil {
		sol.BgOperation.Do(t)
	} else {
		t.Fill(t.Bounds(), color.Black, screen.Src)
	}
	if sol.BgRectOperation != nil {
		sol.BgRectOperation.Do(t)
	}
	for _, op := range sol.FigureOperations {
		op.Do(t)
	}
	return false
}

func (sol *StatefulOperationList) Update(o StateTweaker) {
	o.SetState(sol)
}


var UpdateOp = updateOp{}

type updateOp struct{}

func (op updateOp) Do(t screen.Texture) bool { return true }


type OperationFunc func(t screen.Texture)

func (f OperationFunc) Do(t screen.Texture) bool {
	f(t)
	return false
}


type OperationFill struct {
	Color color.Color
}

func (op OperationFill) Do(t screen.Texture) bool {
	t.Fill(t.Bounds(), op.Color, screen.Src)
	return false
}

func (op OperationFill) SetState(sol *StatefulOperationList) {
	sol.BgOperation = op
}

type RelativePoint struct {
	X float64
	Y float64
}

func (p RelativePoint) ToAbs(size image.Point) image.Point {
	return image.Point{
		X: int(p.X * float64(size.X)),
		Y: int(p.Y * float64(size.Y)),
	}
}

type OperationBGRect struct {
	Min RelativePoint
	Max RelativePoint
}

func (op OperationBGRect) Do(t screen.Texture) bool {
	minAbs := op.Min.ToAbs(t.Size())
	maxAbs := op.Max.ToAbs(t.Size())

	rect := image.Rect(minAbs.X, minAbs.Y, maxAbs.X, maxAbs.Y)
	t.Fill(rect, color.Black, draw.Src)
	return false
}

func (op OperationBGRect) SetState(sol *StatefulOperationList) {
	sol.BgRectOperation = op
}

type OperationFigure struct {
	Center RelativePoint
}

func (op OperationFigure) Do(t screen.Texture) bool {
	centerAbs := op.Center.ToAbs(t.Size())
	x := centerAbs.X
	y := centerAbs.Y

	len := 100
	thickness := 10
	color := color.RGBA{R: 255, G: 255, A: 255}

	// Вертикальна лінія
	verticalLine := image.Rect(x-thickness/2, y-len/2, x+thickness/2, y+len/2)
	t.Fill(verticalLine, color, draw.Src)

	// Горизонтальна лінія
	horizontalLine := image.Rect(x-len/2, y-thickness/2, x+len/2, y+thickness/2)
	t.Fill(horizontalLine, color, draw.Src)

	return false
}

func (op OperationFigure) SetState(sol *StatefulOperationList) {
	sol.FigureOperations = append(sol.FigureOperations, &op)
}

type MoveTweaker struct {
	Offset RelativePoint
}

func (t MoveTweaker) SetState(sol *StatefulOperationList) {
	for _, op := range sol.FigureOperations {
		op.Center.X += t.Offset.X
		op.Center.Y += t.Offset.Y
	}
}

type ResetTweaker struct{}

func (op ResetTweaker) SetState(sol *StatefulOperationList) {
	sol.BgOperation = nil
	sol.BgRectOperation = nil
	sol.FigureOperations = []*OperationFigure{}
}