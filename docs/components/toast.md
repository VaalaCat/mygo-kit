---
title: Toast
description: Notifications in the top-right corner of the window, using a ShowToast call and a ToastHost.
---

# Toast

Notifications in the top-right corner of the window, with the newest at the top. There are three kinds: information, success, and error. Each can include a description and an action button.

<Preview page="Toast" :height="300">

<<< @/../cmd/gallery/toast.go

</Preview>

## Usage

```go
import kit "github.com/VaalaCat/mygo-kit"
```

Call `ToastHost` once in each window's view to display notifications. Usually, place it at the start of the view function with `defer` so it draws last:

```go
func view(c *ui.Context) {
	kit.ApplyTheme(c, kit.OneLight())
	defer kit.ToastHost(c)
	// Draw the interface
}
```

Call `ShowToast` wherever you want to send a notification:

```go
if kit.Button(c, "Export").Clicked() {
	kit.ShowToast(c, kit.ToastOpt{Title: "Session exported", Kind: kit.ToastInfo})
}
```

Notifications appear immediately even when the window is idle. You do not need to request a redraw yourself.

## Examples

### An action button

A nonempty `Action` shows a text button. Clicking it calls `OnAction`, then closes the notification.

```go
kit.ShowToast(c, kit.ToastOpt{
	Title:    "Session archived",
	Kind:     kit.ToastSuccess,
	Action:   "Undo",
	OnAction: undoArchive,
})
```

### An error notification

Error notifications stay visible for 8 seconds, longer than other notifications.

```go
kit.ShowToast(c, kit.ToastOpt{Title: "Connection lost", Description: "The runtime is not responding.", Kind: kit.ToastError})
```

## API

### ToastHost

`func ToastHost(c *ui.Context) ui.Element`

Draws notifications in the top-right corner of the window, below the title bar. Each notification is 360 wide and includes an icon for its kind, a title, an optional description, an action, and a close button. When the window uses a [TitleBar](/components/title-bar), notifications stack downward starting 8 below it. Without a TitleBar, they start 8 below the top of the window. Call once per window.

### ShowToast

`func ShowToast(c *ui.Context, opt ToastOpt)`

The window must already have a `ToastHost`.

### ToastOpt

- **Title** `string`: The title.
- **Description** `string`: A description in the secondary color. Can be empty.
- **Kind** `ToastKind`: The notification kind. The zero value (an empty string) uses the information style.
- **Action** `string`: The action button text. The button appears when this is nonempty.
- **OnAction** `func()`: Called when the action button is clicked, after which the notification closes.

### ToastKind

- **ToastInfo**: Information, with an accent-colored icon.
- **ToastSuccess**: Success, with a green icon.
- **ToastError**: Error, with a danger-colored icon. Stays visible for 8 seconds.
