package wordlist

import (
	"bufio"
	"os"
	"sync"
)

func DoWordlistBruteForce(filePath string, file *os.File, tryDecrypt func(string, string), count *int, threads int) {
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var wg sync.WaitGroup

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for scanner.Scan() {
				password := scanner.Text()
				tryDecrypt(filePath, password)
				*count++
			}
		}()
	}

	wg.Wait()
}
