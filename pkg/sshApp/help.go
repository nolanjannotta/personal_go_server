package sshApp

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m *model) UpdateHelpPage() {

	options := getHelp(m.pageName)

	var optsStr string

	for _, opts := range options {
		optsStr += "• " + opts + "\n"
	}

	optsStr = "\n\n" + "HELP" + "\n\n" + optsStr + "• h back"
	m.content = optsStr

}

func getHelp(page string) (options []string) {

	switch page {
	case "index":
		options = []string{"↑/scroll move up", "↓/scroll move down", "↓/scroll move down", "ctrl+c/q quit"}
	case "snake":
		options = []string{"↑/scroll move up", "↓/scroll move down", "h home", "ctrl+c/q quit"}
	case "calculator":
		options = []string{"↑/scroll move up", "↓/scroll move down", "h home", "ctrl+c/q quit"}
	case "email":
		options = []string{"ctrl+s send email", "ctrl+b home", "ctrl+c quit", "tab next text box", "shift+tab previous text box"}

	}
	return
}

func (m *model) RenderFooter() string {
	options := getHelp(m.pageName)

	var optsStr string

	for i, opts := range options {
		if i == len(options)-1 {
			optsStr += opts
		} else {
			optsStr += opts + "  •  "
		}

	}

	if lipgloss.Width(optsStr) >= m.width-5 {
		optsStr = "ctrl+h help"
	}

	repeats := m.width - lipgloss.Width(optsStr) - 5
	if repeats < 0 {
		repeats = 0
	}

	return m.footerStyle.Render(
		fmt.Sprint(
			optsStr,
			strings.Repeat(" ", repeats),
			fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100),
		))

}
