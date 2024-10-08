package cogs

import (
	"context"
	"testing"

	"github.com/maguro-alternative/remake_bot/repository"
	"github.com/maguro-alternative/remake_bot/testutil/mock"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVcSignal(t *testing.T) {
	ctx := context.Background()
	testUser := &discordgo.User{
		ID:       "11",
		Username: "testuser",
		Avatar:   "a_",
		Bot:      false,
	}
	testUser2 := &discordgo.User{
		ID:       "22",
		Username: "testuser2",
		Avatar:   "a_",
		Bot:      false,
	}
	testUser3 := &discordgo.User{
		ID:       "33",
		Username: "testuser3",
		Avatar:   "a_",
		Bot:      false,
	}
	beforeGuildId := "111"
	afterGuildId := "222"
	beforeChannelId := "1111"
	afterChannelId := "2222"
	beforeChannelId2 := "1112"
	afterChannelId2 := "2223"
	beforeSendChannelId := "11111"
	afterSendChannelId := "22222"

	discordState := discordgo.NewState()
	err := discordState.GuildAdd(&discordgo.Guild{
		ID: afterGuildId,
		Channels: []*discordgo.Channel{
			{
				ID:       afterChannelId,
				Name:     "after_test_vc",
				Position: 1,
				Type:     discordgo.ChannelTypeGuildVoice,
			},
			{
				ID:       afterSendChannelId,
				Name:     "after_test_text",
				Position: 2,
				Type:     discordgo.ChannelTypeGuildText,
			},
			{
				ID:       afterChannelId2,
				Name:     "after_test_vc2",
				Position: 3,
				Type:     discordgo.ChannelTypeGuildVoice,
			},
		},
	})
	require.NoError(t, err)

	err = discordState.GuildAdd(&discordgo.Guild{
		ID: beforeGuildId,
		Channels: []*discordgo.Channel{
			{
				ID:       beforeChannelId,
				Name:     "before_test_vc",
				Position: 1,
				Type:     discordgo.ChannelTypeGuildVoice,
			},
			{
				ID:       beforeSendChannelId,
				Name:     "vefore_test_text",
				Position: 2,
				Type:     discordgo.ChannelTypeGuildText,
			},
			{
				ID:       beforeChannelId2,
				Name:     "vefore_test_vc2",
				Position: 3,
				Type:     discordgo.ChannelTypeGuildVoice,
			},
		},
	})
	require.NoError(t, err)

	t.Run("正常系(通話開始)", func(t *testing.T) {
		discordState.Guilds[0].VoiceStates = []*discordgo.VoiceState{
			{
				GuildID:   afterGuildId,
				ChannelID: afterChannelId,
				Member: &discordgo.Member{
					User: testUser,
				},
				SelfStream: false,
				SelfVideo:  false,
			},
		}
		messages, err := onVoiceStateUpdateFunc(
			ctx,
			&repository.RepositoryFuncMock{
				GetVcSignalNgUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
				GetVcSignalNgRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
				GetVcSignalChannelAllColumnByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) (*repository.VcSignalChannelAllColumn, error) {
					return &repository.VcSignalChannelAllColumn{
						VcChannelID:     vcChannelID,
						GuildID:         afterGuildId,
						SendSignal:      true,
						SendChannelID:   afterSendChannelId,
						JoinBot:         false,
						EveryoneMention: false,
					}, nil
				},
				GetVcSignalMentionUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
				GetVcSignalMentionRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
			},
			&mock.SessionMock{
				ChannelMessageSendFunc: func(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error) {
					return &discordgo.Message{
						Content: content,
					}, nil
				},
				ChannelMessageSendEmbedFunc: func(channelID string, embed *discordgo.MessageEmbed, options ...discordgo.RequestOption) (*discordgo.Message, error) {
					return &discordgo.Message{
						Embeds: []*discordgo.MessageEmbed{embed},
					}, nil
				},
			},
			discordState,
			&discordgo.VoiceStateUpdate{
				VoiceState: &discordgo.VoiceState{
					GuildID:   afterGuildId,
					ChannelID: afterChannelId,
					Member: &discordgo.Member{
						User: testUser,
					},
					SelfStream: false,
					SelfVideo:  false,
				},
			},
		)
		assert.NoError(t, err)
		assert.Len(t, messages, 2)
		assert.Equal(t, messages[0].Content, "現在1人 <@11> が after_test_vcに入室しました。")
		assert.Equal(t, messages[1].Embeds[0].Title, "通話開始")
		assert.Equal(t, messages[1].Embeds[0].Description, "<#2222>")
		assert.Equal(t, messages[1].Embeds[0].Author.Name, "testuser")
		assert.Equal(t, messages[1].Embeds[0].Author.IconURL, "https://cdn.discordapp.com/avatars/11/a_.gif?size=64")
	})

	t.Run("正常系(通話終了)", func(t *testing.T) {
		messages, err := onVoiceStateUpdateFunc(
			ctx,
			&repository.RepositoryFuncMock{
				GetVcSignalNgUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
				GetVcSignalNgRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
				GetVcSignalChannelAllColumnByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) (*repository.VcSignalChannelAllColumn, error) {
					return &repository.VcSignalChannelAllColumn{
						VcChannelID:     vcChannelID,
						GuildID:         beforeGuildId,
						SendSignal:      true,
						SendChannelID:   beforeSendChannelId,
						JoinBot:         false,
						EveryoneMention: false,
					}, nil
				},
				GetVcSignalMentionUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
				GetVcSignalMentionRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
			},
			&mock.SessionMock{
				ChannelMessageSendFunc: func(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error) {
					return &discordgo.Message{
						Content: content,
					}, nil
				},
				ChannelMessageSendEmbedFunc: func(channelID string, embed *discordgo.MessageEmbed, options ...discordgo.RequestOption) (*discordgo.Message, error) {
					return &discordgo.Message{
						Embeds: []*discordgo.MessageEmbed{embed},
					}, nil
				},
			},
			discordState,
			&discordgo.VoiceStateUpdate{
				VoiceState: &discordgo.VoiceState{
					GuildID:   "",
					ChannelID: "",
					Member: &discordgo.Member{
						User: testUser,
					},
					SelfStream: false,
					SelfVideo:  false,
				},
				BeforeUpdate: &discordgo.VoiceState{
					GuildID:   beforeGuildId,
					ChannelID: beforeChannelId,
					Member: &discordgo.Member{
						User: testUser,
					},
					SelfStream: false,
					SelfVideo:  false,
				},
			},
		)
		assert.NoError(t, err)
		assert.Len(t, messages, 3)
		assert.Equal(t, messages[0].Content, "現在0人 <@11> が before_test_vcから退室しました。")
		assert.Equal(t, messages[1].Content, "通話が終了しました。")
		assert.Equal(t, messages[2].Embeds[0].Title, "通話終了")
	})

	t.Run("正常系(ボイスチャンネル移動で移動前サーバー通話終了と移動先サーバー通話開始)", func(t *testing.T) {
		discordState.Guilds[0].VoiceStates = []*discordgo.VoiceState{
			{
				GuildID:   afterGuildId,
				ChannelID: afterChannelId,
				Member: &discordgo.Member{
					User: testUser,
				},
				SelfStream: false,
				SelfVideo:  false,
			},
		}
		discordState.Guilds[1].VoiceStates = []*discordgo.VoiceState{}
		messages, err := onVoiceStateUpdateFunc(
			ctx,
			&repository.RepositoryFuncMock{
				GetVcSignalNgUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
				GetVcSignalNgRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
				GetVcSignalChannelAllColumnByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) (*repository.VcSignalChannelAllColumn, error) {
					if vcChannelID == afterChannelId {
						return &repository.VcSignalChannelAllColumn{
							VcChannelID:     vcChannelID,
							GuildID:         afterGuildId,
							SendSignal:      true,
							SendChannelID:   afterSendChannelId,
							JoinBot:         false,
							EveryoneMention: false,
						}, nil
					}
					return &repository.VcSignalChannelAllColumn{
						VcChannelID:     vcChannelID,
						GuildID:         beforeGuildId,
						SendSignal:      true,
						SendChannelID:   beforeSendChannelId,
						JoinBot:         false,
						EveryoneMention: false,
					}, nil
				},
				GetVcSignalMentionUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
				GetVcSignalMentionRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
			},
			&mock.SessionMock{
				ChannelMessageSendFunc: func(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error) {
					return &discordgo.Message{
						Content: content,
					}, nil
				},
				ChannelMessageSendEmbedFunc: func(channelID string, embed *discordgo.MessageEmbed, options ...discordgo.RequestOption) (*discordgo.Message, error) {
					return &discordgo.Message{
						Embeds: []*discordgo.MessageEmbed{embed},
					}, nil
				},
			},
			discordState,
			&discordgo.VoiceStateUpdate{
				VoiceState: &discordgo.VoiceState{
					GuildID:   afterGuildId,
					ChannelID: afterChannelId,
					Member: &discordgo.Member{
						User: testUser,
					},
					SelfStream: false,
					SelfVideo:  false,
				},
				BeforeUpdate: &discordgo.VoiceState{
					GuildID:   beforeGuildId,
					ChannelID: beforeChannelId,
					Member: &discordgo.Member{
						User: testUser,
					},
					SelfStream: false,
					SelfVideo:  false,
				},
			},
		)
		assert.NoError(t, err)
		assert.Len(t, messages, 5)
		assert.Equal(t, messages[0].Content, "現在0人 <@11> が before_test_vcから退室しました。")
		assert.Equal(t, messages[1].Content, "通話が終了しました。")
		assert.Equal(t, messages[2].Embeds[0].Title, "通話終了")
		assert.Equal(t, messages[3].Content, "現在1人 <@11> が after_test_vcに入室しました。")
		assert.Equal(t, messages[4].Embeds[0].Title, "通話開始")
		assert.Equal(t, messages[4].Embeds[0].Description, "<#2222>")
		assert.Equal(t, messages[4].Embeds[0].Author.Name, "testuser")
		assert.Equal(t, messages[4].Embeds[0].Author.IconURL, "https://cdn.discordapp.com/avatars/11/a_.gif?size=64")
	})

	t.Run("正常系(サーバー内でのボイスチャンネル移動)", func(t *testing.T) {
		discordState.Guilds[0].VoiceStates = []*discordgo.VoiceState{
			{
				GuildID:   afterGuildId,
				ChannelID: afterChannelId,
				Member: &discordgo.Member{
					User: testUser,
				},
				SelfStream: false,
				SelfVideo:  false,
			},
		}
		messages, err := onVoiceStateUpdateFunc(
			ctx,
			&repository.RepositoryFuncMock{
				GetVcSignalNgUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
				GetVcSignalNgRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
				GetVcSignalChannelAllColumnByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) (*repository.VcSignalChannelAllColumn, error) {
					if vcChannelID == afterChannelId {
						return &repository.VcSignalChannelAllColumn{
							VcChannelID:     vcChannelID,
							GuildID:         afterGuildId,
							SendSignal:      true,
							SendChannelID:   afterSendChannelId,
							JoinBot:         false,
							EveryoneMention: false,
						}, nil
					}
					return &repository.VcSignalChannelAllColumn{
						VcChannelID:     vcChannelID,
						GuildID:         beforeGuildId,
						SendSignal:      true,
						SendChannelID:   beforeSendChannelId,
						JoinBot:         false,
						EveryoneMention: false,
					}, nil
				},
				GetVcSignalMentionUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
				GetVcSignalMentionRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
			},
			&mock.SessionMock{
				ChannelMessageSendFunc: func(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error) {
					return &discordgo.Message{
						Content: content,
					}, nil
				},
				ChannelMessageSendEmbedFunc: func(channelID string, embed *discordgo.MessageEmbed, options ...discordgo.RequestOption) (*discordgo.Message, error) {
					return &discordgo.Message{
						Embeds: []*discordgo.MessageEmbed{embed},
					}, nil
				},
			},
			discordState,
			&discordgo.VoiceStateUpdate{
				VoiceState: &discordgo.VoiceState{
					GuildID:   afterGuildId,
					ChannelID: afterChannelId,
					Member: &discordgo.Member{
						User: testUser,
					},
					SelfStream: false,
					SelfVideo:  false,
				},
				BeforeUpdate: &discordgo.VoiceState{
					GuildID:   afterGuildId,
					ChannelID: afterChannelId2,
					Member: &discordgo.Member{
						User: testUser,
					},
					SelfStream: false,
					SelfVideo:  false,
				},
			},
		)
		assert.NoError(t, err)
		assert.Len(t, messages, 3)
		assert.Equal(t, messages[0].Content, "現在0人 <@11> が after_test_vc2から退室しました。")
		assert.Equal(t, messages[1].Content, "現在1人 <@11> が after_test_vcに入室しました。")
		assert.Equal(t, messages[2].Embeds[0].Title, "通話開始")
	})

	t.Run("正常系(サーバー内での2人以降のボイスチャンネル入室)", func(t *testing.T) {
		discordState.Guilds[0].VoiceStates = []*discordgo.VoiceState{
			{
				GuildID:   afterGuildId,
				ChannelID: afterChannelId,
				Member: &discordgo.Member{
					User: testUser,
				},
				SelfStream: false,
				SelfVideo:  false,
			},
			{
				GuildID:   afterGuildId,
				ChannelID: afterChannelId,
				Member: &discordgo.Member{
					User: testUser2,
				},
				SelfStream: false,
				SelfVideo:  false,
			},
		}
		messages, err := onVoiceStateUpdateFunc(
			ctx,
			&repository.RepositoryFuncMock{
				GetVcSignalNgUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
				GetVcSignalNgRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
				GetVcSignalChannelAllColumnByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) (*repository.VcSignalChannelAllColumn, error) {
					if vcChannelID == afterChannelId {
						return &repository.VcSignalChannelAllColumn{
							VcChannelID:     vcChannelID,
							GuildID:         afterGuildId,
							SendSignal:      true,
							SendChannelID:   afterSendChannelId,
							JoinBot:         false,
							EveryoneMention: false,
						}, nil
					}
					return &repository.VcSignalChannelAllColumn{
						VcChannelID:     vcChannelID,
						GuildID:         beforeGuildId,
						SendSignal:      true,
						SendChannelID:   beforeSendChannelId,
						JoinBot:         false,
						EveryoneMention: false,
					}, nil
				},
				GetVcSignalMentionUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
				GetVcSignalMentionRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
			},
			&mock.SessionMock{
				ChannelMessageSendFunc: func(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error) {
					return &discordgo.Message{
						Content: content,
					}, nil
				},
				ChannelMessageSendEmbedFunc: func(channelID string, embed *discordgo.MessageEmbed, options ...discordgo.RequestOption) (*discordgo.Message, error) {
					return &discordgo.Message{
						Embeds: []*discordgo.MessageEmbed{embed},
					}, nil
				},
			},
			discordState,
			&discordgo.VoiceStateUpdate{
				VoiceState: &discordgo.VoiceState{
					GuildID:   afterGuildId,
					ChannelID: afterChannelId,
					Member: &discordgo.Member{
						User: testUser,
					},
					SelfStream: false,
					SelfVideo:  false,
				},
			},
		)
		assert.NoError(t, err)
		assert.Len(t, messages, 1)
		assert.Equal(t, messages[0].Content, "現在2人 <@11> が after_test_vcに入室しました。")
	})

	t.Run("正常系(サーバー内での2人以降のボイスチャンネル退室)", func(t *testing.T) {
		discordState.Guilds[0].VoiceStates = []*discordgo.VoiceState{
			{
				GuildID:   afterGuildId,
				ChannelID: afterChannelId,
				Member: &discordgo.Member{
					User: testUser,
				},
				SelfStream: false,
				SelfVideo:  false,
			},
		}
		messages, err := onVoiceStateUpdateFunc(
			ctx,
			&repository.RepositoryFuncMock{
				GetVcSignalNgUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
				GetVcSignalNgRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
				GetVcSignalChannelAllColumnByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) (*repository.VcSignalChannelAllColumn, error) {
					if vcChannelID == afterChannelId {
						return &repository.VcSignalChannelAllColumn{
							VcChannelID:     vcChannelID,
							GuildID:         afterGuildId,
							SendSignal:      true,
							SendChannelID:   afterSendChannelId,
							JoinBot:         false,
							EveryoneMention: false,
						}, nil
					}
					return &repository.VcSignalChannelAllColumn{
						VcChannelID:     vcChannelID,
						GuildID:         beforeGuildId,
						SendSignal:      true,
						SendChannelID:   beforeSendChannelId,
						JoinBot:         false,
						EveryoneMention: false,
					}, nil
				},
				GetVcSignalMentionUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
				GetVcSignalMentionRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
			},
			&mock.SessionMock{
				ChannelMessageSendFunc: func(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error) {
					return &discordgo.Message{
						Content: content,
					}, nil
				},
				ChannelMessageSendEmbedFunc: func(channelID string, embed *discordgo.MessageEmbed, options ...discordgo.RequestOption) (*discordgo.Message, error) {
					return &discordgo.Message{
						Embeds: []*discordgo.MessageEmbed{embed},
					}, nil
				},
			},
			discordState,
			&discordgo.VoiceStateUpdate{
				VoiceState: &discordgo.VoiceState{
					GuildID:   "",
					ChannelID: "",
					Member: &discordgo.Member{
						User: testUser,
					},
					SelfStream: false,
					SelfVideo:  false,
				},
				BeforeUpdate: &discordgo.VoiceState{
					GuildID:   afterGuildId,
					ChannelID: afterChannelId,
					Member: &discordgo.Member{
						User: testUser,
					},
					SelfStream: false,
					SelfVideo:  false,
				},
			},
		)
		assert.NoError(t, err)
		assert.Len(t, messages, 1)
		assert.Equal(t, messages[0].Content, "現在1人 <@11> が after_test_vcから退室しました。")
	})

	t.Run("正常系(サーバー内での2人以降のボイスチャンネル移動)", func(t *testing.T) {
		discordState.Guilds[0].VoiceStates = []*discordgo.VoiceState{
			{
				GuildID:   afterGuildId,
				ChannelID: afterChannelId,
				Member: &discordgo.Member{
					User: testUser,
				},
				SelfStream: false,
				SelfVideo:  false,
			},
			{
				GuildID:   afterGuildId,
				ChannelID: afterChannelId,
				Member: &discordgo.Member{
					User: testUser2,
				},
				SelfStream: false,
				SelfVideo:  false,
			},
			{
				GuildID:   afterGuildId,
				ChannelID: afterChannelId2,
				Member: &discordgo.Member{
					User: testUser3,
				},
				SelfStream: false,
				SelfVideo:  false,
			},
		}
		messages, err := onVoiceStateUpdateFunc(
			ctx,
			&repository.RepositoryFuncMock{
				GetVcSignalNgUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
				GetVcSignalNgRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
				GetVcSignalChannelAllColumnByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) (*repository.VcSignalChannelAllColumn, error) {
					if vcChannelID == afterChannelId {
						return &repository.VcSignalChannelAllColumn{
							VcChannelID:     vcChannelID,
							GuildID:         afterGuildId,
							SendSignal:      true,
							SendChannelID:   afterSendChannelId,
							JoinBot:         false,
							EveryoneMention: false,
						}, nil
					}
					return &repository.VcSignalChannelAllColumn{
						VcChannelID:     vcChannelID,
						GuildID:         beforeGuildId,
						SendSignal:      true,
						SendChannelID:   beforeSendChannelId,
						JoinBot:         false,
						EveryoneMention: false,
					}, nil
				},
				GetVcSignalMentionUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
				GetVcSignalMentionRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
			},
			&mock.SessionMock{
				ChannelMessageSendFunc: func(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error) {
					return &discordgo.Message{
						Content: content,
					}, nil
				},
				ChannelMessageSendEmbedFunc: func(channelID string, embed *discordgo.MessageEmbed, options ...discordgo.RequestOption) (*discordgo.Message, error) {
					return &discordgo.Message{
						Embeds: []*discordgo.MessageEmbed{embed},
					}, nil
				},
			},
			discordState,
			&discordgo.VoiceStateUpdate{
				VoiceState: &discordgo.VoiceState{
					GuildID:   afterGuildId,
					ChannelID: afterChannelId,
					Member: &discordgo.Member{
						User: testUser,
					},
					SelfStream: false,
					SelfVideo:  false,
				},
				BeforeUpdate: &discordgo.VoiceState{
					GuildID:   afterGuildId,
					ChannelID: afterChannelId2,
					Member: &discordgo.Member{
						User: testUser,
					},
					SelfStream: false,
					SelfVideo:  false,
				},
			},
		)
		assert.NoError(t, err)
		assert.Len(t, messages, 2)
		assert.Equal(t, messages[0].Content, "現在1人 <@11> が after_test_vc2から退室しました。")
		assert.Equal(t, messages[1].Content, "現在2人 <@11> が after_test_vcに入室しました。")
	})

	t.Run("正常系(カメラ配信)", func(t *testing.T) {
		discordState.Guilds[0].VoiceStates = []*discordgo.VoiceState{
			{
				GuildID:   afterGuildId,
				ChannelID: afterChannelId,
				Member: &discordgo.Member{
					User: testUser,
				},
				SelfStream: false,
				SelfVideo:  true,
			},
		}
		messages, err := onVoiceStateUpdateFunc(
			ctx,
			&repository.RepositoryFuncMock{
				GetVcSignalNgUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
				GetVcSignalNgRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
				GetVcSignalChannelAllColumnByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) (*repository.VcSignalChannelAllColumn, error) {
					return &repository.VcSignalChannelAllColumn{
						VcChannelID:     vcChannelID,
						GuildID:         afterGuildId,
						SendSignal:      true,
						SendChannelID:   afterSendChannelId,
						JoinBot:         false,
						EveryoneMention: false,
					}, nil
				},
				GetVcSignalMentionUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
				GetVcSignalMentionRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
			},
			&mock.SessionMock{
				ChannelMessageSendFunc: func(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error) {
					return &discordgo.Message{
						Content: content,
					}, nil
				},
				ChannelMessageSendEmbedFunc: func(channelID string, embed *discordgo.MessageEmbed, options ...discordgo.RequestOption) (*discordgo.Message, error) {
					return &discordgo.Message{
						Embeds: []*discordgo.MessageEmbed{embed},
					}, nil
				},
			},
			discordState,
			&discordgo.VoiceStateUpdate{
				VoiceState: &discordgo.VoiceState{
					GuildID:   afterGuildId,
					ChannelID: afterChannelId,
					Member: &discordgo.Member{
						User: testUser,
					},
					SelfStream: false,
					SelfVideo:  true,
				},
				BeforeUpdate: &discordgo.VoiceState{
					GuildID:   afterGuildId,
					ChannelID: afterChannelId,
					Member: &discordgo.Member{
						User: testUser,
					},
					SelfStream: false,
					SelfVideo:  false,
				},
			},
		)
		assert.NoError(t, err)
		assert.Len(t, messages, 2)
		assert.Equal(t, messages[0].Content, "<@11> がafter_test_vcでカメラ配信を開始しました。")
		assert.Equal(t, messages[1].Embeds[0].Title, "カメラ配信")
		assert.Contains(t, messages[1].Embeds[0].Description, "testuser")
		assert.Contains(t, messages[1].Embeds[0].Description, "<#2222>")
		assert.Equal(t, messages[1].Embeds[0].Author.Name, "testuser")
		assert.Equal(t, messages[1].Embeds[0].Author.IconURL, "https://cdn.discordapp.com/avatars/11/a_.gif?size=64")
	})

	t.Run("正常系(カメラ配信終了)", func(t *testing.T) {
		discordState.Guilds[0].VoiceStates = []*discordgo.VoiceState{
			{
				GuildID:   afterGuildId,
				ChannelID: afterChannelId,
				Member: &discordgo.Member{
					User: testUser,
				},
				SelfStream: true,
				SelfVideo:  false,
			},
		}
		messages, err := onVoiceStateUpdateFunc(
			ctx,
			&repository.RepositoryFuncMock{
				GetVcSignalNgUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
				GetVcSignalNgRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
				GetVcSignalChannelAllColumnByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) (*repository.VcSignalChannelAllColumn, error) {
					return &repository.VcSignalChannelAllColumn{
						VcChannelID:     vcChannelID,
						GuildID:         afterGuildId,
						SendSignal:      true,
						SendChannelID:   afterSendChannelId,
						JoinBot:         false,
						EveryoneMention: false,
					}, nil
				},
				GetVcSignalMentionUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
				GetVcSignalMentionRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
			},
			&mock.SessionMock{
				ChannelMessageSendFunc: func(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error) {
					return &discordgo.Message{
						Content: content,
					}, nil
				},
				ChannelMessageSendEmbedFunc: func(channelID string, embed *discordgo.MessageEmbed, options ...discordgo.RequestOption) (*discordgo.Message, error) {
					return &discordgo.Message{
						Embeds: []*discordgo.MessageEmbed{embed},
					}, nil
				},
			},
			discordState,
			&discordgo.VoiceStateUpdate{
				VoiceState: &discordgo.VoiceState{
					GuildID:   afterGuildId,
					ChannelID: afterChannelId,
					Member: &discordgo.Member{
						User: testUser,
					},
					SelfStream: false,
					SelfVideo:  false,
				},
				BeforeUpdate: &discordgo.VoiceState{
					GuildID:   afterGuildId,
					ChannelID: afterChannelId,
					Member: &discordgo.Member{
						User: testUser,
					},
					SelfStream: false,
					SelfVideo:  true,
				},
			},
		)
		assert.NoError(t, err)
		assert.Len(t, messages, 1)
		assert.Equal(t, messages[0].Content, "<@11> がカメラ配信を終了しました。")
	})

	t.Run("正常系(画面共有開始)", func(t *testing.T) {
		discordState.Guilds[0].VoiceStates = []*discordgo.VoiceState{
			{
				GuildID:   afterGuildId,
				ChannelID: afterChannelId,
				Member: &discordgo.Member{
					User: testUser,
				},
				SelfStream: false,
				SelfVideo:  false,
				SelfMute:   false,
				SelfDeaf:   false,
			},
		}
		messages, err := onVoiceStateUpdateFunc(
			ctx,
			&repository.RepositoryFuncMock{
				GetVcSignalNgUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
				GetVcSignalNgRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
				GetVcSignalChannelAllColumnByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) (*repository.VcSignalChannelAllColumn, error) {
					return &repository.VcSignalChannelAllColumn{
						VcChannelID:     vcChannelID,
						GuildID:         afterGuildId,
						SendSignal:      true,
						SendChannelID:   afterSendChannelId,
						JoinBot:         false,
						EveryoneMention: false,
					}, nil
				},
				GetVcSignalMentionUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
				GetVcSignalMentionRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
			},
			&mock.SessionMock{
				ChannelMessageSendFunc: func(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error) {
					return &discordgo.Message{
						Content: content,
					}, nil
				},
				ChannelMessageSendEmbedFunc: func(channelID string, embed *discordgo.MessageEmbed, options ...discordgo.RequestOption) (*discordgo.Message, error) {
					return &discordgo.Message{
						Embeds: []*discordgo.MessageEmbed{embed},
					}, nil
				},
			},
			discordState,
			&discordgo.VoiceStateUpdate{
				VoiceState: &discordgo.VoiceState{
					GuildID:   afterGuildId,
					ChannelID: afterChannelId,
					Member: &discordgo.Member{
						User: testUser,
					},
					SelfStream: true,
					SelfVideo:  false,
					SelfMute:   false,
					SelfDeaf:   false,
				},
				BeforeUpdate: &discordgo.VoiceState{
					GuildID:   afterGuildId,
					ChannelID: afterChannelId,
					Member: &discordgo.Member{
						User: testUser,
					},
					SelfStream: false,
					SelfVideo:  false,
					SelfMute:   false,
					SelfDeaf:   false,
				},
			},
		)
		assert.NoError(t, err)
		assert.Len(t, messages, 2)
		assert.Equal(t, messages[0].Content, "<@11> がafter_test_vcで画面共有を開始しました。")
		assert.Equal(t, messages[1].Embeds[0].Title, "画面共有")
	})

	t.Run("正常系(アクティビティ画面共有開始)", func(t *testing.T) {
		discordState.Guilds[0].VoiceStates = []*discordgo.VoiceState{
			{
				GuildID:   afterGuildId,
				ChannelID: afterChannelId,
				Member: &discordgo.Member{
					User: testUser,
				},
				SelfStream: false,
				SelfVideo:  false,
				SelfMute:   false,
				SelfDeaf:   false,
			},
		}
		discordState.Guilds[0].Presences = []*discordgo.Presence{
			{
				User: &discordgo.User{
					ID:       testUser.ID,
					Username: testUser.Username,
					Avatar:   testUser.Avatar,
					Bot:      false,
				},
				Activities: []*discordgo.Activity{
					{
						Name:          "シノビマスター 閃乱カグラ NEW LINK",
						Type:          discordgo.ActivityTypeStreaming,
						URL:           "https://dmg.hpgames.jp/shinomas/index.html",
						ApplicationID: "123456789012345678",
						Assets: discordgo.Assets{
							LargeImageID: "shinomas",
							LargeText:    "シノビマスター 閃乱カグラ NEW LINK",
						},
					},
				},
			},
		}
		messages, err := onVoiceStateUpdateFunc(
			ctx,
			&repository.RepositoryFuncMock{
				GetVcSignalNgUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
				GetVcSignalNgRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
				GetVcSignalChannelAllColumnByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) (*repository.VcSignalChannelAllColumn, error) {
					return &repository.VcSignalChannelAllColumn{
						VcChannelID:     vcChannelID,
						GuildID:         afterGuildId,
						SendSignal:      true,
						SendChannelID:   afterSendChannelId,
						JoinBot:         false,
						EveryoneMention: false,
					}, nil
				},
				GetVcSignalMentionUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
				GetVcSignalMentionRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
			},
			&mock.SessionMock{
				ChannelMessageSendFunc: func(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error) {
					return &discordgo.Message{
						Content: content,
					}, nil
				},
				ChannelMessageSendEmbedFunc: func(channelID string, embed *discordgo.MessageEmbed, options ...discordgo.RequestOption) (*discordgo.Message, error) {
					return &discordgo.Message{
						Embeds: []*discordgo.MessageEmbed{embed},
					}, nil
				},
			},
			discordState,
			&discordgo.VoiceStateUpdate{
				VoiceState: &discordgo.VoiceState{
					GuildID:   afterGuildId,
					ChannelID: afterChannelId,
					UserID:    testUser.ID,
					Member: &discordgo.Member{
						User: testUser,
					},
					SelfStream: true,
					SelfVideo:  false,
					SelfMute:   false,
					SelfDeaf:   false,
				},
				BeforeUpdate: &discordgo.VoiceState{
					GuildID:   afterGuildId,
					ChannelID: afterChannelId,
					UserID:    testUser.ID,
					Member: &discordgo.Member{
						User: testUser,
					},
					SelfStream: false,
					SelfVideo:  false,
					SelfMute:   false,
					SelfDeaf:   false,
				},
			},
		)
		assert.NoError(t, err)
		assert.Len(t, messages, 2)
		assert.Equal(t, messages[0].Content, "<@11> がafter_test_vcで「シノビマスター 閃乱カグラ NEW LINK」を配信開始しました。")
		assert.Equal(t, messages[1].Embeds[0].Title, "配信タイトル:シノビマスター 閃乱カグラ NEW LINK")
		assert.Contains(t, messages[1].Embeds[0].Description, "testuser")
		assert.Contains(t, messages[1].Embeds[0].Description, "<#2222>")
		assert.Equal(t, messages[1].Embeds[0].Author.Name, "testuser")
		assert.Equal(t, messages[1].Embeds[0].Author.IconURL, "https://cdn.discordapp.com/avatars/11/a_.gif?size=64")
		assert.Equal(t, messages[1].Embeds[0].Image.URL, "https://cdn.discordapp.com/app-assets/123456789012345678/shinomas.png")
	})

	t.Run("正常系(画像なしアクティビティ画面共有開始)", func(t *testing.T) {
		discordState.Guilds[0].VoiceStates = []*discordgo.VoiceState{
			{
				GuildID:   afterGuildId,
				ChannelID: afterChannelId,
				Member: &discordgo.Member{
					User: testUser,
				},
				SelfStream: false,
				SelfVideo:  false,
				SelfMute:   false,
				SelfDeaf:   false,
			},
		}
		discordState.Guilds[0].Presences = []*discordgo.Presence{
			{
				User: &discordgo.User{
					ID:       testUser.ID,
					Username: testUser.Username,
					Avatar:   testUser.Avatar,
					Bot:      false,
				},
				Activities: []*discordgo.Activity{
					{
						Name:          "Devil May Cry 5 Special Edition",
						Type:          discordgo.ActivityTypeStreaming,
						ApplicationID: "123456789012345678",
					},
				},
			},
		}
		messages, err := onVoiceStateUpdateFunc(
			ctx,
			&repository.RepositoryFuncMock{
				GetVcSignalNgUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
				GetVcSignalNgRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
				GetVcSignalChannelAllColumnByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) (*repository.VcSignalChannelAllColumn, error) {
					return &repository.VcSignalChannelAllColumn{
						VcChannelID:     vcChannelID,
						GuildID:         afterGuildId,
						SendSignal:      true,
						SendChannelID:   afterSendChannelId,
						JoinBot:         false,
						EveryoneMention: false,
					}, nil
				},
				GetVcSignalMentionUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
				GetVcSignalMentionRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
			},
			&mock.SessionMock{
				ChannelMessageSendFunc: func(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error) {
					return &discordgo.Message{
						Content: content,
					}, nil
				},
				ChannelMessageSendEmbedFunc: func(channelID string, embed *discordgo.MessageEmbed, options ...discordgo.RequestOption) (*discordgo.Message, error) {
					return &discordgo.Message{
						Embeds: []*discordgo.MessageEmbed{embed},
					}, nil
				},
			},
			discordState,
			&discordgo.VoiceStateUpdate{
				VoiceState: &discordgo.VoiceState{
					GuildID:   afterGuildId,
					ChannelID: afterChannelId,
					UserID:    testUser.ID,
					Member: &discordgo.Member{
						User: testUser,
					},
					SelfStream: true,
					SelfVideo:  false,
					SelfMute:   false,
					SelfDeaf:   false,
				},
				BeforeUpdate: &discordgo.VoiceState{
					GuildID:   afterGuildId,
					ChannelID: afterChannelId,
					UserID:    testUser.ID,
					Member: &discordgo.Member{
						User: testUser,
					},
					SelfStream: false,
					SelfVideo:  false,
					SelfMute:   false,
					SelfDeaf:   false,
				},
			},
		)
		assert.NoError(t, err)
		assert.Len(t, messages, 2)
		assert.Equal(t, messages[0].Content, "<@11> がafter_test_vcで「Devil May Cry 5 Special Edition」を配信開始しました。")
		assert.Equal(t, messages[1].Embeds[0].Title, "配信タイトル:Devil May Cry 5 Special Edition")
		assert.Contains(t, messages[1].Embeds[0].Description, "testuser")
		assert.Contains(t, messages[1].Embeds[0].Description, "<#2222>")
		assert.Equal(t, messages[1].Embeds[0].Author.Name, "testuser")
		assert.Equal(t, messages[1].Embeds[0].Author.IconURL, "https://cdn.discordapp.com/avatars/11/a_.gif?size=64")
	})

	t.Run("正常系(画面共有終了)", func(t *testing.T) {
		discordState.Guilds[0].VoiceStates = []*discordgo.VoiceState{
			{
				GuildID:   afterGuildId,
				ChannelID: afterChannelId,
				Member: &discordgo.Member{
					User: testUser,
				},
				SelfStream: true,
				SelfVideo:  false,
				SelfMute:   false,
				SelfDeaf:   false,
			},
		}
		messages, err := onVoiceStateUpdateFunc(
			ctx,
			&repository.RepositoryFuncMock{
				GetVcSignalNgUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
				GetVcSignalNgRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
				GetVcSignalChannelAllColumnByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) (*repository.VcSignalChannelAllColumn, error) {
					return &repository.VcSignalChannelAllColumn{
						VcChannelID:     vcChannelID,
						GuildID:         afterGuildId,
						SendSignal:      true,
						SendChannelID:   afterSendChannelId,
						JoinBot:         false,
						EveryoneMention: false,
					}, nil
				},
				GetVcSignalMentionUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
				GetVcSignalMentionRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
			},
			&mock.SessionMock{
				ChannelMessageSendFunc: func(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error) {
					return &discordgo.Message{
						Content: content,
					}, nil
				},
				ChannelMessageSendEmbedFunc: func(channelID string, embed *discordgo.MessageEmbed, options ...discordgo.RequestOption) (*discordgo.Message, error) {
					return &discordgo.Message{
						Embeds: []*discordgo.MessageEmbed{embed},
					}, nil
				},
			},
			discordState,
			&discordgo.VoiceStateUpdate{
				VoiceState: &discordgo.VoiceState{
					GuildID:   afterGuildId,
					ChannelID: afterChannelId,
					UserID:    testUser.ID,
					Member: &discordgo.Member{
						User: testUser,
					},
					SelfStream: false,
					SelfVideo:  false,
					SelfMute:   false,
					SelfDeaf:   false,
				},
				BeforeUpdate: &discordgo.VoiceState{
					GuildID:   afterGuildId,
					ChannelID: afterChannelId,
					UserID:    testUser.ID,
					Member: &discordgo.Member{
						User: testUser,
					},
					SelfStream: true,
					SelfVideo:  false,
					SelfMute:   false,
					SelfDeaf:   false,
				},
			},
		)
		assert.NoError(t, err)
		assert.Len(t, messages, 1)
		assert.Equal(t, messages[0].Content, "<@11> が画面共有を終了しました。")
	})

	t.Run("正常系(画面共有時退室&通話終了)", func(t *testing.T) {
		discordState.Guilds[0].VoiceStates = []*discordgo.VoiceState{
			{
				GuildID:   afterGuildId,
				ChannelID: afterChannelId,
				Member: &discordgo.Member{
					User: testUser,
				},
				SelfStream: true,
				SelfVideo:  false,
				SelfMute:   false,
				SelfDeaf:   false,
			},
		}
		messages, err := onVoiceStateUpdateFunc(
			ctx,
			&repository.RepositoryFuncMock{
				GetVcSignalNgUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
				GetVcSignalNgRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
				GetVcSignalChannelAllColumnByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) (*repository.VcSignalChannelAllColumn, error) {
					return &repository.VcSignalChannelAllColumn{
						VcChannelID:     vcChannelID,
						GuildID:         afterGuildId,
						SendSignal:      true,
						SendChannelID:   afterSendChannelId,
						JoinBot:         false,
						EveryoneMention: false,
					}, nil
				},
				GetVcSignalMentionUserIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
				GetVcSignalMentionRoleIDsByVcChannelIDFunc: func(ctx context.Context, vcChannelID string) ([]string, error) {
					return []string{}, nil
				},
			},
			&mock.SessionMock{
				ChannelMessageSendFunc: func(channelID string, content string, options ...discordgo.RequestOption) (*discordgo.Message, error) {
					return &discordgo.Message{
						Content: content,
					}, nil
				},
				ChannelMessageSendEmbedFunc: func(channelID string, embed *discordgo.MessageEmbed, options ...discordgo.RequestOption) (*discordgo.Message, error) {
					return &discordgo.Message{
						Embeds: []*discordgo.MessageEmbed{embed},
					}, nil
				},
			},
			discordState,
			&discordgo.VoiceStateUpdate{
				VoiceState: &discordgo.VoiceState{
					GuildID:   "",
					ChannelID: "",
					UserID:    testUser.ID,
					Member: &discordgo.Member{
						User: testUser,
					},
					SelfStream: true,
					SelfVideo:  false,
				},
				BeforeUpdate: &discordgo.VoiceState{
					GuildID:   beforeGuildId,
					ChannelID: beforeChannelId,
					UserID:    testUser.ID,
					Member: &discordgo.Member{
						User: testUser,
					},
					SelfStream: true,
					SelfVideo:  false,
				},
			},
		)
		assert.NoError(t, err)
		assert.Len(t, messages, 3)
		assert.Equal(t, messages[0].Content, "現在0人 <@11> が before_test_vcから退室しました。")
		assert.Equal(t, messages[1].Content, "通話が終了しました。")
		assert.Equal(t, messages[2].Embeds[0].Title, "通話終了")
	})

	t.Run("正常系(ミュートの場合何もしない)", func(t *testing.T) {
		messages, err := onVoiceStateUpdateFunc(
			ctx,
			&repository.RepositoryFuncMock{},
			&mock.SessionMock{},
			discordState,
			&discordgo.VoiceStateUpdate{
				VoiceState: &discordgo.VoiceState{
					GuildID:   afterGuildId,
					ChannelID: afterChannelId,
					UserID:    testUser.ID,
					Member: &discordgo.Member{
						User: testUser,
					},
					SelfStream: false,
					SelfVideo:  false,
					SelfMute:   true,
					SelfDeaf:   false,
				},
				BeforeUpdate: &discordgo.VoiceState{
					GuildID:   afterGuildId,
					ChannelID: afterChannelId,
					UserID:    testUser.ID,
					Member: &discordgo.Member{
						User: testUser,
					},
					SelfStream: false,
					SelfVideo:  false,
					SelfMute:   false,
					SelfDeaf:   false,
				},
			},
		)
		assert.NoError(t, err)
		assert.Len(t, messages, 0)
	})
}
