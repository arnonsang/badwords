package usecase

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/arnonsang/badwords/assets"
)

type NxWord struct {
	Status int      `json:"status"`
	Words  []string `json:"word"`
}

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

type Replaced struct {
	Status   int      `json:"status"`
	Sentence string   `json:"sentence"`
	Count    int      `json:"count"`
	BadWords []string `json:"badWords"`
	Replaced string   `json:"replaced"`
}

type BadWordsUseCase interface {
	GetWord(n int) NxWord
	GetWords() WordList
	CheckWord(word string) Word
	CheckSentence(sentence string) Sentence
	Replacer(sentence string) Replaced
}

type badWordsUseCase struct{}

func NewBadWordsUseCase() BadWordsUseCase {
	return &badWordsUseCase{}
}

func (uc *badWordsUseCase) GetWord(n int) NxWord {
	wordCount := len(assets.BadWords)
	if n > wordCount {
		n = wordCount
	}
	words := make([]string, n)
	for i := 0; i < n; i++ {
		words[i] = assets.BadWords[i]
	}
	return NxWord{Status: http.StatusOK, Words: words}
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
	lowercaseWord := strings.ToLower(word)
	for _, badWord := range assets.BadWords {
		if lowercaseWord == badWord {
			return true
		}
	}
	return false
}

func (uc *badWordsUseCase) detectBadWords(sentence string) []string {
	badWords := []string{}
	lowercaseSentence := strings.ToLower(sentence)
	for _, word := range assets.BadWords {
		re := regexp.MustCompile(`\b` + regexp.QuoteMeta(strings.ToLower(word)) + `\b`)
		if re.MatchString(lowercaseSentence) {
			badWords = append(badWords, word)
		}
	}
	return badWords
}

func (uc *badWordsUseCase) countBadWords(sentence string) int {
	badWords := uc.detectBadWords(sentence)
	return len(badWords)
}

func (uc *badWordsUseCase) Replacer(sentence string) Replaced {
	badWords := uc.detectBadWords(sentence)
	replaced := sentence
	for _, word := range badWords {
		re := regexp.MustCompile(`(?i)\b` + regexp.QuoteMeta(word) + `\b`)
		replaced = re.ReplaceAllString(replaced, strings.Repeat("*", len(word)))
	}
	return Replaced{Status: http.StatusOK, Sentence: sentence, Count: len(badWords), BadWords: badWords, Replaced: replaced}
}
