package main

import (
    "fmt"
    "strconv"
    "golang.org/x/crypto/ssh/terminal"
)

type (
    Canvas struct {
        width int
        height int
        characters []string
    }

    Color uint32
)

const (
    fg Color = 7
    Black Color = 0
    Red Color = 1
    Green Color = 2
    Yellow Color = 3
    Blue Color = 4
    Magenta Color = 5
    Cyan Color = 6
    White Color = 7

    bg Color = 56
    BgBlack Color = 0
    BgRed Color = 8
    BgGreen Color = 16
    BgYellow Color = 24
    BgBlue Color = 32
    BgMagenta Color = 40
    BgCyan Color = 48
    BgWhite Color = 56

    Bold Color = 64
    Dim Color = 128
    Italic Color = 256
    Underline Color = 512
    Blink Color = 1024
)

func NewScreenCanvas() Canvas {
    width, height := getTerminalSize()
    return NewCanvas(width, height)
}

func NewCanvas(width int, height int) Canvas {
    characters := make([]string, height*width)
    for i := range characters {
        characters[i] = " "
    }
    return Canvas{width, height, characters}
}

func getTerminalSize() (int, int) {
    width, height, _ := terminal.GetSize(0)
    return width, height
}

func (canvas *Canvas) PrintToTerm() {
    fmt.Print("\033[?25l\033[0;0H")
    for l:=0; l<canvas.height; l++ {
        line := ""
        for c:=0; c<canvas.width; c++ {
            line += canvas.characters[l*canvas.width+c]
        }
        if l < canvas.height - 1 {
            fmt.Println(line)
        } else {
            fmt.Print(line)
        }
    }
}

func (canvas *Canvas) Blit(x int, y int, otherCanvas Canvas) {
    // TODO : check if it fits
    for l:=0; l<otherCanvas.height; l++ {
        targetY := y+otherCanvas.height-l-1
        for c:=0; c<otherCanvas.width; c++ {
            canvas.AddString(x+c, targetY, otherCanvas.characters[l*otherCanvas.width+c])
        }
    }
}

func (canvas *Canvas) AddString(x int, y int, str string) {
    dx := 0
    for _, char := range str {
        canvas.addChar(x+dx, y, string(char))
        dx += 1
    }
}

func (canvas *Canvas) addChar(x int, y int, char string) {
    canvas.characters[(canvas.height-y-1)*canvas.width + x] = char
}

func (canvas *Canvas) AddColoredString(x int, y int, color Color, str string) {
    dx := 0
    for _, char := range str {
        canvas.addColoredChar(x+dx, y, color, string(char))
        dx += 1
    }
}

func (canvas *Canvas) addColoredChar(x int, y int, color Color, char string) {
    canvas.characters[(canvas.height-y-1)*canvas.width + x] = getColorAnsiCode(color) + char + "\033[0m"
}

func getColorAnsiCode(color Color) string {
    fgColor := color & fg
    bgColor := color & bg >> 3
    bold := (color & Bold) != 0
    dim := (color & Dim) != 0
    italic := (color & Italic) != 0
    underline := (color & Underline) != 0
    blink := (color & Blink) != 0

    ansiCode := "\033["
    ansiCode += "3"+strconv.Itoa(int(fgColor))
    ansiCode += ";4"+strconv.Itoa(int(bgColor))
    if bold {
        ansiCode += ";1"
    }
    if dim {
        ansiCode += ";2"
    }
    if italic {
        ansiCode += ";3"
    }
    if underline {
        ansiCode += ";4"
    }
    if blink {
        ansiCode+= ";5"
    }
    ansiCode += "m"
    return ansiCode
}
