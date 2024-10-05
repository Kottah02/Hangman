package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode"
)

const maxAttempts = 7

var hangmanStages = []string{
	`
  +---+
  |   |
      |
      |
      |
      |
=========
`,
	`
  +---+
  |   |
  O   |
      |
      |
      |
=========
`,
	`
  +---+
  |   |
  O   |
  |   |
      |
      |
=========
`,
	`
  +---+
  |   |
  O   |
 /|   |
      |
      |
=========
`,
	`
  +---+
  |   |
  O   |
 /|\  |
      |
      |
=========
`,
	`
  +---+
  |   |
  O   |
 /|\  |
 /    |
      |
=========
`,
	`
  +---+
  |   |
  O   |
 /|\  |
 / \  |
      |
=========
`,
}

func main() {
	fmt.Println(" =====Bienvenue dans le jeu du pendu!===== ")
	fmt.Println("Devinez le mot avant que l'homme ne soit pendu!")

	// Afficher le menu de difficulté
	difficulty, err := chooseDifficulty()
	if err != nil {
		fmt.Println("Erreur:", err)
		return
	}

	err = startGame(difficulty)
	if err != nil {
		fmt.Println("Erreur:", err)
	}
}

func startGame(difficulty string) error {
	word, err := getRandomWord(difficulty)
	if err != nil {
		return err
	}

	// Initialize the discovered slice with underscores
	discovered := make([]rune, len(word))
	for i := range discovered {
		discovered[i] = '_'
	}

	// Select a random position and reveal that letter
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(word))
	discovered[randomIndex] = rune(word[randomIndex])

	usedLetters := map[rune]bool{}
	// Mark the revealed letter as already used
	usedLetters[rune(word[randomIndex])] = true

	attempts := 0

	for attempts < maxAttempts {
		printGameState(discovered, usedLetters, attempts)

		guess := getUserGuess()

		// Vérifier si l'utilisateur devine le mot entier
		if len(guess) > 1 {
			if guess == word {
				fmt.Println("---------------------------------------------")
				fmt.Println("! Félicitations, vous avez deviné le mot:", word)
				return nil
			} else {
				fmt.Println("Mauvaise supposition! Vous perdez 2 tentatives.")
				attempts += 2
			}
		} else {
			// Si c'est une seule lettre
			runeGuess := rune(guess[0])
			if usedLetters[runeGuess] {
				fmt.Println("Vous avez déjà utilisé cette lettre.")
				continue
			}

			usedLetters[runeGuess] = true

			if strings.ContainsRune(word, runeGuess) {
				updateDiscovered(discovered, word, runeGuess)
				if strings.Join(stringSlice(discovered), "") == word {
					fmt.Println("---------------------------------------------")
					fmt.Println("! Félicitations, vous avez deviné le mot:", word)
					return nil
				}
			} else {
				attempts++
			}
		}
	}

	printGameState(discovered, usedLetters, attempts)
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
	index := attempts
	if index >= len(hangmanStages) {
		index = len(hangmanStages) - 1
	}
	fmt.Println(hangmanStages[index]) // Affiche l'art ASCII correspondant au nombre d'essais
	fmt.Printf("Mot: %s\n", strings.Join(stringSlice(discovered), " "))
	fmt.Printf("Lettres utilisées: %s\n", strings.Join(mapKeysToSlice(usedLetters), " "))
	fmt.Printf("Essais restants: %d\n", maxAttempts-attempts)
}

func getUserGuess() string {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Entrez une lettre ou devinez le mot entier: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// Si l'utilisateur entre une seule lettre, on vérifie que c'est une lettre valide
		if len(input) == 1 && unicode.IsLetter(rune(input[0])) {
			return input
		}

		// Si l'utilisateur entre un mot entier, on accepte aussi
		if len(input) > 1 {
			return input
		}

		fmt.Println("Veuillez entrer une lettre valide ou deviner le mot entier.")
	}
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

// Fonction pour afficher le menu et choisir la difficulté
func chooseDifficulty() (string, error) {
	fmt.Println("Choisissez la difficulté du jeu:")
	fmt.Println("1. Facile")
	fmt.Println("2. Moyen")
	fmt.Println("3. Difficile")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	switch input {
	case "1":
		return "facile", nil
	case "2":
		return "moyen", nil
	case "3":
		return "difficile", nil
	default:
		return "", errors.New("choix invalide")
	}
}

// Fonction pour obtenir un mot aléatoire selon la difficulté
func getRandomWord(difficulty string) (string, error) {
	var filename string
	switch difficulty {
	case "facile":
		filename = "words_facile.txt"
	case "moyen":
		filename = "words_moyen.txt"
	case "difficile":
		filename = "words_difficile.txt"
	default:
		return "", errors.New("difficulté inconnue")
	}

	file, err := os.Open(filename)
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
