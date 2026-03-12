package tui

import "charm.land/lipgloss/v2"

// ── Color Palette ───────────────────────────────────────────────────
//
// Tailwind-inspired slate/violet/cyan palette for a modern feel.

var (
	// Primary accent — violet
	colorViolet    = lipgloss.Color("#a78bfa")
	colorVioletDim = lipgloss.Color("#7c3aed")

	// Secondary accent — cyan
	colorCyan = lipgloss.Color("#22d3ee")

	// Semantic
	colorGreen    = lipgloss.Color("#34d399")
	colorGreenDim = lipgloss.Color("#059669")
	colorRed      = lipgloss.Color("#fb7185")
	colorYellow   = lipgloss.Color("#fbbf24")
	colorAmberDim = lipgloss.Color("#92400e")

	// Neutral scale (slate)
	colorWhite     = lipgloss.Color("#f1f5f9")
	colorTextSub   = lipgloss.Color("#94a3b8")
	colorTextMuted = lipgloss.Color("#64748b")
	colorTextFaint = lipgloss.Color("#475569")
	colorBorder    = lipgloss.Color("#334155")
)

// ── Style Definitions ───────────────────────────────────────────────

var (
	// ── Header ──────────────────────────────────────────────────
	logoStyle = lipgloss.NewStyle().
			Foreground(colorViolet).
			Bold(true)

	// ── Breadcrumb stepper ──────────────────────────────────────
	bcActiveStyle  = lipgloss.NewStyle().Foreground(colorCyan).Bold(true)
	bcDoneStyle    = lipgloss.NewStyle().Foreground(colorGreen)
	bcPendingStyle = lipgloss.NewStyle().Foreground(colorTextFaint)
	bcSepStyle     = lipgloss.NewStyle().Foreground(colorTextFaint)

	// ── Section heading ─────────────────────────────────────────
	sectionStyle = lipgloss.NewStyle().
			Foreground(colorWhite).
			Bold(true)

	// ── Card container ──────────────────────────────────────────
	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Padding(1, 2)

	// ── List items ──────────────────────────────────────────────
	selectedStyle = lipgloss.NewStyle().Foreground(colorCyan).Bold(true)
	normalStyle   = lipgloss.NewStyle().Foreground(colorTextSub)
	descStyle     = lipgloss.NewStyle().Foreground(colorTextMuted)

	// ── Cursor ──────────────────────────────────────────────────
	cursorStyle = lipgloss.NewStyle().Foreground(colorCyan).Bold(true)

	// ── Checkboxes ──────────────────────────────────────────────
	checkOnStyle  = lipgloss.NewStyle().Foreground(colorCyan)
	checkOffStyle = lipgloss.NewStyle().Foreground(colorTextFaint)

	// ── Step status indicators ──────────────────────────────────
	successStyle = lipgloss.NewStyle().Foreground(colorGreen)
	errorStyle   = lipgloss.NewStyle().Foreground(colorRed)
	warnStyle    = lipgloss.NewStyle().Foreground(colorYellow)
	runStyle     = lipgloss.NewStyle().Foreground(colorViolet)
	pendStyle    = lipgloss.NewStyle().Foreground(colorTextFaint)

	// ── Log message ─────────────────────────────────────────────
	logStyle = lipgloss.NewStyle().
			Foreground(colorTextMuted).
			Italic(true)

	// ── Help bar ────────────────────────────────────────────────
	helpKeyStyle  = lipgloss.NewStyle().Foreground(colorTextSub).Bold(true)
	helpDescStyle = lipgloss.NewStyle().Foreground(colorTextMuted)

	// ── Utility ─────────────────────────────────────────────────
	dimStyle   = lipgloss.NewStyle().Foreground(colorTextMuted)
	faintStyle = lipgloss.NewStyle().Foreground(colorTextFaint)
)
