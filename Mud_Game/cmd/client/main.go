package main

import (
	"bytes"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var myFont *text.GoTextFace                          //указатель на объект шрифта из пакета text/v2
var cwet = color.RGBA{R: 64, G: 224, B: 208, A: 255} //цвет бирюзовый для шрифта

// Она представляет игру.
// Ebiten требует, чтобы у нас был объект с методами Update, Draw, Layout.
type Game struct{}

// Специальная функция Go, которая автоматически вызывается до main()
func init() {
	//загружаем шрифт
	fontData, err := os.ReadFile("assets/fonts/Roboto-Regular.ttf")
	if err != nil {
		log.Fatal("Не удалось загрузить шрифт:", err)
	}

	//Создаем источник шрифта
	src, err := text.NewGoTextFaceSource(bytes.NewReader(fontData))
	if err != nil {
		log.Fatal("Не удалось распарсить шрифт:", err)
	}

	//создаем шрифт размером 16
	myFont = &text.GoTextFace{
		Source: src, //сам шрифт (данные)
		Size:   16,
	}
}

// Здесь обрабатывается логика(60 кад/сек)
func (g *Game) Update() error {
	return nil
}

// Здесь рисуется всё(60 кад/сек)
func (g *Game) Draw(screen *ebiten.Image) {
	DrawCenteredText(screen, "Привет! Это будущая MUDка!", myFont, 640, nil)
	DrawText(screen, "Городская площадь", myFont, 40.0, 100.0, color.White)
}

// Это нужно для масштабирования. Если окно растянут — Ebiten сам отмасштабирует.
func (g *Game) Layout(outsideWidth, outsideHeigh int) (int, int) {
	return 640, 480

}

func main() {
	ebiten.SetWindowSize(640, 480)    //создание окна +размер окна
	ebiten.SetWindowTitle("MUD Game") //заголовок окна

	//запускает игровой цикл:
	// 60 раз в секунду вызывает Update
	//60 раз в секунду вызывает Draw
	//Game НЕ создает проект — проект создает RunGame. Game — это объект с логикой, который RunGame использует.
	if err := ebiten.RunGame(&Game{}); err != nil { //создаем экземпляр нашей структуры и передаем указатель.
		log.Fatal(err)
	}
}

// ввод текста по центру чуть выше
func DrawCenteredText(screen *ebiten.Image, str string, font *text.GoTextFace, screenW int, clr color.Color) {
	if clr == nil {
		clr = cwet
	}
	textWidth, _ := text.Measure(str, font, 0)
	x := (float64(screenW) - textWidth) / 2
	y := 50.0
	//op - Как рисовать	Позиция + цвет
	op := &text.DrawOptions{}         ////опции для рисования
	op.GeoM.Translate(x, y)           //позиция
	op.ColorScale.ScaleWithColor(clr) //цвет
	text.Draw(screen, str, font, op)  //сама надпись
}

// вывод текста
func DrawText(screen *ebiten.Image, str string, font *text.GoTextFace, x, y float64, clr color.Color) {
	if clr == nil {
		clr = cwet
	}
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(clr)
	text.Draw(screen, str, font, op)

}
