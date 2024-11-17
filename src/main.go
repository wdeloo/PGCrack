package main

import (
	"PGCrack/src/modes/incremental"
	"PGCrack/src/modes/random"
	"PGCrack/src/modes/wordlist"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
)

var paramHelpErrorMessage = fmt.Sprintf("\nexecute [ %s --help ] to print usage", os.Args[0])

var startTime time.Time = time.Now()
var count int = 0
var prevCount int = 0
var foundPasswords []string

func exitWithError(errorMessage string) {
	fmt.Fprintln(os.Stderr, errorMessage)
	fmt.Fprintln(os.Stderr, paramHelpErrorMessage)
	os.Exit(1)
}

func catchCtrlC() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-signals

	fmt.Print("\r")
	printStatus()
	fmt.Print("\n[x] Bruteforce interrupted\n\n")
	printSummary()
	fmt.Print("\033[?25h") // show cursor
	os.Exit(1)
}

func tryDecrypt(file string, password string) {
	encFile, _ := os.ReadFile(file)

	pgp := crypto.PGP()

	decHandle, _ := pgp.Decryption().Password([]byte(password)).New()
	_, err := decHandle.Decrypt(encFile, crypto.Bytes)

	if err == nil {
		fmt.Printf("\n[+] Password found: \033[7m%s\033[0m\n", password)
		foundPasswords = append(foundPasswords, password)
	}
}

func fileExists(fileName string) bool {
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err)
}

func printSummary() {
	fmt.Println("----------------------------------------------------------")

	fmt.Printf("              TOTAL TIME: %d hours %d minutes %d seconds\n", int(time.Since(startTime).Abs().Hours()), int(time.Since(startTime).Abs().Minutes())%60, int(time.Since(startTime).Abs().Seconds())%60)
	fmt.Printf("  TOTAL TESTED PASSWORDS: %d\n", count)
	fmt.Printf("           AVERAGE SPEED: %s passwords/second\n", fmt.Sprintf("%.2f", (float64(count)/float64(time.Since(startTime).Abs().Milliseconds()))*1000))

	fmt.Println("")

	if len(foundPasswords) == 0 {
		fmt.Println("      NO PASSWORDS FOUND")
	} else {
		for _, password := range foundPasswords {
			fmt.Printf("          FOUND PASSWORD: \033[7m%s\033[0m\n", password)
		}
	}

	fmt.Println("----------------------------------------------------------")
}

func printStatus() {
	fmt.Printf("UPTIME: %dh %dm %ds · TESTED PASSWORDS: %d · SPEED: %dp/s     \r",
		int(time.Since(startTime).Abs().Hours()),
		int(time.Since(startTime).Abs().Minutes())%60,
		int(time.Since(startTime).Abs().Seconds())%60,
		count,
		count-prevCount)
}

func statusBar() {
	fmt.Print("\033[?25l") // hide cursor
	for {
		printStatus()
		prevCount = count

		time.Sleep(1 * time.Second)
	}
}

func getMode(modes map[string]any) string {
	var mode string
	for k, v := range modes {
		if v != "" && v != false {
			if mode == "" {
				mode = k
			} else {
				exitWithError("only one mode can be specified")
			}
		}
	}
	return mode
}

func getFileName(args []string) string {
	if len(flag.Args()) == 0 {
		exitWithError("missing argument: encrypted \".gpg\" file")
	}

	if len(flag.Args()) > 1 {
		exitWithError("too many arguments")
	}

	filePath := args[0]

	if !fileExists(filePath) {
		exitWithError(fmt.Sprintf("%s: no such file or directory\n", filePath))
	}

	return filePath
}

func checkFlags(mode string, threads int, length int, minLength int, maxLength int, characters string) {
	switch mode {
	case "random":
		if length < 1 {
			exitWithError("missing or error in parameter \"-l\": password length is required for random mode and must be at least 1")
		}
		if minLength != 0 || maxLength != 0 {
			exitWithError("error in parameter \"--min/max-length\": password min/max length can not be used in random mode, use length (-l) instead")
		}
	case "wordlist":
		if length != 0 {
			exitWithError("error in parameter \"-l\": cannot set length in wordlist mode")
		}
		if minLength != 0 || maxLength != 0 {
			exitWithError("error in parameter \"--min/max-length\": cannot set password min/max length in wordlist mode")
		}
		if characters != "" {
			exitWithError("error in parameter \"-c\": cannot set charset in wordlist mode")
		}
	case "incremental":
		if length != 0 && (minLength != 0 || maxLength != 0) {
			exitWithError("error in parameter \"-l, --min/max-length\": cannot set password length and password min/max length at the same time")
		}
		if length < 0 {
			exitWithError("error in parameter \"-l\": password length must be at least 1")
		}
		if minLength < 0 || maxLength < 0 {
			exitWithError("error in parameter \"--min/max-length\": password min/max length must be at least 1")
		}
		if minLength > maxLength && maxLength != 0 {
			exitWithError("error in parameter \"--min/max-length\": min length cannot be greater than max length")
		}
	}

	if mode == "" {
		exitWithError("error: no mode specified")
	}

	if threads < 1 {
		exitWithError("error in parameter \"-t\": number of threads must be at least 1")
	}
}

