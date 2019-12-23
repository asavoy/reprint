# reprint 📖

Cleans styling in EPUB ebooks.

## Why?

Many ebooks have terrible, obnoxious or unnecessary styling.

For example: oversized margins, excessive line breaks, non-default fonts and 
sizes, peculiar spacing.

This is undesirable, because it distracts from the reading experience. It also
overrides styling that has been optimised by ebook reading apps and user 
settings.

Thus, **reprint** tries to fix as much of this as possible.

## Usage

```
reprint source.epub fixed.epub
```

## Design goals

**Optimise for the Apple Books app**

- Simply because it's the only app I use

**Replace styling with good defaults**

- It's too hard to repair styles automatically
- This assumes the content semantics are reasonable

**Preserve the metadata**

- It's necessary for library management

**Don't optimise the ebook internals**

- As long at the output works, the internals don't matter to the user
