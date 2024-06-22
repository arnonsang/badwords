package usecase

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/arnonsang/badwords/assets"
)

type WordList struct {
	Status int      `json:"status"`
	Count  int      `json:"count"`
	Words  []string `json:"words"`
}

type Word struct {
	Status int    `json:"status"`
	Word   string `json:"word"`
	IsBad  bool   `json:"isBad"`
}

type Sentence struct {
	Status   int    `json:"status"`
	Sentence string `json:"sentence"`
	Count    int    `json:"count"`
	BadWords []string
}

type BadWordsUseCase interface {
	GetWords() WordList
	CheckWord(word string) Word
	CheckSentence(sentence string) Sentence
}

type badWordsUseCase struct{}

func NewBadWordsUseCase() BadWordsUseCase {
	return &badWordsUseCase{}
}

func (uc *badWordsUseCase) GetWords() WordList {
	wordCount := len(assets.BadWords)
	return WordList{Status: http.StatusOK, Count: wordCount, Words: assets.BadWords}
}

func (uc *badWordsUseCase) CheckWord(word string) Word {
	return Word{Status: http.StatusOK, Word: word, IsBad: uc.isBadWord(word)}
}

func (uc *badWordsUseCase) CheckSentence(sentence string) Sentence {
	return Sentence{Status: http.StatusOK, Sentence: sentence, Count: uc.countBadWords(sentence), BadWords: uc.detectBadWords(sentence)}
}

func (uc *badWordsUseCase) isBadWord(word string) bool {
	for _, badWord := range assets.BadWords {
		if word == badWord {
			return true
		}
	}
	return false
}

func (uc *badWordsUseCase) detectBadWords(sentence string) []string {
	badWords := []string{}
	for _, word := range assets.BadWords {
		re := regexp.MustCompile(`\b` + regexp.QuoteMeta(word) + `\b`)
		if re.MatchString(strings.ToLower(sentence)) {
			badWords = append(badWords, word)
		}
	}
	return badWords
}

func (uc *badWordsUseCase) countBadWords(sentence string) int {
	badWords := uc.detectBadWords(sentence)
	return len(badWords)
}
