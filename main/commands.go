// display content file
// render markdown
// choose files
// on close file
package main

import (
	"bubble-text-editor/functions"

	tea "github.com/charmbracelet/bubbletea"
)

type ResultMsg struct {
	Result string
	Err    error
}

func EditFile(filename string) tea.Cmd {
	//open text editor
	return func() tea.Msg {
		res, err := functions.ReadFile((filename))

		return ResultMsg{Result: res, Err: err}
	}
}

func RenderMarkdown(input string) string {
	return setTextToBold.Render(input)
}
