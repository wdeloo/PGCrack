package main

import (
	"PGCrack/src/modes/incremental"
	"PGCrack/src/modes/random"
	"PGCrack/src/modes/wordlist"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

var startTime time.Time = time.Now()
var count int = 0
var prevCount int = 0
var foundPasswords []string

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
	cmd := exec.Command("bash", "-c", "echo '"+password+"' | gpg --batch --passphrase-fd 0 --decrypt "+file)
	_, err := cmd.Output()

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
				fmt.Fprintln(os.Stderr, "Only one mode can be specified")
				fmt.Fprintf(os.Stderr, "\nExecute [ %s --help ] to print usage\n", os.Args[0])
				os.Exit(1)
			}
		}
	}
	return mode
}

func getFileName(args []string) string {
	if len(flag.Args()) == 0 {
		fmt.Fprintln(os.Stderr, "Missing argument: encrypted \".gpg\" file")
		fmt.Fprintf(os.Stderr, "\nExecute [ %s --help ] to print usage\n", os.Args[0])
		os.Exit(1)
	}

	if len(flag.Args()) > 1 {
		fmt.Fprintln(os.Stderr, "Too many arguments")
		fmt.Fprintf(os.Stderr, "\nExecute [ %s --help ] to print usage\n", os.Args[0])
		os.Exit(1)
	}

	filePath := args[0]

	if !fileExists(filePath) {
		fmt.Fprintf(os.Stderr, "%s: No such file or directory\n", filePath)
		fmt.Fprintf(os.Stderr, "\nExecute [ %s --help ] to print usage\n", os.Args[0])
		os.Exit(1)
	}

	return filePath
}

func checkFlags(mode string, threads int, length int, minLength int, maxLength int) {
	switch mode {
	case "random":
		if length < 1 {
			fmt.Fprintln(os.Stderr, "Missing or error in parameter \"-l\": password length is required and must be at least 1")
			fmt.Fprintf(os.Stderr, "\nExecute [ %s --help ] to print usage\n", os.Args[0])
			os.Exit(1)
		}
		if minLength != 0 || maxLength != 0 {
			fmt.Fprintln(os.Stderr, "Error in parameter \"--min/max-length\": password min/max length can not be used in random mode, use length (-l) instead")
			fmt.Fprintf(os.Stderr, "\nExecute [ %s --help ] to print usage\n", os.Args[0])
			os.Exit(1)
		}
	case "wordlist":
		if length != 0 {
			fmt.Fprintln(os.Stderr, "Error in parameter \"-l\": cannot set length in wordlist mode")
			fmt.Fprintf(os.Stderr, "\nExecute [ %s --help ] to print usage\n", os.Args[0])
			os.Exit(1)
		}
		if minLength != 0 || maxLength != 0 {
			fmt.Fprintln(os.Stderr, "Error in parameter \"--min/max-length\": cannot set password min/max length in wordlist mode")
			fmt.Fprintf(os.Stderr, "\nExecute [ %s --help ] to print usage\n", os.Args[0])
			os.Exit(1)
		}
	case "incremental":
		if length != 0 && (minLength != 0 || maxLength != 0) {
			fmt.Fprintln(os.Stderr, "Error in parameter \"-l, --min/max-length\": cannot set password length and min/max password length at the same time")
			fmt.Fprintf(os.Stderr, "\nExecute [ %s --help ] to print usage\n", os.Args[0])
			os.Exit(1)
		}
		if length < 0 {
			fmt.Fprintln(os.Stderr, "Error in parameter \"-l\": password length must be at least 1")
			fmt.Fprintf(os.Stderr, "\nExecute [ %s --help ] to print usage\n", os.Args[0])
			os.Exit(1)
		}
		if minLength < 0 || maxLength < 0 {
			fmt.Fprintln(os.Stderr, "Error in parameter \"--min/max-length\": password min/max length must be at least 1")
			fmt.Fprintf(os.Stderr, "\nExecute [ %s --help ] to print usage\n", os.Args[0])
			os.Exit(1)
		}
		if minLength > maxLength && maxLength != 0 {
			fmt.Fprintln(os.Stderr, "Error in parameter \"--min/max-length\": min length cannot be greater than max length")
			fmt.Fprintf(os.Stderr, "\nExecute [ %s --help ] to print usage\n", os.Args[0])
			os.Exit(1)
		}
	}

	if mode == "" {
		fmt.Fprintln(os.Stderr, "No mode specified")
		fmt.Fprintf(os.Stderr, "\nExecute [ %s --help ] to print usage\n", os.Args[0])
		os.Exit(1)
	}

	if threads < 1 {
		fmt.Fprintln(os.Stderr, "Error in parameter \"-t\": number of threads must be at least 1")
		fmt.Fprintf(os.Stderr, "\nExecute [ %s --help ] to print usage\n", os.Args[0])
		os.Exit(1)
	}
}

func main() {
	go catchCtrlC()

	randomMode := flag.Bool("r", false, "")
	wordlistMode := flag.String("w", "", "")
	incrementalMode := flag.Bool("i", false, "")
	threads := flag.Int("t", 1, "")
	length := flag.Int("l", 0, "")
	minLength := flag.Int("min-length", 0, "")
	maxLength := flag.Int("max-length", 0, "")
	help := flag.Bool("help", false, "")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\nExecute [ %s --help ] to print usage\n", os.Args[0])
	}

	showHelp := func() {
		fmt.Printf("Usage: %s [mode] [parameters] encrypted.gpg\n", os.Args[0])
		fmt.Println("")
		fmt.Println("# Modes")
		fmt.Println("  -w <file>  Wordlist: Bruteforce using a wordlist")
		fmt.Printf("     Example: %s -w wordlist.txt [parameters] encrypted.gpg\n", os.Args[0])
		fmt.Println("")
		fmt.Println("  -r         Random: Bruteforce using random passwords of a given length (-l)")
		fmt.Printf("     Example: %s -r -l <num> [parameters] encrypted.gpg\n", os.Args[0])
		fmt.Println("")
		fmt.Println("# Parameters")
		fmt.Println("  -t <num>   Threads: Number of threads bruteforcing simultaneously")
		fmt.Println("  -l <num>   Length: Length of the password (random mode)")
		fmt.Println("")
		fmt.Println("  --help     Help: Show this help pannel")
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
	checkFlags(mode, *threads, *length, *minLength, *maxLength)

	chars := []rune{
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'.', '-', '_', ',',
	}

	go statusBar()

	switch mode {
	case "random":
		random.DoRandomBruteForce(filePath, chars, *length, tryDecrypt, &count, *threads)

	case "wordlist":
		file, err := os.Open(*wordlistMode)
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "%s: No such file or directory\n", *wordlistMode)
			fmt.Fprintf(os.Stderr, "\nExecute [ %s --help ] to print usage\n", os.Args[0])
			os.Exit(1)
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
