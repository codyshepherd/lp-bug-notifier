# lp-bug-notifier

The purpose of this application is to be a minimal, terminal-based tool for getting updates on Launchpad bugs of interest. 

1. Add a bug to follow with `add <bug #>`, and the application will poll LP for updates on that bug.
2. Bug title and message age are updated and displayed and updated every five minutes
3. Drop a bug you no longer care about with `drop <bug #>`

That's it.

## Install and Run Instructions

Assuming your $GOPATH and $PATH are configured correctly (a la direnv), from top level in the project
(i.e. `ls` reveals `src/` and `README.md`):

`go run main`

### Planned features

- Adjust refresh rate
- Nicer display, something ncurses-y
- Show bug status, latest message, and/or bug page link
