# reprint 📖

Cleans styling in EPUB ebooks.

## Why?

Many ebooks have terrible, obnoxious or unnecessary styling.

For example: oversized margins, excessive line breaks, non-default fonts and
sizes, peculiar spacing.

This is undesirable because it distracts from the reading experience.

Thus, **reprint** tries to fix as much of this as possible.

## Install

**reprint** requires [Go](https://golang.org/doc/install) to be installed.

```
go get -u github.com/asavoy/reprint
```

To find out where **reprint** was installed, you can run
`go list -f {{.Target}} github.com/asavoy/reprint`.

For `reprint` to be used globally, add that directory to your `$PATH`
environment var.

## Usage

```
reprint source.epub fixed.epub
```

## Design goals

**Optimise for the Apple Books app**

- Simply because it's the only app I use
- Consider both iPhone and iPad

**Remove custom styling, in favor of built-in defaults**

- The built-in styles have been optimised by the app and user settings, don't
  override them
- Assumes the content semantics are reasonable
- Justified text is bad for readability

**Add styling only when it improves readability**

- Tables and asides don't have built-in styling
- Avoid page breaks between images and their captions

**Fix common markup problems**

- Use of `<br>` for spacing between paragraphs
- Use of `<blockquote>` for unnecessary margins

**Preserve metadata when needed for library management**

- Don't bother with EPUB metadata that isn't utilised by apps

**Don't optimise the ebook internals**

- As long at the output works, the internals don't matter to the reader
- Don't bother making stylesheets DRY
