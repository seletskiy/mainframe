# Events

To receive events client should first subscribe using `subscribe` command,
which is described in [COMMANDS.md](COMMANDS.md).

All events has following format:

```
event tick: 123 kind: "<event kind>" <event specific data>
```

`<event kind>` is same as name which was used with `subscribe` command.


## `resize`: emitted on window resize

```
event tick: 123 kind: "resize" columns: 80 rows: 20 width: 1024 height: 768
```

* `columns` and `rows` specify width and height in glyphs;
* `width` and `height` specify width and height in pixels;


## `keyboard`: emitted on keyboard events like key press

```
event tick: 123 kind: "keyboard" symbol: "Ctrl+Q" (press|release|repeat) code: 234 [ctrl] [shift] [alt] [super]
```

Keyboard events never consider current active layout and symbol will be always
reported in US layout.

To read input text use `input` subscription. For example, this event can't be
used to capture `!`, because pressing `Shift+1` produce `symbol: "Shift+1'`
event, not `symbol: "!"`.

* `symbol` can contain:
    * single key like `Q`;
    * key with modifier like `Ctrl+Q`;
    * key with several modifiers like `Ctrl+Shift+Q`;
    * maximum available modifier `Ctrl+Shift+Alt+Super` (in that order);
    * single modifier key like `LShift` or `RCtrl`;
    * special key like `Esc` or `F12`;
    * special key with modifiers like `Ctrl+Shift+Esc`;
* `press` flag will be specified at key down event;
* `release` flag will be specified at key up event;
* `repeat` flag will be specified ater key down event and when key is hold
  longer than repeat interval set in OS;
* `ctrl`, `shift`, `alt` and `super` will be specified when corresponding
  modifier key was hold during event; for example, you can receive
  `symbol: "Ctrl+Q" ctrl`;
* `code` will be equal to `0` if `release` event was produced by loosing focus;
* it is possible to capture `LShift+RShift` like events; see examples above;

### Examples

#### Single keypress of key `Q`

```
event tick: 1537483348933897 kind: "keyboard" symbol: "Q" press code: 24
event tick: 1537483348967395 kind: "keyboard" symbol: "Q" release code: 24
```

#### Holding of key `Q`

```

event tick: 1537484086799050 kind: "keyboard" symbol: "Q" press code: 24
event tick: 1537484086915845 kind: "keyboard" symbol: "Q" repeat code: 24
event tick: 1537484086915845 kind: "keyboard" symbol: "Q" repeat code: 24
event tick: 1537484086932732 kind: "keyboard" symbol: "Q" repeat code: 24
event tick: 1537484086932732 kind: "keyboard" symbol: "Q" release code: 24
```

#### Single keypress of `Ctrl+C`

```
event tick: 1537483485868994 kind: "keyboard" symbol: "LCtrl" press code: 66
event tick: 1537483486536455 kind: "keyboard" symbol: "Ctrl+C" press code: 54 ctrl
event tick: 1537483486586198 kind: "keyboard" symbol: "Ctrl+C" release code: 54 ctrl
event tick: 1537483487036929 kind: "keyboard" symbol: "LCtrl" release code: 66 ctrl
```

Note, that `LCtrl` key press and following release was captured too.

#### Single keypress of `LShift+RShift`

```
event tick: 1537483734544364 kind: "keyboard" symbol: "LShift" press code: 37
event tick: 1537483735579292 kind: "keyboard" symbol: "RShift" press code: 108 shift
event tick: 1537483735629280 kind: "keyboard" symbol: "RShift" release code: 108 shift
event tick: 1537483736512712 kind: "keyboard" symbol: "LShift" release code: 37 shift
```

Note, that actual `LShift+RShift` key press event is recorded as `RShift`
symbol with `shift` modifier.

So, if you want to capture `LShift+RShift`, look for `symbol: "RShift" press
shift` events.

But if you want to capture `RShift+LShift`, look for `symbol: "LShift" press
shift`


## `input`: emitted on keyboard events which produce printable characters

```
event tick: 123 kind: "input" char: "q" [ctrl] [shift] [alt] [super]
```

* `char` will refer to layout specific typed character taking into account
  modifier key, e.g. if you will press key `Q`, you will received `char: "q"`,
  but if you will press `Shift+Q`, you will receive `char: "Q" shift`
