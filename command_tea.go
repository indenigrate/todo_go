package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/enescakir/emoji"
)

const listHeight = 11

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

// Define a new enum-like type for the state of the application
type state int

const (
	chooseItemState state = iota
	getInputState
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	list        list.Model
	choice      string
	input       string
	state       state // Tracks the current state (choosing item or inputting data)
	inputPrompt string
	quitting    bool
	todos       *Todos
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		keypress := msg.String()

		switch m.state {
		// When in the item choosing state
		case chooseItemState:
			switch keypress {
			case "q", "ctrl+c":
				m.quitting = true
				return m, tea.Quit

			case "enter":
				// Get the selected item and move to input state
				i, ok := m.list.SelectedItem().(item)
				if ok {
					m.choice = string(i)
					m.state = getInputState
					switch m.choice {
					case "Add":
						m.inputPrompt = fmt.Sprintf("Enter the title of the new task: ")
					case "Delete", "Toggle":
						m.inputPrompt = fmt.Sprintf("Enter the index to %s : ", m.choice)
					case "Edit":
						m.inputPrompt = fmt.Sprintf("Enter the index and the new title asa INDEX:NEW_TITLE :")
					case "List":
						m.quitting = true
						return m, tea.Quit
					}
				}
				return m, nil
			}

		// When in the input state
		case getInputState:
			switch keypress {
			case "backspace":
				m.input = m.input[:len(m.input)-1]
			case "ctrl+c":
				m.quitting = true
				return m, tea.Quit
			case "enter":
				// Call the function with the input and reset
				switch m.choice {
				case "Add":
					m.todos.add(m.input)
				case "Delete":
					num, _ := strconv.Atoi(m.input)
					m.todos.delete(num - 1)
				case "Toggle":
					num, _ := strconv.Atoi(m.input)
					m.todos.toggle(num - 1)
				case "Edit":
					parts := strings.SplitN(m.input, ":", 2)
					if len(parts) != 2 {
						fmt.Println("Error, invalid format for edit. Please use id:new_title")
						os.Exit(1)
					}
					index, err := strconv.Atoi(parts[0])
					if err != nil {
						fmt.Println("Error, invalid format for edit. Please use id:new_title")
						os.Exit(1)
					}
					m.todos.edit(index-1, parts[1])
				}
				// m.quitting = true
				m.input = ""

				m.state = chooseItemState

				// return m, tea.Quit
				return m, nil

			default:
				// Collect input from the user
				m.input += keypress
				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	// If quitting, return quit message
	if m.quitting {
		// return quitTextStyle.Render("Closed TODO!")
		return m.todos.renderTodos()
	}

	// If in the input state, prompt for input
	if m.state == getInputState {
		return fmt.Sprintf("%s%s", m.inputPrompt, m.input)
	}

	// If in the list selection state, show the list
	return "\n" + m.list.View()

}

func run_command(todos *Todos) {
	items := []list.Item{
		item("List"),
		item("Add"),
		item("Toggle"),
		item("Delete"),
		item("Edit"),
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "What do you want to do?"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := model{list: l, state: chooseItemState, todos: todos}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

// In the Todos struct, add a method to render todos as a string
func (todos *Todos) renderTodos() string {
	var output strings.Builder
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")) // Magenta color
	output.WriteString(headerStyle.Render(fmt.Sprintf("%-5s %-40s %-15s %-35s %-40s", "#", "Title", "Completed", "         Created At", "      Completed At")))
	output.WriteString(fmt.Sprintln(""))
	output.WriteString(strings.Repeat("-", 140) + "\n")

	for index, todo := range *todos {
		completed := "   " + emoji.CrossMark.String()
		completedAt := ""
		createdAt := todo.CreatedAt.Format(time.RFC1123)
		if todo.Completed {
			completed = "   " + emoji.CheckMark.String()
			if todo.CompletedAt != nil {
				completedAt = todo.CompletedAt.Format(time.RFC1123)
			}
			createdAt = "" + createdAt
		}

		if len(todo.Title) <= 40 {
			output.WriteString(fmt.Sprintf("%-5d %-40s %-15s %-35s %-40s\n", index+1, todo.Title, completed, createdAt, completedAt))
		} else {
			//split string
			s := splitInParts(todo.Title)
			output.WriteString(fmt.Sprintf("%-5d %-40s %-15s %-35s %-40s\n", index+1, s[0], completed, createdAt, completedAt))
			for i := 1; i < len(s); i++ {
				// output.WriteString(fmt.Sprintf("%-5s %-40s %-15s %-35s %-40s\n", "", s[i], "", "", ""))
				output.WriteString(fmt.Sprintf("%-5s %-40s\n", "  ", s[i]))
			}
		}
	}

	return output.String()
}

func splitInParts(s string) []string {
	n := len(s) / 40
	ans := make([]string, 0)
	for i := 0; i < n; i++ {
		ans = append(ans, s[i*40:40*(i+1)])

	}
	ans = append(ans, s[40*n:])
	return ans
}
