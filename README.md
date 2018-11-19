# lp-bug-notifier

The purpose of this application is to be a minimal, terminal-based tool for getting updates on Launchpad bugs of interest. 

1. Add a bug to follow with `add <bug #>`, and the application will poll LP for updates on that bug.
2. Bug title and message age are pulled and displayed, and updated every five minutes
3. Drop a bug you no longer care about with `drop <bug #>`

That's it.

## Install and Run Instructions

### Install direnv

1. `$ sudo apt install direnv`
2. add the following line to your .bashrc:
    `eval "$(direnv hook bash)"`
3. `cd` to project directory
4. `$ echo 'layout go' > .envrc`
    This should prompt the following message:
    > direnv: error .envrc is blocked. Run `direnv allow` to approve its content.
5. `direnv allow`

### Install project & dependencies

`go get main`

### Run

`go run main`

### Planned features

- Adjust refresh rate
- Nicer display, something ncurses-y
- Show bug status, latest message, and/or bug page link
