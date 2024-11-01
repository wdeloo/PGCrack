package random

import "math/rand/v2"

func getRandomPassword(chars []rune, length int) string {
	var password []rune
	for i := 0; i < length; i++ {
		randomCharacter := chars[rand.IntN(len(chars))]
		password = append(password, randomCharacter)
	}
	return string(password)
}

func DoRandomBruteForce(filePath string, chars []rune, length int, tryDecrypt func(string, string), count *int, threads int) {
	for i := 0; i < threads; i++ {
		go func() {
			for {
				password := getRandomPassword(chars, length)
				tryDecrypt(filePath, password)
				*count++
			}
		}()
	}

	select {}
}
