package component

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
)

type SaveImageCallback func() ([]byte, error)

func NewSaveImageDialog(window fyne.Window, callback SaveImageCallback) *dialog.FileDialog {
	dlg := dialog.NewFileSave(func(uc fyne.URIWriteCloser, e error) {
		if uc == nil || e != nil {
			log.Println("save-image: user closed the dialog or unexpected error", e)
			return
		}
		imageBytes, err := callback()
		if err != nil {
			log.Println("save-image: can not save the image", err)
			return
		}
		uc.Write(imageBytes)
	}, window)

	dlg.SetFilter(storage.NewExtensionFileFilter([]string{".jpg", ".jpeg", ".png"}))

	return dlg
}
