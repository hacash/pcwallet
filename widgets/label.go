package widgets

import (
	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func NewTextWrapWordLabel(text string) *widget.Label {
	lb := widget.NewLabel(text)
	lb.Wrapping = fyne.TextWrapBreak
	return lb
}
