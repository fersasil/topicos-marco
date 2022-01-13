package main

import (
	"fmt"
)

type AnsiColors struct {
	black   string
	red     string
	green   string
	yellow  string
	blue    string
	magenta string
	cyan    string
	white   string
}

type AnsiEffects struct {
	reset     string // \u001b[0m
	bold      string // \u001b[1m
	underline string // \u001b[4m
	reversed  string // \u001b[7m
}

type AnsiFontStyle struct {
	effects         AnsiEffects
	fontColor       AnsiColors
	backgroundColor AnsiColors
}

var AnsiStyle = AnsiFontStyle{
	effects: AnsiEffects{
		reset:     "\u001b[0m",
		bold:      "\u001b[1m",
		underline: "\u001b[4m",
		reversed:  "\u001b[7m",
	},
	fontColor: AnsiColors{
		black:   "\u001b[30m",
		red:     "\u001b[31m",
		green:   "\u001b[32m",
		yellow:  "\u001b[33m",
		blue:    "\u001b[34m",
		magenta: "\u001b[35m",
		cyan:    "\u001b[36m",
		white:   "\u001b[37m",
	},
	backgroundColor: AnsiColors{
		black:   "\u001b[40m",
		red:     "\u001b[41m",
		green:   "\u001b[42m",
		yellow:  "\u001b[43m",
		blue:    "\u001b[44m",
		magenta: "\u001b[45m",
		cyan:    "\u001b[46m",
		white:   "\u001b[47m",
	},
}

func testCli() {
	fmt.Println(AnsiStyle.effects.underline)
	fmt.Printf("%s%s%sPão de batata%s\n",
		AnsiStyle.effects.bold,
		AnsiStyle.fontColor.red,
		AnsiStyle.backgroundColor.white,
		AnsiStyle.effects.reset)
}
