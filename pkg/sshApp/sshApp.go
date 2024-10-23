package sshApp

import (
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/charmbracelet/bubbles/help"
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
	term        string
	profile     string
	width       int
	height      int
	renderer    *lipgloss.Renderer
	footerStyle lipgloss.Style
	// helpStyle   lipgloss.Style

	emailLayout lipgloss.Style
	viewport    viewport.Model

	nameInput textinput.Model
	fromInput textinput.Model
	msgInput  textarea.Model

	textInputIndex  int
	content         string
	pageName        string
	emailSuccessMsg string
	markdown        map[string]string
	footer          help.Model
}

func SetUp() *ssh.Server {

	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(fmt.Sprint(home, "/.ssh/personal-server")),
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

func Start(s *ssh.Server, done chan<- os.Signal) {
	log.Info("Starting SSH server", "host", host, "port", port)

	if err := s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not start server", "error", err)
		done <- nil
	}

}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {

	// fingerprint := uuid.New().String()

	pty, _, _ := s.Pty()

	renderer := bubbletea.MakeRenderer(s)

	emailLayout := renderer.NewStyle().Width(pty.Window.Width).Align(lipgloss.Center)

	footerStyle := renderer.NewStyle().Foreground(lipgloss.Color("#4f5258"))

	m := model{
		renderer:       renderer,
		term:           pty.Term,
		profile:        renderer.ColorProfile().Name(),
		width:          pty.Window.Width,
		height:         pty.Window.Height,
		emailLayout:    emailLayout,
		footerStyle:    footerStyle,
		msgInput:       textarea.New(),
		nameInput:      textinput.New(),
		fromInput:      textinput.New(),
		textInputIndex: 0,
		pageName:       "index",
		viewport:       viewport.New(pty.Window.Width, pty.Window.Height-2),
		footer:         help.New(),
	}

	m.LoadMarkdown()

	m.nameInput.Placeholder = "your name"
	m.fromInput.Placeholder = "your email"

	m.nameInput.PlaceholderStyle = m.renderer.NewStyle().Foreground(lipgloss.Color("240"))
	m.fromInput.PlaceholderStyle = m.renderer.NewStyle().Foreground(lipgloss.Color("240"))
	m.nameInput.Cursor.Style = m.renderer.NewStyle().Background(lipgloss.AdaptiveColor{Light: "255", Dark: "0"})
	m.fromInput.Cursor.Style = m.renderer.NewStyle().Background(lipgloss.AdaptiveColor{Light: "255", Dark: "0"})

	m.msgInput.FocusedStyle.Placeholder = m.renderer.NewStyle().Foreground(lipgloss.Color("240"))
	m.msgInput.BlurredStyle.Placeholder = m.renderer.NewStyle().Foreground(lipgloss.Color("240"))
	m.msgInput.Cursor.Style = m.renderer.NewStyle().Background(lipgloss.AdaptiveColor{Light: "255", Dark: "0"})

	m.msgInput.Placeholder = "say hi!"
	m.msgInput.SetWidth(pty.Window.Width / 2)
	m.msgInput.SetHeight(10)

	return m, []tea.ProgramOption{tea.WithAltScreen(), tea.WithMouseAllMotion()}
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
		cmds = append(cmds, m.UpdateMarkdownPages(msg))
	case "email":

		cmds = append(cmds, m.UpdateEmailPage(msg))

	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:

		m.height = msg.Height
		m.width = msg.Width
		m.emailLayout = m.renderer.NewStyle().Width(m.width).Align(lipgloss.Center)
		m.viewport = viewport.New(m.width, m.height-2)
	}

	m.viewport, cmd = m.viewport.Update(msg)
	m.viewport.SetContent(m.content)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return fmt.Sprint(
		m.viewport.View(),
		"\n\n",
		m.RenderFooter(),
	)
}
