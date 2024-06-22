package presentation

import (
	"net/http"
	"strconv"

	"github.com/arnonsang/badwords/usecase"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	e       *echo.Echo
	usecase usecase.BadWordsUseCase
}

func NewServer(usecase usecase.BadWordsUseCase) *Server {
	return &Server{
		e:       echo.New(),
		usecase: usecase,
	}
}

func (s *Server) Start(address string) error {
	s.setupMiddleware()
	s.setupRoutes()
	return s.e.Start(address)
}

func (s *Server) setupMiddleware() {
	s.e.Use(middleware.Logger())
	s.e.Use(middleware.Recover())
}

func (s *Server) setupRoutes() {
	s.e.Static("/", "static")
	s.e.GET("/api/word", s.getWord)
	s.e.GET("/api/words", s.getWords)
	s.e.GET("/api/word/:word", s.checkWord)
	s.e.GET("/api/sentence/:sentence", s.checkSentence)
	s.e.GET("/healthz", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})
}

func (s *Server) getWord(c echo.Context) error {
	nStr := c.QueryParam("n")
	n, err := strconv.Atoi(nStr)
	if err != nil {
		n = 1
	}
	if n < 1 {
		c.JSON(http.StatusBadRequest, "n must be greater than 0")
	}
	return c.JSON(http.StatusOK, s.usecase.GetWord(n))
}

func (s *Server) getWords(c echo.Context) error {
	return c.JSON(http.StatusOK, s.usecase.GetWords())
}

func (s *Server) checkWord(c echo.Context) error {
	word := c.Param("word")
	return c.JSON(http.StatusOK, s.usecase.CheckWord(word))
}

func (s *Server) checkSentence(c echo.Context) error {
	sentence := c.Param("sentence")
	return c.JSON(http.StatusOK, s.usecase.CheckSentence(sentence))
}
