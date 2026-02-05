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
	termWidth, termHight int
	heading              = lipgloss.NewStyle().Bold(true).Margin(1, 0)
	notChoosen           = lipgloss.NewStyle().Bold(true).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("12")).Width(14).Padding(1)
	choosen              = lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Bold(true).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("5")).Width(14).Padding(1)
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
		termWidth, termHight = msg.Width, msg.Height
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
	s := heading.Render("Where do you want to GO?")

	for i, choice := range m.choices {
		cursor := " "
		var line string
		if m.cursor == i {
			cursor = "❯"
			line = choosen.Render(fmt.Sprintf("%s %s", cursor, choice))
		} else {
			line = notChoosen.Render(fmt.Sprintf("%s %s", cursor, choice))
		}
		s = lipgloss.JoinVertical(lipgloss.Center, s, line)
	}

	s = lipgloss.JoinVertical(lipgloss.Center, s, heading.Render("Press q to quit."))

	textWidth, textHeight := lipgloss.Size(s)
	marginW, marginH := (termWidth-textWidth)/2, (termHight-textHeight)/2

	return lipgloss.NewStyle().Margin(marginH, marginW).Render(s)
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
