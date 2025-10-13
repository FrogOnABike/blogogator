# Blogogator

A Go/Postgres-based RSS browser

Here’s my take on the [Boot.dev](http://boot.dev) Blog Aggregator project - although I changed up the name they'd suggested from gator :)

It’s a CLI-based app, built in Golang with a Postgres backend, which allows for multiple users to follow RSS feeds and view the post titles in the terminal.
## Requirements
### Golang

The app was written with Go v1.25, so it will require that version at a minimum to run on your computer.  
Installation instructions can be found [here](https://go.dev/doc/install)
### Postgres

This is also required as the backend database for the app.  
Installation instructions for it can be found [here](https://www.postgresql.org/download/)
Once installed you will need to create a database called "gator"

Once you have Postgres installed you’ll need to create a database for the app to use and obtain
### Setup config file

The app expects a .gatorconfig.json file in your home directory which will need to initially contain the following contents:  

`{`

  `"db_url": "protocol://username:password@localhost:5432/gator"`

`}`

You would need to replace some sections in the above connection string depending on your OS and setup. MacOS users will probably just need username to connect to the db, Linux/WSL users would need to use "postrgres:password" where password was the one set during the Postgres installation.

## Installation

Once you have the repo cloned locally, you can use [go install](https://pkg.go.dev/cmd/go#hdr-Compile_and_install_packages_and_dependencies) to install the program locally
## Commands

The basic structure for use is:

`blogogator *command* **arguments**`

Not all commands need arguments :)
Here is the current selection, along with expected arguments.

login **username** - Login as given user

register **username** - Register given user (required before login!)

reset - Will reset ALL users, feeds etc. USE CAREFULLY!!

users - List all registered users, along with displaying the currently logged in user

addfeed **feed name** *url* - Will add a given RSS feed to the database and store it with "feed name"
*Adding a feed will also follow it for you

feeds - List of ALL feeds in the database, along with which user added them

follow *url* - Follow the given RSS feed

following - Display all feeds you currently follow

unfollow *url* - Stops following the given RSS feed

agg **time interval** - Will start retrieving all feeds for the logged in user at **time interval** periods. This can be specified in the following formats - s for seconds, m for minutes h for hours, so 15m would be "every 15 minutes"

browse **optional number** - Will retrieve the newest X posts if a number is given, otherwise displays the 2 most recent posts from the database

Possible future improvements:

Extending the Project
You've done all the required steps, but if you'd like to make this project your own, here are some ideas:

Add sorting and filtering options to the browse command
Add pagination to the browse command
Add concurrency to the agg command so that it can fetch more frequently
Add a search command that allows for fuzzy searching of posts
Add bookmarking or liking posts
Add a TUI that allows you to select a post in the terminal and view it in a more readable format (either in the terminal or open in a browser)
Add an HTTP API (and authentication/authorization) that allows other users to interact with the service remotely
Write a service manager that keeps the agg command running in the background and restarts it if it crashes

