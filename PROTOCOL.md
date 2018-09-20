# Text Protocol Design

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

See [COMMANDS.md](COMMANDS.md) for examples of this format.

## General

All successfull responses for commands that does not return value start from
`ok`.

All successfull responses contain `tick: 123` which specify tick on
which response was generated.

All error messages will be in form `err msg: "text"`.

## Ticks

Tick is UNIX timestamp with microsecond precision.

Tick is same for all windows on same render interation.
