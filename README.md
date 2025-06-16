# Trolling telnet (ft. rickrolld)

```sh
telnet macnugget.org
```

If your machine doesn't have telnet installed, good for you. You can also use
Netcat to connect.

```sh
nc macnugget.org 23
```

# What is this? Whose fault is this?

I honestly don't remember what exactly was going through my head in 2013 when I
wrote the original version of this service. I imagine that I felt then -- as I
do now -- that rickrolld represents the zenith of [Rickrolling] humor.

rickrolld was written in Tcl 8.5 and was bespoke to my favored-at-the-time
FreeBSD environment. The script relied on [inetd], [tcpwrappers], and
[tcllauncher] to function. In 2022 I decommissioned the last of my those
FreeBSD servers and lost the enviroment I needed to keep rickrolld running.
Sadly, rickrolld died that day.

## Version 2.0 Rewrite

In June 2025 I rewrote the server in Golang and got it containerized for an
easier deployment in modern infrastructure. Special thanks to [Michael Hazell]
for reminding me to get it finished.

[The code from back then] is a time capsule from a different era. From before
Linux and containers took over the Internet. It's old enough to have been
3-Clause BSD licensed. I remember having strong feelings about that at the
time. I changed it to MIT as part of the rewrite.

# Description

This service will bind port 23 and listen as a telnet server. Incoming
connections will be textually serenaded.

# Local Operation

You can build and run locally straight from this repo if you have Golang installed.

```sh
make run
```

# Containerized

Images are on dockerhub at
[nugget/rickrolld](https://hub.docker.com/r/nugget/rickrolld)

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
RICKROLL_LISTEN_ADDR=:23

# path to the lyrics file
RICKROLL_LYRICS_FILENAME=lyrics.dat
```

[The code from back then]: https://github.com/nugget/rickroll/tree/02f031511578bc33fd5b3df10f857620042bc857
[Rickrolling]: https://en.wikipedia.org/wiki/Rickrolling
[inetd]: https://man.freebsd.org/cgi/man.cgi?inetd
[tcpwrappers]: https://en.wikipedia.org/wiki/TCP_Wrappers
[tcllauncher]: https://github.com/flightaware/tcllauncher
[Michael Hazell]: https://github.com/Techman
