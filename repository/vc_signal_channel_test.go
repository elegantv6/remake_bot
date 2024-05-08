package repository

import (
	"context"
	"testing"

	"github.com/maguro-alternative/remake_bot/bot/config"
	"github.com/maguro-alternative/remake_bot/pkg/db"
	"github.com/maguro-alternative/remake_bot/testutil/fixtures"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsertVcSignalChannel(t *testing.T) {
	ctx := context.Background()
	t.Run("ChannelIDを追加できること", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_channel")

		repo := NewRepository(tx)
		err = repo.InsertVcSignalChannel(ctx, "123456789", "987654321", "1234567890")
		assert.NoError(t, err)

		var channels []VcSignalChannelAllColumn
		err = tx.SelectContext(ctx, &channels, "SELECT * FROM vc_signal_channel")
		assert.NoError(t, err)
		assert.Len(t, channels, 1)
		assert.Equal(t, "123456789", channels[0].VcChannelID)
		assert.Equal(t, "987654321", channels[0].GuildID)
		assert.Equal(t, "1234567890", channels[0].SendChannelID)
	})

	t.Run("ChannelIDが重複している場合はエラーは返さず挿入しないこと", func(t *testing.T) {
		dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
		assert.NoError(t, err)
		defer cleanup()
		tx, err := dbV1.BeginTxx(ctx, nil)
		assert.NoError(t, err)

		defer tx.RollbackCtx(ctx)

		tx.ExecContext(ctx, "DELETE FROM vc_signal_channel")

		repo := NewRepository(tx)
		err = repo.InsertVcSignalChannel(ctx, "123456789", "987654321", "1234567890")
		assert.NoError(t, err)

		err = repo.InsertVcSignalChannel(ctx, "123456789", "987654321", "1234567890")
		assert.NoError(t, err)

		var channels []VcSignalChannelAllColumn
		err = tx.SelectContext(ctx, &channels, "SELECT * FROM vc_signal_channel")
		assert.NoError(t, err)
		assert.Len(t, channels, 1)
		assert.Equal(t, "123456789", channels[0].VcChannelID)
		assert.Equal(t, "987654321", channels[0].GuildID)
		assert.Equal(t, "1234567890", channels[0].SendChannelID)
	})
}

func TestGetVcSignalChannel(t *testing.T) {
	ctx := context.Background()
	dbV1, cleanup, err := db.NewDB(ctx, config.DatabaseName(), config.DatabaseURL())
	require.NoError(t, err)
	defer cleanup()
	tx, err := dbV1.BeginTxx(ctx, nil)
	require.NoError(t, err)

	defer tx.RollbackCtx(ctx)

	tx.ExecContext(ctx, "DELETE FROM vc_signal_channel")

	f := &fixtures.Fixture{DBv1: tx}
	f.Build(t,
		fixtures.NewVcSignalChannel(ctx, func(v *fixtures.VcSignalChannel) {
			v.VcChannelID = "111"
			v.GuildID = "1111"
			v.SendChannelID = "11111"
		}),
		fixtures.NewVcSignalChannel(ctx, func(v *fixtures.VcSignalChannel) {
			v.VcChannelID = "222"
			v.GuildID = "2222"
			v.SendChannelID = "22222"
		}),
		fixtures.NewVcSignalChannel(ctx, func(v *fixtures.VcSignalChannel) {
			v.VcChannelID = "333"
			v.GuildID = "3333"
			v.SendChannelID = "33333"
		}),
	)

	repo := NewRepository(tx)
	t.Run("ボイスチャンネルの情報を取得できること", func(t *testing.T) {
		vcSignalChannel, err := repo.GetVcSignalChennelAllColumn(ctx, "111")
		assert.NoError(t, err)

		assert.Equal(t, true, vcSignalChannel.SendSignal)
		assert.Equal(t, false, vcSignalChannel.JoinBot)
		assert.Equal(t, true, vcSignalChannel.EveryoneMention)
	})

	t.Run("ボイスチャンネルの情報が存在しない場合はエラーを返すこと", func(t *testing.T) {
		vcSignalChannel, err := repo.GetVcSignalChennelAllColumn(ctx, "444")
		assert.Error(t, err)
		assert.Nil(t, vcSignalChannel)
	})
}
