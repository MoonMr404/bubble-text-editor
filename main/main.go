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

// CREATE DISTINCT VIEWS
// Nested Struct>
type editingView struct {
	editingArea textarea.Model
	content     string
	//undoStack
}
type model struct {
	textArea    textarea.Model
	textInput   textinput.Model
	footer      string
	files       []string
	currentFile string
	width       int
	height      int
	cursor      int
	isEditing   bool
}

func initialModel() model {
	ta := textarea.New()
	ta.Focus()
	ta.ShowLineNumbers = true

	// editing cursor style

	// needs to be dynamic
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle().Background(lipgloss.Color("#087249"))
	ta.FocusedStyle.CursorLineNumber = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Bold(true)
	
	ti := textinput.New()
	ti.Placeholder = "Nome del nuovo file..."

	// Caricamento file da argomenti
	if len(os.Args) > 1 {
		fn := os.Args[1]
		content, _ := os.ReadFile(fn)
		ta.SetValue(string(content))
		return model{textArea: ta, textInput: ti, isEditing: true, currentFile: fn}
	}

	// Caricamento lista file
	files, _ := os.ReadDir("C:\\Users\\Francesco\\Desktop\\Appunti\\") //Win vers
	// files, _ := os.ReadDir("/Users/francesco/Desktop/Appunti/") macos vers
	var fileList []string

	for _, f := range files {
		if !f.IsDir() {
			fileList = append(fileList, f.Name())
		}

	}

	return model{textArea: ta, textInput: ti, files: fileList, isEditing: false}
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
		return m, nil

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlQ {
			return m, tea.Quit
		}

		// Gestione Input Nuovo File
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
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}

		// Gestione Lista File
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
			case "n":
				return m, m.textInput.Focus()
			case "enter":
				if len(m.files) > 0 {
					nomeFile := m.files[m.cursor]
					content, _ := os.ReadFile(nomeFile)
					m.textArea.SetValue(string(content))
					m.textArea.SetCursor(0)
					m.currentFile = nomeFile
					m.isEditing = true
					m.textArea.Focus()
				}
			case tea.KeyBackspace.String():
				//TODO fix file removing in other folders
				os.Remove(m.files[m.cursor])
				m.files = append(m.files[:m.cursor], m.files[m.cursor+1:]...)
				return m, nil
			}
			return m, nil
		}

		// Gestione Editor
		switch msg.Type {
		//Save commands

		case tea.KeyCtrlS:
			val := m.textArea.Value()
			functions.UpdateFile(m.currentFile, val)
			// m.footer = "Saved successfully"
			// m.textArea.SetValue("Saved successfully")
			return m, nil
		case tea.KeyTab:
			m.textArea.InsertString("  ")
			return m, nil
		case tea.KeyEsc:
			m.isEditing = false
			m.textArea.Blur()
			return m, nil
		}
	}

	if m.isEditing {
		m.textArea, cmd = m.textArea.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	//New file managing
	if m.textInput.Focused() {
		content := fmt.Sprintf("Create new file\n\n%s\n\n(enter: confirm, esc: cancel)", m.textInput.View())
		return appStyle.Render(boxStyle.Render(content))
	}

	//Editing file UI
	if m.isEditing {

		halfWidth := (m.width / 2) - 4
		m.textArea.SetWidth(halfWidth)
		m.textArea.SetHeight(m.height - 15)
		// Stile di sfondo comune
		bgStyle := lipgloss.NewStyle().Background(lipgloss.Color("#087249"))

		// leftView := normalItemStyle.Width(halfWidth).Height(m.height - 10).Render(m.textArea.View())

		//Define style in styles.go
		// Rendering cursor?

		leftView := normalItemStyle.
			Width((m.width / 2) - 4).
			Height(m.height - 10).
			Padding(2).
			Render(m.textArea.View())

		// pass leftVew buffer for rendering?
		rightView := bgStyle.Bold(true).Width(halfWidth).Height(m.height - 10).Padding(2).Render(m.textArea.Value())

		// Unione orizzontale
		joinedViews := lipgloss.JoinHorizontal(lipgloss.Top, leftView, rightView)

		mainContent := bgStyle.Width(m.width - 4).Render(joinedViews)

		header := fmt.Sprintf("Editing: %s\n\n", m.currentFile)
		footer := "\n\n(esc: back | ctrl+s: save | ctrl+q: quit)"
		//Footer style
		footer = lipgloss.NewStyle().Background(lipgloss.Color("#707208")).Width(m.width).Height(5).Render(footer)

		return appStyle.Render(boxStyle.Render(header + mainContent + footer))
	}

	// 3. Vista Lista File
	s := "Select file\n\n"
	for i, file := range m.files {
		if m.cursor == i {
			s += selectedItemStyle.Render("> "+file) + "\n"
		} else {
			s += fmt.Sprintf("  %s\n", file)
		}
	}
	s += "\n(enter: open | n: new file | backspace: delete | ctrl+q: quit)"

	return appStyle.Render(boxStyle.Render(s))
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Errore: %v", err)
		os.Exit(1)
	}
}
