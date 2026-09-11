# ascii-art-web

## Description

A web version of the `ascii-art` project. Instead of running the converter
from a terminal, the program starts an HTTP server with a simple web page:
the user types text, picks a banner (`standard`, `shadow`, or `thinkertoy`),
and gets the ASCII-art result back on the same page.

The conversion logic itself (loading a banner file and rendering text with
it) is unchanged from the original `ascii-art` project — this project only
adds an HTTP layer around it.


## Usage

Requires Go (built and tested with Go 1.24).

```bash
go run .
```

Then open `http://localhost:8080` in a browser. Type some text, choose a
banner, and click "Generate".

## Implementation details

### Routes

- `GET /` — serves the main page (`templates/index.html`) with the input
  form. Any other path returns `404`.
- `POST /ascii-art` — reads `text` and `banner` from the submitted form,
  converts the text, and re-renders the same page with the result filled
  in.

### Reused conversion logic

- `banner/banner.go` — loads a banner (embedded at compile time with
  `go:embed`, so there is no risk of a missing file at runtime) and parses
  it into a `map[rune][]string` (8 lines per character).
- `render/render.go` — turns input text into the final ASCII-art string
  using that map, one line at a time, handling `\n` as a line break and an
  empty line as a blank output line.

Both packages are copied unchanged from the original `ascii-art` CLI
project, so the web version produces byte-identical output for the same
input.

### Status codes

| Situation | Status |
|---|---|
| Successful page load / successful conversion | 200 |
| Wrong HTTP method, unknown banner name, unsupported character, text too long | 400 |
| Unknown route, or the HTML template is missing | 404 |
| Unexpected panic anywhere in a handler | 500 (via a `recover` middleware) |

### Notes

- A literal `\n` typed into the (single-line) text field is converted to a
  real line break before rendering, matching the behavior of the original
  CLI tool.
- Only Go's standard library is used (`net/http`, `html/template`, plus the
  reused `banner`/`render` packages).
