package tui

import (
	"fmt"
	"strings"
	"sync/atomic"

	"charm.land/bubbles/v2/progress"
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"cyberstart/internal/installer"
)

// ── Phases ──────────────────────────────────────────────────────────

type phase int

const (
	phasePreset  phase = iota // 프리셋 선택
	phaseSteps                // 설치 항목 체크
	phaseInstall              // 설치 진행
	phaseDone                 // 완료 요약
)

// ── Step status ─────────────────────────────────────────────────────

type stepStatus int

const (
	statusPending stepStatus = iota
	statusRunning
	statusDone
	statusFailed
)

// ── Types ───────────────────────────────────────────────────────────

type presetDef struct {
	name    string
	desc    string
	stepIDs []string
}

type stepState struct {
	step    installer.Step
	enabled bool
	status  stepStatus
	errMsg  string
}

// Thread-safe string for live log display from install goroutines.
type atomicString struct{ v atomic.Value }

func (s *atomicString) Store(val string) { s.v.Store(val) }
func (s *atomicString) Load() string {
	if v := s.v.Load(); v != nil {
		return v.(string)
	}
	return ""
}

// ── Messages ────────────────────────────────────────────────────────

type installDoneMsg struct {
	index int
	err   error
}

// ── Model ───────────────────────────────────────────────────────────

type model struct {
	phase    phase
	width    int
	height   int
	quitting bool

	// Preset phase
	presets      []presetDef
	presetCursor int

	// Steps phase
	steps      []stepState
	stepCursor int

	// Install phase
	currentStep int
	currentLog  *atomicString
	spinner     spinner.Model
	progress    progress.Model
	doneCount   int
	totalCount  int
}

// ── Constructor ─────────────────────────────────────────────────────

func NewModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(colorViolet)

	p := progress.New(
		progress.WithColors(colorVioletDim, colorCyan),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)

	presets := []presetDef{
		{
			name:    "기본",
			desc:    "Chrome, VSCode, Bat2Exe",
			stepIDs: []string{"chrome", "vscode", "bat2exe"},
		},
		{
			name:    "Kollus",
			desc:    "Chrome, Kollus Agent, VSCode, Bat2Exe",
			stepIDs: []string{"chrome", "kollus", "vscode", "bat2exe"},
		},
		{
			name:    "PlaynPlay",
			desc:    "Chrome, PlaynPlay, VSCode, Bat2Exe",
			stepIDs: []string{"chrome", "playnplay", "vscode", "bat2exe"},
		},
	}

	allSteps := installer.AllSteps()
	steps := make([]stepState, len(allSteps))
	for i, st := range allSteps {
		steps[i] = stepState{step: st}
	}

	return model{
		phase:      phasePreset,
		presets:    presets,
		steps:      steps,
		spinner:    s,
		progress:   p,
		currentLog: &atomicString{},
	}
}

// ── Tea interface ───────────────────────────────────────────────────

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Global handlers — animation ticks must always be processed.
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		pw := m.innerWidth() - 8
		if pw < 20 {
			pw = 20
		}
		m.progress.SetWidth(pw)
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case progress.FrameMsg:
		var cmd tea.Cmd
		m.progress, cmd = m.progress.Update(msg)
		return m, cmd

	case tea.KeyPressMsg:
		if msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
		if msg.String() == "q" && m.phase != phaseInstall && m.phase != phaseDone {
			m.quitting = true
			return m, tea.Quit
		}
	}

	// Phase-specific handlers.
	switch m.phase {
	case phasePreset:
		return m.updatePreset(msg)
	case phaseSteps:
		return m.updateSteps(msg)
	case phaseInstall:
		return m.updateInstall(msg)
	case phaseDone:
		return m.updateDone(msg)
	}

	return m, nil
}

// ── Width helpers ───────────────────────────────────────────────────

func (m model) outerWidth() int {
	w := m.width - 4
	if m.width == 0 {
		w = 56 // sensible default before first WindowSizeMsg
	}
	if w > 62 {
		w = 62
	}
	if w < 36 {
		w = 36
	}
	return w
}

