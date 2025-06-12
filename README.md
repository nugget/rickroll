# Trolling telnet.

You can `telnet macnugget.org` to see it in action.

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

# Version 2.0 Rewrite

In June 2025 I rewrote the sever in Golang and got it containerized for easy
deployment in a manner that's better aligned with current best practices.

This service (and docker container) will bind port 23 and listen as a telnet server.
Incoming connections will be handled and enjoy the song of our people.

# Installation

You can run it locally from this repo if you have a relatively modern Golang installed.

```sh
make run
```

Or you can run my published [rickrolld](https://hub.docker.com/r/nugget/rickrolld) container:

```sh
docker run nugget/rickrolld
```

Or use the `docker-compse.yml` file from this repo.

[Rickrolling]: https://en.wikipedia.org/wiki/Rickrolling
[inetd]: https://man.freebsd.org/cgi/man.cgi?inetd
[tcpwrappers]: https://en.wikipedia.org/wiki/TCP_Wrappers
[tcllauncher]: https://github.com/flightaware/tcllauncher
