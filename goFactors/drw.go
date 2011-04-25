package main

import (
	"exp/draw"
	"exp/draw/x11"
	"fmt"
	"image"
	"rand"
)

const (
	src = 3478912056
)

var rnd *rand.Rand

var div int = 2

var nonDivOff bool = false

var pic *image.Gray16

func randomizeWindow(window draw.Window) {
	s := window.Screen()

	colorScaler := 0xFFFF/div

	for i := range pic.Pix {
		if x := ((i%pic.Stride * (i / pic.Stride)) % div); x == 0 {
			pic.Pix[i].Y = 0xFFFF
		} else {
			if nonDivOff {
				pic.Pix[i].Y = 0
			} else{
				pic.Pix[i].Y = uint16(x * colorScaler)
			}
		}
		//pic.Set(i % pic.Stride, i / pic.Stride, image.Gray16Color{0xFFFF})
	}
	draw.Draw(s, s.Bounds(), pic, image.ZP)

/*	for i := 0; i < bounds.Dy(); i++ {
		for j := 0; j < bounds.Dx(); j++ {
			var c image.Gray16Color
			if x:= (j * i)%div; x == 0 {
				c = image.Gray16Color{0xFFFF}
			} else {
				if nonDivOff {
					c = image.Gray16Color{0}
				} else {
					c = image.Gray16Color{uint16(x*(0xFFFF/div))}
				}
			}
//image.GrayColor{uint8(rnd.Intn(256))}
			s.Set(bounds.Min.X + j,
				bounds.Min.Y + i,
				c)
		}
	} */
	window.FlushImage()
}

func main() {
	defer func() {
		if x := recover(); x!= nil {
			fmt.Printf("Run time panic: %v", x)
		}
	}()

	rnd = rand.New(rand.NewSource(src))
	window, err := x11.NewWindowDisplay(":1")

	b := (window.Screen()).Bounds()

	pic = image.NewGray16(b.Dx(), b.Dy())

	if err != nil {
		fmt.Println(err)
		return;
	}

	evChan := window.EventChan()
	fmt.Println("Initialized loop")

	Mainloop: for {
		e := <-evChan
		me, ok := e.(draw.MouseEvent)

		if ok{
			switch me.Buttons {
			case 1:
				fmt.Printf("Pos: (%v)\n", me.Loc)
			}
		} else {
			ke, ok := e.(draw.KeyEvent)
			if ok {
				switch ke.Key {
				case 'q':
					break Mainloop
				case 'r':
					randomizeWindow(window)
				case 'i':
					//Up arrow
					if(div < 1024){
						div++
					}
					fmt.Printf("Increased Div to %v\n", div)
				case 'k':
					//Down arrow
					if(div > 1) {
						div--
					}
					fmt.Printf("Decreased div to %v\n", div)
				case 'o':
					nonDivOff = !nonDivOff
					randomizeWindow(window)
				}
			}
		}
	}

	fmt.Println("Exiting")
	window.Close()
}