func (m model) innerWidth() int {
	return m.outerWidth() - 6 // 2 border chars + 4 padding chars
}

// ── Preset Phase ────────────────────────────────────────────────────

func (m model) updatePreset(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyPressMsg); ok {
		switch msg.String() {
		case "up", "k":
			if m.presetCursor > 0 {
				m.presetCursor--
			}
		case "down", "j":
			if m.presetCursor < len(m.presets)-1 {
				m.presetCursor++
			}
		case "enter":
			preset := m.presets[m.presetCursor]
			enabled := make(map[string]bool)
			for _, id := range preset.stepIDs {
				enabled[id] = true
			}
			for i := range m.steps {
				m.steps[i].enabled = enabled[m.steps[i].step.ID]
			}
			m.phase = phaseSteps
			m.stepCursor = 0
		}
	}
	return m, nil
}

// ── Steps Phase ─────────────────────────────────────────────────────

func (m model) updateSteps(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyPressMsg); ok {
		switch msg.String() {
		case "up", "k":
			if m.stepCursor > 0 {
				m.stepCursor--
			}
		case "down", "j":
			if m.stepCursor < len(m.steps)-1 {
				m.stepCursor++
			}
		case "space":
			m.steps[m.stepCursor].enabled = !m.steps[m.stepCursor].enabled
		case "a":
			allEnabled := true
			for _, s := range m.steps {
				if !s.enabled {
					allEnabled = false
					break
				}
			}
			for i := range m.steps {
				m.steps[i].enabled = !allEnabled
			}
		case "enter":
			return m.beginInstall()
		case "esc":
			m.phase = phasePreset
		}
	}
	return m, nil
}

// ── Install Phase ───────────────────────────────────────────────────

func (m model) beginInstall() (model, tea.Cmd) {
	m.phase = phaseInstall
	m.totalCount = 0
	for _, s := range m.steps {
		if s.enabled {
			m.totalCount++
		}
	}
	m.doneCount = 0

	m.currentStep = m.nextEnabled(-1)
	if m.currentStep < 0 {
		m.phase = phaseDone
		return m, nil
	}

	m.steps[m.currentStep].status = statusRunning
	m.currentLog.Store("")

	return m, tea.Batch(
		m.spinner.Tick,
		m.runStep(m.currentStep),
	)
}

func (m model) runStep(idx int) tea.Cmd {
	step := m.steps[idx].step
	log := m.currentLog
	return func() tea.Msg {
		err := step.Run(func(msg string) { log.Store(msg) })
		return installDoneMsg{index: idx, err: err}
	}
}

func (m model) updateInstall(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(installDoneMsg); ok {
		if msg.err != nil {
			m.steps[msg.index].status = statusFailed
			m.steps[msg.index].errMsg = msg.err.Error()
		} else {
			m.steps[msg.index].status = statusDone
		}
		m.doneCount++

		pct := float64(m.doneCount) / float64(m.totalCount)
		pCmd := m.progress.SetPercent(pct)

		next := m.nextEnabled(msg.index)
		if next < 0 {
			installer.CleanupTempDir()
			m.phase = phaseDone
			return m, pCmd
		}

		m.currentStep = next
		m.steps[next].status = statusRunning
		m.currentLog.Store("")
		return m, tea.Batch(m.runStep(next), pCmd)
	}
	return m, nil
}

func (m model) nextEnabled(after int) int {
	for i := after + 1; i < len(m.steps); i++ {
		if m.steps[i].enabled {
			return i
		}
	}
	return -1
}

// ── Done Phase ──────────────────────────────────────────────────────

func (m model) updateDone(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(tea.KeyPressMsg); ok {
		m.quitting = true
		return m, tea.Quit
	}
	return m, nil
}

// ════════════════════════════════════════════════════════════════════
//  VIEW
// ════════════════════════════════════════════════════════════════════

// ── Helpers ─────────────────────────────────────────────────────────

