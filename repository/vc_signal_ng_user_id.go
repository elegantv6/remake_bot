package repository

import (
	"context"
)

type VcSignalNgUserAllColumn struct {
	VcChannelID string `db:"vc_channel_id"`
	GuildID     string `db:"guild_id"`
	UserID      string `db:"user_id"`
}

func (r *Repository) InsertVcSignalNgUser(ctx context.Context, vcChannelID, guildID, userID string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO vc_signal_ng_user_id (
			vc_channel_id,
			guild_id,
			user_id
		) VALUES (
			$1,
			$2,
			$3
		)
	`, vcChannelID, guildID, userID)
	return err
}

func (r *Repository) GetVcSignalNgUsersByChannelIDAllColumn(ctx context.Context, channelID string) ([]*VcSignalNgUserAllColumn, error) {
	var ngUserIDs []*VcSignalNgUserAllColumn
	err := r.db.SelectContext(ctx, &ngUserIDs, `
		SELECT
			*
		FROM
			vc_signal_ng_user_id
		WHERE
			vc_channel_id = $1
	`, channelID)
	if err != nil {
		return nil, err
	}
	return ngUserIDs, nil
}

func (r *Repository) DeleteVcNgUserByChannelID(ctx context.Context, vcChannelID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM
			vc_signal_ng_user_id
		WHERE
			vc_channel_id = $1
	`, vcChannelID)
	return err
}
func (r *Repository) DeleteVcNgUserByUserID(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM
			vc_signal_ng_user_id
		WHERE
			user_id = $1
	`, userID)
	return err
}
