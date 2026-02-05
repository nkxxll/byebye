package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	termWidth, termHeight int
)

const (
	Suspend   = iota
	Lock      = iota
	Logout    = iota
	Sleep     = iota
	Shutdown  = iota
	Restart   = iota
	Hibernate = iota
)

type model struct {
	cursor   int
	choices  []string
	selected map[int]struct{}
}

func end(choice string) {
	action := strings.ToLower(choice)
	if err := executeAction(action); err != nil {
		log.Printf("action %s failed: %v", choice, err)
	}
}

func executeAction(action string) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	env := detectEnvironment()
	wm := detectWindowManager(env)
	commands := cfg.getCommands(wm, action, env.Display)
	if len(commands) == 0 {
		return fmt.Errorf("no commands configured for %s", action)
	}

	return runCommands(action, commands)
}

func runCommands(action string, commands []string) error {
	var lastErr error
	for _, command := range commands {
		if strings.TrimSpace(command) == "" {
			continue
		}
		err := runShell(command)
		if err != nil {
			lastErr = err
			continue
		}
		lastErr = nil
		if shouldStopAfterSuccess(action, command) {
			return nil
		}
	}
	return lastErr
}

func runShell(command string) error {
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func shouldStopAfterSuccess(action, command string) bool {
	trimmed := strings.TrimSpace(command)
	if strings.HasSuffix(trimmed, "&") {
		return false
	}
	if action == "suspend" || action == "hibernate" {
		return true
	}
	return true
}

func initialModel() model {
	return model{
		choices: []string{"Suspend", "Lock", "Logout", "Sleep", "Shutdown", "Restart", "Hibernate"},

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(nil, tea.EnterAltScreen)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		termWidth, termHeight = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case strconv.Itoa(Suspend + 1):
			end("Suspend")
			return m, tea.Quit
		case strconv.Itoa(Lock + 1):
			end("Lock")
			return m, tea.Quit
		case strconv.Itoa(Logout + 1):
			end("Logout")
			return m, tea.Quit
		case strconv.Itoa(Sleep + 1):
			end("Sleep")
			return m, tea.Quit
		case strconv.Itoa(Shutdown + 1):
			end("Shutdown")
			return m, tea.Quit
		case strconv.Itoa(Restart + 1):
			end("Restart")
			return m, tea.Quit
		case strconv.Itoa(Hibernate + 1):
			end("Hibernate")
			return m, tea.Quit
		case "enter":
			end(m.choices[m.cursor])
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	// Use terminal ANSI colors for better theme compatibility
	dimColor := lipgloss.ANSIColor(8)       // bright black / dim
	accentColor := lipgloss.ANSIColor(4)    // blue
	highlightColor := lipgloss.ANSIColor(6) // cyan

	// Responsive: use 2 columns if terminal is small (< 40 height or < 60 width)
	useColumns := termHeight < 40 || termWidth < 60

	// Styles
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(accentColor).
		MarginBottom(1)

	subtitleStyle := lipgloss.NewStyle().
		Foreground(dimColor).
		Italic(true)

	itemWidth := 16
	if useColumns && termWidth < 50 {
		itemWidth = 12
	}

	selectedStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(highlightColor).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(highlightColor).
		Width(itemWidth).
		Align(lipgloss.Center).
		Padding(0, 1)

	normalStyle := lipgloss.NewStyle().
		Foreground(dimColor).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(dimColor).
		Width(itemWidth).
		Align(lipgloss.Center).
		Padding(0, 1)

	// Build option items
	var items []string
	for i, choice := range m.choices {
		label := fmt.Sprintf("%d %s", i+1, choice)
		if m.cursor == i {
			items = append(items, selectedStyle.Render(label))
		} else {
			items = append(items, normalStyle.Render(label))
		}
	}

	// Layout options
	var optionsBlock string
	if useColumns {
		// Two-column layout
		var rows []string
		for i := 0; i < len(items); i += 2 {
			if i+1 < len(items) {
				row := lipgloss.JoinHorizontal(lipgloss.Top, items[i], "  ", items[i+1])
				rows = append(rows, row)
			} else {
				rows = append(rows, items[i])
			}
		}
		optionsBlock = lipgloss.JoinVertical(lipgloss.Center, rows...)
	} else {
		// Single-column layout
		optionsBlock = lipgloss.JoinVertical(lipgloss.Center, items...)
	}

	// Compose final view
	title := titleStyle.Render("Where do you want to GO?")
	hint := subtitleStyle.Render("↑/↓ navigate • enter select • q quit")

	content := lipgloss.JoinVertical(lipgloss.Center, title, "", optionsBlock, "", hint)

	// Center in terminal
	contentWidth, contentHeight := lipgloss.Size(content)
	marginH := max((termHeight-contentHeight)/2, 0)
	marginW := max((termWidth-contentWidth)/2, 0)

	return lipgloss.NewStyle().Margin(marginH, marginW).Render(content)
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
