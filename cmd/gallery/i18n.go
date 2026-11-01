package main

import "github.com/egoist/mygo/ui"

// t 返回当前界面语言下的示例文字：中文界面查 zhTexts，查不到时显示英文原文。
// 画廊的示例文字都用英文原文作参数，例如 t(c, "Actions")。
func t(c *ui.Context, english string) string {
	// 1. 只有中文界面查表。
	if c.Locale() == "zh-Hans" {
		if chinese, ok := zhTexts[english]; ok {
			return chinese
		}
	}
	return english
}

// zhTexts 按英文原文保存画廊示例文字的简体中文译文。
var zhTexts = map[string]string{
	// Overview
	"A focused native workspace": "专注的原生工作区",
	"Continuous panels, compact toolbars and a clear type hierarchy.\nReusable Go controls, without another application runtime.": "连贯的面板、紧凑的工具栏与清晰的文字层级。\n可复用的 Go 控件，无需额外的应用运行时。",
	"Preview":            "预览",
	"Source":             "源码",
	"Details":            "详细信息",
	"Open the workspace": "打开工作区",
	"Documentation":      "文档",

	// Button
	"Actions":                               "操作",
	"Ghost":                                 "幽灵按钮",
	"Primary":                               "主要操作",
	"Delete":                                "删除",
	"Disabled":                              "已禁用",
	"Icon and tool buttons":                 "图标按钮和工具按钮",
	"Settings":                              "设置",
	"Add":                                   "添加",
	"Workspace":                             "工作区",
	"IconButtonTip (hover for the tooltip)": "IconButtonTip（悬停显示提示）",
	"Copy":                                  "复制",
	"Refresh":                               "刷新",
	"Compact icon button (20)":              "紧凑图标按钮（20）",
	"Close":                                 "关闭",
	"Flush buttons align with the text line above": "无内边距按钮与上方文字对齐",
	"Provider list":      "Provider 列表",
	"Link":               "链接",
	"MyGo documentation": "MyGo 文档",

	// Input
	"TextInput":               "单行文本输入",
	"Name":                    "名称",
	"A project name":          "项目名称",
	"TextArea":                "多行文本输入",
	"Description":             "描述",
	"Multi-line text and IME": "多行文本与输入法",
	"Invalid":                 "无效输入",
	"Invalid name":            "无效名称",
	"ReadOnly (a plain selectable value, aligned with editable fields)": "ReadOnly（可选中的纯文本值，与可编辑字段对齐）",
	"Provider ID":                     "Provider ID",
	"Name (editable, for comparison)": "名称（可编辑，用于对比）",
	"Leading icon":                    "前置图标",
	"Search":                          "搜索",
	"Leading icon with long text (the text is clipped beside the icon)": "前置图标与长文本（文本在图标旁裁剪）",
	"Directory": "目录",
	"Borderless (the surrounding box decides the look)": "无边框（外观由外层容器决定）",
	"Message": "消息",
	"Filter":  "筛选",

	// Select
	"System":                         "跟随系统",
	"Toolbar select":                 "工具栏选择器",
	"Appearance":                     "外观",
	"Bordered select":                "带边框的选择器",
	"Theme":                          "主题",
	"Descriptions and group headers": "描述与分组标题",
	"Local":                          "本地",
	"Model":                          "模型",
	"Custom trigger text (options read just the level)": "自定义触发器文本（选项只显示级别）",
	"Auto":        "自动",
	"Low":         "低",
	"Medium":      "中",
	"High":        "高",
	"Reasoning: ": "推理：",
	"Reasoning":   "推理",

	// Checkbox
	"Auto-save": "自动保存",

	// Tabs
	"General":           "常规",
	"Editor":            "编辑器",
	"Advanced":          "高级",
	"General settings":  "常规设置",
	"Editor settings":   "编辑器设置",
	"Advanced settings": "高级设置",

	// Disclosure
	"More settings": "更多设置",
	"Native keyboard behavior, independently styled.": "原生键盘交互，独立外观样式。",

	// ListRow
	"Single line · icon · selection":      "单行 · 图标 · 选择",
	"Refactor the parser":                 "重构解析器",
	"Fix flaky test":                      "修复不稳定测试",
	"Explain this error":                  "解释此错误",
	"Two lines (44) · trailing":           "双行（44）· 尾部内容",
	"Add retry to the client":             "为客户端添加重试",
	"openvaala · 3 minutes ago":           "openvaala · 3 分钟前",
	"Update the changelog":                "更新变更日志",
	"docs · yesterday":                    "文档 · 昨天",
	"Unread":                              "未读",
	"Draft release notes":                 "起草发布说明",
	"docs · 2 days ago":                   "文档 · 2 天前",
	"Depth (16 per level)":                "层级（每级缩进 16）",
	"Actions on hover / selected / focus": "悬停、选中或聚焦时显示操作",
	"Hover this row":                      "悬停此行",
	"Or select it":                        "或选中此行",
	"Tab into it":                         "按 Tab 聚焦此行",
	"local":                               "本地",
	"More":                                "更多",
	"Trailing text on the title baseline · disabled": "尾部文字与标题基线对齐 · 禁用状态",
	"running":     "运行中",
	"Unavailable": "不可用",

	// Dialog
	"Dialog":        "对话框",
	"Open dialog":   "打开对话框",
	"Native dialog": "原生对话框",
	"Focus stays inside this dialog. Escape returns it to the trigger.": "焦点保持在对话框内。按 Escape 将焦点返回触发按钮。",
	"Done": "完成",
	"DialogWith (cancel + confirm, Enter confirms, Esc cancels)": "DialogWith（取消与确认，Enter 确认，Esc 取消）",
	"Confirm":          "确认",
	"Danger":           "危险",
	"With input":       "带输入框",
	"With tabs":        "带选项卡",
	"Read-only values": "只读值",
	"Last result: ":    "上次结果：",
	"Compact context":  "压缩上下文",
	"Older turns are summarized to free up the context window.": "汇总较早的对话轮次，释放上下文空间。",
	"Compact":                "压缩",
	"Delete conversation":    "删除对话",
	"This cannot be undone.": "此操作无法撤销。",
	"A title is required.":   "请输入标题。",
	"Rename":                 "重命名",
	"Save":                   "保存",
	"Title":                  "标题",
	"Provider details":       "Provider 详细信息",
	"Popover and tooltip":    "弹出浮层与工具提示",
	"Popover":                "弹出浮层",
	"Anchored content":       "锚定的内容",
	"Hover for tooltip":      "悬停显示提示",
	"A short hint":           "简短提示",
	"(none)":                 "（无）",
	"compact":                "压缩",
	"delete":                 "删除",
	"rename: ":               "重命名：",

	// Kbd
	"Open the command palette": "打开命令面板",
	"New conversation":         "新建对话",

	// Badge
	"Neutral":    "中性",
	"Accent":     "强调",
	"Success":    "成功",
	"Warning":    "警告",
	"Connected":  "已连接",
	"Connecting": "连接中",
	"Failed":     "失败",

	// Spinner
	"Thinking…": "思考中…",

	// Progress
	"ProgressBar (4 high)": "ProgressBar（高度 4）",
	"Context usage":        "上下文用量",
	"Context window":       "上下文窗口",

	// StatusBar
	"Toggle sidebar": "切换侧边栏",
	"Terminal":       "终端",
	"Outline":        "大纲",

	// TitleBar
	"Window content": "窗口内容",
	"The title bar drags the window and keeps clear of system buttons.": "通过标题栏拖动窗口，内容避开系统按钮。",

	// Sheet
	"Open left sheet":                     "打开左侧抽屉",
	"Open right sheet":                    "打开右侧抽屉",
	"Conversations":                       "对话",
	"Click the backdrop or press Escape.": "点击遮罩或按 Escape 关闭。",
	"A drawer from the right edge.":       "从右侧滑出的抽屉。",

	// Toast
	"Info":                            "信息",
	"Conversation exported":           "对话已导出",
	"Success with action":             "带操作的成功通知",
	"Conversation archived":           "对话已归档",
	"Undo":                            "撤销",
	"Error":                           "错误",
	"Connection lost":                 "连接已断开",
	"The runtime stopped responding.": "运行时已停止响应。",

	// CommandPalette
	"Open command palette": "打开命令面板",
	"Last command: ":       "上次命令：",
	"Conversation":         "对话",
	"Archive conversation": "归档对话",
	"View":                 "视图",
	"Toggle terminal":      "切换终端",
	"Open settings":        "打开设置",
	"Application":          "应用",

	// Workspace
	"Back to gallery": "返回组件库",
	"CONVERSATIONS":   "对话列表",
	"Drag the dividers; narrow the window below 900 to get sheets.": "拖动分隔线；窗口宽度小于 900 时显示抽屉。",
	"Files touched by the agent":                                    "智能体修改的文件",
	"Model, usage and context":                                      "模型、用量与上下文",

	// Images
	"Grid preview": "网格预览",

	// Virtual list
	"Item %05d":       "项目 %05d",
	"Variable height": "可变高度",

	// MessageScroller
	"%d messages · top reached %d times": "%d 条消息 · 已到达顶部 %d 次",
	"Start streaming":                    "开始流式输出",
	"Stop streaming":                     "停止流式输出",

	// ActivityCard
	"Inspect":                            "查看详情",
	" · lines 1-120":                     " · 第 1–120 行",
	"Model call failed":                  "模型调用失败",
	"review the diff of the last commit": "审查上次提交的差异",
	"Turn interrupted":                   "本轮已中断",
	"Thinking":                           "思考中",
	"Reasoning in progress.":             "正在推理。",
	"Model changed from A to B":          "模型已从 A 切换为 B",

	// Terminal
	"%d rows × %d cols · resized %d times": "%d 行 × %d 列 · 已调整大小 %d 次",
	"Run progress":                         "运行进度演示",
	"Mouse reporting: off":                 "鼠标上报：关",
	"Mouse reporting: on":                  "鼠标上报：开",

	// Callout
	"Model request failed": "模型请求失败",
	"The provider returned 429: rate limit exceeded. Retry in a few seconds.": "Provider 返回 429：请求超出速率限制。请稍后重试。",
	"Resend":                 "重新发送",
	"Context is nearly full": "上下文即将用尽",
	"86% of the window is used; older turns may be compacted.": "上下文窗口已使用 86%；较早的对话轮次可能会被压缩。",
	"Conversation forked": "已创建对话分支",
	"Open":                "打开",
	"Dismiss":             "关闭",

	// AttachmentChip
	" (no thumbnail)": "（无缩略图）",

	// Shimmer
	"Generating a title for this conversation": "正在为此对话生成标题",

	// MorphText
	"Reading files":   "正在读取文件",
	"Reading config":  "正在读取配置",
	"Running tests":   "正在运行测试",
	"Running 3 tests": "正在运行 3 项测试",
	"Next status":     "下一个状态",

	// Switch
	"With label (the whole row toggles; Space toggles while focused)": "带标签（点击整行切换；聚焦时按 Space 切换）",
	"Reduce motion":                   "减少动态效果",
	"Notify when done":                "完成时通知",
	"Track only (name it with Label)": "仅开关（用 Label 设置名称）",
	"Disabled (framework opacity)":    "已禁用（框架控制透明度）",
	"On":                              "开启",
	"Off":                             "关闭",

	// Segmented
	"Text only, hugging content (arrows move the choice)": "纯文本，宽度适应内容（方向键切换选项）",
	"Light":                            "浅色",
	"Dark":                             "深色",
	"Two segments, second selected":    "两个选项，选中第二项",
	"Archive":                          "归档",
	"With icons and tooltips":          "带图标与工具提示",
	"Executor and directory":           "执行器与目录",
	"Connected MCP servers":            "已连接的 MCP 服务器",
	"Skill":                            "Skill",
	"Skills this conversation can use": "此对话可用的 Skill",
	"Sections":                         "分区",
	"Full width, equal segments":       "占满宽度，等宽选项",
	"Full width":                       "占满宽度",
	"One":                              "一",
	"Two":                              "二",
	"Three":                            "三",

	// Field
	"Field: label above, description or error below (16 between fields)": "字段：标签在上，说明或错误在下（字段间距为 16）",
	"Used in config.yaml and as the model prefix.":                       "用于 config.yaml，并作为模型前缀。",
	"Hidden while there is an error.":                                    "出现错误时隐藏。",
	"Only letters, digits and dashes.":                                   "仅限字母、数字和连字符。",
	"Base URL":                                                           "基础 URL",
	"Protocol":                                                           "协议",
	"Wire format of the provider.":                                       "Provider 的通信格式。",
	"Inherit settings":                                                   "继承设置",
	"Use the default provider's headers.":                                "使用默认 Provider 的请求头。",
	"FieldMessage":                                                       "字段提示",
	"A description in muted text.":                                       "以弱化文字显示的说明。",
	"A warning in the warning color.":                                    "以警告色显示的警告。",
	"An error in the danger color.":                                      "以危险色显示的错误。",
	"InlineHint":                                                         "行内提示",
	"Loading…":                                                           "加载中…",
	"No conversations yet":                                               "暂无对话",

	// Settings form
	"Revert":                               "还原",
	"Apply":                                "应用",
	"System follows the operating system.": "此外观选项会跟随操作系统。",
	"Color theme":                          "配色主题",
	"Turns off transitions and spinners.":  "关闭过渡动画和加载动画。",
	"Notify when a turn finishes":          "对话轮次结束时通知",
	"Notify":                               "通知",
	"Model provider":                       "模型 Provider",
	"Include the /v1 suffix.":              "请包含 /v1 后缀。",

	// Menu
	"Fork":    "创建分支",
	"Export":  "导出",
	"Remote":  "远程",
	"Sandbox": "沙箱",
	"MenuButton (Up/Down, Enter, Esc, first letter)": "菜单按钮（Up/Down、Enter、Esc、首字母）",
	"Last chosen: ": "上次选择：",
	"ContextMenu (right click, Shift+F10 or the Menu key while focused)": "右键菜单（右键点击，或聚焦时按 Shift+F10 或 Menu 键）",
	"Right click here": "在此右键点击",
	"ContextMenu on Selectable text (drag to select, then right click: Copy is built in)": "可选择文本上的右键菜单（拖动选择后右键点击：内置复制操作）",
	"Select any part of this text, then right click.":                                     "选择这段文字中的任意内容，然后右键点击。",
	"The menu is kit's own, not the system's.":                                            "此菜单由组件库提供，并非系统菜单。",

	// CompletionMenu
	"Items":   "选项",
	"Loading": "加载中",
	"Empty":   "空",
	"State":   "状态",
	"   (type / to open; Up/Down, Enter or Tab, Esc)":        "   （输入 / 打开；Up/Down、Enter 或 Tab、Esc）",
	"Message, / for skills":                                  "输入消息，输入 / 选择 Skill",
	"No matching skills":                                     "没有匹配的 Skill",
	"Review the pending changes and list the problems found": "审查待提交的改动并列出发现的问题",
	"Restructure code without changing its behavior":         "重构代码并保持行为不变",
	"Draft release notes from the commit history":            "根据提交历史起草发布说明",
	"Rename a symbol everywhere":                             "重命名符号及其所有引用",
	"Read the code and answer a question":                    "阅读代码并回答问题",
	"Rebase the branch":                                      "变基当前分支",
	"Retry the last tool call":                               "重试上次工具调用",
	"Revert a commit":                                        "还原一次提交",
	"Rewrite the selected text":                              "重写所选文本",
	"Restore a deleted file":                                 "恢复已删除的文件",
	"Could not load skills: connection refused":              "无法加载 Skill：连接被拒绝",

	// Gallery
	"mygo-kit  /  Native workspace": "mygo-kit  /  原生工作区",
	"◐ Theme":                       "◐ 主题",
	"COMPONENTS":                    "组件",
	"● Ready":                       "● 就绪",
	"Native UI  ·  ":                "原生 UI  ·  ",
}
