package main

import (
	"errors"
	"strings"
)

func validateChirp(content *ChirpInput) error {
	const MAX_CHIRP_LENGTH = 140
	if len(content.Body) > MAX_CHIRP_LENGTH {
		return errors.New("Invalid body length")
	}

	// Cleaning content
	ammendedChirp := content
	ammendedChirp.Body = matchAndFixBanWords(content.Body)

	return nil
}

func matchAndFixBanWords(body string) string {
	var BAD_WORDS = []string{"", "kerfuffle", "sharbert", "fornax"}
	lowerCaseWords := strings.Split(body, " ")
	var cleanBody []string
	for _, word := range lowerCaseWords {
		for _, bad := range BAD_WORDS {
			if strings.ToLower(word) == bad {
				word = "****"
			}
		}
		cleanBody = append(cleanBody, word)
	}

	return strings.Join(cleanBody, " ")
}
