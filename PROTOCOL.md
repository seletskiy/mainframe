# Protocol

Protocol is:
* full-duplex
* request-response.

Protocol can be used in text or binary form. Only text form is supported now.

## Text Protocol

### Format

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

### General

All successfull responses for commands starts from `ok`.

All error messages will be in form `error message: "text"`.

### Ticks

Some commands return `tick` value in their response.

Tick is UNIX timestamp with microsecond precision.

Tick is generated on render event for each window separately.
