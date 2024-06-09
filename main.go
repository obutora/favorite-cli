package main

import (
	"fmt"
	"hagakun/service"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle   = focusedStyle
	noStyle       = lipgloss.NewStyle()
	focusedButton = focusedStyle.Render("[ 登録 ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("登録"))
)

type item struct {
	title string
	desc  string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	focusIndex   int
	inputs       []textinput.Model
	isAddMode    bool
	isDeleteMode bool
	list         list.Model
}

func initialModel() model {
	items := service.ReadItems()

	var models = []list.Item{}
	for _, e := range items {
		models = append(models, item{
			title: e.Title,
			desc:  e.Desc,
		})
	}

	m := model{
		list:         list.New(models, list.NewDefaultDelegate(), 0, 0),
		inputs:       make([]textinput.Model, 2),
		isAddMode:    false,
		isDeleteMode: false,
	}

	m.list.Title = "Bookmarks"

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 1024

		switch i {
		case 0:
			t.Placeholder = "Name: ブックマークを表す端的な説明を入力してください"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Address: https://example.com"
		}

		m.inputs[i] = t
	}

	return m
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "ctrl+a":
			m.isAddMode = true
			m.list.Title = "Add Bookmark"

		case "ctrl+d":
			m.isDeleteMode = true
			m.list.Title = "Delete Bookmark"

		case "enter":
			if m.isAddMode && m.focusIndex == len(m.inputs) {
				fmt.Printf("name: %s, address: %s\n", m.inputs[0].Value(), m.inputs[1].Value())
				service.AddItem(m.inputs[0].Value(), m.inputs[1].Value())
				fmt.Printf("add: %s\n", m.inputs[0].Value())
				return m, tea.Quit
			}

			if m.isDeleteMode {
				service.DeleteItem(m.list.SelectedItem().(item).Title())
				fmt.Printf("delete: %s\n", m.list.SelectedItem().(item).Title())
				return m, tea.Quit
			}

			if !m.isAddMode && !m.isDeleteMode {
				service.OpenURL(m.list.SelectedItem().(item).Description())
				return m, tea.Quit
			}

		// Set focus to next input
		case "tab", "shift+tab":
			s := msg.String()

			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	if m.isAddMode {
		cmd := m.updateInputs(msg)
		return m, cmd
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m model) View() string {
	if m.isAddMode {
		var b strings.Builder

		for i := range m.inputs {
			b.WriteString(m.inputs[i].View())
			if i < len(m.inputs)-1 {
				b.WriteRune('\n')
			}
		}

		button := &blurredButton
		if m.focusIndex == len(m.inputs) {
			button = &focusedButton
		}
		fmt.Fprintf(&b, "\n\n%s\n\n", *button)

		// b.WriteString(helpStyle.Render("cursor mode is "))
		// b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
		// b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))
		return b.String()
	}
	return docStyle.Render(m.list.View())
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
