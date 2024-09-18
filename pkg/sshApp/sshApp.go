package sshApp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/mail"
	"os"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/nolanjannotta/personal_go_server/pkg/httpServer"
)

const (
	host = "0.0.0.0"
	port = "23234"
)

var md map[string]string // markdown content

var pageLoadErr string = "A mysterious error occurred while loading this page. Please let me know through twitter or warpcast :)"
var emailErr string = "A mysterious error occurred. Please try again or contact me through twitter or warpcast :)"

type model struct {
	term            string
	profile         string
	width           int
	height          int
	pStyle          lipgloss.Style
	borderStyle     lipgloss.Style
	titleStyle      lipgloss.Style
	footerStyle     lipgloss.Style
	viewport        viewport.Model
	textArea        textarea.Model
	textInput       []textinput.Model
	textInputIndex  int
	content         string //whats being displayed currently
	pageName        string
	emailSuccessMsg string
}

func SetUp() *ssh.Server {

	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath("./.ssh/nolanj"),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	return s

}

func Start(s *ssh.Server, done chan os.Signal) {
	log.Info("Starting SSH server", "host", host, "port", port)

	if err := s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not start server", "error", err)
		done <- nil
	}

}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	// This should never fail, as we are using the activeterm middleware.
	pty, _, _ := s.Pty()

	renderer := bubbletea.MakeRenderer(s)
	pStyle := renderer.NewStyle().
		Foreground(lipgloss.Color("10"))

	borderStyle := renderer.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		PaddingRight(4).
		PaddingLeft(1).
		Width(pty.Window.Width - 2).
		Align(lipgloss.Left).
		TabWidth(0)

	titleStyle := renderer.NewStyle().
		Foreground(lipgloss.Color("10")).
		Width(pty.Window.Width).
		Bold(true)

	footerStyle := renderer.NewStyle().Foreground(lipgloss.Color("8")).Align(lipgloss.Left)

	loadMarkdown()

	m := model{
		term:           pty.Term,
		profile:        renderer.ColorProfile().Name(),
		width:          pty.Window.Width,
		height:         pty.Window.Height,
		titleStyle:     titleStyle,
		borderStyle:    borderStyle,
		pStyle:         pStyle,
		footerStyle:    footerStyle,
		textArea:       textarea.New(),
		textInput:      make([]textinput.Model, 2),
		textInputIndex: 0,
		pageName:       "index",
		viewport:       viewport.New(pty.Window.Width, pty.Window.Height-3),
		content:        md["index"],
	}

	m.textInput[0] = textinput.New()
	m.textInput[1] = textinput.New()

	m.textInput[0].Placeholder = "your email"
	m.textInput[1].Placeholder = "your name"

	m.textArea.Placeholder = "say hi!"
	m.viewport.YPosition = 2
	m.viewport.HighPerformanceRendering = false
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

func loadMarkdown() error {

	md = make(map[string]string)
	index, err1 := os.ReadFile("../../web/markdown/index.md")
	snake, err2 := os.ReadFile("../../web/markdown/snake.md")
	calculator, err3 := os.ReadFile("../../web/markdown/calculator.md")

	md["index"] = string(index)
	md["snake"] = string(snake)
	md["calculator"] = string(calculator)

	if err1 != nil {
		md["index"] = pageLoadErr
	}
	if err2 != nil {
		md["snake"] = pageLoadErr
	}
	if err3 != nil {
		md["calculator"] = pageLoadErr
	}

	return nil
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, 3)

	m.textInput[0], cmds[0] = m.textInput[0].Update(msg)
	m.textInput[1], cmds[1] = m.textInput[1].Update(msg)
	m.textArea, cmds[2] = m.textArea.Update(msg)

	return tea.Batch(cmds...)
}

func (m *model) sendEmail() {

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

func (m model) Init() tea.Cmd {

	return nil
}

func updateMarkdownPages(m *model, msg tea.Msg) []tea.Cmd {
	var (
		// cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return append(cmds, tea.Quit)
		case "s":
			m.pageName = "snake"
			m.content = md[m.pageName]
			m.viewport.GotoTop()
		case "c":
			m.pageName = "calculator"
			m.content = md[m.pageName]
			m.viewport.GotoTop()
		case "h":
			m.pageName = "index"
			m.content = md[m.pageName]
			m.viewport.GotoTop()

		case "e":
			emailCmds := updateEmailPage(m, msg)
			m.viewport.GotoTop()
			cmds = append(cmds, emailCmds...)
			m.pageName = "email"

		}

	}
	return cmds

}

func updateEmailPage(m *model, msg tea.Msg) []tea.Cmd {
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
		// m.textInput[0], cmd = m.textInput[0].Update(msg)

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

	// cmds := make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+b":
			m.textArea.Blur()
			m.textArea.Reset()
			m.textInput[0].Blur()
			m.textInput[0].Reset()
			m.textInput[1].Blur()
			m.textInput[1].Reset()
			m.textInputIndex = 0

			m.pageName = "index"
			m.content = md[m.pageName]
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
			m.sendEmail()
			// fmt.Println(emailErr)

		}
	}

	// if !m.textArea.Focused() {

	// 	cmds = append(cmds, m.textArea.Focus())
	// }
	// m.textArea, cmd = m.textArea.Update(msg)
	// cmds = append(cmds, cmd)

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

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch m.pageName {
	case "index", "snake", "calculator":
		markdownCmds := updateMarkdownPages(&m, msg)
		cmds = append(cmds, markdownCmds...)
	case "email":
		// m.textInput[0].Focus()
		emailCmds := updateEmailPage(&m, msg)
		cmds = append(cmds, emailCmds...)

	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		m.viewport = viewport.New(msg.Width, msg.Height-3)
	}

	m.viewport.SetContent(m.content)
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return fmt.Sprint(
		m.viewport.View(),
		// m.content,

		"\n\n",
		m.footerStyle.Render(fmt.Sprintf("---------[arrow keys or mouse]🠢 scroll---------[q]🠢 quit---------progress🠢 %3.f%%---------", m.viewport.ScrollPercent()*100)))

	// return m.borderStyle.Render(m.titleStyle.Render(space+"Nolan Jannotta") + "\n\n\n\n" + m.txtStyle.Render(p1) + "\n\n" + m.quitStyle.Render("Press 'q' to quit\n"))
}
