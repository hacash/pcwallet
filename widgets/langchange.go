package widgets

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

// map[string]string{"en": "", "zh": ""}
// langChangeManager.SetText(accshow, map[string]string{"en": "", "zh": ""})

type LangItem struct {
	words map[string]string
}

type LangChangeManager struct {
	currentLangName string // en zh
	objs            []interface{}
	langs           []*LangItem
}

func NewLangChangeManager() *LangChangeManager {
	return &LangChangeManager{
		currentLangName: "en",
		objs:            make([]interface{}, 0),
		langs:           make([]*LangItem, 0),
	}
}

func (l *LangChangeManager) NewTextWrapWordLabel(texts map[string]string) *widget.Label {
	lb := widget.NewLabel(texts["en"])
	lb.Wrapping = fyne.TextWrapBreak
	// add ary
	l.objs = append(l.objs, lb)
	l.langs = append(l.langs, &LangItem{words: texts})
	// return ok
	return lb
}

func (l *LangChangeManager) NewEntrySetPlaceHolder(texts map[string]string) *widget.Entry {
	lb := widget.NewEntry()
	lb.SetPlaceHolder("- " + texts["en"] + " -")
	// add ary
	l.objs = append(l.objs, lb)
	l.langs = append(l.langs, &LangItem{words: texts})
	// return ok
	return lb
}

func (l *LangChangeManager) NewCardSetTitle(texts map[string]string, content fyne.CanvasObject) *widget.Card {
	lb := widget.NewCard(texts["en"], "", content)
	// add ary
	l.objs = append(l.objs, lb)
	l.langs = append(l.langs, &LangItem{words: texts})
	// return ok
	return lb
}

func (l *LangChangeManager) NewButton(texts map[string]string, tapped func()) *widget.Button {
	lb := widget.NewButton(texts["en"], tapped)
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
