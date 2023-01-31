package main

import (
	"bufio"
	"bytes"
	"fmt"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"sort"
)

func addLabel(img *image.RGBA, x, y int, label string) {
	col := color.RGBA{R: 0, G: 0, B: 0, A: 255}
	point := fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}

type handler struct{}

func (h *handler) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	lines := make([]string, 0)
	lines = append(lines, fmt.Sprintf("%s %s", r.Method, r.URL.String()))
	lines = append(lines, "")

	// header
	headerLines := make([]string, 0)
	headerLines = append(headerLines, fmt.Sprintf("%s: %s", "Host", r.Host))
	for k, v := range r.Header {
		for _, item := range v {
			headerLines = append(headerLines, fmt.Sprintf("%s: %s", k, item))
		}
	}
	sort.Strings(headerLines)
	for _, item := range headerLines {
		lines = append(lines, item)
	}
	log.Print(lines)
	lines = append(lines, "")

	width := 0
	for _, str := range lines {
		if len(str) > width {
			width = len(str)
		}
	}
	if width < 100 {
		width = 100
	}

	// body
	lineWidth := width * 4 / 5
	body := bufio.NewReader(r.Body)
	for {
		line, _, err := body.ReadLine()
		if err != nil {
			break
		}

		if len(line) > lineWidth {
			lineReader := bytes.NewReader(append(line, 10))
			for {
				buffer := make([]byte, lineWidth)
				read, err := lineReader.Read(buffer)
				if err != nil {
					break
				}
				if read > 0 {
					lines = append(lines, string(buffer[:read]))
				}
			}
		} else {
			lines = append(lines, string(append(line, 10)))
		}
	}

	// print
	width = width*7 + 100
	height := len(lines)*15 + 100

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, image.White)
		}
	}
	for line, str := range lines {
		addLabel(img, 50, 50+line*15, str)
	}

	err := png.Encode(writer, img)
	if err != nil {
		return
	}
}

func main() {
	handler := &handler{}
	http.ListenAndServe(":8090", handler)
}
