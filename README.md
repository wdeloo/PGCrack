# PGCrack

PGCrack is a tool made in go to **bruteforce** symmetrically **encrypted** `GPG` files.

⚠️ PGCrack is currently in a **beta phase**.

## Installation

### Download PGCrack Binary

Navigate to [releases](https://github.com/wdeloo/PGCrack/releases) and download the **latest version**.

### Build PGCrack from Source

#### Clone the repository

```
git clone --depth 1 https://github.com/wdeloo/PGCrack.git
```

#### Build it

```
cd PGCrack
mkdir dist
go build -o ./dist/pgcrack ./src/main.go
```

The **compiled binary** will be in `dist` directory.

## Usage

```
./pgcrack [mode] [parameters] encrypted.gpg
```

### Modes

#### Wordlist
`-w wordlist` Get the passwords from a wordlist

*Example: `./pgcrack -w wordlist.txt [parameters] encrypted.gpg`*

#### Random
`-r` Use random passwords of a given length (`-l`)

*Example: `./pgcrack -r -l <num> [parameters] encrypted.gpg`*

### Parameters
`-t number-of-threads` (*default: 1*)

`-l password-length` (**required for random mode**)
