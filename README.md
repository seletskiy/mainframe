# `mainframe` — not a typical terminal

Typical terminals are actually terminal emulators that abide specifications
from eighties, where control functions are represented by escape sequences
which can only be controlled synchronously in one thread and requires moving
cursor around to display text in different parts of screen. Text that displayed
on screen is not separated by control characters and `cat`-ing binary file will
corrupt that kind of terminals.

![Grand Scheme](diagram.png)

`mainframe` proposes different design. Terminal can be controlled via
unix-socket and multiple clients can control what is displayed on the screen
via text-based or binary-based protocol. So terminal does not process `stdin`,
`stdout` or `stderr` of attached programs at all and all terminal manipulations
should be done via socket connections.

# State of development

`mainframe` is in very early development stage.

Terminal currently supports only:

- [x] load simple bitmap font from image;

- [x] starting in daemon mode (`-L` flag);

- [x] render any number of windows from single instance;

- [x] executing specified commands (`-E` flag) with automatic window creating
  and attaching running command by opening socket and passing it as file
  descriptor no. `3`;

- [x] displaying specified text at specified coords;

Next:

- [ ] clarify error messages;

- [ ] define common protocol messages;
