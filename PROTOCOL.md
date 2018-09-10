# Protocol Design

Protocol is:
* full-duplex
* request-response.

## Format

* Protocol is line oriented: messages are separated by `\n`.
* Message consists of message tag followed by optional body.
* Message tag consists of `[a-z_]`.
* Message body consist of white-space delimited key or key-value pairs.
* Keys consists of `[a-z_]`.
* Key-value pairs delimited with `:`.
* Values can be following types:
    * integer (signed);
    * color in form `#rgb` or `#rrggbb`;
    * double-quoted string where `"\\" = \` and `"\"" = "`

## General

All messages from server prefixed to ease parsing:

* `=` for responses to commands;
* `!` for events;
* `#` for errors;

All successfull responses for commands that does not return value start from
`=ok`.

All successfull responses contain `tick: 123` which specify tick on
which response was generated.

All error messages will be in form `#err msg: "text"`.

## Legend

* `[x]` — `x` or nothing.
* `(x|y)` — `x` or `y`.
* `([x] [y])` — `x` or `y` or `x y`.

## `client` → `server`

### `clear`: clear cells on screen

#### Request

```
clear [x: 1 y: 2 [width: 80] [height: 20]]
```

* `clear` without arguments will clear entire screen;
* `clear x: 1 y: 2` will clear single cell at position `(1; 2)`;
* full form will clear specified area;

#### Response

```
=ok [offscreen]
```

* `offscreen` flag will be in response if request attempts to clear cells
   outside of screen;

### `put`: assign text, fg or bg color to cells on screen

#### Request

```
put x: 1 y: 2 [width: 80] [height: 20] ([fg: #ff0] [bg: #f00] [text: "text string"]) [tick: 123] [exclusive]
```

* `put` can change foreground, background or text at once or one-by-one;
* `tick: 123` schedule `put` command on specified terminal tick;
* `exclusive` will mark area changed by `put` command; if another command will
  change cell at specified coordinates, then all marked area will
  be cleared first;
* if `width` and `height` is given, then entire area will be changed; text
  will be wrapped to subsequent lines when single line will exceed given `width`;
* if only `width` is given, `height` is assumed to be `1`;
* if only `height` is given, `width` is assumed to be `1`;
* if `width` and `height` is not specified, then `height` is assumed to be `1`
  and `width` to be equal to length of given `text`; if `text` is not given,
  then `put` will only change single specified cell;

#### Response

```
=ok tick: 123 [offscreen] [overflow]
```

* `offscreen` flag will be in response if request attempts to modify cells
  outside of screen;
* `overflow` flag will be in response if given `text` can't be fit in specified
  area;

### `subscribe`: subscribe on given events

Format:

```
subscribe ([keyboard] [resize])
```

* `keyboard` subscription will allow to receive `keyup`, `keydown` and
  `keypress` events;
* `resize` will allow to receive window size change events;
* after subscription it is still possible to send other commands;
* see next section for list of events that server will send back to client;

## `server` → `client`

### `keyboard` events

### Format

```
!event tick: 123 kind: ("keyup"|"keydown"|"keypress") key: "x" [shift] [alt] [ctrl]
```

### `resize` event

```
!event tick: 123 kind: "resize" width: 80 height: 20
```
