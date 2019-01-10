# Commands

* [Session commands](#window-commands)
   * [`get`: retrieve various mainframe information](#get)
      * [Request](#get-request)
      * [Response](#get-response)
* [Window commands](#window-commands)
   * [`open`: opens new window and binds session to it](#open)
      * [Request](#open-request)
      * [Response](#open-response)
   * [`reshape`: change window position and/or size](#reshape)
      * [Request](#reshape-request)
      * [Response](#reshape-response)
* [Output commands](#output-commands)
   * [`clear`: clear cells on screen](#clear)
      * [Request](#clear-request)
      * [Response](#clear-response)
   * [`put`: assign text, fg or bg color to cells on screen](#put)
      * [Request](#put-request)
      * [Response](#put-response)
      * [Example: draw vim-like line numbers column](#put-example-1)
   * [`subscribe`: subscribe on given events](#subscribe)
      * [Request](#subscribe-request)
      * [Response](#subscribe-response)

## Legend

Following syntax is used to describe commands:

* `[x]` — `x` or nothing.
* `(x|y)` — `x` or `y`.
* `[x|y]` — `x` or `y` or nothing.
* `([x] [y])` — `x` or `y` or `x y`.
* `<x>...` — placeholder for any arguments.

## <a id="session-commands"> Session commands

### <a id="get"> `get`: retrieve various mainframe information

#### <a id="get-request"> Request

```
get <options>...
```

See following sections for each option.

### <a id="get-font"> `get font`: retrieve font information

#### <a id="get-font-request"> Request

```
get font
```

#### <a id="get-font-response"> Response

```
ok width: 8 height: 18
```

## <a id="window-commands"> Window commands

### <a id="open"> `open`: opens new window and binds session to it

#### <a id="open-request"> Request

```
open [width: 640 height: 480|columns: 80 rows: 20] [x: 1 y: 2] [title: "string"] [raw] [hidden] [fixed] [bare] [floating]
```

* when used in new open connection to mainframe `open` will bind created window
  to this connection, so all further commands can be used without need to
  specify window ID.
* `open` without arguments will open window with default size;
* `width` and `height` can be used to specify window size in pixels;
* `columns` and `rows` can be used to specify window size using current font size;
* if tiling WM is used, then size and position arguments will be ignored
  unless window is made floating via WM configuration or `raw` option is
  specified;
* `raw` creates window that completely ignored by WM;
* `floating` hint can be ignored by WM and have no effect;
* `hidden` creates window that should be shown to be displayed;
* `bare` specify that window should be created without any WM decorations (e.g. no borders);

#### <a id="open-args"> Arguments

| Argument | Type   | Description                                             |
| :------- | :---   | :----------                                             |
| width    | int    | Width of window in pixels.                              |
| height   | int    | Height of window in pixels.                             |
| columns  | int    | Width of window in columns (based on font width).       |
| rows     | int    | Height of window in rows (based on font height).        |
| x        | int    | Onscreen position of window in pixels.                  |
| y        | int    | Onscreen position of window in pixels.                  |
| title    | string | Title for window.                                       |
| raw      | bool   | Create window that is not managed by WM.                |
| hidden   | bool   | Create hidden window that need to be shown with `show`. |
| fixed    | bool   | Create fixed size window.                               |
| bare     | bool   | Create window without WM decorations.                   |
| floating | bool   | Create floating window (WM specific).                   |

#### <a id="open-response"> Response

```
ok id: 123
```

* `id` can be used for other window manipulation commands;
* after window is open all further commands will operate on this window by
  default;
* if client closes connection, all open windows that were opened via this
  connection will be closed;

---

### <a id="reshape"> `reshape`: change window position and/or size

#### <a id="reshape-request"> Request

```
reshape ([width: 640 height: 480|columns: 80 rows: 20] [x: 1 y: 2])
```

* `reshape` can be used to move and resize window in single command;

#### <a id="reshape-args"> Arguments

| Argument | Type | Description                                                    |
| :------- | :--- | :----------                                                    |
| width    | int  | Target window width in pixels.                                 |
| height   | int  | Target window height in pixels.                                |
| columns  | int  | Target window width in cell columns (based on font width).     |
| rows     | int  | Target window height size in cell rows (based on font height). |
| x        | int  | Onscreen position of window in pixels.                         |
| y        | int  | Onscreen position of window in pixels.                         |

#### <a id="reshape-response"> Response

```
ok
```

## <a id="output-commands"> Output commands

### <a id="clear"> `clear`: clear cells on screen

#### <a id="clear-request"> Request

```
clear [x: 1 y: 2 [columns: 80] [rows: 20]]
```

* `clear` without arguments will clear entire screen;
* `clear x: 1 y: 2` will clear single cell at position `(1; 2)`;
* full form will clear specified area;

#### <a id="clear-args"> Arguments

| Argument | Type | Description                                        |
| :------- | :--- | :----------                                        |
| x        | int  | Column coordinate of first cell to clear.          |
| y        | int  | Row coordinate of first cell to clear.             |
| columns  | int  | Amount of columns to clear (including first cell). |
| rows     | int  | Amount of rows to clear (including first cell).    |

#### <a id="clear-response"> Response

```
ok [offscreen]
```

* `offscreen` flag will be in response if request attempts to clear cells
   outside of screen;

---

### <a id="put"> `put`: assign text, fg or bg color to cells on screen

#### <a id="put-request"> Request

```
put x: 1 y: 2 [columns: 80] [rows: 20] ([fg: #ff0] [bg: #f00] [text: "string"]) [tick: 123] [exclusive]
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
* text may contain `\n` to put following text to next row;

#### <a id="put-args"> Arguments

| Argument  | Type   | Description                                                                                                               |
| :-------  | :---   | :----------                                                                                                               |
| x         | int    | Column coordinate of first cell.                                                                                          |
| y         | int    | Row coordinate of first cell.                                                                                             |
| columns   | int    | Maximum number of columns this operation will touch.                                                                      |
| rows      | int    | Maximum number of rows this operation can touch.                                                                          |
| fg        | color  | New foreground color for cells (e.g. text color).                                                                         |
| bg        | color  | New background color for cells.                                                                                           |
| text      | string | Text to put in cells. Text will be wrapped to next row if rows specified or trimmed otherwise.                            |
| exclusive | bool   | Mark region of cells `(x, y, x+columns, y+rows)` as exclusive, which will be cleared when cell `(x, y)` will be modified. |
| tick      | int    | *Not implemented.*                                                                                                        |

#### <a id="put-response"> Response

```
ok [offscreen] [overflow]
```

* `offscreen` flag will be in response if request attempts to modify cells
  outside of screen;
* `overflow` flag will be in response if given `text` can't be fit in specified
  area;

#### <a id="put-example-1"> Example: draw vim-like line numbers column

```
put x: 0 y: 0 columns: 2 rows: 15 text: " 1 2 3 4 5 6 7 8 9101112131415" bg: #333
```

Note: because of columns is specified to `2` given `text` will be wrapped after
every `2` characters.

---

### <a id="subscribe"> `subscribe`: subscribe on given events

#### <a id="subscribe-request"> Request

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

#### <a id="subscribe-response"> Response

```
ok
```
