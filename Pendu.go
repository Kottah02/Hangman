package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

const maxAttempts = 6

func main() {
	fmt.Println("=====Bienvenue dans le jeu du pendu!=====")

	err := startGame()
	if err != nil {
		fmt.Println("Erreur:", err)
	}
}

func startGame() error {
	word, err := getRandomWord()
	if err != nil {
		return err
	}

	discovered := make([]rune, len(word))
	for i := range discovered {
		discovered[i] = '_'
	}

	usedLetters := map[rune]bool{}
	attempts := 0

	for attempts < maxAttempts {
		printGameState(discovered, usedLetters, attempts)

		guess := getUserGuess()
		if usedLetters[guess] {
			fmt.Println("Vous avez déjà utilisé cette lettre.")
			continue
		}

		usedLetters[guess] = true

		if strings.ContainsRune(word, guess) {
			updateDiscovered(discovered, word, guess)
			if strings.Join(stringSlice(discovered), "") == word {
				fmt.Println("Félicitations, vous avez deviné le mot:", word)
				return nil
			}
		} else {
			attempts++
		}
	}

	fmt.Println("Vous avez perdu! Le mot était:", word)
	return nil
}

func updateDiscovered(discovered []rune, word string, guess rune) {
	for i, letter := range word {
		if letter == guess {
			discovered[i] = guess
		}
	}
}

func printGameState(discovered []rune, usedLetters map[rune]bool, attempts int) {
	fmt.Printf("\nMot: %s\n", strings.Join(stringSlice(discovered), " "))
	fmt.Printf("Lettres utilisées: %s\n", strings.Join(mapKeysToSlice(usedLetters), " "))
	fmt.Printf("Essais restants: %d\n", maxAttempts-attempts)
}

func getUserGuess() rune {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Entrez une lettre: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if len(input) != 1 {
		fmt.Println("Veuillez entrer une seule lettre.")
		return getUserGuess()
	}
	return rune(input[0])
}
func mapKeysToSlice(m map[rune]bool) []string {
	var keys []string
	for k := range m {
		keys = append(keys, string(k))
	}
	return keys
}

func stringSlice(runes []rune) []string {
	var result []string
	for _, r := range runes {
		result = append(result, string(r))
	}
	return result
}

func getRandomWord() (string, error) {
	file, err := os.Open("words.txt")
	if err != nil {
		return "", errors.New("impossible d'ouvrir le fichier des mots")
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, strings.TrimSpace(scanner.Text()))
	}

	if len(words) == 0 {
		return "", errors.New("le fichier des mots est vide")
	}

	rand.Seed(time.Now().UnixNano())
	return words[rand.Intn(len(words))], nil
}
