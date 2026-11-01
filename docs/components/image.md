---
title: Image
description: Displays caller-provided images with loading and error states; ImagePreview supports fitting to the window, zooming, and scrolling.
---

# Image

Displays images. This library does not read from disk or download images. Your application supplies image data, loading state, errors, and retry actions through `ImageResource`.

<Preview page="Images" :height="420">

<<< @/../cmd/gallery/image.go

</Preview>

## Usage

```go
import kit "github.com/VaalaCat/mygo-kit"
```

```go
kit.Image(c, kit.ImageResource{Source: source, Alt: "Architecture diagram"})
```

`Source` is MyGo's `ui.ImageSource`; for example, `ui.NewBitmap(img)` creates one from a Go `image.Image`. Images keep their original aspect ratio and never exceed the container's width.

## Examples

### Loading and errors

When `Source` is `nil`, the component displays a message, taking the first available item in this order: `Error` text, loading, then `Alt`. If none is available, it shows “Image unavailable”. When `Error` is present and `Retry` is set, a “Retry image” button appears below the message.

```go
kit.Image(c, kit.ImageResource{Alt: "Architecture diagram", Loading: true})

kit.Image(c, kit.ImageResource{
	Alt:   "Architecture diagram",
	Error: "Download failed: 404",
	Retry: func() { download(url) },
})
```

### Preview

`ImagePreview` is a preview area with a toolbar containing “Fit to window”, “Original size”, “−”, and “+”. The image can scroll in both directions. You own the zoom state.

```go
var preview = kit.ImagePreviewState{Fit: true} // Owned by the application

kit.ImagePreview(c, resource, &preview).Grow(1)
```

## API

### Image

`func Image(c *ui.Context, image ImageResource) ui.Element`

An image with a corner radius of 4. `Alt` is its accessible name.

### ImageResource

- **Source** `ui.ImageSource`: Image data. `nil` means it is not yet available.
- **Alt** `string`: Alternative text, also used as the accessible name.
- **Loading** `bool`: Loading in progress. Shows “Loading image…” when `Source` is `nil`.
- **Error** `string`: Description of a loading failure. Takes precedence over `Loading`.
- **Retry** `func()`: Called when the “Retry image” button is clicked. The button is shown only when `Error` is non-empty.

### ImagePreview

`func ImagePreview(c *ui.Context, image ImageResource, state *ImagePreviewState) ui.Element`

The toolbar title is `image.Alt`. “−” and “+” change the zoom by a factor of 1.25 each time, within a range of 0.1 to 8.

- **image** `ImageResource`: Image to preview.
- **state** `*ImagePreviewState`: Zoom state, owned by the caller.

### ImagePreviewState

- **Fit** `bool`: Fits the image to the available window size when true.
- **Zoom** `float32`: Zoom factor, used when `Fit` is false. Values less than or equal to 0 are treated as 1.
