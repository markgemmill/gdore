package main

import (
	"fmt"
	"gdore/broker"
	"gdore/environ"
	"gdore/scraper"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var messanger *broker.MessageBroker = broker.NewMessageBroker()

type ApplicationParameters struct {
	user      binding.String
	password  binding.String
	region    binding.String
	documents binding.String
	progress  binding.String
	button    *widget.Button
	display   *widget.Label
}

func Execute(params *ApplicationParameters, env environ.Environ) func() {
	return func() {
		// get the input values
		userValue, _ := params.user.Get()
		passValue, _ := params.password.Get()
		regionValue, _ := params.region.Get()
		docsValue, _ := params.documents.Get()
		params.button.Disable()
		params.button.SetText("running...")
		openFile := ""

		go scraper.RunScraper(userValue, passValue, regionValue, docsValue, env)

		go func() {
			jobCompleted := false
			ticker := time.NewTicker(time.Second)
			for range ticker.C {
				if jobCompleted {
					break
				}
				fyne.Do(func() {
					msg, okay := messanger.Receive()
					if !okay {
						return
					}
					if msg.Done || msg.Err != nil {
						jobCompleted = true
						openFile = msg.File
					}
					params.progress.Set(msg.Message)
					if msg.Err != nil {
					}
				})
			}
			fyne.Do(func() {
				params.button.Enable()
				params.button.SetText("GO")
				if openFile != "" {
					environ.OpenFile(openFile)
				}
			})
		}()
	}
}

func Gui(env environ.Environ) {
	os.Setenv("FYNE_SCALE", "1.5")
	application := app.New()
	window := application.NewWindow("Guillaume's Helper")
	window.SetFixedSize(true)
	window.Resize(fyne.NewSize(400, 400))

	params := ApplicationParameters{
		user:      binding.NewString(),
		password:  binding.NewString(),
		region:    binding.NewString(),
		documents: binding.NewString(),
		progress:  binding.NewString(),
	}

	if env.DefaultUser != "" {
		params.user.Set(env.DefaultUser)
	}

	userLabel := widget.NewLabel("User")
	passLabel := widget.NewLabel("Password")
	// separatorLabel := widget.NewLabel("")
	regionLabel := widget.NewLabel("Region")
	docLabel := widget.NewLabel("Documents")

	params.button = widget.NewButton("GO", Execute(&params, env))
	params.display = widget.NewLabelWithData(params.progress)
	params.display.Alignment = fyne.TextAlignCenter
	params.display.TextStyle = fyne.TextStyle{
		Bold: true,
	}

	// user id widget
	userInput := widget.NewEntryWithData(params.user)
	userInput.OnSubmitted = func(s string) {
		if s != env.DefaultUser {
			env.DefaultUser = s
			fmt.Println("reset config default user")
			environ.DumpConfig(env)
		}
	}

	// password widget
	passInput := widget.NewEntryWithData(params.password)
	passInput.Password = true

	// separator := widget.NewSeparator()

	regions := []string{
		"National",
		"Atlantic",
		"Ontario",
		"West",
		"Quebec",
	}

	// region widget
	regionInput := widget.NewSelectWithData(regions, params.region)

	// documents widget
	docInput := widget.NewEntryWithData(params.documents)
	docInput.MultiLine = true
	docInput.SetMinRowsVisible(10)

	// loginForm := container.New(
	// 	layout.NewFormLayout(),
	// 	userLabel,
	// 	userInput,
	// 	passLabel,
	// 	passInput,
	// )

	searchForm := container.New(
		layout.NewFormLayout(),
		userLabel,
		userInput,
		passLabel,
		passInput,
		regionLabel,
		regionInput,
		docLabel,
		docInput,
	)

	execution := container.NewVBox(
		params.display,
		params.button,
	)

	content := container.NewBorder(
		nil,
		execution,
		nil,
		nil,
		searchForm,
	)

	window.SetContent(content)
	window.ShowAndRun()
}
