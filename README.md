# Element Call Bots

Quick and dirty tool that emulates users joining the given Element Call URL by using the playwright-go framework for the end-to-end testing.

# Usage

## Install

Assuming that Go is installed on your system, download the repository and install `playwright-go`:

```sh
git clone git@github.com:vector-im/sfu-load-tester-go.git && go run github.com/playwright-community/playwright-go/cmd/playwright install --with-deps
```

## Run

You can run it without any parameters to see the available options:

```sh
go run *.go
```

Example of usage:

```sh
go run *.go -call-url https://pr805--element-call.netlify.app/room/\#testcall:call.ems.host -num-bots 3 -headless
```

# How Does It Work?

Currently, the tool works by assuming that there is a list of pre-registered
bot users (hardcoded at the moment). Alternatively, we could register the users
on the fly (initial version of this tool did so), but both things are subject
to rate limitting and it seems like registering users on the fly requires much
more time to really spawn bots.

There is a limit of max. 20 bots at the moment. It's unlikely that a regular
user machine would be able to spawn more though.

This is in no way a tool that is ready for production usage.

Press `CTRL+C` to stop the tool after running.
