package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

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
)

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		// fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	p := tea.NewProgram(model{ready: false}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		// fmt.Printf("There's been an error: %v", err)
		os.Exit(1)
	}
}

type custResp struct {
	resp *http.Response
	Body []byte
}

type model struct {
	ready    bool
	url      textinput.Model
	viewport viewport.Model
	resp     custResp
}

type respMsg custResp
type errorMsg error

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) headerView() string {
	title := titleStyle.Render("Mr. Pager")
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
			m.url = textinput.New()
			m.url.Placeholder = "https://httpbin.org/anything"
			m.url.SetValue("https://httpbin.org/anything")
			// m.url.Focus()
			m.url.CharLimit = 2048
			m.url.Width = 60

			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
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
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+z":
			return m, tea.Suspend
		case "u":
			m.url.Focus()
			return m, nil
		case "s":
			return m, sendRequest(m.url.Value())
		}
	case errorMsg:
		m.viewport.SetContent(msg.Error())
		return m, nil
	case respMsg:
		m.resp = custResp(msg)

		bodyStr := string(m.resp.Body)
		// fmt.Println("bodyStr: ", bodyStr)

		m.viewport.SetContent(bodyStr)

		return m, nil
	}

	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf(
		"Request URL:\n\n%s\n\nResponse:\n\n%s",
		m.url.View(),
		m.viewport.View(),
	)
}

func sendRequest(url string) tea.Cmd {
	// fmt.Println("Sending GET request to:", url)
	resp, err := http.Get(url)
	if err != nil {
		// fmt.Println("Error sending GET request:", err)
		return func() tea.Msg {
			return errorMsg(err)
		}
	}
	defer resp.Body.Close()

	// fmt.Println("Response status:", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// fmt.Println("Error sending GET request:", err)
		return func() tea.Msg {
			return errorMsg(err)
		}
	}

	return func() tea.Msg {
		return respMsg(custResp{resp: resp, Body: body})
	}
}
