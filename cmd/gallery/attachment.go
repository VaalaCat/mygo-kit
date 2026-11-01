package main

import (
	kit "github.com/VaalaCat/mygo-kit"
	"github.com/egoist/mygo/ui"
)

var attachmentNames = []string{"screenshot-2026-10-31.png", "design.png", "a-rather-long-attachment-file-name-that-gets-truncated.png"}

func init() {
	register(page{name: "AttachmentChip", order: 44, view: func(c *ui.Context) {
		theme := kit.GetTheme(c)
		pageBody(c, func() {
			ui.Row(c).Gap(theme.Space(2)).Wrap().Children(func() {
				for i, name := range attachmentNames {
					kit.AttachmentChip(c.Key(name), name, galleryImage.Source, func() {
						attachmentNames = append(attachmentNames[:i], attachmentNames[i+1:]...)
					})
				}
				kit.AttachmentChip(c, "notes.txt"+t(c, " (no thumbnail)"), nil, nil)
			})
		})
	}})
}
