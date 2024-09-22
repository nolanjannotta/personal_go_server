package sshApp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/nolanjannotta/personal_go_server/pkg/httpServer"
)

func (m *model) ResetEmailForm() {
	m.msgInput.Blur()
	m.nameInput.Blur()
	m.fromInput.Blur()

	m.msgInput.Reset()
	m.nameInput.Reset()
	m.fromInput.Reset()
	m.textInputIndex = 0

}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 3)

	m.nameInput, cmds[0] = m.nameInput.Update(msg)
	m.fromInput, cmds[1] = m.fromInput.Update(msg)
	m.msgInput, cmds[2] = m.msgInput.Update(msg)

	return tea.Batch(cmds...)
}

func (m *model) RenderEmailPage() {
	m.content = m.emailLayout.Render(
		fmt.Sprint(
			"\n\nCONTACT FORM\n\n" +

				lipgloss.JoinVertical(
					lipgloss.Left,
					"name: "+m.nameInput.View()+"\n",
					"email: "+m.fromInput.View()+"\n",
					"message: "+"\n",
					m.msgInput.View(),
					m.emailSuccessMsg),
		))
}

func (m *model) UpdateEmailPage(msg tea.Msg) tea.Cmd {
	var (
		cmds []tea.Cmd
	)

	switch m.textInputIndex {
	case 0:
		if !m.nameInput.Focused() {
			cmds = append(cmds, m.nameInput.Focus())
		}
		m.fromInput.Blur()
		m.msgInput.Blur()

	case 1:
		if !m.fromInput.Focused() {
			cmds = append(cmds, m.fromInput.Focus())
		}
		m.nameInput.Blur()
		m.msgInput.Blur()

	case 2:
		if !m.msgInput.Focused() {
			cmds = append(cmds, m.msgInput.Focus())
		}
		m.fromInput.Blur()
		m.nameInput.Blur()

	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+b":

			m.pageName = "index"
			m.content = m.markdown["index"]
			m.emailSuccessMsg = ""
			return nil
		case "tab":

			m.textInputIndex++
			if m.textInputIndex > 2 {
				m.textInputIndex = 2
			}
		case "shift+tab":

			m.textInputIndex--
			if m.textInputIndex < 0 {
				m.textInputIndex = 0
			}
		case "ctrl+c":
			return tea.Quit
		case "ctrl+s":
			fmt.Println("SENDING EMAIL")
			m.SendEmail()

		}
	}

	cmds = append(cmds, m.updateInputs(msg))

	m.RenderEmailPage()

	return tea.Batch(cmds...)
}

func (m *model) SendEmail() {
	var emailErr string = "A mysterious error occurred. Please try again or contact me through twitter or warpcast :)"

	// validate that all fields are filled out
	if m.nameInput.Value() == "" || m.fromInput.Value() == "" || m.msgInput.Value() == "" {
		m.emailSuccessMsg = "ERROR: please fill out all fields before sending an email"
		return
	}

	// check if return email address is valid
	_, err := mail.ParseAddress(m.fromInput.Value())
	if err != nil {
		m.emailSuccessMsg = "ERROR: invalid email address detected, please try again."
		return
	}

	email := httpServer.Email{
		From: m.fromInput.Value(),
		Name: m.nameInput.Value(),
		Msg:  m.msgInput.Value(),
	}

	b, err := json.Marshal(email)
	if err != nil {
		m.emailSuccessMsg = emailErr
		return
	}

	request, err := http.NewRequest("POST", "http://localhost:8080/email", bytes.NewBuffer(b))
	if err != nil {
		m.emailSuccessMsg = emailErr
		return
	}

	client := &http.Client{}
	res, err := client.Do(request)

	if err != nil {
		m.emailSuccessMsg = emailErr
		return
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		m.emailSuccessMsg = emailErr
		return
	}
	m.emailSuccessMsg = "Email sent successfully! I'll get back to you as soon as I can. Thanks for reaching out!"

}
