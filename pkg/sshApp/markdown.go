package sshApp

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) UpdateMarkdownPages(msg tea.Msg) tea.Cmd {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return tea.Quit
		case "s":
			m.pageName = "snake"
			m.content = m.markdown[m.pageName]
			m.viewport.GotoTop()
		case "c":
			m.pageName = "calculator"
			m.content = m.markdown[m.pageName]
			m.viewport.GotoTop()
		case "h":
			m.pageName = "index"
			m.content = m.markdown[m.pageName]
			m.viewport.GotoTop()
		case "e":
			m.pageName = "email"
			m.RenderEmailPage()
			m.viewport.GotoTop()
			return m.nameInput.Focus()
		case "ctrl+h":
			m.UpdateHelpPage()

		}

	}

	return nil

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
