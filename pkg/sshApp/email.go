package sshApp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/nolanjannotta/personal_go_server/pkg/httpServer"
)

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 3)

	m.textInput[0], cmds[0] = m.textInput[0].Update(msg)
	m.textInput[1], cmds[1] = m.textInput[1].Update(msg)
	m.textArea, cmds[2] = m.textArea.Update(msg)

	return tea.Batch(cmds...)
}

func (m *model) UpdateEmailPage(msg tea.Msg) []tea.Cmd {
	var (
		cmds []tea.Cmd
	)

	cmds = append(cmds, m.updateInputs(msg))

	switch m.textInputIndex {
	case 0:
		if !m.textInput[0].Focused() {
			cmds = append(cmds, m.textInput[0].Focus())
		}

		m.textInput[1].Blur()
		m.textArea.Blur()
	case 1:
		if !m.textInput[1].Focused() {
			cmds = append(cmds, m.textInput[1].Focus())
		}
		m.textInput[0].Blur()
		m.textArea.Blur()

	case 2:
		if !m.textArea.Focused() {
			cmds = append(cmds, m.textArea.Focus())
		}
		m.textInput[0].Blur()
		m.textInput[1].Blur()

	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+b":
			m.textArea.Blur()
			m.textInput[0].Blur()
			m.textInput[1].Blur()
			m.textArea.Reset()
			m.textInput[0].Reset()
			m.textInput[1].Reset()
			m.textInputIndex = 0

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
			return append(cmds, tea.Quit)
		case "ctrl+s":
			fmt.Println("SENDING EMAIL")
			m.SendEmail()
			// fmt.Println(emailErr)

		}
	}

	m.content = fmt.Sprint(
		"CONTACT FORM" + "\n\n" +
			"email: " + m.textInput[0].View() + "\n\n" +
			"name: " + m.textInput[1].View() + "\n\n" +
			"message: " + "\n" +
			m.textArea.View() + "\n\n" +
			"'ctrl+s' to send email\n'ctrl+b' to go back\n'ctrl+c' to quit\n'tab' for next text box\n'shift+tab' for previous text box\n\n\n" +
			m.emailSuccessMsg)

	return cmds
}

func (m *model) SendEmail() {
	var emailErr string = "A mysterious error occurred. Please try again or contact me through twitter or warpcast :)"

	// validate that all fields are filled out
	if m.textInput[0].Value() == "" || m.textInput[1].Value() == "" || m.textArea.Value() == "" {
		m.emailSuccessMsg = "ERROR: please fill out all fields before sending an email"
		return
	}

	// check if return email address is valid
	_, err := mail.ParseAddress(m.textInput[0].Value())
	if err != nil {
		m.emailSuccessMsg = "ERROR: invalid email address detected, please try again."
		return
	}

	email := httpServer.Email{
		From: m.textInput[0].Value(),
		Name: m.textInput[1].Value(),
		Msg:  m.textArea.Value(),
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
		// panic(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		m.emailSuccessMsg = emailErr
		return
		// fmt.Println(res.Status)
	}
	m.emailSuccessMsg = "Email sent successfully! I'll get back to you as soon as I can. Thanks for reaching out!"

}
