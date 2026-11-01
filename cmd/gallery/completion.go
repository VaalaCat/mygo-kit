package main

import (
	"strings"

	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

var (
	completionText  = "/re"
	completionState = 0
	completionPick  = "(none)"
)

var completionSkills = []kit.CompletionItem{
	{Icon: kit.IconBookOpen, Title: "/review", Description: "Review the pending changes and list the problems found"},
	{Icon: kit.IconBookOpen, Title: "/refactor", Description: "Restructure code without changing its behavior"},
	{Icon: kit.IconBookOpen, Title: "/release-notes", Description: "Draft release notes from the commit history"},
	{Icon: kit.IconBookOpen, Title: "/rename", Description: "Rename a symbol everywhere"},
	{Icon: kit.IconBookOpen, Title: "/research", Description: "Read the code and answer a question"},
	{Icon: kit.IconBookOpen, Title: "/rebase", Description: "Rebase the branch"},
	{Icon: kit.IconBookOpen, Title: "/retry", Description: "Retry the last tool call"},
	{Icon: kit.IconBookOpen, Title: "/revert", Description: "Revert a commit"},
	{Icon: kit.IconBookOpen, Title: "/rewrite", Description: "Rewrite the selected text"},
	{Icon: kit.IconBookOpen, Title: "/restore", Description: "Restore a deleted file"},
}

func init() {
	register(page{name: "CompletionMenu", order: 55, view: func(c *ui.Context) {
		theme := kit.GetTheme(c)
		completionPickLabel := completionPick
		if completionPickLabel == "(none)" {
			completionPickLabel = t(c, "(none)")
		}
		ui.Column(c).Grow(1).Padding(24).Gap(16).Children(func() {
			kit.Segmented(c.Key("state"), &completionState, []string{t(c, "Items"), t(c, "Loading"), t(c, "Empty"), t(c, "Error")}).Label(t(c, "State")).AlignSelf(ui.Start)
			ui.Text(c, t(c, "Last chosen: ")+completionPickLabel+t(c, "   (type / to open; Up/Down, Enter or Tab, Esc)")).TextColor(theme.Colors.Muted)
			ui.Spacer(c)
			// 1. 输入框在页面底部，补全浮层浮在它上方，输入区高度不变。
			box := ui.Column(c.Key("composer")).Padding(theme.Space(2)).Radius(theme.Metrics.RadiusLarge).Border(1, theme.Colors.Border).Background(theme.Colors.Editor)
			var input ui.Element
			box.Children(func() {
				input = kit.TextAreaWith(c.Key("composer-input"), &completionText, kit.InputOpt{Borderless: true, Placeholder: t(c, "Message, / for skills")}).Label(t(c, "Message")).MinHeight(theme.Space(12))
			})
			opt := kit.CompletionOpt{Open: strings.HasPrefix(completionText, "/"), Empty: t(c, "No matching skills")}
			switch completionState {
			case 0:
				query := strings.TrimPrefix(completionText, "/")
				for _, item := range completionSkills {
					if strings.Contains(item.Title, query) {
						item.Description = t(c, item.Description)
						opt.Items = append(opt.Items, item)
					}
				}
			case 1:
				opt.Loading = true
			case 3:
				opt.Error = t(c, "Could not load skills: connection refused")
			}
			if i := kit.CompletionMenu(c, box, input, opt); i >= 0 {
				completionPick = opt.Items[i].Title
				completionText = opt.Items[i].Title + " "
			}
		})
	}})
}
