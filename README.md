# ascii-art-web

A small web front-end for the [`ascii-art`](../ascii-art) project. Instead of
running the converter from a terminal, this starts a Go web server: the user
types some text, picks a banner, and gets the ASCII-art result back in the
browser.

## Description

The page has a text field and a choice of three banners — `standard`,
`shadow`, and `thinkertoy`. Submitting the form sends the text to the server,
which converts it and re-renders the same page with the result shown below
the form.

The conversion itself (loading a banner file, turning text into ASCII art)
is not reimplemented here — it is the exact same `banner` and `render`
packages from the original `ascii-art` CLI project, reused unchanged. This
project only adds the HTTP layer around them.

## Authors

- Mohammed Zuhair ([@mozuhair](https://learn.reboot01.com/git/mozuhair))

## Usage

Requires Go (built and tested with Go 1.24). From the project root:

```bash
go run .
```

Then open `http://localhost:8090` in a browser, type some text, choose a
banner, and click **Generate**.

## Implementation details

### Routes

| Route | Method | Behavior |
|---|---|---|
| `/` | GET | Serves the form (`templates/index.html`). Any other path is `404`. |
| `/ascii-art` | POST | Reads `text` and `banner` from the form, converts the text, and re-renders the page with the result (or an error message). |
| `/static/*` | GET | Serves `style.css` and any other file placed under `static/`. |

### Algorithm

1. `banner.Load(name)` reads the chosen banner file (embedded into the
   binary at compile time with `go:embed`) and parses it into a map of
   `rune → 8 lines of ASCII art`.
2. `render.Render(text, chars)` walks the input line by line (splitting on
   `\n`), looks up each character's 8-line block, and writes them out row by
   row so the characters line up side by side.
3. The result is injected into `templates/index.html` through Go's
   `html/template` package, which also escapes the input safely.

### Status codes

| Situation | Code |
|---|---|
| Page loads / conversion succeeds | `200` |
| Wrong HTTP method, invalid banner, text too long, or unsupported character | `400` |
| Unknown route, or a missing template file | `404` |
| An unexpected panic anywhere in a handler | `500` (caught by `recoverFromPanic`) |

### Notes

- A literal `\n` typed into the text field is converted to a real line
  break before rendering, matching the original CLI's behavior.
- Only Go's standard library is used — no external packages.