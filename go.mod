module github.com/syleron/426c

go 1.24.0

toolchain go1.24.4

require (
	github.com/AllenDang/giu v0.14.1
	github.com/ProtonMail/gopenpgp v1.0.0
	github.com/boltdb/bolt v1.3.1
	github.com/charmbracelet/bubbletea v1.3.6
	github.com/charmbracelet/lipgloss v1.1.0
	github.com/gdamore/tcell v1.3.0
	github.com/labstack/gommon v0.3.0
	github.com/olekukonko/tablewriter v0.0.4
	github.com/prometheus/client_golang v1.19.1
	github.com/rivo/tview v0.0.0-20191129065140-82b05c9fb329
	go.etcd.io/bbolt v1.4.2
)

require (
	github.com/AllenDang/cimgui-go v1.3.2-0.20250409185506-6b2ff1aa26b5 // indirect
	github.com/AllenDang/go-findfont v0.0.0-20200702051237-9f180485aeb8 // indirect
	github.com/ProtonMail/go-mime v0.0.0-20190521135552-09454e3dbe72 // indirect
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/charmbracelet/colorprofile v0.2.3-0.20250311203215-f60798e515dc // indirect
	github.com/charmbracelet/x/ansi v0.9.3 // indirect
	github.com/charmbracelet/x/cellbuf v0.0.13-0.20250311204145-2c3ea96c31dd // indirect
	github.com/charmbracelet/x/term v0.2.1 // indirect
	github.com/erikgeiser/coninput v0.0.0-20211004153227-1c3628e74d0f // indirect
	github.com/faiface/mainthread v0.0.0-20171120011319-8b78f0a41ae3 // indirect
	github.com/gdamore/encoding v1.0.0 // indirect
	github.com/gucio321/glm-go v0.0.0-20241029220517-e1b5a3e011c8 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.1 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-colorable v0.1.2 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-localereader v0.0.1 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/mazznoer/csscolorparser v0.1.6 // indirect
	github.com/muesli/ansi v0.0.0-20230316100256-276c6243b2f6 // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/muesli/termenv v0.16.0 // indirect
	github.com/napsy/go-css v1.0.0 // indirect
	github.com/pkg/browser v0.0.0-20210911075715-681adbf594b8 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.48.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/sahilm/fuzzy v0.1.1 // indirect
	github.com/sirupsen/logrus v1.4.2 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.0.1 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	golang.design/x/hotkey v0.4.1 // indirect
	golang.design/x/mainthread v0.3.0 // indirect
	golang.org/x/crypto v0.23.0 // indirect
	golang.org/x/image v0.27.0 // indirect
	golang.org/x/sync v0.15.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/eapache/queue.v1 v1.1.0 // indirect
)

replace golang.org/x/crypto => github.com/ProtonMail/crypto v0.0.0-20190427044656-efb430e751f2

replace github.com/boltdb/bolt => go.etcd.io/bbolt v1.3.9
