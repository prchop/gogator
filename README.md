# `gogator`

**gogator** is a simple RSS feed aggregator written in Go.

##  Features

* Add RSS feeds from across the internet to be collected
* Store the collected posts in a PostgreSQL database
* Follow and unfollow RSS feeds that other users have added
* View summaries of the aggregated posts in the terminal

## Get The Source Code

Download or clone the repository:

- Clone it using Git:

```bash
git clone https://github.com/prchop/gogator.git
```

Or download it as a ZIP from GitHub and extract it.

## Installation

### 1. Install Go

You’ll need Go 1.25 or higher.

Check if you already have it:

```bash
go version
```

If not installed, visit: https://go.dev/dl/ and follow the installation instructions for your OS.

After installation, ensure Go binaries are available:

* Make sure `$(go env GOPATH)/bin or $GOBIN` is included:

```bash
echo $PATH
```

* Add the go binary path to the `$PATH`

> go binary default path is `$HOME/go/bin or %USERPROFILE%\go\bin.`

```bash
export PATH="$PATH:$(go env GOPATH)/bin
```

* Then install using `go install`:

```bash
cd gogator
go install ./cmd/gogator
```

### 2. Install PostgreSQL

Check if you already have it:

```bash
psql --version
```

If not installed, visit: https://www.postgresql.org/download/ and follow the installation instructions for your OS.

## Configuration

### 1. Run migrations (using Goose)

* Install [goose](https://github.com/pressly/goose?tab=readme-ov-file#install)

```bash
cd gogator/sql/schema
goose postgres "postgres://username:password@localhost:5432/dbname" up
```

### 2. Configure the config file

This project uses a `.gogatorconfig.json` file in `$HOME` directory. Before running the app, create file in your `$HOME` directory and set similarly to `.gogatorconfig.json.example`. You can left the `username` blank, but you must set the `db_url` field with your PostgreSQL connection string:

```
postgres://username:password@host:port/dbname
```

Example:

```
postgres://postgres:secret123@localhost:5432/gator
```

## Usage

### Running commands

Use the following syntax:

```bash
gogator <command> [args...]
```

### Available commands

#### Authentication

| Command               | Description                                        |
| --------------------- | -------------------------------------------------- |
| `register <username>` | Register a new user                                |
| `login <username>`    | Log in as an existing user                         |
| `users`               | List all users and show who is currently logged in |

#### Feeds

| Command                | Description                                    |
| ---------------------- | ---------------------------------------------- |
| `feeds`                | List all available feeds                       |
| `addfeed <name> <url>` | Add a new feed to the aggregator and follow it |
| `follow <url>`         | Follow an existing feed                        |
| `unfollow <url>`       | Unfollow a feed                                |
| `following`            | Show all feeds the current user is following   |

#### Aggregation

| Command              | Description                                                     |
| -------------------- | --------------------------------------------------------------- |
| `agg <intervals>`    | Start the background aggregator that fetches feeds at intervals |
| `browse <limit>`     | Browse aggregated posts (default limit to 2)                    |

## Example

```bash
gogator register bob
gogator login bob
gogator addfeed "Hacker News RSS" "https://hnrss.org/newest"
gogator follow "https://hnrss.org/newest"
gogator unfollow "https://hnrss.org/newest"
gogator following
gogator agg 2s
```

## Notes

* Make sure PostgreSQL is running before you start the app.
* The aggregator will periodically fetch and store new posts in the background.
* You can modify configuration inside `.gogatorconfig.json`.

## Technologies Used

* [Go](https://go.dev/) — Backend and CLI
* [PostgreSQL](https://www.postgresql.org/) — Database
* [Goose](https://github.com/pressly/goose) — Migrations
* [sqlc](https://docs.sqlc.dev/en/latest/index.html) — Type-safe SQL to Go code
