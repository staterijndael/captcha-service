package service

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"math"
	"math/rand"
)

type CaptchaService struct {
}

func NewCaptchaService() *CaptchaService {
	return &CaptchaService{}
}

func (c *CaptchaService) CreateCaptcha(word string, width int, height int) (*image.RGBA, error) {
	background := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(background, background.Bounds(), &image.Uniform{C: color.White},
		image.Point{}, draw.Src)
	fontFace, err := freetype.ParseFont(goregular.TTF)
	if err != nil {
		return nil, err
	}

	// second layer of random lines creating (upside the background)
	firstLayer := rand.Intn(110-50) + 50
	drawRandLines(background, firstLayer)

	text := []rune(word)
	textBytes := []byte(string(text))

	fontSize := 178.0 - float64(3*len(text))

	lettersFontDrawer := &font.Drawer{
		Dst: background,
		Face: truetype.NewFace(fontFace, &truetype.Options{
			Size:    fontSize,
			Hinting: font.HintingFull,
		}),
	}

	textBounds, _ := lettersFontDrawer.BoundBytes(textBytes)
	textHeight := textBounds.Max.Y - textBounds.Min.Y

	offsetsX := make([]fixed.Int26_6, len(text))
	var offsetsSum fixed.Int26_6

	offsetsY := make([]fixed.Int26_6, len(text))

	measures := make([]fixed.Int26_6, len(text))
	var measuresSum fixed.Int26_6
	for i := 0; i < len(text); i++ {
		if i != 0 {
			offsetsX[i] = fixed.Int26_6(rand.Intn(2500-2100) + 2100)
			offsetsSum += offsetsX[i]
		}

		offsetsY[i] = fixed.Int26_6(rand.Intn(800+800) - 800)

		measures[i] = lettersFontDrawer.MeasureBytes([]byte(string(text[i])))
		measuresSum += measures[i]
	}

	currentX := (fixed.I(background.Rect.Max.X) - measuresSum + offsetsSum) / 2
	currentY := fixed.I((background.Rect.Max.Y)-textHeight.Ceil())/2 + fixed.I(textHeight.Ceil())

	for i, symb := range text {
		xOffset := offsetsX[i]
		yOffset := offsetsY[i]

		symbBytes := []byte(string(symb))
		currentSymbMeasure := measures[i]
		err = drawRandLetter(lettersFontDrawer, currentX-xOffset, currentY+yOffset, symbBytes)
		if err != nil {
			return nil, err
		}

		currentX += currentSymbMeasure - xOffset
	}

	imageWaveCurvature(background)

	// second layer of random lines creating (upside the background)
	secondLayerLinesCount := rand.Intn(110-50) + 50
	drawRandLines(background, secondLayerLinesCount)

	return background, nil
}

func drawRandLines(canvas *image.RGBA, count int) {
	for i := 0; i < count; i++ {
		randX1 := rand.Intn(canvas.Rect.Max.X-canvas.Rect.Min.X) + canvas.Rect.Min.X
		randX2 := rand.Intn(canvas.Rect.Max.X-canvas.Rect.Min.X) + canvas.Rect.Min.X

		randY1 := rand.Intn(canvas.Rect.Max.Y-canvas.Rect.Min.Y) + canvas.Rect.Min.Y
		randY2 := rand.Intn(canvas.Rect.Max.Y-canvas.Rect.Min.Y) + canvas.Rect.Min.Y

		randColor := getRandColor()

		drawLine(canvas, randColor, randX1, randX2, randY1, randY2)
	}
}

// bresenham algorithm
func drawLine(canvas *image.RGBA, color color.Color, x1 int, x2 int, y1 int, y2 int) {
	dx, dy := x2-x1, y2-y1
	a := float64(dy) / float64(dx)
	b := int(float64(y1) - a*float64(x1))

	canvas.Set(x1, y1, color)
	for x := x1 + 1; x <= x2; x++ {
		y := int(a*float64(x)) + b
		canvas.Set(x, y, color)
	}
}

