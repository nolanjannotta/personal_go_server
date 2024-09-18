package sshApp

import (
	"errors"
	"fmt"
	"net"
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
	gossh "golang.org/x/crypto/ssh"
)

const (
	host = "0.0.0.0"
	port = "23234"
)



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
	content         string
	pageName        string
	emailSuccessMsg string
	markdown        map[string]string
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
		wish.WithPublicKeyAuth(func(_ ssh.Context, key ssh.PublicKey) bool {
			return true
		}),
		wish.WithKeyboardInteractiveAuth(func(ctx ssh.Context, challenger gossh.KeyboardInteractiveChallenge) bool {
			return true
		}),
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
	}
	m.LoadMarkdown()
	m.textInput[0] = textinput.New()
	m.textInput[1] = textinput.New()

	m.textInput[0].Placeholder = "your email"
	m.textInput[1].Placeholder = "your name"

	m.textArea.Placeholder = "say hi!"
	m.viewport.YPosition = 2
	m.viewport.HighPerformanceRendering = false
	return m, []tea.ProgramOption{tea.WithAltScreen()}
}

func (m model) Init() tea.Cmd {

	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch m.pageName {

	case "index", "snake", "calculator":
		markdownCmds := m.UpdateMarkdownPages(msg)
		cmds = append(cmds, markdownCmds...)
	case "email":
		// m.textInput[0].Focus()
		emailCmds := m.UpdateEmailPage(msg)
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
}
