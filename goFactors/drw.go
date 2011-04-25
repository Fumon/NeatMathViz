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
var wholes = 1

var nonDivOff bool = false

var pic *image.Gray16
var picStride int

var window draw.Window

func randomizeWindow() {
	s := window.Screen()

	colorScaler := 0xFFFF/div/2

	//Overlay
	totalSide := div * wholes
	st :=  600 / totalSide
	boundary := st * div * wholes

	for i := range pic.Pix {
		var setting uint16
		x, y := i % picStride, i / picStride
		if x > boundary || y > boundary {
			pic.Pix[i].Y = 0
		} else {
			x, y = x / st, y / st
			if z := x * y % div; z == 0 {
				setting = 0xFFFF
			} else {
				if nonDivOff {
					setting = 0x0000
				} else {
					setting = uint16(z * colorScaler)
				}
			}
		}
		pic.Pix[i].Y = setting
		/*if x := ((i%pic.Stride * (i / pic.Stride)) % div); x == 0 {
			pic.Pix[i].Y = 0xFFFF
		} else {
			if nonDivOff {
				pic.Pix[i].Y = 0
			} else{
				pic.Pix[i].Y = uint16(x * colorScaler)
			}
		}*/
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
	/*defer func() {
		if x := recover(); x!= nil {
			fmt.Printf("Run time panic: %v", x)
		}
	}()
	*/
	rnd = rand.New(rand.NewSource(src))
	win, err := x11.NewWindowDisplay(":1")

	window = win
	b := (window.Screen()).Bounds()

	pic = image.NewGray16(b.Dx(), b.Dy())
	picStride = pic.Stride

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
					randomizeWindow()
				case 'f':
					if(div < 1024) {
						div++
					}
					fmt.Printf("Increased Div to %v\n", div)
					randomizeWindow()
				case 'b':
					if(div > 1) {
						div--
					}
					fmt.Printf("Decreased div to %v\n", div)
					randomizeWindow()
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
					randomizeWindow()
				case 'w':
					if wholes < 30 {
						wholes++
					}
					fmt.Printf("Wholes++: %v", wholes)
				case 'W':
					if wholes > 1{
						wholes--
					}
					fmt.Printf("Wholes--: %v", wholes)
				}
			}
		}
	}

	fmt.Println("Exiting")
	window.Close()
}
