# Element Call Bots

Emulates users joining the given Element Call URL.

The only things that you need to specify are:

1. Element Call URL.
2. Room name.
3. The amount of bots you want to spawn.

Note that the joining process may be slow due to the rate-limiting of the
home-server since each attempt to join a session essentially registers a new
bot user (this is due to the fact that Element Call works on top of Matrix that
does not have a notion of an anonymous or guest participant that does not
require the registration).
