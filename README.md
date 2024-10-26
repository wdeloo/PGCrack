# PGCrack

PGCrack is a tool made in go to **bruteforce** symmetrically **encrypted** `GPG` files.

PGCrack is currently in a **beta phase**, it only supports **random** passwords for the momment.

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

`./PGCrack [parameters] encrypted.gpg`

### Parameters
`-l password-length` (**required**)

`-t number-of-threads` (*optional*)