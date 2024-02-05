package cogs

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (h *CogHandler) OnVoiceStateUpdate(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
	fmt.Println("OnVoiceStateUpdate")
	fmt.Println(vs)
}
