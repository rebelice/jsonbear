package app

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	zone "github.com/lrstanley/bubblezone"
)

const (
	initialInputs = 2
	maxInputs     = 6
	minInputs     = 1
	helpHeight    = 5
)

type Model struct {
	Inputs []textarea.Model
	Help   help.Model
	KeyMap keyMap
	Focus  int

	Height int
	Width  int
}

var (
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))

	cursorLineStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("57")).
			Foreground(lipgloss.Color("230"))

	placeholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("238"))

	endOfBufferStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("235"))

	FocusedPlaceholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("99"))

	FocusedBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("238"))

	blurredBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.HiddenBorder())
)

func newTextarea() textarea.Model {
	t := textarea.New()
	t.Prompt = ""
	t.Placeholder = "Type something"
	t.ShowLineNumbers = true
	t.Cursor.Style = cursorStyle
	t.FocusedStyle.Placeholder = FocusedPlaceholderStyle
	t.BlurredStyle.Placeholder = placeholderStyle
	t.FocusedStyle.CursorLine = cursorLineStyle
	t.FocusedStyle.Base = FocusedBorderStyle
	t.BlurredStyle.Base = blurredBorderStyle
	t.FocusedStyle.EndOfBuffer = endOfBufferStyle
	t.BlurredStyle.EndOfBuffer = endOfBufferStyle
	t.KeyMap.DeleteWordBackward.SetEnabled(false)
	t.KeyMap.LineNext = key.NewBinding(key.WithKeys("down"))
	t.KeyMap.LinePrevious = key.NewBinding(key.WithKeys("up"))
	t.Blur()
	return t
}

type keyMap = struct {
	next, generate, quit key.Binding
}

func NewModel() Model {
	zone.NewGlobal()

	m := Model{
		Inputs: make([]textarea.Model, 2),
		Help:   help.New(),
		KeyMap: keyMap{
			next: key.NewBinding(
				key.WithKeys("ctrl+n"),
				key.WithHelp("ctrl+n", "next"),
			),
			generate: key.NewBinding(
				key.WithKeys("ctrl+enter"),
				key.WithHelp("ctrl+enter", "generate"),
			),
			quit: key.NewBinding(
				key.WithKeys("esc", "ctrl+c"),
				key.WithHelp("esc", "quit"),
			),
		},
	}
	for i := range m.Inputs {
		m.Inputs[i] = newTextarea()
	}
	m.Inputs[m.Focus].Focus()
	return m
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.quit):
			for i := range m.Inputs {
				m.Inputs[i].Blur()
			}
			return m, tea.Quit
		case key.Matches(msg, m.KeyMap.next):
			m.Inputs[m.Focus].Blur()
			m.Focus++
			if m.Focus > len(m.Inputs)-1 {
				m.Focus = 0
			}
			cmd := m.Inputs[m.Focus].Focus()
			cmds = append(cmds, cmd)
		}
	case tea.WindowSizeMsg:
		m.Height = msg.Height
		m.Width = msg.Width
	}
	m.sizeInputs()

	// Update all textareas
	for i := range m.Inputs {
		newModel, cmd := m.Inputs[i].Update(msg)
		m.Inputs[i] = newModel
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) sizeInputs() {
	for i := range m.Inputs {
		m.Inputs[i].SetWidth(m.Width / len(m.Inputs))
		m.Inputs[i].SetHeight(m.Height - helpHeight)
	}
}

func (m Model) View() string {
	help := m.Help.ShortHelpView([]key.Binding{
		m.KeyMap.next,
		m.KeyMap.generate,
		m.KeyMap.quit,
	})

	var views []string
	for i := range m.Inputs {
		views = append(views, m.Inputs[i].View())
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, views...) + "\n\n" + help
}