func indent(s string, n int) string {
	pad := strings.Repeat(" ", n)
	lines := strings.Split(s, "\n")
	for i := range lines {
		lines[i] = pad + lines[i]
	}
	return strings.Join(lines, "\n")
}

func helpBar(entries ...string) string {
	var parts []string
	for i := 0; i+1 < len(entries); i += 2 {
		parts = append(parts, helpKeyStyle.Render(entries[i])+" "+helpDescStyle.Render(entries[i+1]))
	}
	return strings.Join(parts, faintStyle.Render("  ·  "))
}

func (m model) enabledList() []stepState {
	var r []stepState
	for _, s := range m.steps {
		if s.enabled {
			r = append(r, s)
		}
	}
	return r
}

// ── Breadcrumb ──────────────────────────────────────────────────────

func (m model) breadcrumb() string {
	names := []string{"프리셋", "항목", "설치", "완료"}
	idx := int(m.phase)

	var parts []string
	for i, name := range names {
		var s string
		if m.phase == phaseDone || i < idx {
			s = bcDoneStyle.Render("✓") + " " + bcDoneStyle.Render(name)
		} else if i == idx {
			s = bcActiveStyle.Render("●") + " " + bcActiveStyle.Render(name)
		} else {
			s = bcPendingStyle.Render("○") + " " + bcPendingStyle.Render(name)
		}
		parts = append(parts, s)
	}

	return strings.Join(parts, bcSepStyle.Render(" ── "))
}

// ── Main View ───────────────────────────────────────────────────────

func (m model) View() tea.View {
	if m.quitting {
		return tea.NewView("")
	}

	var b strings.Builder

	// Header
	b.WriteString("\n")
	b.WriteString("  " + logoStyle.Render("✦ cyberstart"))
	b.WriteString("\n\n")

	// Breadcrumb stepper
	b.WriteString("  " + m.breadcrumb())
	b.WriteString("\n\n")

	// Phase content
	switch m.phase {
	case phasePreset:
		m.viewPreset(&b)
	case phaseSteps:
		m.viewSteps(&b)
	case phaseInstall:
		m.viewInstall(&b)
	case phaseDone:
		m.viewDone(&b)
	}

	v := tea.NewView(b.String())
	v.AltScreen = true
	return v
}

// ── Preset view ─────────────────────────────────────────────────────

func (m model) viewPreset(b *strings.Builder) {
	b.WriteString("  " + sectionStyle.Render("설치 프리셋을 선택하세요."))
	b.WriteString("\n\n")

	// Build card content
	var c strings.Builder
	for i, p := range m.presets {
		if i > 0 {
			c.WriteString("\n")
		}

		cur := "  "
		dot := bcPendingStyle.Render("○")
		name := normalStyle.Render(p.name)
		desc := descStyle.Render(p.desc)

		if i == m.presetCursor {
			cur = cursorStyle.Render("❯ ")
			dot = bcActiveStyle.Render("●")
			name = selectedStyle.Render(p.name)
		}

		c.WriteString(cur + dot + " " + name + "\n")
		c.WriteString("    " + desc)

		if i < len(m.presets)-1 {
			c.WriteString("\n")
		}
	}

	card := cardStyle.Width(m.outerWidth()).Render(c.String())
	b.WriteString(indent(card, 2))
	b.WriteString("\n\n")

	b.WriteString("  " + helpBar("↑↓", "이동", "enter", "선택", "q", "종료"))
	b.WriteString("\n")
}

// ── Steps view ──────────────────────────────────────────────────────