func main() {
	go catchCtrlC()

	randomMode := flag.Bool("r", false, "")
	wordlistMode := flag.String("w", "", "")
	incrementalMode := flag.Bool("i", false, "")
	threads := flag.Int("t", 1, "")
	charset := flag.String("c", "", "")
	length := flag.Int("l", 0, "")
	minLength := flag.Int("min-length", 0, "")
	maxLength := flag.Int("max-length", 0, "")
	help := flag.Bool("help", false, "")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, paramHelpErrorMessage)
	}

	showHelp := func() {
		fmt.Printf("Usage: %s [mode] [parameters] encrypted.gpg\n", os.Args[0])
		fmt.Println("")
		fmt.Println("# Modes")
		fmt.Println("  -w <file>           Wordlist: Get the passwords from a wordlist")
		fmt.Printf("     Usage: %s -w wordlist.txt [parameters] encrypted.gpg\n", os.Args[0])
		fmt.Println("")
		fmt.Println("  -i                  Incremental: Bruteforce using all possible combinations")
		fmt.Printf("     Usage: %s -i [parameters] encrypted.gpg\n", os.Args[0])
		fmt.Println("")
		fmt.Println("  -r                  Random: Bruteforce using random passwords of a given length (-l)")
		fmt.Printf("     Usage: %s -r -l <num> [parameters] encrypted.gpg\n", os.Args[0])
		fmt.Println("")
		fmt.Println("# Parameters")
		fmt.Println("  -t <num>            Threads: Number of threads bruteforcing simultaneously")
		fmt.Println("")
		fmt.Println("  -c <str>            Charset: What characters will be used to compose the passwords (random/incremental mode)")
		fmt.Println("     Default: all letters (capitalized and uncapitalized) and digits")
		fmt.Println("     Example: [ -c \"abc123\" ] will use 'a', 'b', 'c', '1', '2' and '3' characters")
		fmt.Println("")
		fmt.Println("  -l <num>            Length: Length of the password (random/incremental mode)")
		fmt.Println("")
		fmt.Println("  --min-length <num>  Minimum Length: Minimum length of the password (incremental mode)")
		fmt.Println("  --max-length <num>  Maximum Length: Maximum length of the password (incremental mode)")
		fmt.Println("")
		fmt.Println("  --help              Help: Show this help pannel")
	}

	flag.Parse()

	if *help {
		showHelp()
		os.Exit(0)
	}

	modes := map[string]any{
		"random":      *randomMode,
		"wordlist":    *wordlistMode,
		"incremental": *incrementalMode,
	}

	mode := getMode(modes)

	filePath := getFileName(flag.Args())
	checkFlags(mode, *threads, *length, *minLength, *maxLength, *charset)

	var chars []rune
	if *charset == "" {
		chars = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	} else {
		chars = []rune(*charset)
	}

	go statusBar()

	switch mode {
	case "random":
		random.DoRandomBruteForce(filePath, chars, *length, tryDecrypt, &count, *threads)

	case "wordlist":
		file, err := os.Open(*wordlistMode)
		if os.IsNotExist(err) {
			exitWithError(fmt.Sprintf("%s: no such file or directory\n", *wordlistMode))
		}

		wordlist.DoWordlistBruteForce(filePath, file, tryDecrypt, &count, *threads)

	case "incremental":
		if *length != 0 {
			*minLength = *length
			*maxLength = *length
		}

		incremental.DoIncrementalBruteForce(filePath, chars, *minLength, *maxLength, &count, *threads, tryDecrypt)
	}

	printStatus()
	fmt.Print("\n[x] Bruteforce concluded\n\n")
	printSummary()
	fmt.Print("\033[?25h") // show cursor
}
