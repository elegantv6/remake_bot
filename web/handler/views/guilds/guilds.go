package guilds

import (
	"context"
	"encoding/json"
	"html/template"
	"log/slog"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/maguro-alternative/remake_bot/web/shared/ctxvalue"

	"github.com/maguro-alternative/remake_bot/web/components"
	"github.com/maguro-alternative/remake_bot/web/service"
	"github.com/maguro-alternative/remake_bot/web/shared/model"
)

type userGuild struct {
	ID          string                   `json:"id"`
	Name        string                   `json:"name"`
	Icon        string                   `json:"icon"`
	Owner       bool                     `json:"owner"`
	Permissions int64                    `json:"permissions"`
	Features    []discordgo.GuildFeature `json:"features"`
}

type GuildsViewHandler struct {
	indexService *service.IndexService
}

func NewGuildsViewHandler(indexService *service.IndexService) *GuildsViewHandler {
	return &GuildsViewHandler{
		indexService: indexService,
	}
}

func (g *GuildsViewHandler) Index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	discordLoginUser, err := ctxvalue.DiscordUserFromContext(ctx)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "Discordログインユーザー取得に失敗しました。", "エラー: ", err.Error())
		return
	}
	// Lineの認証情報なしでもアクセス可能なためエラーレスポンスは出さない
	lineSession, err := ctxvalue.LineUserFromContext(ctx)
	if err != nil {
		lineSession = &model.LineOAuthSession{}
	}

	var matchGuilds []discordgo.UserGuild
	var client = *g.indexService.Client
	botGuilds, err := g.indexService.DiscordSession.UserGuilds(100, "", "", true, discordgo.WithClient(&client))
	if err != nil {
		slog.ErrorContext(ctx, "Botサーバー取得に失敗しました。", "エラー: ", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	userGuilds, err := getUserGuilds(discordLoginUser.Token, client)
	if err != nil {
		slog.ErrorContext(ctx, "ユーザーサーバー取得に失敗しました。", "エラー: ", err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	for _, botGuild := range botGuilds {
		for _, userGuild := range userGuilds {
			if botGuild.ID == userGuild.ID {
				matchGuilds = append(matchGuilds, userGuild)
				break
			}
		}
	}

	accountVer := strings.Builder{}
	accountVer.WriteString(components.CreateDiscordAccountVer(discordLoginUser.User))
	accountVer.WriteString(components.CreateLineAccountVer(lineSession.User))
	htmlGuildBuilders := strings.Builder{}
	for _, guild := range matchGuilds {
		if guild.Icon == "" {
			htmlGuildBuilders.WriteString(`
			<a href="/guild/` + guild.ID + `">
				<img src="/static/img/discord-icon.jpg" alt="` + guild.Name + `">
				<li>` + guild.Name + `</li>
			</a><br>
			`)
			continue
		}
		htmlGuildBuilders.WriteString(`
		<a href="/guild/` + guild.ID + `">
			<img src="https://cdn.discordapp.com/icons/` + guild.ID + `/` + guild.Icon + `.png" alt="` + guild.Name + `">
			<li>` + guild.Name + `</li>
		</a><br>
		`)
	}
	data := struct {
		Title       string
		AccountVer  template.HTML
		JsScriptTag template.HTML
		Guilds      template.HTML
	}{
		Title:      "サーバー一覧",
		AccountVer: template.HTML(accountVer.String()),
		Guilds:     template.HTML(htmlGuildBuilders.String()),
	}
	tmpl := template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/views/guilds/guilds.html"))
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		slog.ErrorContext(ctx, "template execute error: "+err.Error())
	}
}

// discordgo.UserGuildをそのまま使用すると、jsonデコード時にエラーが発生するため、userGuildを使用する
func getUserGuilds(token string, client http.Client) ([]discordgo.UserGuild, error) {
	url := "https://discord.com/api/users/@me/guilds"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var guilds []userGuild
	if err := json.NewDecoder(resp.Body).Decode(&guilds); err != nil {
		return nil, err
	}
	var userGuilds []discordgo.UserGuild
	for _, guild := range guilds {
		userGuilds = append(userGuilds, discordgo.UserGuild{
			ID:          guild.ID,
			Name:        guild.Name,
			Icon:        guild.Icon,
			Owner:       guild.Owner,
			Permissions: guild.Permissions,
			Features:    guild.Features,
		})
	}
	return userGuilds, nil
}
