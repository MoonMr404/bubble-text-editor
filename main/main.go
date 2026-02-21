package main

import (
	"bubble-text-editor/functions"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	textArea  textarea.Model
	textInput textinput.Model
	files     []string
	width     int
	cursor    int
	isEditing bool
	height    int
}

func initialModel() model {
	// Inizializza SEMPRE una textarea base per evitare nil pointer

	ta := textarea.New()
	ta.Focus()
	ta.ShowLineNumbers = true

	ti := textinput.New()
	ti.Placeholder = "Nome del nuovo file..."
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle().
		Background(lipgloss.Color("#597db1"))

	ta.FocusedStyle.CursorLineNumber = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6696d9")).
		Bold(true)

	if len(os.Args) > 1 {
		fn := os.Args[1]
		content, _ := os.ReadFile(fn)
		ta.SetValue(string(content))
		ta.Focus()

		return model{
			textArea:  ta,
			textInput: ti,
			isEditing: true,
		}
	}

	// Modalità lista
	files, _ := os.ReadDir(".")
	var fileList []string
	for _, f := range files {
		if !f.IsDir() {
			fileList = append(fileList, f.Name())
		}
	}

	return model{
		textArea:  ta,
		textInput: ti,
		files:     fileList,
		isEditing: false,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.textArea.SetWidth(m.width - 4)
		m.textArea.SetHeight(m.height - 6)
		return m, nil

	case tea.KeyMsg:
		// 1. TASTI GLOBALI
		if msg.Type == tea.KeyCtrlQ {
			return m, tea.Quit
		}

		//managing new file 
		if m.textInput.Focused() {
			switch msg.Type {
			case tea.KeyEnter:
				fn := functions.CreateNewFile(m.textInput.Value())
				m.files = append(m.files, fn)
				m.textInput.Blur()
				m.textInput.Reset()
				return m, nil
			case tea.KeyEsc:
				m.textInput.Blur()
				m.textInput.Reset()
				return m, nil
			}
			// Aggiorna l'input mentre scrivi
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}

		//List logic
		if !m.isEditing {
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.files)-1 {
					m.cursor++
				}
			//create new file
			case "n":
				return m, m.textInput.Focus()

			//open editor
			case "enter":
				if len(m.files) > 0 {
					nomeFile := m.files[m.cursor]
					content, _ := os.ReadFile(nomeFile)
					m.textArea.SetValue(string(content))
					m.textArea.SetCursor(0)
					m.textArea.Focus()
					m.isEditing = true
				}
			case tea.KeyBackspace.String():
				os.Remove(m.files[m.cursor])
				m.files = append(m.files[:m.cursor], m.files[m.cursor+1:]...)
				return m, nil
			}
			return m, nil
		}

		// 4. LOGICA EDITOR (Solo se isEditing è true)
		switch msg.Type {
		case tea.KeyTab:
			m.textArea.InsertString("    ")
			return m, nil
		case tea.KeyEscape:
			m.isEditing = false
			m.textArea.Blur()
			return m, nil
		}
	}

	// Aggiorna la textarea se siamo in editing
	if m.isEditing {
		m.textArea, cmd = m.textArea.Update(msg)
	}
	return m, cmd
}
func (m model) View() string {
	if m.textInput.Focused() {
		content := fmt.Sprintf(
			"Crea Nuovo File\n\n%s\n\n%s",
			m.textInput.View(),
			"(enter per confermare, esc per annullare)",
		)
		return appStyle.Render(boxStyle.Render(content))
	}
	// 1. Se siamo in modalità editor, mostriamo la TextArea
	if m.isEditing {
		content := fmt.Sprintf(
			"Editing file...\n\n%s\n\n%s",
			m.textArea.View(),
			"(ctrl+q to quit, esc to go back)",
		)

		return appStyle.Render(
			boxStyle.Render(content),
		)
	}

	// 2. Se NON siamo in editing, mostriamo la lista dei file
	s := "Select file\n\n"
	for i, file := range m.files {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, file)
	}
	s += "\n(`Enter` to open file, `ctrl+q` to quit, `n` to create a new file)"

	return appStyle.Render(
		boxStyle.Render(s),
	)
}

func main() {
	p := tea.NewProgram(
		initialModel(),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Errore: %v", err)
		os.Exit(1)
	}
}
