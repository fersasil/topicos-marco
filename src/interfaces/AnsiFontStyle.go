package main

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
