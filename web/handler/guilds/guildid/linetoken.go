package guildid

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/maguro-alternative/remake_bot/web/handler/guilds/guildid/internal"
	"github.com/maguro-alternative/remake_bot/web/service"
)

type GuildIdHandler struct {
	IndexService *service.IndexService
}

func NewGuildIdHandler(indexService *service.IndexService) *GuildIdHandler {
	return &GuildIdHandler{
		IndexService: indexService,
	}
}

func (g *GuildIdHandler) Index(w http.ResponseWriter, r *http.Request) {
	//       7
	// /guild/{guildId:[0-9]+}/linetoken
	categoryPositions := make(map[string]internal.DiscordChannel)
	guildId := r.URL.String()[7:strings.Index(r.URL.String(), "/linetoken")]
	guild, err := g.IndexService.DiscordSession.State.Guild(guildId)
	if err != nil {
		http.Error(w, "Not get guild id", http.StatusInternalServerError)
		return
	}
	repo := internal.NewRepository(g.IndexService.DB)
	for _, channel := range guild.Channels {
		if channel.Type != discordgo.ChannelTypeGuildCategory {
			continue
		}
		categoryPositions[channel.ID] = internal.DiscordChannel{
			Name:     channel.Name,
			Position: channel.Position,
		}
	}
	//[categoryPosition]map[channelPosition]channelName
	channelsInCategory := make(map[int]map[int]internal.DiscordChannelSelect, len(categoryPositions)+1)
	for _, channel := range guild.Channels {
		if channel.Type == discordgo.ChannelTypeGuildForum {
			continue
		}
		if channel.Type == discordgo.ChannelTypeGuildCategory {
			categoryPosition := categoryPositions[channel.ID]
			channelsInCategory[categoryPosition.Position] = make(map[int]internal.DiscordChannelSelect)
			continue
		}
		typeIcon := "🔊"
		if channel.Type == discordgo.ChannelTypeGuildText {
			typeIcon = "📝"
		}
		categoryPosition := categoryPositions[channel.ParentID]
		channelsInCategory[categoryPosition.Position][channel.Position] = internal.DiscordChannelSelect{
			ID:   channel.ID,
			Name: fmt.Sprintf("%s:%s:%s", categoryPosition.Name, typeIcon, channel.Name),
		}
	}
	lineBot, err := repo.GetLineBot(r.Context(), guildId)
	if err != nil && err.Error() == "sql: no rows in result set" {
		err = repo.InsertLineBot(r.Context(), &internal.LineBot{
			GuildID: guildId,
			DefaultChannelID: guild.SystemChannelID,
			DebugMode: false,
		})
		if err != nil {
			http.Error(w, "line_bot:"+err.Error(), http.StatusInternalServerError)
			return
		}
		err = repo.InsertLineBotIv(r.Context(), &internal.LineBotIv{
			GuildID: guildId,
		})
		if err != nil {
			http.Error(w, "line_bot_iv:"+err.Error(), http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	htmlSelectChannels := ``
	for _, channels := range channelsInCategory {
		for _, channelSelect := range channels {
			if lineBot.DefaultChannelID == channelSelect.ID {
				htmlSelectChannels += fmt.Sprintf(`<option value="%s" selected>%s</option>`, channelSelect.ID, channelSelect.Name)
				continue
			}
			htmlSelectChannels += fmt.Sprintf(`<option value="%s">%s</option>`, channelSelect.ID, channelSelect.Name)
		}
	}
	data := struct {
		guildID  string
		chennels string
	}{
		guildID:  guildId,
		chennels: htmlSelectChannels,
	}
	t := template.Must(template.New("linetoken.html").ParseFiles("web/templates/views/guilds/linetoken.html"))
	t.Execute(os.Stdout, data)
}
