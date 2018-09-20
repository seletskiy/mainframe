# Commands

Following syntax is used to describe commands:

* `[x]` — `x` or nothing.
* `(x|y)` — `x` or `y`.
* `([x] [y])` — `x` or `y` or `x y`.

## Output commands

### `clear`: clear cells on screen

#### Request

```
clear [x: 1 y: 2 [columns: 80] [rows: 20]]
```

* `clear` without arguments will clear entire screen;
* `clear x: 1 y: 2` will clear single cell at position `(1; 2)`;
* full form will clear specified area;

#### Response

```
ok [offscreen]
```

* `offscreen` flag will be in response if request attempts to clear cells
   outside of screen;

### `put`: assign text, fg or bg color to cells on screen

#### Request

```
put x: 1 y: 2 [columns: 80] [rows: 20] ([fg: #ff0] [bg: #f00] [text: "text string"]) [tick: 123] [exclusive]
```

* `put` can change foreground, background or text at once or one-by-one;
* `tick: 123` schedule `put` command on specified terminal tick;
* `exclusive` will mark area changed by `put` command; if another command will
  change cell at specified coordinates, then all marked area will
  be cleared first;
* if `columns` and `rows` is given, then entire area will be changed; text
  will be wrapped to subsequent lines when single line will exceed given `columns`;
* if only `columns` is given, `rows` is assumed to be `1`;
* if only `rows` is given, `columns` is assumed to be `1`;
* if `columns` and `rows` is not specified, then `rows` is assumed to be `1`
  and `columns` to be equal to length of given `text`; if `text` is not given,
  then `put` will only change single specified cell;

#### Response

```
ok [offscreen] [overflow]
```

* `offscreen` flag will be in response if request attempts to modify cells
  outside of screen;
* `overflow` flag will be in response if given `text` can't be fit in specified
  area;

##### Example: draw vim-like line-number column

```
put x: 0 y: 0 columns: 2 rows: 15 text: " 1 2 3 4 5 6 7 8 9101112131415" bg: #333
```

### `subscribe`: subscribe on given events

#### Request

```
subscribe ([keyboard] [input] [resize])
```

* `keyboard` subscription will allow to receive `press`, `release` and
  `repeat` events, but without notion of keyboard layout;
* `input` subscription will allow to receive text input with notion of
  keyboard layout;
* `resize` will allow to receive window size change events;
* after subscription it is still possible to send other commands;
* see [EVENTS.md](EVENTS.md) for full description of received events;

#### Response

```
ok
```
