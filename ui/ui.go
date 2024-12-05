package ui

import (
	"fmt"
	"log/slog"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/4aykovski/text-entropia-randomizer/lib"
)

func Run() {
	a := app.New()
	w := a.NewWindow("randomizer")
	w.Resize(fyne.NewSize(800, 600))

	w.SetContent(buildContent(w))
	w.ShowAndRun()
}

func buildContent(w fyne.Window) fyne.CanvasObject {
	fileLabel := widget.NewLabel("Choose file")
	fileEntry := widget.NewEntry()
	fileEntry.SetPlaceHolder("file path")

	fileButton := widget.NewButton("File", func() {
		dialog.ShowFileOpen(func(closer fyne.URIReadCloser, err error) {
			if err != nil {
				slog.Error("can't open file", slog.String("error", err.Error()))
			}

			if closer == nil {
				slog.Warn("nothing chosen")

				return
			}

			if closer.URI().Extension() != ".txt" {
				slog.Error("wrong file type", slog.String("extension", closer.URI().Extension()))

				return
			}

			inputPath := closer.URI().Path()
			fileEntry.SetText(inputPath)
			err = closer.Close()
		}, w)
	})

	lengthLabel := widget.NewLabel("Length")
	lengthEntry := widget.NewEntry()
	lengthEntry.SetPlaceHolder("length")

	shiftLabel := widget.NewLabel("Shift")
	shiftEntry := widget.NewEntry()
	shiftEntry.SetPlaceHolder("shift")

	infinite := widget.NewProgressBarInfinite()
	infinite.Hide()

	entropyLabel := widget.NewLabel("Сгенерированная строка на основе энтропии")
	entropyEntry := widget.NewEntry()

	shiftOutputLabel := widget.NewLabel("Сгенерированная строка со сдвигом на основе энтропии")
	shiftOutputEntry := widget.NewEntry()

	runButton := widget.NewButton("Сгенерировать", func() {
		defer infinite.Hide()
		infinite.Show()

		lengthString := lengthEntry.Text
		length, err := strconv.Atoi(lengthString)
		if err != nil {
			dialog.ShowError(fmt.Errorf("неверная длина"), w)
			return
		}

		shiftString := shiftEntry.Text
		shift, err := strconv.Atoi(shiftString)
		if err != nil {
			dialog.ShowError(fmt.Errorf("неверный сдвиг"), w)
			return
		}

		entropyOutput, shiftOutput := lib.GenerateTwoString(fileEntry.Text, length, shift)

		entropyEntry.SetText(entropyOutput)
		shiftOutputEntry.SetText(shiftOutput)
	})

	fileInput := container.New(layout.NewVBoxLayout(), fileLabel, fileEntry, fileButton)

	lengthInput := container.New(layout.NewVBoxLayout(), lengthLabel, lengthEntry)
	shiftInput := container.New(layout.NewVBoxLayout(), shiftLabel, shiftEntry)

	entropyOutput := container.New(layout.NewVBoxLayout(), entropyLabel, entropyEntry)
	shiftOutput := container.New(layout.NewVBoxLayout(), shiftOutputLabel, shiftOutputEntry)

	content := container.New(
		layout.NewVBoxLayout(),
		fileInput,
		lengthInput,
		shiftInput,
		runButton,
		infinite,
		entropyOutput,
		shiftOutput,
	)

	return content
}
