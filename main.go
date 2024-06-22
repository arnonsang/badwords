package main

import (
	"github.com/arnonsang/badwords/presentation"
	"github.com/arnonsang/badwords/usecase"
)

func main() {
	server := presentation.NewServer(usecase.NewBadWordsUseCase())
	server.Start(":8089")
}
