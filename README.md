Break Out Of Hareka's Underground (Boohu) is a roguelike game mainly inspired from DCSS and its tavern, with some ideas from Brogue, but aiming for very short games, almost no character building, and a simplified inventory.

*Every year, a hero is chosen by the elders to collect medicinal simellas plants in the Underground. This year, the honor fell on you. You have heard rumours of vile creatures, which you may fight or try to escape. Various items will help you along the way, in your search for a way back to your village.*

![Boohu introduction screen](https://raw.githubusercontent.com/anaseto/boohu-pages/master/intro-screen.png)

Screenshot
----------

[![Introduction Screeshot](https://raw.githubusercontent.com/anaseto/boohu-pages/master/screenshot.png)](https://anaseto.github.io/boohu-pages/)

Follow the link on the image to see a little introduction screencast and some
short animations.

Install
-------

You can download binaries on the [releases
page](https://github.com/anaseto/boohu/releases).

You can also build from source by following these steps:

+ Install the [go compiler](https://golang.org/).
+ Set `$GOPATH` variable (for example `export GOPATH=$HOME/go`).
+ Add `$GOPATH/bin` to your `$PATH` (for example `export PATH="$PATH:$GOPATH/bin"`).
+ Use the command `go get -u github.com/anaseto/boohu`.
  
The `boohu` command should now be available.

The only dependency outside of the go standard library is the lightweight
curses-like library [termbox-go](https://github.com/nsf/termbox-go), which is
installed automatically by the previous `go get` command.

*Portability note.* If you happen to experience input problems, try adding
option `--tags tcell` or `--tags ansi` to the `go get` command. The first will use
[tcell](https://github.com/gdamore/tcell) instead of termbox-go, and requires
cgo on some platforms, but is more portable. The second will work on POSIX
systems with a `stty` command.

Colors
------

If the default colors do not display nicely on your terminal emulator, you can
use the `-s` option: `boohu -s` to use the 16-color palette, which
will display nicely if the [solarized](http://ethanschoonover.com/solarized)
palette is used. Configurations are available for most terminal emulators, otherwise, colors may have to be configured manually to one's liking in
the terminal emulator options.

Basic Survival Tips
-------------

+ Position yourself to fight one enemy at a time whenever possible.
+ Fight far enough from unknown terrain if you can: combat is noisy, more
  monsters will come if they hear it.
+ Use your potions, projectiles and rods. With experience you learn when you
  can spare them, and when you should use them, but dying with a potion of heal
  wounds in the inventory should never be a thing.
+ Avoid dead-ends if you do not have any means of escape, such as potions of
  teleportation, unless no better options are available.
+ Use *pillar dancing*: sometimes you can turn around a block several times to
  avoid being killed while replenishing your HP.
+ You do not have to kill every monster. You want, though, to find as many items
  as you can, but survival comes first.
+ Use doors and dense foliage to break line of sight with monsters.
