package incremental

import (
	"sync"
)

func DoIncrementalBruteForce(filePath string, chars []rune, minLength int, maxLength int, count *int, threads int, tryDecrypt func(string, string)) {
	continueGenerating := func(length int, maxLength int) bool {
		if maxLength == 0 {
			return true
		} else {
			return length <= maxLength
		}
	}

	if minLength == 0 {
		minLength = 1
	}

	for length := minLength; continueGenerating(length, maxLength); length++ {
		passwords := make(chan string)
		var wg sync.WaitGroup

		var generate func(string, int)
		generate = func(prefix string, n int) {
			if n == 0 {
				passwords <- prefix
				return
			}

			for _, char := range chars {
				generate(prefix+string(char), n-1)
			}
		}

		go func() {
			generate("", length)
			close(passwords)
		}()

		for i := 0; i < threads; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for password := range passwords {
					tryDecrypt(filePath, password)
					*count++
				}
			}()
		}

		wg.Wait()
	}
}
