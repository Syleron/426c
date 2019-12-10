package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"time"
)

type UserInterface struct {
	textView *tview.TextView
}

func (ui *UserInterface) setup() {
	app := tview.NewApplication()

	// Setup the text view
	ui.textView = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			app.Draw()
		})
	ui.textView.SetBorder(true)

	// Setup the chat input form
	chatInput := tview.NewInputField().
		SetLabel(" " + "Syleron" + " ").
		SetFieldWidth(0)
	chatInput.SetDoneFunc(func(key tcell.Key) {
		//words := strings.Fields(chatInput.GetText())
		//command := &Commands{}
		//if len(words) > 0 {
		//	if exists := command.Send(words[0], words); !exists {
				ui.addMessage("white", chatInput.GetText(), "Syleron")
		//	}
		//}
		chatInput.SetText("")
	})
	chatInput.SetBorder(true)

	// Setup the flex grid
	flex := tview.NewFlex().
		AddItem(
			tview.NewFlex().
				SetDirection(tview.FlexRow).
				AddItem(
					ui.textView,
					0,
					3,
					false,
				).
				AddItem(
					chatInput,
					3,
					1,
					true,
				),
			0,
			2,
			true,
		)

	// **
	// **
	// **

	ui.addText(`
 ___ ___ ___     ___ ___ ___ ___ ___ _   
| | |_  |  _|___|  _|  _|  _|_  |  _| |_ 
|_  |  _| . |  _| . |  _| . |_  | . | . |
  |_|___|___|___|___|_| |___|___|___|___|
                       Security. Privacy.
	`)

	ui.addText("[yellow]This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.\n")

	ui.addWarningMessage("testing123")
	ui.addInfoMessage("asdadssssss")
	ui.addDebugMessage("this is a debug message")
	ui.addErrorMessage("this is an error message")

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}


func (ui *UserInterface) addWarningMessage(message string) {
	ui.addMessage("yellow", message, ".System.")
}

func (ui *UserInterface) addDebugMessage(message string) {
	ui.addMessage("orange", message, ".System.")
}

func (ui *UserInterface) addErrorMessage(message string) {
	ui.addMessage("red", message, ".System.")
}

func (ui *UserInterface) addInfoMessage(message string) {
	ui.addMessage("blue", message, ".System.")
}

func (ui *UserInterface) addMessage(color, message, username string) {
	if ui.textView == nil {
		return
	}
	ui.textView.Write(
		[]byte(fmt.Sprintf("["+color+"][%s][%s] " + message + "\n", time.Now().Format("15:04"), username)),
	)
}

/**
Add message to text view
Note: this will add a new line automatically
*/
func (ui *UserInterface) addText(message string) {
	if ui.textView == nil {
		return
	}
	ui.textView.Write([]byte(message + "\n"))
}
