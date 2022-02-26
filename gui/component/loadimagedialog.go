package component

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
)

type LoadImageCallback func(path string)

func NewLoadImageDialog(window fyne.Window, callback LoadImageCallback) *dialog.FileDialog {
	dlg := dialog.NewFileOpen(func(fileURI fyne.URIReadCloser, err error) {
		if fileURI == nil || err != nil {
			return
		}

		callback(fileURI.URI().Path())
	}, window)

	dlg.SetFilter(storage.NewExtensionFileFilter([]string{".jpg", ".jpeg", ".png"}))

	return dlg
}
