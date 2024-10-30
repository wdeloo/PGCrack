package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand/v2"
	"os"
	"os/exec"
	"time"
)

var startTime time.Time = time.Now()
var count int = 0
var prevCount int = 0

var finished bool

func tryDecrypt(file string, password string) {
	cmd := exec.Command("bash", "-c", "echo "+password+" | gpg --batch --passphrase-fd 0 --decrypt "+file)
	_, err := cmd.Output()

	if err == nil {
		fmt.Printf("\n[+] Password found: %s\n", password)
	}
}

func fileExists(fileName string) bool {
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err)
}

func statusBar() {
	for {
		fmt.Printf("UPTIME: %dh %dm %ds · TESTED PASSWORDS: %d · SPEED: %dp/s     \r",
			int(time.Since(startTime).Abs().Hours()),
			int(time.Since(startTime).Abs().Minutes())%60,
			int(time.Since(startTime).Abs().Seconds())%60,
			count,
			count-prevCount)

		prevCount = count

		if finished {
			fmt.Println("\n[x] Bruteforce concluded")
			return
		}

		time.Sleep(1 * time.Second)
	}
}

func doRandomBruteForce(filePath string, chars []rune, length int) {
	getRandomPassword := func(chars []rune, length int) string {
		var password []rune
		for i := 0; i < length; i++ {
			randomCharacter := chars[rand.IntN(len(chars))]
			password = append(password, randomCharacter)
		}
		return string(password)
	}

	for {
		password := getRandomPassword(chars, length)
		tryDecrypt(filePath, password)
		count++
	}
}

func doWordlistBruteForce(filePath string, scanner *bufio.Scanner) {
	for scanner.Scan() {
		password := scanner.Text()
		tryDecrypt(filePath, password)
		count++
	}
	if !scanner.Scan() {
		finished = true
	}
}

func getMode(modes map[string]any) string {
	var mode string
	for k, v := range modes {
		if v != "" && v != false {
			if mode == "" {
				mode = k
			} else {
				fmt.Fprintln(os.Stderr, "Only one mode can be specified")
				fmt.Fprintf(os.Stderr, "\nExecute: %s --help to print usage\n", os.Args[0])
				os.Exit(1)
			}
		}
	}
	return mode
}

func getFileName(args []string) string {
	filePath := args[0]

	if len(flag.Args()) == 0 {
		fmt.Fprintln(os.Stderr, "Missing argument: encrypted \".gpg\" file")
		fmt.Fprintf(os.Stderr, "\nExecute: %s --help to print usage\n", os.Args[0])
		os.Exit(1)
	}

	if len(flag.Args()) > 1 {
		fmt.Fprintln(os.Stderr, "Too many arguments")
		fmt.Fprintf(os.Stderr, "\nExecute: %s --help to print usage\n", os.Args[0])
		os.Exit(1)
	}

	if !fileExists(filePath) {
		fmt.Fprintf(os.Stderr, "%s: No such file or directory\n", filePath)
		fmt.Fprintf(os.Stderr, "\nExecute: %s --help to print usage\n", os.Args[0])
		os.Exit(1)
	}

	return filePath
}

func checkFlags(mode string, threads int, length int) {
	switch mode {
	case "random":
		if length < 1 {
			fmt.Fprintln(os.Stderr, "Missing or error in parameter \"-l\": password length is required and must be at least 1")
			fmt.Fprintf(os.Stderr, "\nExecute: %s --help to print usage\n", os.Args[0])
			os.Exit(1)
		}
	case "wordlist":
		if length != 0 {
			fmt.Fprintln(os.Stderr, "Error in parameter \"-l\": cannot set length in wordlist mode")
			fmt.Fprintf(os.Stderr, "\nExecute: %s --help to print usage\n", os.Args[0])
			os.Exit(1)
		}
	}

	if mode == "" {
		fmt.Fprintln(os.Stderr, "No mode specified")
		fmt.Fprintf(os.Stderr, "\nExecute: %s --help to print usage\n", os.Args[0])
		os.Exit(1)
	}

	if threads < 1 {
		fmt.Fprintln(os.Stderr, "Error in parameter \"-t\": number of threads must be at least 1")
		fmt.Fprintf(os.Stderr, "\nExecute: %s --help to print usage\n", os.Args[0])
		os.Exit(1)
	}
}

func main() {
	threads := flag.Int("t", 1, "")

	random := flag.Bool("r", false, "")
	length := flag.Int("l", 0, "")
	help := flag.Bool("help", false, "")

	wordlist := flag.String("w", "", "Specify password wordlist")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\nExecute: %s --help to print usage\n", os.Args[0])
	}

	showHelp := func() {
		fmt.Printf("Usage: %s [mode] [parameters] encrypted.gpg\n", os.Args[0])
		fmt.Println("")
		fmt.Println("# Parameters")
		fmt.Println("  -t <num>   Threads: Number of threads bruteforcing simultaneously")
		fmt.Println("  -l <num>   Length: Length of the password (random mode)")
		fmt.Println("")
		fmt.Println("  --help     Help: Show this help pannel")
		fmt.Println("")
		fmt.Println("# Modes")
		fmt.Println("  -w <file>  Wordlist: Bruteforce using a wordlist")
		fmt.Printf("     Example: %s -w wordlist.txt [parameters] encrypted.gpg\n", os.Args[0])
		fmt.Println("")
		fmt.Println("  -r         Random: Bruteforce using random passwords of a given length (-l)")
		fmt.Printf("     Example: %s -r -l <num> [parameters] encrypted.gpg\n", os.Args[0])
	}

	flag.Parse()

	if *help {
		showHelp()
		os.Exit(0)
	}

	modes := map[string]any{
		"random":   *random,
		"wordlist": *wordlist,
	}

	mode := getMode(modes)

	filePath := getFileName(flag.Args())
	checkFlags(mode, *threads, *length)

	chars := []rune{
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'.', '-', '_', ',',
	}

	fmt.Print("\033[?25l")
	switch mode {
	case "random":
		for i := 0; i < *threads; i++ {
			go doRandomBruteForce(filePath, chars, *length)
		}
	case "wordlist":
		file, err := os.Open(*wordlist)
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "%s: No such file or directory\n", *wordlist)
			fmt.Fprintf(os.Stderr, "\nExecute: %s --help to print usage\n", os.Args[0])
			os.Exit(1)
		}

		defer file.Close()

		scanner := bufio.NewScanner(file)
		for i := 0; i < *threads; i++ {
			go doWordlistBruteForce(filePath, scanner)
		}
	}

	statusBar()
}
