# Element Call Bots

Quick and dirty tool that emulates users joining the given Element Call URL by using the playwright-go framework for the end-to-end testing.

Use `--help` flag to understand the usage.

Currently, the tool works by assuming that there is a list of pre-registered
bot users (hardcoded at the moment). Alternatively, we could register the users
on the fly (initial version of this tool did so), but both things are subject
to rate limitting and it seems like registering users on the fly requires much
more time to really spawn bots.

There is a limit of max. 20 bots at the moment. It's unlikely that a regular user machine would be able to spawn more though.

This is in no way a tool that is ready for production usage.

Press `CTRL+C` to stop the tool after running.
