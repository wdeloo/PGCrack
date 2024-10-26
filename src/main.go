package main

import (
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

func tryDecrypt(file string, password string) {
	cmd := exec.Command("bash", "-c", "echo "+password+" | gpg --batch --passphrase-fd 0 --decrypt "+file)
	_, err := cmd.Output()

	if err == nil {
		fmt.Printf("\nPassword found: %s\n", password)
	}
}

func getRandomPassword(chars []rune, length int) string {
	var password []rune
	for i := 0; i < length; i++ {
		randomCharacter := chars[rand.IntN(len(chars))]
		password = append(password, randomCharacter)
	}
	return string(password)
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
		time.Sleep(1 * time.Second)
	}
}

func doBruteForce(filePath string, chars []rune, length int) {
	for {
		password := getRandomPassword(chars, length)
		tryDecrypt(filePath, password)
		count++
	}
}

func main() {
	length := flag.Int("l", 0, "Password length (required)")
	threads := flag.Int("t", 1, "Number of threads running simultaniously")

	flag.Usage = func() {
		fmt.Printf("Usage: %s [parameters] encrypted.gpg\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Fprintln(os.Stderr, "Missing argument: encrypted \".gpg\" file")
		os.Exit(1)
	}

	if len(flag.Args()) > 1 {
		fmt.Fprintln(os.Stderr, "Too many arguments")
		os.Exit(1)
	}

	filePath := flag.Args()[0]
	if !fileExists(filePath) {
		fmt.Fprintf(os.Stderr, "%s: No such file or directory\n", filePath)
		os.Exit(1)
	}

	if *threads < 1 {
		fmt.Fprintln(os.Stderr, "Error in parameter \"-t\": number of threads must be at least 1")
		os.Exit(1)
	}

	if *length < 1 {
		fmt.Fprintln(os.Stderr, "Missing or error in parameter \"-l\": password length is required and must be at least 1")
		os.Exit(1)
	}

	chars := []rune{
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'.', '-', '_', ',',
	}

	fmt.Print("\033[?25l")
	for i := 0; i < *threads; i++ {
		go doBruteForce(filePath, chars, *length)
	}

	statusBar()
}
