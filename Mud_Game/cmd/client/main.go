package main

import (
	"bytes"
	"image"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// неизменяемые значения
const (
	frameWidth  = 1280
	frameHeight = 720
	frameCols   = 7  //колонок в спрайте
	frameRows   = 13 // строк
	totalFrames = 91 //всего слайдов
	screenW     = 1280
	screenH     = 720

	btnX = 540
	btnY = 450
	btnW = 200
	btnH = 60
)

var (
	myFont      *text.GoTextFace                            //указатель на объект шрифта из пакета text/v2
	cwet        = color.RGBA{R: 64, G: 224, B: 208, A: 255} //цвет бирюзовый для шрифта
	spriteSheet *ebiten.Image                               //спрайт лист картинки
)

// Game struct представляет игру.
// Ebiten требует, чтобы у нас был объект с методами Update, Draw, Layout.
type Game struct {
	frameIndex     int    //Хранить текущий кадр анимации.
	tickCount      int    //  счетчик тиков
	playerName     string //имя игрока
	cursorVisible  bool   //виден ли курсор сейчас
	cursorTimer    int    //таймер мигания
	backspaceTimer int    //для того чтобы удалял 1 символ или при зажатии больше
	gameState      string //строка состояния "login" или "game"
}

// Специальная функция Go, которая автоматически вызывается до main()
func init() {
	//загружаем шрифт
	fontData, err := os.ReadFile("assets/fonts/CyrillicPixel7-LPeg.ttf")
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
		Size:   32,
	}

	//загрузка спрайт-листа
	spriteSheet, _, err = ebitenutil.NewImageFromFile("assets/images/start.png")
	if err != nil {
		log.Fatal("Не удалось загрузить спрайт:", err)
	}

}

func main() {
	ebiten.SetWindowSize(1280, 720)   //создание окна +размер окна
	ebiten.SetWindowTitle("MUD Game") //заголовок окна

	//запускает игровой цикл:
	// 60 раз в секунду вызывает Update
	//60 раз в секунду вызывает Draw
	//Game НЕ создает проект — проект создает RunGame. Game — это объект с логикой, который RunGame использует.
	if err := ebiten.RunGame(&Game{
		gameState: "login",
	}); err != nil { //создаем экземпляр нашей структуры и передаем указатель.
		log.Fatal(err)
	}
}

// Здесь обрабатывается логика(60 кад/сек) фона
func (g *Game) Update() error {
	///////////////////////////////Каждый тик (60 раз/сек) — следующий кадр. Цикл. меняется слайд
	g.tickCount++
	if g.tickCount >= 2 {
		g.tickCount = 0
		g.frameIndex++
		if g.frameIndex >= totalFrames {
			g.frameIndex = 0
		}
	}

	//////////////////////////////мигание курсора////////////////////////////////////
	g.cursorTimer++
	if g.cursorTimer >= 30 {
		g.cursorTimer = 0
		g.cursorVisible = !g.cursorVisible //переключаем
	}

	///////////////////////////////ввод текста///////////////////////////////////////////
	//возвращает список символов, которые игрок ввел за один тик
	//ch — это один символ (руна, rune).
	for _, ch := range ebiten.AppendInputChars(nil) {
		if len([]rune(g.playerName)) >= 10 {
			break //выходим больше не добавляем
		}
		//В компьютере каждая буква — это число (код в ASCII):
		//'a' = 97
		//'z' = 122
		//'A' = 65
		//'Z' = 90
		//'0' = 48
		//'9' = 57
		if (ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '_' {
			g.playerName += string(ch)
		}

	}

	///////////////////////////////////обработка Backspace///////////////////////////
	//Проверяем, нажата ли клавиша Backspace прямо сейчас.
	if ebiten.IsKeyPressed(ebiten.KeyBackspace) {
		g.backspaceTimer++
		if g.backspaceTimer == 1 || g.backspaceTimer > 30 {
			if len(g.playerName) > 0 {
				//убираем последний символ(руну)
				runa := []rune(g.playerName)
				g.playerName = string(runa[:len(runa)-1])
			}
		}
	} else {
		//Чтобы следующее нажатие снова удаляло сразу (счетчик обнулился).
		g.backspaceTimer = 0 //сброс при отпускании
	}

	///////////////////////////////////обработка клика по кнопке(только на экране логина)//////////
	if g.gameState == "login" && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		//отрезок для клика(прямоугольник)
		if mx >= btnX && mx <= btnX+btnW && my >= btnY && my <= btnY+btnH {
			if g.playerName != "" {
				g.gameState = "game"
			}
		}
	}

	return nil

}

