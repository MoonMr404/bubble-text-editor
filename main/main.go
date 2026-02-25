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
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle().Background(lipgloss.Color("#597db1"))
	ta.FocusedStyle.CursorLineNumber = lipgloss.NewStyle().Foreground(lipgloss.Color("#6696d9")).Bold(true)

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
	files, _ := os.ReadDir(".")
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
	// 1. Vista Input Nome File
	if m.textInput.Focused() {
		content := fmt.Sprintf("Crea Nuovo File\n\n%s\n\n(enter: conferma, esc: annulla)", m.textInput.View())
		return appStyle.Render(boxStyle.Render(content))
	}

	// 2. Vista Editor Split-Screen
	if m.isEditing {
		// Calcolo larghezze per le due colonne
		halfWidth := (m.width / 2) - 6
		m.textArea.SetWidth(halfWidth)
		m.textArea.SetHeight(m.height - 12)

		// Colonna Sinistra (Editor)
		leftView := normalItemStyle.Width(halfWidth).Render(m.textArea.View())

		// Colonna Destra (Anteprima Bold)
		rightView := setTextToBold.Width(halfWidth).Render(setTextToBold.Render(m.textArea.Value()))

		// Unione orizzontale
		mainContent := lipgloss.JoinHorizontal(lipgloss.Top, leftView, rightView)

		header := fmt.Sprintf("Editing: %s\n\n", m.currentFile)
		m.footer = "\n\n(esc: back | ctrl+s: save | ctrl+q: quit)"

		return appStyle.Render(boxStyle.Render(header + mainContent + m.footer))
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
