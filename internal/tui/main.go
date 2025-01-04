package tui

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const useHighPerformanceRenderer = false

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

func init() {
	now := time.Now()
	logDir := "logs/"
	logFile := fmt.Sprint(now.Format("2006-01-02"), ".log")
	logFilePath := fmt.Sprint(logDir, logFile)

	_, err := os.Stat(logDir)
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.Mkdir(logDir, 0755)
	}

	f, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	// defer f.Close()

	// log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(f)
	log.SetLevel(log.DebugLevel)
}

func MainView() {
	log.Info("Starting the application...")

	p := tea.NewProgram(model{ready: false}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("There's been an error: %v", err)
	}
}

// OpenApi Structs
type Collection struct {
	BasePath string
	Host     string
	Schemes  []string
	Protocol string

	Tags      []Tag
	Endpoints []Endpoint
}

type Tag struct {
	Name        string
	Description string
}

type HTTPMethod int32

const (
	GET HTTPMethod = iota
	POST
	PUT
	DELETE
	PATCH
	HEAD
	OPTIONS
	CONNECT
	TRACE
)

var httpMethodName = map[HTTPMethod]string{
	GET:     "get",
	POST:    "post",
	PUT:     "put",
	DELETE:  "delete",
	PATCH:   "patch",
	HEAD:    "head",
	OPTIONS: "options",
	CONNECT: "connect",
	TRACE:   "trace",
}

func (h HTTPMethod) String() string {
	return httpMethodName[h]
}

type Endpoint struct {
	Path   string
	Method HTTPMethod
}

type custResp struct {
	resp     *http.Response
	Body     []byte
	Duration time.Duration
}

type model struct {
	ready            bool
	list             list.Model
	listKeys         *listKeyMap
	listDelegateKeys *delegateKeyMap
	url              textinput.Model
	viewport         viewport.Model
	resp             custResp
}

type respMsg custResp
type errorMsg error

type item struct {
	title       string
	description string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.description }
func (i item) FilterValue() string { return i.title }

