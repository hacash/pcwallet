package widgets

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

func NewTextWrapWordLabel(text string) *widget.Label {
	lb := widget.NewLabel(text)
	lb.Wrapping = fyne.TextWrapBreak
	return lb
}
