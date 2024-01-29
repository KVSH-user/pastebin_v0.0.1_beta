package pastebin

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	resp "pastebin/internal/lib/api/response"
	"pastebin/internal/lib/random"
	"pastebin/internal/storage/postgres"
)

type Reader struct {
	Text string `json:"text"`
}

type Request struct {
	Text        string `json:"text,omitempty"`
	Alias       string `json:"alias,omitempty"`
	AliasForDel string `json:"alias_for_del,omitempty"`
	OnlyOne     bool   `json:"only_one,omitempty"`
}

type Response struct {
	resp.Response
	Alias       string `json:"alias,omitempty"`
	AliasForDel string `json:"alias_for_del,omitempty"`
	Text        string `json:"text,omitempty"`
	OnlyOne     bool   `json:"only_one,omitempty"`
}

type PastebinSaver interface {
	SavePastebin(text string, alias string, aliasForDel string, onlyOne bool) error
}

type PastebinDel interface {
	DelPastebin(aliasForDel string) error
}

type PastebinRead interface {
	ReadPastebin(alias string) (*postgres.Pastebin, error)
}

func New(log *slog.Logger, pastebinSaver PastebinSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.handlers.pastebin.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", err)

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString()
		}

		aliasForDel := req.AliasForDel
		if aliasForDel == "" {
			aliasForDel = random.NewRandomString() + "_del"
		}

		err = pastebinSaver.SavePastebin(req.Text, alias, aliasForDel, req.OnlyOne)
		if err != nil {
			log.Info("failed to add pastebin: ", err)

			render.JSON(w, r, resp.Error("failed to add pastebin"))

			return
		}

		log.Info("pastebin added")

		responseOK(w, r, alias, aliasForDel, req.OnlyOne)
	}
}

func Del(log *slog.Logger, pastebinDel PastebinDel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.handlers.pastebin.Del"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		aliasForDel := chi.URLParam(r, "aliasForDel")
		if aliasForDel == "" {
			log.Error("aliasForDel is empty")

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		err := pastebinDel.DelPastebin(aliasForDel)
		if err != nil {
			log.Info("failed to delete pastebin: ", err)

			render.JSON(w, r, resp.Error("failed to delete pastebin"))

			return
		}

		log.Info("pastebin deleted")

		render.JSON(w, r, "Pastebin deleted")
	}
}

func Read(log *slog.Logger, pastebinRead PastebinRead, pastebinDel PastebinDel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.handlers.pastebin.Read"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Error("alias is empty")

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		result, err := pastebinRead.ReadPastebin(alias)
		if err != nil {
			log.Info("failed to read pastebin: ", err)

			render.JSON(w, r, resp.Error("failed to read pastebin"))

			return
		}

		log.Info("pastebin read")

		render.JSON(w, r, result.Text)

		if result.OnlyOne == true {
			err = pastebinDel.DelPastebin(result.AliasForDel)
		}
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, alias string, aliasForDel string, onlyOne bool) {
	render.JSON(w, r, Response{
		Response:    resp.OK(),
		Alias:       alias,
		AliasForDel: aliasForDel,
		OnlyOne:     onlyOne,
	})
}
