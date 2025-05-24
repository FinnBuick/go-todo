package main

import "github.com/rivo/tview"

type UI struct {
	app            *tview.Application
	mainFlex       *tview.Flex
	list           *tview.List
	helpText       *tview.TextView
	inputField     *tview.InputField
	isShowingInput bool
}
