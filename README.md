# Giterm
[![Build Status](https://travis-ci.org/bharath-srinivas/giterm.svg?branch=master)](https://travis-ci.org/bharath-srinivas/giterm)
[![Go Report Card](https://goreportcard.com/badge/github.com/bharath-srinivas/giterm)](https://goreportcard.com/report/github.com/bharath-srinivas/giterm)

A terminal client for Github :sunglasses:

## Installation
### Installing from source
If you want to build the application yourself, do the following from within your `$GOPATH`:
```bash
$ go get -u github.com/bharath-srinivas/giterm
$ cd $GOPATH/src/github.com/bharath-srinivas/giterm
$ go install ./cmd/giterm
```

## Setup
Giterm requires github `personal access token` to work. For more information on how to create a
personal access token, you can refer to the official documentation [here](https://help.github.com/en/github/authenticating-to-github/creating-a-personal-access-token-for-the-command-line).

When running giterm for the first time, you have to set your personal access token before using
the application.

1. Create you personal access token by clicking [here](https://github.com/settings/tokens).
Refer to the [Required scopes](#required-scopes) section when creating the token
3. Copy your personal access token and keep it ready
4. Now set the token by executing the following command:
```bash
$ giterm --token [your_personal_access_token]
```
This will create a new config file at `~/.giterm/config.json`.

### Required scopes
Giterm requires the following scopes to be selected when creating a personal access token for
it to function properly.

| Scope             | Description                                     |
| ----------------- | ----------------------------------------------- |
| `repo`            | Full control of private repositories            |
| `read:packages`   | Download packages from github package registry  |
| `read:org`        | Read org and team membership, read org projects |
| `notifications`   | Access notifications                            |
| `user`            | Update all user data                            |
| `read:discussion` | Read team discussions                           |

### Feeds (optional)
If you wish to view your github news feeds along with other data, you have to set your
private github news feed URL.

To get your private github news feed URL, execute the following command:
```bash
$ curl -u "username" https://api.github.com/feeds | grep -oP '(?<="current_user_url": ")[^"]*'
```

Running the above command will ask you to enter your github password. After getting the URL,
you can set it using the following command:
```bash
$ giterm --feeds-url [your_github_feed_url]
```

## Usage

Once everything is setup properly, you can run the application with the following command:
```bash
$ giterm
```

## Default key mappings

| Key              | Action            |
| ---------------- | ----------------- |
| `Ctrl-n`         | Next page         |
| `Ctrl-p`         | Previous page     |
| `Tab`            | Next widget       |
| `Ctrl-Tab`       | Previous widget   |
| `Ctrl-r`         | Refresh page      |
| `Left`           | Scroll left       |
| `Right`          | Scroll right      |
| `Up/Page up`     | Scroll up         |
| `Down/Page down` | Scroll down       |
| `Home`           | Scroll to top     |
| `End`            | Scroll to bottom  |
| `Enter`          | Confirm           |
| `Ctrl-q`         | Quit              |