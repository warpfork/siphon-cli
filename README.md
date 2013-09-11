siphon-cli
==========

Siphon is a program for creating psuedoterminals.  The two core commands are `siphon host` (which hosts a new process, say, `bash`), and `siphon attach`, which attaches the input and output of your current terminal to the hosted process.
  
Siphon works over a unix socket, so it can work on any systems that share a filesystem that supports unix sockets.  (For example, you can create a [docker](https://www.docker.io/) container with a shared directory mounted between host and container, then use Siphon to create new shells inside the container.)

Siphon can serve much the same role as`ssh`, spawning new terminals (or other processes) in response to connections, but without any of the complexity of networking or authentication.  Siphon works over unix socket files, so you can create as many Siphon hosts as desired, and authentication is simply a function of the permissions on that socket file.

Siphon is also comparable to the role of the well-known `screen` command, but whereas `screen` requires full visibility to your process tree for... reasons, Siphon works purely over a unix socket, thus being usable in a wider range of scenarios.  Siphon is also drastically less complex than `screen` or `tmux`; Siphon isn't trying to be a window manager, it's just a passthrough for a terminal.  `screen` and `tmux` can be used *inside* Siphon if you desire their features but need a terminal over a socket.



Command reference
-----------------

Siphon has three commands: `siphon host`, `siphon attach`, and `siphon daemon`.  You can append `--help` to `siphon` or any of its subcommands (i.e., `siphon host --help`) to see a list of options and their descriptions.


### siphon host

`siphon host` hosts a new process in a new psuedoterminal, opening a new socket which can be used by `siphon attach` to provide input and output to the hosted process.

By default, `siphon host` spawns a new shell (`/bin/sh`) in the current working directory, and makes a unix socket file called `siphon.sock`.

#### options

 * `-C` lets you specify the command to host.  For example, `siphon host -C /bin/bash` hosts a bash shell instead.  Any command that can be found on your $PATH is acceptable.
 * `-L` lets you specify where the socket should be created.  For example, `siphon host -L unix:///var/siphon.sock` forces the socket to be located in `/var` instead of the current directory.


### siphon attach

`siphon attach` connects to the socket created by `siphon host`, and shuttles input from your terminal to the hosted process, and output from the hosted process back to your terminal.

Multiple `siphon attach` commands can connect to a single `siphon host`.  Input from each of them will be accepted, and output will be copied to all.  The size of the host terminal may change when additional terminals are attached; this may cause confusing behavior, and so while supported, is not necessarily recommended.

#### options

 * `-L` lets you specify where to look for a socket to attach to.  The format is exactly the same as for the `-L` option on `siphon host`.


### siphon daemon

`siphon daemon` allows new hosts to be spawned.  It also creates a unix socket, which you connect to using `siphon attach` exactly the same way.  The difference is, instead of hosting a process directly, `siphon daemon` launches a new `siphon host` every time a siphon client attaches, and transparently redirects the client to the new host.

This can be used like the ssh daemon that's probably already on your server: every time you connect, you get a new shell.

#### options
 * `-L` lets you specify where the socket should be created.  The format is exactly the same as for the `-L` option on `siphon host`.
 * other options are passed through to the `siphon host` command:
   * `-C` lets you specify the command to host.  Passed on literally to a `siphon host` command.
   * `-H` lets you specify where the host socket should be created.  This is a pattern; the default is 'unix://siphon.#####.sock', which at runtime will be turned into 'unix://siphon.43560.sock' or some other random sequence.

`siphon daemon` quite literally exec's a new `siphon host` process, so `siphon` must be on your $PATH for daemon mode to work correctly.



Building
--------

Run `./go.build.sh`.  Assuming you have bash, a working install of `go`, and are on a *nix system of some kind, you should be good.



Usage example, from clone to execution
--------------------------------------

```
git clone https://github.com/heavenlyhash/siphon-cli
cd siphon-cli
./go.build.sh
./siphon host &
./siphon attach
```

You now have a working shell contained in a psuedoterminal.

Try using the `tty` command before `siphon attach`, and again afterwards.  You should see different answers!  These are the names of your psuedoterminals.



Usage example, with docker
--------------------------

```
# assuming you already have the siphon binary in this directory,
#  since it was built in the previous example
docker run -v $PWD:/shared -d ubuntu bash -c "cd /shared/ && ./siphon host"
siphon attach
```

In this example, we're just using `siphon host`, which means when that process exists, the docker container will exit.  This isn't much different than just running the shell as the docker command directly.  The magic part is if you want do be able to do multiple attaches, just switch `siphon host` for `siphon attach`, and great success ensues.

Much easier to attach a debug terminal this way than sucking an sshd server into your container and then bothering with key management, isn't it?  ;)

Note: when using `siphon daemon`, the socket names are passed to the client quite literally -- which means if your paths are different from depending on if it's the daemon's or the client's perspective, you will want to be sure to use relative paths instead of absolutes, so that the relative path string coincides on the same real location for the daemon and client.



Debug mode
----------

Additional logging statements can be enabled by setting an environment variable called "DEBUG".  Set it to "*" to enable all logging.

In bash, turning on debug mode for the siphon host process looks like this:

```
DEBUG='*' siphon host
```

(Note the single quotes!  In bash, single quotes will capture the asterisk exactly.  Double quotes will not!  Double quotes will cause the asterisk to be expanded to a list of files in the current directory.  This will not cause Siphon's debug statements to be enabled!)



Siphon as a Library
-------------------

Core features of Siphon are also available as a pure library (no main methods, no args parsing, just the good stuff): https://github.com/heavenlyhash/siphon

siphon-cli (this repo) is as thin as possible of a wrapper on top of the siphon library.  It adds args parsing, and also daemon mode.  (The siphon library contains support for redirects, so you can implement your own variations of "daemon mode" using the library.)



License
-------

Siphon is distributed under the Apache v2 license.  Enjoy freely.


