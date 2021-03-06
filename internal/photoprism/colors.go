package photoprism

import (
	"fmt"
	"image/color"
	"log"

	"github.com/disintegration/imaging"
	"github.com/lucasb-eyer/go-colorful"
)

type MaterialColor uint16
type MaterialColors []MaterialColor

const ColorSampleSize = 3

const (
	Black MaterialColor = iota
	Brown
	Grey
	White
	Purple
	Indigo
	Blue
	Cyan
	Teal
	Green
	Lime
	Yellow
	Amber
	Orange
	Red
	Pink
)

var materialColorNames = map[MaterialColor]string{
	Black:  "black",  // 0
	Brown:  "brown",  // 1
	Grey:   "grey",   // 2
	White:  "white",  // 3
	Purple: "purple", // 4
	Indigo: "indigo", // 5
	Blue:   "blue",   // 6
	Cyan:   "cyan",   // 7
	Teal:   "teal",   // 8
	Green:  "green",  // 9
	Lime:   "lime",   // A
	Yellow: "yellow", // B
	Amber:  "amber",  // C
	Orange: "orange", // D
	Red:    "red",    // E
	Pink:   "pink",   // F
}

var materialColorWeight = map[MaterialColor]uint16{
	Black:  2,
	Brown:  1,
	Grey:   2,
	White:  2,
	Purple: 5,
	Indigo: 3,
	Blue:   3,
	Cyan:   4,
	Teal:   4,
	Green:  3,
	Lime:   5,
	Yellow: 5,
	Amber:  5,
	Orange: 5,
	Red:    5,
	Pink:   5,
}

func (c MaterialColor) Name() string {
	return materialColorNames[c]
}

func (c MaterialColor) Hex() string {
	return fmt.Sprintf("%X", c)
}

func (c MaterialColors) Hex() (result string) {
	for _, materialColor := range c {
		result += materialColor.Hex()
	}

	return result
}

var materialColorMap = map[color.RGBA]MaterialColor{
	{0x00, 0x00, 0x00, 0xff}: Black,
	{0x79, 0x55, 0x48, 0xff}: Brown,
	{0x9E, 0x9E, 0x9E, 0xff}: Grey,
	{0xFF, 0xFF, 0xFF, 0xff}: White,
	{0x9c, 0x27, 0xb0, 0xff}: Purple,
	{0x3F, 0x51, 0xB5, 0xff}: Indigo,
	{0x21, 0x96, 0xF3, 0xff}: Blue,
	{0x00, 0xBC, 0xD4, 0xff}: Cyan,
	{0x00, 0x96, 0x88, 0xff}: Teal,
	{0x4C, 0xAF, 0x50, 0xff}: Green,
	{0xCD, 0xDC, 0x39, 0xff}: Lime,
	{0xFF, 0xEB, 0x3B, 0xff}: Yellow,
	{0xFF, 0xC1, 0x07, 0xff}: Amber,
	{0xFF, 0x98, 0x00, 0xff}: Orange,
	{0xf4, 0x43, 0x36, 0xff}: Red,
	{0xe9, 0x1e, 0x63, 0xff}: Pink,
}

func colorfulToMaterialColor(actualColor colorful.Color) (result MaterialColor) {
	var distance = 1.0

	for colorRGBA, materialColor := range materialColorMap {
		colorColorful, _ := colorful.MakeColor(colorRGBA)
		currentDistance := colorColorful.DistanceLab(actualColor)

		if distance >= currentDistance {
			distance = currentDistance
			result = materialColor
		}
	}

	return result
}

// Colors returns color information for a media file.
func (m *MediaFile) Colors() (colors MaterialColors, mainColor MaterialColor, err error) {
	jpeg, err := m.GetJpeg()

	if err != nil {
		log.Printf("can't find jpeg: %s", err.Error())

		return colors, mainColor, err
	}

	img, err := imaging.Open(jpeg.GetFilename(), imaging.AutoOrientation(true))

	if err != nil {
		log.Printf("can't open jpeg: %s", err.Error())

		return colors, mainColor, err
	}

	img = imaging.Resize(img, ColorSampleSize, ColorSampleSize, imaging.Box)

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	colorCount := make(map[MaterialColor]uint16)
	var mainColorCount uint16

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			rgbColor, _ := colorful.MakeColor(color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)})
			materialColor := colorfulToMaterialColor(rgbColor)
			colors = append(colors, materialColor)

			if _, ok := colorCount[materialColor]; ok == true {
				colorCount[materialColor] += materialColorWeight[materialColor]
			} else {
				colorCount[materialColor] = materialColorWeight[materialColor]
			}

			if colorCount[materialColor] > mainColorCount {
				mainColorCount = colorCount[materialColor]
				mainColor = materialColor
			}

		}
	}

	return colors, mainColor, nil
}