// Здесь рисуется всё(60 кад/сек)
func (g *Game) Draw(screen *ebiten.Image) {
	//ФОН - пока что всегда
	g.drawFrame(screen, g.frameIndex)

	if g.gameState == "login" {
		g.drawLoginScreen(screen)
	} else {
		g.drawGameScreen(screen)
	}

}

// Это нужно для масштабирования. Если окно растянут — Ebiten сам отмасштабирует.
func (g *Game) Layout(outsideWidth, outsideHeigh int) (int, int) {
	return 1280, 720

}

// Вырезать кадр из спрайт-листа и нарисовать.
func (g *Game) drawFrame(screen *ebiten.Image, frameIndex int) {
	//8 / 7 = 1 (целая часть) остаток не берем
	//8 - (1 × 7) = 1 — это остаток = колонка

	col := frameIndex % frameCols //Вычисляем колонку кадра в спрайт-листе.

	//8 / 7 = 1
	row := frameIndex / frameCols //Вычисляем строку кадра в спрайт-листе.

	sx := col * frameWidth  //Вычисляем x-координату в спрайт-листе.
	sy := row * frameHeight //Вычисляем y-координату в спрайт-листе.

	//Вырезаем прямоугольник из спрайт-листа.

	//image.Rect(x1, y1, x2, y2) — прямоугольник

	//sx, sy — верхний левый угол

	//sx+frameWidth, sy+frameHeight — нижний правый угол

	//SubImage(...) — "дай мне часть картинки"

	//.(*ebiten.Image) — превращаем в Ebiten-картинку
	frame := spriteSheet.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image)

	//Рисуем вырезанный кадр на экран.
	screen.DrawImage(frame, nil)
}

// ввод текста по центру чуть выше(screenW всегда 1280)
func DrawCenteredText(screen *ebiten.Image, str string, font *text.GoTextFace, screenW int, y float64, clr color.Color) {
	if clr == nil {
		clr = cwet
	}
	textWidth, _ := text.Measure(str, font, 0)
	x := (float64(screenW) - textWidth) / 2
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

// экран авторизации
func (g *Game) drawLoginScreen(screen *ebiten.Image) {
	//ЗАГОЛОВОК
	DrawCenteredText(screen, "Привет! Это будущая MUDка!", myFont, 1280, 100, cwet)
	//ПОДСКАЗКА
	DrawCenteredText(screen, "Введите имя персонажа", myFont, 1280, 320, color.White)
	//ПОЛЕ ВВОДА
	vector.FillRect(screen, 490, 360, 300, 50, color.RGBA{40, 40, 60, 50}, true)
	//КНОПКА ВОЙТИ
	vector.FillRect(screen, btnX, btnY, btnW, btnH, color.RGBA{60, 120, 120, 200}, false)
	DrawCenteredText(screen, "Войти", myFont, 1280, 465, color.White)

	//текст внутри поля(ник игрока, пока пустой)
	if g.playerName == "" {
		//Подсказка - плейсхолдер
		if g.cursorVisible {
			DrawText(screen, "_", myFont, 510, 375, color.RGBA{150, 150, 150, 255})
		}
	} else {
		//Введеный ник + курсор
		text := g.playerName
		if g.cursorVisible {
			text += "_"
		}
		DrawText(screen, g.playerName, myFont, 510, 375, color.White)
	}
}

func (g *Game) drawGameScreen(screen *ebiten.Image) {
	DrawCenteredText(screen, "Ты в игре!", myFont, 1280, 100, cwet)
	DrawCenteredText(screen, "Игрок: "+g.playerName, myFont, 1280, 200, color.White)
}
