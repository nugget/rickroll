# Trolling telnet.

You can `telnet macnugget.org` to see it in action.

If your host doesn't have the telnet client installed, good for you. You can also
use Netcat to access the service with the command `nc macnugget.org 23`

# What is this? Whose fault is this?

I honestly don't remember what was going through my head in 2013 when I wrote
the original version of this server in Tcl 8.5 and set it up to run on one of
my FreeBSD servers. I'm sure I felt then -- as I do now -- that it represents
the pinnacle of [Rickrolling] humor.

In 2022 I decommissioned the last of my FreeBSD servers which had the required
enviroment to host a Tcl script which relied on [inetd], [tcpwrappers], and
[tcllauncher] to function. Sadly, `telnet macnugget.org` died that day. Only
the code remained as an archeological examination of how we used to do things
back in the olden days of UNIX before the rise of Linux and containers.

It was a different world [back then](https://github.com/nugget/rickroll/tree/02f031511578bc33fd5b3df10f857620042bc857).

It's old enough to have been 3-Clause BSD licensed. I remember I had strong
feelings about that at the time. I changed it to MIT as part of the rewrite.

# Version 2.0 Rewrite

In June 2025 I rewrote the server in Golang and got it containerized for an easier
deployment in modern infrastructure.

This service will bind port 23 and listen as a telnet server. Incoming
connections will be textually serenaded.

# Local Operation

You can build and run locally straight from this repo if you have Golang installed.

```sh
make run
```

# Containerized

Images are on dockerhub at [nugget/rickrolld](https://hub.docker.com/r/nugget/rickrolld)

```sh
docker run -p 23:23 nugget/rickrolld
```

Or use the `docker-compse.yml` file from this repo.

# Building the container itself

You can build/tag the rickrolld container locally tagged as rickrolld:dev

```sh
make container
```

```sh
make runcontainer
```

# Configuration

These environment variables can be set to override default values:

```.env
# string form of address (for example, "192.0.2.1:25", "[2001:db8::1]:80")
RICKROLL=LISTEN_ADDR=:23

# path to the lyrics file
RICKROLL_LYRICS_FILENAME=lyrics.dat
```

[Rickrolling]: https://en.wikipedia.org/wiki/Rickrolling
[inetd]: https://man.freebsd.org/cgi/man.cgi?inetd
[tcpwrappers]: https://en.wikipedia.org/wiki/TCP_Wrappers
[tcllauncher]: https://github.com/flightaware/tcllauncher
