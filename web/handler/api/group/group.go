package group

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/maguro-alternative/remake_bot/repository"

	"github.com/maguro-alternative/remake_bot/web/handler/api/group/internal"
	"github.com/maguro-alternative/remake_bot/web/service"
)

type LineGroupHandler struct {
	IndexService    *service.IndexService
	Repo            repository.RepositoryFunc
}

func NewLineGroupHandler(
	indexService *service.IndexService,
	repo repository.RepositoryFunc,
) *LineGroupHandler {
	return &LineGroupHandler{
		IndexService:    indexService,
		Repo:            repo,
	}
}

func (g *LineGroupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		slog.ErrorContext(ctx, "/api/group Method Not Allowed")
		return
	}
	var lineGroupJson internal.LineBotJson
	if err := json.NewDecoder(r.Body).Decode(&lineGroupJson); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		slog.ErrorContext(ctx, "jsonの読み取りに失敗しました:"+err.Error())
		return
	}
	if err := lineGroupJson.Validate(); err != nil {
		http.Error(w, "Unprocessable Entity", http.StatusUnprocessableEntity)
		slog.ErrorContext(ctx, "jsonのバリデーションに失敗しました:"+err.Error())
		return
	}
	guildId := r.PathValue("guildId")
	err := g.Repo.UpdateLineBot(ctx, &repository.LineBot{
		GuildID:          guildId,
		DefaultChannelID: lineGroupJson.DefaultChannelID,
	})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "DBの更新に失敗しました:"+err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("OK")
}
