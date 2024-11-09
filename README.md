# PGCrack

PGCrack is a tool made in go to **bruteforce** symmetrically **encrypted** `GPG` files.

⚠️ PGCrack is currently in a **beta phase**.

# Installation

## Download PGCrack Binary

Navigate to [releases](https://github.com/wdeloo/PGCrack/releases) and download the **latest version**.

## Build PGCrack from Source

### Clone the repository

```
git clone --depth 1 https://github.com/wdeloo/PGCrack.git
```

### Build it

```
cd PGCrack
mkdir dist
go build -o ./dist/pgcrack ./src/main.go
```

The **compiled binary** will be in `dist` directory.

# Usage

```
./pgcrack [mode] [parameters] encrypted.gpg
```

## Modes

### Wordlist
Get the passwords from a wordlist

*Usage: `./pgcrack -w wordlist.txt [parameters] encrypted.gpg`*

### Incremental
Bruteforce using all possible combinations

*Usage: `./pgcrack -i [parameters] encrypted.gpg`*

### Random
Use random passwords of a given length (`-l`)

*Usage: `./pgcrack -r -l <num> [parameters] encrypted.gpg`*

## Parameters

### Threads (*optional*)
Number of threads bruteforcing simultaneously

*Usage: `-t <number-of-threads>`*

### Charset (*optional for incremental and random modes*)
What characters will be used to compose the passwords

*Usage: `-c <string-containing-characters>`*

*Example: `-c "abc123"` will use 'a', 'b', 'c', '1', '2' and '3' characters*

### Length (**required for random** and *optional for incremental* modes)
Passwords length

*Usage: `-l <length>`*

### Minimum/Maximum Length (*optional for incremental mode*)
Minimum/Maximum length of the passwords

*Usage: `--min-length <length>` / `--max-length <length>`*