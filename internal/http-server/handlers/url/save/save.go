package save

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/lib/random"
)

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

// TODO:move to config
const aliasLength = 6

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers/url/save"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var request Request

		err := render.DecodeJSON(r.Body, &request)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, response.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", request))

		if err := validator.New().Struct(request); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("failed to validate request", sl.Err(err))

			render.JSON(w, r, response.ValidationError(validateErr))

			return
		}

		alias := request.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}
	}
}
