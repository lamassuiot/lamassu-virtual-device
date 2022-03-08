package deviceview

import (
	"github.com/gdamore/tcell/v2"
	"github.com/go-kit/log"
	"github.com/lamassuiot/lamassu-virtual-device/pkg/observer"
	"github.com/rivo/tview"
)

func GetItem(logger log.Logger, deviceData *observer.DeviceState, app *tview.Application) tview.Primitive {

	deviceForm := tview.NewForm().
		AddInputField("Device ID", deviceData.Config.Device.Id, 50, nil, func(text string) {
			deviceData.Config.Device.Id = text
		}).
		AddInputField("Device Alias", deviceData.Config.Device.Alias, 50, nil, func(text string) {
			deviceData.Config.Device.Alias = text
		})

	pages := tview.NewPages()

	activePage := 1
	pages.AddPage("enroll-page", GetEnrollItem(logger, deviceData, app), true, activePage == 1)
	pages.AddPage("re-enroll-page", GetReenrollItem(logger, deviceData, app), true, activePage == 2)
	pages.AddPage("reset-page", GetReenrollItem(logger, deviceData, app), true, activePage == 2)
	pages.AddPage("aws-page", GetAwsItem(logger, deviceData, app), true, activePage == 2)

	actionsList := tview.NewList().
		ShowSecondaryText(false)

	actionsList.AddItem("Enroll", "", '1', func() {
		activePage = 1
		pages.SwitchToPage("enroll-page")
	})
	actionsList.AddItem("Reenroll", "", '2', func() {
		activePage = 2
		pages.SwitchToPage("re-enroll-page")
	})
	actionsList.AddItem("Reset Bootstrap", "", '3', func() {
		activePage = 3
		pages.SwitchToPage("reset-page")
	})
	actionsList.AddItem("AWS Integration", "", '4', func() {
		activePage = 4
		pages.SwitchToPage("aws-page")
	})

	test := tview.NewBox().
		SetBorder(false).
		SetDrawFunc(func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
			// Draw a horizontal line across the middle of the box.
			centerY := y + height/2
			for cx := x + 1; cx < x+width-1; cx++ {
				screen.SetContent(cx, centerY, tview.BoxDrawingsLightHorizontal, nil, tcell.StyleDefault.Foreground(tcell.NewHexColor(0xffffff)))
			}
			// Write some text along the horizontal line.
			tview.Print(screen, " Operation ", x+1, centerY, width-2, tview.AlignCenter, tcell.NewHexColor(0xF9F1A5))
			// Space for other content.
			return x + 1, centerY + 1, width - 2, height - (centerY + 1 - y)
		})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(deviceForm, 5, 0, false).
		AddItem(test, 2, 0, false).
		AddItem(pages, 0, 1, false).
		AddItem(actionsList, 4, 0, true)

	flex.SetBorder(true)
	return flex
}
