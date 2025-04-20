package scraper

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/cascax/colorthief-go"
	"github.com/rs/zerolog/log"
	_ "golang.org/x/image/webp"
)

func hexColor(c color.Color) string {
	rgba := color.RGBAModel.Convert(c).(color.RGBA)
	return fmt.Sprintf("#%.2x%.2x%.2x", rgba.R, rgba.G, rgba.B)
}
func ParseCover(coverData []byte) (format string, colorStr string, err error) {
	r := bytes.NewReader(coverData)
	img, format, err := image.Decode(r)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Not valid image")
		return
	}
	baseColor, err := colorthief.GetColor(img)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Cant get base color")
		return
	}
	colorStr = hexColor(baseColor)
	log.Logger.Trace().Str("Format", format).Str("color", colorStr).Msg("Cover parsed")
	return
}