type listKeyMap struct {
	toggleSpinner    key.Binding
	toggleTitleBar   key.Binding
	toggleStatusBar  key.Binding
	togglePagination key.Binding
	toggleHelpMenu   key.Binding
	insertItem       key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		insertItem: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add item"),
		),
		toggleSpinner: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "toggle spinner"),
		),
		toggleTitleBar: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "toggle title"),
		),
		toggleStatusBar: key.NewBinding(
			key.WithKeys("S"),
			key.WithHelp("S", "toggle status"),
		),
		togglePagination: key.NewBinding(
			key.WithKeys("P"),
			key.WithHelp("P", "toggle pagination"),
		),
		toggleHelpMenu: key.NewBinding(
			key.WithKeys("H"),
			key.WithHelp("H", "toggle help"),
		),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) headerView() string {
	title := titleStyle.Render("Response Body:")
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m model) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	if m.url.Focused() {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				return m, tea.Quit
			case "enter":
				return m, sendRequest(m.url.Value())
			}
			m.url, cmd = m.url.Update(msg)
			return m, cmd
		}
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			var (
				delegateKeys = newDelegateKeyMap()
				listKeys     = newListKeyMap()
			)

			// Make initial list of items
			items := []list.Item{
				item{
					title:       "test 1",
					description: "test 1 description",
				},
				item{
					title:       "test 2",
					description: "test 2 description",
				},
				item{
					title:       "test 3",
					description: "test 3 description",
				},
			}

			// Setup list
			delegate := newItemDelegate(delegateKeys)
			m.list = list.New(items, delegate, msg.Width/3, msg.Height)
			m.list.Title = "Groceries"
			m.list.Styles.Title = titleStyle
			m.list.AdditionalFullHelpKeys = func() []key.Binding {
				return []key.Binding{
					listKeys.toggleSpinner,
					listKeys.insertItem,
					listKeys.toggleTitleBar,
					listKeys.toggleStatusBar,
					listKeys.togglePagination,
					listKeys.toggleHelpMenu,
				}
			}
			m.listKeys = listKeys
			m.listDelegateKeys = delegateKeys

			m.url = textinput.New()
			m.url.Placeholder = "https://httpbin.org/anything"
			// m.url.SetValue("https://httpbin.org/anything")
			// m.url.Focus()
			m.url.CharLimit = 2048
			m.url.Width = 60

			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight-lipgloss.Height(m.url.View())-lipgloss.Height("Placeholder"))
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			m.viewport.SetContent("No Data")
			m.ready = true

			// This is only necessary for high performance rendering, which in
			// most cases you won't need.
			//
			// Render the viewport one line below the header.
			m.viewport.YPosition = headerHeight + 1
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

		if useHighPerformanceRenderer {
			// Render (or re-render) the whole viewport. Necessary both to
			// initialize the viewport and when the window is resized.
			//
			// This is needed for high-performance rendering only.
			cmds = append(cmds, viewport.Sync(m.viewport))
		}
	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, m.listKeys.toggleSpinner):
			cmd := m.list.ToggleSpinner()
			return m, cmd

		case key.Matches(msg, m.listKeys.toggleTitleBar):
			v := !m.list.ShowTitle()
			m.list.SetShowTitle(v)
			m.list.SetShowFilter(v)
			m.list.SetFilteringEnabled(v)
			return m, nil

		case key.Matches(msg, m.listKeys.toggleStatusBar):
			m.list.SetShowStatusBar(!m.list.ShowStatusBar())
			return m, nil

		case key.Matches(msg, m.listKeys.togglePagination):
			m.list.SetShowPagination(!m.list.ShowPagination())
			return m, nil

		case key.Matches(msg, m.listKeys.toggleHelpMenu):
			m.list.SetShowHelp(!m.list.ShowHelp())
			return m, nil

		case key.Matches(msg, m.listKeys.insertItem):
			m.listDelegateKeys.remove.SetEnabled(true)
			newId := len(m.list.Items()) + 1
			newItem := item{
				title:       fmt.Sprintf("test %d", newId),
				description: fmt.Sprintf("test %d description", newId),
			}
			insCmd := m.list.InsertItem(0, newItem)
			statusCmd := m.list.NewStatusMessage(statusMessageStyle("Added " + newItem.Title()))
			return m, tea.Batch(insCmd, statusCmd)
		}
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+z":
			return m, tea.Suspend
		case "u":
			m.url.Focus()
			return m, nil
		case "s":
			url := m.url.Value()
			if url == "" {
				url = m.url.Placeholder
			}
			return m, sendRequest(url)
		}
	case errorMsg:
		m.viewport.SetContent(msg.Error())
		return m, nil
	case respMsg:
		m.resp = custResp(msg)

		bodyStr := fmt.Sprint(m.resp.Duration, "\n") + string(m.resp.Body)
		log.Debug("bodyStr: ", bodyStr)

		m.viewport.SetContent(bodyStr)

		return m, nil
	}

	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var style = lipgloss.NewStyle()

	reqView := lipgloss.JoinVertical(
		lipgloss.Top,
		"Request URL:",
		m.url.View(),
		m.headerView(),
		m.viewport.View(),
		m.footerView(),
	)
	retView := lipgloss.JoinHorizontal(lipgloss.Left, m.list.View(), reqView)

	return style.Render(retView)
}

func sendRequest(url string) tea.Cmd {
	timeBeforeReq := time.Now()
	log.Debug("Sending GET request to:", url)
	resp, err := http.Get(url)
	timeAfterReq := time.Now()
	if err != nil {
		log.Error("Error sending GET request:", err)
		return func() tea.Msg {
			return errorMsg(err)
		}
	}
	defer resp.Body.Close()

	requestDuration := timeAfterReq.Sub(timeBeforeReq)
	log.Debug("Request duration:", requestDuration)

	log.Debug("Response status:", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Error sending GET request:", err)
		return func() tea.Msg {
			return errorMsg(err)
		}
	}

	return func() tea.Msg {
		return respMsg(custResp{resp: resp, Body: body, Duration: requestDuration})
	}
}
