package log

import "github.com/mgutz/ansi"

const banner = `
  ____        ____       ____
 / ___| ___  / ___|___  / ___|___
| |  _ / _ \| |   / _ \| |   / _ \
| |_| | (_) | |__| (_) | |__| (_) |
 \____|\___/ \____\___/ \____\___/
`

func DisplayGoCoCo() {
	stdout.Write([]byte(ansi.Color(banner, "cyan+b")))
}