func drawRandLetter(fontDrawer *font.Drawer, x fixed.Int26_6, y fixed.Int26_6, letterBytes []byte) error {
	fontDrawer.Src = &image.Uniform{
		C: getRandColor(),
	}

	fontDrawer.Dot = fixed.Point26_6{
		X: x,
		Y: y,
	}

	fontDrawer.DrawBytes(letterBytes)

	return nil
}

func imageWaveCurvature(image *image.RGBA) {
	// частота (зернисность изображения)
	ratio1 := float64(300000) / float64(30000)
	ratio2 := float64(300000) / float64(30000)
	ratio3 := float64(300000) / float64(30000)
	ratio4 := float64(300000) / float64(30000)
	// фазы (сдвиг пикселей)
	phase1 := float64(4000000) / float64(1600000)
	phase2 := float64(4000000) / float64(1600000)
	phase3 := float64(4000000) / float64(1600000)
	phase4 := float64(4000000) / float64(1600000)
	// амплитуды (насколько далеко размываются пиксели друга от друга)
	amplitude1 := float64(500) / float64(90)
	amplitude2 := float64(500) / float64(90)

	var x, y float64

	for x = 0; x < float64(image.Rect.Max.X); x++ {
		for y = 0; y < float64(image.Rect.Max.Y); y++ {
			// координаты пикселя-первообраза.
			sx := x + (math.Sin(x*ratio1+phase1)+math.Sin(y*ratio3+phase2))*amplitude1
			sy := y + (math.Sin(x*ratio2+phase3)+math.Sin(y*ratio4+phase4))*amplitude2

			var (
				col      color.Color
				color_x  color.Color
				color_y  color.Color
				color_xy color.Color
			)

			// первообраз за пределами изображения
			if sx < 0 || sy < 0 || sx >= float64(image.Rect.Max.X)-1 || sy >= float64(image.Rect.Max.Y)-1 {
				col = color.White
				color_x = color.White
				color_y = color.White
				color_xy = color.White
			} else { // цвета основного пикселя и его 3-х соседей для лучшего антиалиасинга
				col = image.At(int(sx), int(sy))
				color_x = image.At(int(sx)+1, int(sy))
				color_y = image.At(int(sx), int(sy)+1)
				color_xy = image.At(int(sx)+1, int(sy)+1)
			}

			var newColor color.Color

			var (
				frsx  float64
				frsy  float64
				frsx1 float64
				frsy1 float64
			)

			// сглаживаем только точки, цвета соседей которых отличается
			if col == color_x && col == color_y && col == color_xy {
				newColor = col
			} else {
				frsx = sx - math.Floor(sx) //отклонение координат первообраза от целого
				frsy = sy - math.Floor(sy)
				frsx1 = 1 - frsx
				frsy1 = 1 - frsy

				// вычисление цвета нового пикселя как пропорции от цвета основного пикселя и его соседей
				colR, _, _, _ := col.RGBA()
				colXR, _, _, _ := color_x.RGBA()
				colYR, _, _, _ := color_y.RGBA()
				colXYR, _, _, _ := color_xy.RGBA()

				pixelColR := int(colR >> 8)
				pixelColXR := int(colXR >> 8)
				pixelColYR := int(colYR >> 8)
				pixelColXYR := int(colXYR >> 8)

				newColor = color.RGBA{
					R: uint8(float64(pixelColR)*frsx1*frsy1 + float64(pixelColXR)*frsx*frsy1 + float64(pixelColYR)*frsx1*frsy + float64(pixelColXYR)*frsx*frsy),
					G: uint8(float64(pixelColR)*frsx1*frsy1 + float64(pixelColXR)*frsx*frsy1 + float64(pixelColYR)*frsx1*frsy + float64(pixelColXYR)*frsx*frsy),
					B: uint8(float64(pixelColR)*frsx1*frsy1 + float64(pixelColXR)*frsx*frsy1 + float64(pixelColYR)*frsx1*frsy + float64(pixelColXYR)*frsx*frsy),
					A: 255,
				}
			}

			image.Set(int(x), int(y), newColor)
		}
	}
}

func getRandColor() color.RGBA {
	return color.RGBA{
		R: uint8(rand.Intn(math.MaxUint8)),
		G: uint8(rand.Intn(math.MaxUint8)),
		B: uint8(rand.Intn(math.MaxUint8)),
		A: 255,
	}
}
