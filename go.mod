module github.com/syleron/426c

go 1.13

require (
	github.com/ProtonMail/go-mime v0.0.0-20190521135552-09454e3dbe72 // indirect
	github.com/ProtonMail/gopenpgp v1.0.0
	github.com/boltdb/bolt v1.3.1
	github.com/gdamore/tcell v1.3.0
	github.com/ipfs/go-log v1.0.0
	github.com/labstack/gommon v0.3.0
	github.com/mattn/go-runewidth v0.0.7 // indirect
	github.com/rivo/tview v0.0.0-20191129065140-82b05c9fb329
	golang.org/x/crypto v0.0.0-20191206172530-e9b2fee46413 // indirect
)

replace golang.org/x/crypto => github.com/ProtonMail/crypto v0.0.0-20190427044656-efb430e751f2
