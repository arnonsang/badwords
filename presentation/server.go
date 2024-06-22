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

type SentenceRequest struct {
	Sentence string `json:"sentence" validate:"required"`
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
	s.e.Use(middleware.CORS())
	s.e.Use(middleware.Secure())
	s.e.Use(middleware.RequestID())
	s.e.Use(middleware.CSRF())
	s.e.Use(middleware.Gzip())
	s.e.Use(middleware.Decompress())
	s.e.Use(middleware.BodyLimit("2M"))
	s.e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	s.e.Use(middleware.Timeout())
}

func (s *Server) setupRoutes() {
	s.e.Static("/", "static")
	s.e.GET("/api/word", s.getWord)
	s.e.GET("/api/words", s.getWords)
	s.e.GET("/api/word/:word", s.checkWord)
	s.e.GET("/api/sentence/:sentence", s.checkSentence)
	s.e.POST("/api/sentence", s.checkSentence)
	s.e.GET("/api/replacer/:sentence", s.sentenceReplacer)
	s.e.POST("/api/replacer", s.sentenceReplacer)
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
	var sentenceReq SentenceRequest

	if c.Request().Method == "POST" {
		if err := c.Bind(&sentenceReq); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		}

		if sentenceReq.Sentence == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "body.sentence is required, This service only accept JSON format"})
		}

	} else {
		sentenceReq.Sentence = c.Param("sentence")
	}

	return c.JSON(http.StatusOK, s.usecase.CheckSentence(sentenceReq.Sentence))
}

func (s *Server) sentenceReplacer(c echo.Context) error {
	var sentenceReq SentenceRequest

	if c.Request().Method == "POST" {
		if err := c.Bind(&sentenceReq); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		}

		if sentenceReq.Sentence == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "body.sentence is required, This service only accept JSON format"})
		}

	} else {
		sentenceReq.Sentence = c.Param("sentence")
	}

	return c.JSON(http.StatusOK, s.usecase.Replacer(sentenceReq.Sentence))
}
