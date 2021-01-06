package main

import (
    "strings"
)

var bottomFilledCharacters = []string{
    " ",
    "\u2581",
    "\u2582",
    "\u2583",
    "\u2584",
    "\u2585",
    "\u2586",
    "\u2587",
    "\u2588",
}

var leftBottomQuartersFilledCharacters = []string{
    " ",
    "\u2596",
    "\u258C",
    "\u2599",
    "\u2588",
}

var leftFilledCharacters = []string{
    " ",
    "\u258F",
    "\u258E",
    "\u258D",
    "\u258C",
    "\u258B",
    "\u258A",
    "\u2589",
    "\u2588",
}

const filledCharacter = "\u2588"

func (color1 Color) combineNext(color2 Color) Color{
    fg1 := color1 & fg
    fg2 := color2 & fg
    otherStyles := (color1 & color2) & ^(fg | bg)
    return otherStyles | (fg2<<3) | fg1
}

func (canvas *Canvas) FullWidthVerticalFill(proportions []float64, colors []Color) {
    total := canvas.height*8*canvas.width
    start := 0.
    position := 0
    nextColor := colors[0]
    for i := range proportions {
        proportion := proportions[i]
        size := float64(total)*proportion
        end := int(start + size)
        if end > total {
            end = total
        }

        currentColor := nextColor
        if i+1<len(colors) {
            nextColor = colors[i+1]
        } else {
            nextColor = Black | (currentColor & ^(fg|bg))
        }
        combinedColor := currentColor.combineNext(nextColor)

        for position < (end & ^0x7) {
            x := position>>3 % canvas.width
            y := position>>3 / canvas.width
            canvas.AddColoredString(x,y, combinedColor, filledCharacter)
            position += 8
        }
        if position < end {
            x := position>>3 % canvas.width
            y := position>>3 / canvas.width
            canvas.AddColoredString(x,y, combinedColor, leftFilledCharacters[end-position])
            position += 8
        }

        start = start + size
    }
}

func (canvas *Canvas) BiColorFullWidthVerticalFill(proportion float64, color1 Color, color2 Color) {
    horizontalBars := proportion*float64(canvas.height)
    exactEnd := int(proportion*float64(canvas.height*canvas.width)*8)
    combinedColor := color1.combineNext(color2)
    for y := 0; y<canvas.height; y++ {
        if y < int(horizontalBars) {
            canvas.AddColoredString(0,y, combinedColor, strings.Repeat(filledCharacter, canvas.width))
        } else if y > int(horizontalBars) {
            canvas.AddColoredString(0,y, combinedColor, strings.Repeat(" ", canvas.width))
        } else {
            for x:= 0; x<canvas.width; x++ {
                base := int((horizontalBars - float64(y))*8)
                position := (y*8+base)*canvas.width+x
                if position <= exactEnd {
                    canvas.AddColoredString(x,y,combinedColor, bottomFilledCharacters[base+1])
                } else {
                    canvas.AddColoredString(x,y,combinedColor, bottomFilledCharacters[base])
                }
            }
        }
    }
}
