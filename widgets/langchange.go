package widgets

import (
	fyne "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type LangItem struct {
	words map[string]string
}

type LangChangeManager struct {
	a               fyne.App
	currentLangName string // en zh
	objs            []interface{}
	langs           []*LangItem
}

func NewLangChangeManager(a fyne.App) *LangChangeManager {
	return &LangChangeManager{
		a:               a,
		currentLangName: "en",
		objs:            make([]interface{}, 0),
		langs:           make([]*LangItem, 0),
	}
}

func (l *LangChangeManager) NewWindowAndShow(title map[string]string, windowSize *fyne.Size, content fyne.CanvasObject) fyne.Window {
	w := l.a.NewWindow(title[l.currentLangName])
	w.Resize(*windowSize)

	box := container.NewVBox()
	box.Add(content)

	// 页面翻动
	scroll := container.NewVScroll(box)

	w.SetContent(scroll)
	w.Show()

	// add ary
	l.objs = append(l.objs, w)
	l.langs = append(l.langs, &LangItem{words: title})

	// return ok
	return w
}

func (l *LangChangeManager) NewTextWrapWordLabel(texts map[string]string) *widget.Label {
	lb := widget.NewLabel(texts[l.currentLangName])
	lb.Wrapping = fyne.TextWrapBreak
	// add ary
	l.objs = append(l.objs, lb)
	l.langs = append(l.langs, &LangItem{words: texts})
	// return ok
	return lb
}

func (l *LangChangeManager) NewEntrySetPlaceHolder(texts map[string]string) *widget.Entry {
	lb := widget.NewEntry()
	lb.SetPlaceHolder("- " + texts[l.currentLangName] + " -")

	// add ary
	l.objs = append(l.objs, lb)
	l.langs = append(l.langs, &LangItem{words: texts})
	// return ok
	return lb
}

func (l *LangChangeManager) NewCardSetTitle(texts map[string]string, content fyne.CanvasObject) *widget.Card {
	lb := widget.NewCard(texts[l.currentLangName], "", content)

	// add ary
	l.objs = append(l.objs, lb)
	l.langs = append(l.langs, &LangItem{words: texts})
	// return ok
	return lb
}

func (l *LangChangeManager) NewButton(texts map[string]string, tapped func()) *widget.Button {
	lb := widget.NewButton(texts[l.currentLangName], tapped)

	// add ary
	l.objs = append(l.objs, lb)
	l.langs = append(l.langs, &LangItem{words: texts})
	// return ok
	return lb
}

func (l *LangChangeManager) NewCheck(texts map[string]string, tapped func(bool)) *widget.Check {
	lb := widget.NewCheck(texts[l.currentLangName], tapped)

	// add ary
	l.objs = append(l.objs, lb)
	l.langs = append(l.langs, &LangItem{words: texts})
	// return ok
	return lb
}

func (l *LangChangeManager) SetText(objs interface{}, texts map[string]string) {
	l.setTextEx(objs, texts, false)
}

func (l *LangChangeManager) setTextEx(objs interface{}, texts map[string]string, holdermark bool) {
	name := l.currentLangName
	if obj, ok := objs.(*widget.Label); ok {
		obj.SetText(texts[name])
	} else if obj, ok := objs.(*widget.Entry); ok {
		ttt := texts[name]
		if holdermark {
			ttt = "- " + ttt + " -"
			obj.SetPlaceHolder(ttt)
		} else {
			obj.SetText(ttt)
		}
	} else if obj, ok := objs.(*widget.Button); ok {
		obj.SetText(texts[name])
	} else if obj, ok := objs.(*widget.Card); ok {
		obj.SetTitle(texts[name])
	} else if obj, ok := objs.(fyne.Window); ok {
		obj.SetTitle(texts[name])
	}
}

func (l *LangChangeManager) ChangeLangByName(name string) {
	if name == "zh" || name == "en" {
	} else {
		return // do nothing
	}

	// check
	if len(l.langs) != len(l.objs) {
		panic("len(words) != len(l.objs)")
		return // do nothing
	}

	l.currentLangName = name
	for i, v := range l.objs {
		l.setTextEx(v, l.langs[i].words, true)
	}
}
