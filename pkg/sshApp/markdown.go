package sshApp

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) UpdateMarkdownPages(msg tea.Msg) []tea.Cmd {
	var (
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return append(cmds, tea.Quit)
		case "s":
			m.pageName = "snake"
			m.content = m.markdown["snake"]
			m.viewport.GotoTop()
		case "c":
			m.pageName = "calculator"
			m.content = m.markdown["calculator"]
			m.viewport.GotoTop()
		case "h":
			m.pageName = "index"
			m.content = m.markdown["index"]
			m.viewport.GotoTop()

		case "e":
			emailCmds := m.UpdateEmailPage(msg)
			m.viewport.GotoTop()
			cmds = append(cmds, emailCmds...)
			m.pageName = "email"

		}

	}

	return cmds

}

func (m *model) LoadMarkdown() error {
	var pageLoadErr string = "A mysterious error occurred while loading this page. Please let me know through twitter or warpcast :)"

	m.markdown = make(map[string]string)
	index, err1 := os.ReadFile("../../web/markdown/index.md")
	snake, err2 := os.ReadFile("../../web/markdown/snake.md")
	calculator, err3 := os.ReadFile("../../web/markdown/calculator.md")

	m.markdown["index"] = string(index)
	m.markdown["snake"] = string(snake)
	m.markdown["calculator"] = string(calculator)

	if err1 != nil {
		m.markdown["index"] = pageLoadErr
	}
	if err2 != nil {
		m.markdown["snake"] = pageLoadErr
	}
	if err3 != nil {
		m.markdown["calculator"] = pageLoadErr
	}

	m.content = m.markdown["index"]
	return nil
}
