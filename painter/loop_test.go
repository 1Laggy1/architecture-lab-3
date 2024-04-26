package painter

import (
	"image"
	"image/color"
	"image/draw"
	"testing"

	"golang.org/x/exp/shiny/screen"
)



func TestLoop_Post(t *testing.T) {
	var (
		l  Loop
		tr testReceiver
	)
	l.Receiver = &tr

	l.Start(mockScreen{})
	l.Post(logOp(t, "do white fill", OperationFunc(func(t screen.Texture) {
		OperationFill{Color: color.White}.Do(t)
	})))
	l.Post(logOp(t, "do green fill", OperationFunc(func(t screen.Texture) {
		OperationFill{Color: color.RGBA{0, 255, 0, 255}}.Do(t)
	})))
	l.Post(UpdateOp)

	for i := 0; i < 3; i++ {
		go l.Post(logOp(t, "do green fill", OperationFunc(func(t screen.Texture) {
			OperationFill{Color: color.RGBA{0, 255, 0, 255}}.Do(t)
		})))
	}

	l.StopAndWait()

	if tr.lastTexture == nil {
		t.Fatal("Texture was not updated")
	}
	mt, ok := tr.lastTexture.(*mockTexture)
	if !ok {
		t.Fatal("Unexpected texture", tr.lastTexture)
	}
	if mt.Colors[0] != color.White {
		t.Error("First color is not white:", mt.Colors)
	}
	if len(mt.Colors) != 2 {
		t.Error("Unexpected size of colors:", mt.Colors)
	}

}

func logOp(t *testing.T, msg string, op OperationFunc) OperationFunc {
	return func(tx screen.Texture) {
		t.Log(msg)
		op(tx)
	}
}

type testReceiver struct {
	lastTexture screen.Texture
}

func (tr *testReceiver) Update(t screen.Texture) {
	tr.lastTexture = t
}

type mockScreen struct{}

func (m mockScreen) NewBuffer(size image.Point) (screen.Buffer, error) {
	panic("implement me")
}

func (m mockScreen) NewTexture(size image.Point) (screen.Texture, error) {
	return new(mockTexture), nil
}

func (m mockScreen) NewWindow(opts *screen.NewWindowOptions) (screen.Window, error) {
	panic("implement me")
}

type mockTexture struct {
	Colors []color.Color
}

func (m *mockTexture) Release() {}

func (m *mockTexture) Size() image.Point { return size }

func (m *mockTexture) Bounds() image.Rectangle {
	return image.Rectangle{Max: m.Size()}
}

func (m *mockTexture) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle) {}
func (m *mockTexture) Fill(dr image.Rectangle, src color.Color, op draw.Op) {
	m.Colors = append(m.Colors, src)
}
