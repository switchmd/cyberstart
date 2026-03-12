package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	"cyberstart/internal/tui"
)

func main() {
	p := tea.NewProgram(tui.NewModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "오류: %v\n", err)
		os.Exit(1)
	}
}