func (m model) viewSteps(b *strings.Builder) {
	b.WriteString("  " + sectionStyle.Render("설치할 프로그램을 선택하세요."))
	b.WriteString("\n\n")

	var c strings.Builder
	for i, s := range m.steps {
		cur := "  "
		chk := checkOffStyle.Render("□")
		name := normalStyle.Render(s.step.Name)

		if s.enabled {
			chk = checkOnStyle.Render("■")
		}
		if i == m.stepCursor {
			cur = cursorStyle.Render("❯ ")
			name = selectedStyle.Render(s.step.Name)
		}

		c.WriteString(cur + chk + " " + name)
		if i < len(m.steps)-1 {
			c.WriteString("\n")
		}
	}

	card := cardStyle.Width(m.outerWidth()).Render(c.String())
	b.WriteString(indent(card, 2))
	b.WriteString("\n\n")

	b.WriteString("  " + helpBar("↑↓", "이동", "space", "선택", "a", "전체", "enter", "설치", "esc", "뒤로"))
	b.WriteString("\n")
}

// ── Install view ────────────────────────────────────────────────────

func (m model) viewInstall(b *strings.Builder) {
	b.WriteString("  " + sectionStyle.Render("프로그램을 설치하고 있습니다..."))
	b.WriteString("\n\n")

	var c strings.Builder
	list := m.enabledList()

	for i, s := range list {
		var icon string
		var st lipgloss.Style

		switch s.status {
		case statusPending:
			icon = pendStyle.Render("○")
			st = pendStyle
		case statusRunning:
			icon = m.spinner.View()
			st = runStyle
		case statusDone:
			icon = successStyle.Render("✓")
			st = successStyle
		case statusFailed:
			icon = errorStyle.Render("✗")
			st = errorStyle
		}

		c.WriteString("  " + icon + " " + st.Render(s.step.Name))
		if s.status == statusFailed && s.errMsg != "" {
			c.WriteString("\n    " + dimStyle.Render(s.errMsg))
		}
		if i < len(list)-1 {
			c.WriteString("\n")
		}
	}

	// Animated gradient progress bar + counter
	counter := dimStyle.Render(fmt.Sprintf("  %d / %d", m.doneCount, m.totalCount))
	c.WriteString("\n\n")
	c.WriteString("  " + m.progress.View() + counter)

	// Current log message
	if log := m.currentLog.Load(); log != "" {
		c.WriteString("\n\n")
		c.WriteString("  " + logStyle.Render(log))
	}

	card := cardStyle.
		Width(m.outerWidth()).
		BorderForeground(colorVioletDim).
		Render(c.String())
	b.WriteString(indent(card, 2))
	b.WriteString("\n")
}

// ── Done view ───────────────────────────────────────────────────────

func (m model) viewDone(b *strings.Builder) {
	hasFailure := false
	for _, s := range m.steps {
		if s.enabled && s.status == statusFailed {
			hasFailure = true
			break
		}
	}

	if hasFailure {
		b.WriteString("  " + warnStyle.Render("⚠ 일부 설치에 실패했습니다."))
	} else {
		b.WriteString("  " + successStyle.Render("✓ 모든 프로그램이 설치되었습니다!"))
	}
	b.WriteString("\n\n")

	// Result card
	var c strings.Builder
	list := m.enabledList()

	for i, s := range list {
		var icon, text string

		switch s.status {
		case statusDone:
			icon = successStyle.Render("✓")
			text = normalStyle.Render(s.step.Name)
		case statusFailed:
			icon = errorStyle.Render("✗")
			text = errorStyle.Render(s.step.Name)
		default:
			icon = faintStyle.Render("–")
			text = faintStyle.Render(s.step.Name + " (건너뜀)")
		}

		c.WriteString("  " + icon + " " + text)
		if s.status == statusFailed && s.errMsg != "" {
			c.WriteString("\n    " + dimStyle.Render(s.errMsg))
		}
		if i < len(list)-1 {
			c.WriteString("\n")
		}
	}

	borderColor := colorGreenDim
	if hasFailure {
		borderColor = colorAmberDim
	}

	card := cardStyle.
		Width(m.outerWidth()).
		BorderForeground(borderColor).
		Render(c.String())
	b.WriteString(indent(card, 2))
	b.WriteString("\n\n")

	b.WriteString("  " + dimStyle.Render("Made with ♥ by switchmd"))
	b.WriteString("\n\n")

	b.WriteString("  " + helpBar("아무 키", "종료"))
	b.WriteString("\n")
}
