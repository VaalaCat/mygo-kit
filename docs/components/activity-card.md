---
title: ActivityCard
description: A single-line entry header with expandable details and a status icon, for recording tool calls, thinking, and other steps.
---

# ActivityCard

A 28-high, single-line entry header that can expand when details are available. Commonly used to show an Agent's tool calls, thinking, errors, and other steps.

```
[Icon/arrow] Title  Status  Summary ········ [Actions]
```

<Preview page="ActivityCard" :height="460">

<<< @/../cmd/gallery/activitycard.go

</Preview>

## Usage

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
var open bool // The application owns the expanded state

kit.ActivityCard(c, &open, kit.ActivityCardOpt{
	Icon:    kit.IconTerminal,
	Title:   "bash",
	Summary: "go test ./...",
	Status:  kit.ActivityRunning,
}, nil, func() {
	ui.Text(c, "$ go test ./...")
})
```

The last two arguments are the actions on the right and the expanded details, respectively. With `detail` set to `nil`, the row is a plain text row with no expansion or arrow. Keep a separate `open` for each card.

## Examples

### Actions shown on hover

The icon buttons in `actions` are separate from the expand button. They appear only when the row is hovered or has keyboard focus, while their space stays reserved.

```go
inspect := func() { kit.IconButton(c, kit.IconBug, "Inspect") }
kit.ActivityCard(c, &open, opt, inspect, detail)
```

### Status

`Status` shows an icon after the title. An empty value means success or no status and shows no icon.

```go
kit.ActivityCardOpt{Icon: kit.IconPencil, Title: "apply_patch", Status: kit.ActivityFailed}
```

### Semantic colors

`Tone` changes the title and icon colors.

```go
kit.ActivityCardOpt{Icon: kit.IconTriangleAlert, Title: "Model call failed", Tone: kit.ActivityDanger}
```

### In progress

`Shimmer` makes the title shimmer, for thinking, compaction, or work in progress.

```go
kit.ActivityCardOpt{Icon: kit.IconBrain, Title: "Thinking", Shimmer: true}
```

### Custom content after the title

`Extra` appears immediately after the title, outside the expand button. It can hold interactive content such as a conversation link. Setting it hides `Summary`.

## API

### ActivityCard

`func ActivityCard(c *ui.Context, open *bool, opt ActivityCardOpt, actions func(), detail func()) ui.Element`

Details are indented to align with the title text. The vertical line on the left aligns with the center of the icon, with 12 between blocks.

- **open** `*bool`: Expanded state, owned by the caller.
- **opt** `ActivityCardOpt`: Header row content.
- **actions** `func()`: Function that draws the action buttons on the right. May be `nil`.
- **detail** `func()`: Function that draws the expanded content. `nil` disables expansion.

### ActivityCardOpt

- **Icon** `*ui.SVG`: Type icon. Expandable entries replace it with an arrow on hover, focus, and expansion.
- **Title** `string`: Title.
- **Summary** `string`: Argument summary, truncated to one line, with no hover tooltip and in the placeholder color.
- **Status** `ActivityStatus`: Status. Empty and hidden by default.
- **Tone** `ActivityTone`: Semantic color. Defaults to `ActivityNormal`.
- **Shimmer** `bool`: Makes the title shimmer.
- **Extra** `func()`: Custom content immediately after the title.

### ActivityStatus

- **ActivityRunning**: Running; shows a spinner.
- **ActivityFailed**: Failed; shows a danger-colored icon.
- **ActivityInterrupted**: Interrupted; uses the warning color.
- **ActivityBlocked**: Blocked; uses the warning color.
- **ActivityStopped**: Stopped; uses the warning color.

### ActivityTone

- **ActivityNormal**: Secondary color, changing to the body text color on hover.
- **ActivityWarning**: Warning color.
- **ActivityDanger**: Danger color.

### StatusIcon

`func StatusIcon(c *ui.Context, status ActivityStatus)`

Draws a status icon on its own and returns no element. Success (an empty value) draws nothing. The icon has a localized accessible name, such as “Running”. It can be reused in lists outside `ActivityCard`.
