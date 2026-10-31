package kit

import (
	"embed"

	"github.com/egoist/mygo/ui"
)

//go:embed icons/*.svg
var iconFiles embed.FS

// Lucide 线性图标，统一 1.5 描边，颜色跟随文字。
var (
	IconPanelLeft     = mustIcon("panel-left")
	IconPanelRight    = mustIcon("panel-right")
	IconPanelBottom   = mustIcon("panel-bottom")
	IconPlus          = mustIcon("plus")
	IconSettings      = mustIcon("settings")
	IconSearch        = mustIcon("search")
	IconChevronDown   = mustIcon("chevron-down")
	IconChevronLeft   = mustIcon("chevron-left")
	IconChevronRight  = mustIcon("chevron-right")
	IconX             = mustIcon("x")
	IconEllipsis      = mustIcon("ellipsis")
	IconArrowUp       = mustIcon("arrow-up")
	IconArrowDown     = mustIcon("arrow-down")
	IconStop          = mustIcon("square")
	IconImagePlus     = mustIcon("image-plus")
	IconFolder        = mustIcon("folder")
	IconFold          = mustIcon("fold-vertical")
	IconCopy          = mustIcon("copy")
	IconAlert         = mustIcon("circle-alert")
	IconLoader        = mustIcon("loader-circle")
	IconBot           = mustIcon("bot")
	IconBranch        = mustIcon("git-branch")
	IconTerminal      = mustIcon("terminal")
	IconInfo          = mustIcon("info")
	IconCircleCheck   = mustIcon("circle-check")
	IconCircleX       = mustIcon("circle-x")
	IconCommand       = mustIcon("command")
	IconFiles         = mustIcon("files")
	IconBrain         = mustIcon("brain")
	IconTriangleAlert = mustIcon("triangle-alert")
	IconFileText      = mustIcon("file-text")
	IconPencil        = mustIcon("pencil")
	IconGauge         = mustIcon("gauge")
	IconListTree      = mustIcon("list-tree")
	IconBug           = mustIcon("bug")
	IconActivity      = mustIcon("activity")
	IconWrench        = mustIcon("wrench")
	IconArchive       = mustIcon("archive")
	IconBookOpen      = mustIcon("book-open")
	IconWebhook       = mustIcon("webhook")
	IconGitFork       = mustIcon("git-fork")
	IconReplace       = mustIcon("replace")
	IconServerCog     = mustIcon("server-cog")
	IconBan           = mustIcon("ban")
	IconMessageSquare = mustIcon("message-square")
	IconBraces        = mustIcon("braces")
	IconRotateCcw     = mustIcon("rotate-ccw")
	IconCirclePause   = mustIcon("circle-pause")
	IconCircleStop    = mustIcon("circle-stop")
	IconRefresh       = mustIcon("refresh-cw")
	IconCheck         = mustIcon("check")
	IconListEnd       = mustIcon("list-end")

	IconPanelBottomClose = mustIcon("panel-bottom-close")
	IconPanelRightClose  = mustIcon("panel-right-close")
	IconPlay             = mustIcon("play")
	IconTrash            = mustIcon("trash-2")
	IconDownload         = mustIcon("download")
)

// mustIcon 读取程序内置的图标文件。
func mustIcon(name string) *ui.SVG {
	data, err := iconFiles.ReadFile("icons/" + name + ".svg")
	if err != nil {
		panic(err)
	}
	return ui.MustParseSVG(data)
}
