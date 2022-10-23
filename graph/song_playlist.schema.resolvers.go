package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"

	"github.com/dan6erbond/jamboree-api/graph/generated"
	"github.com/dan6erbond/jamboree-api/graph/model"
	"github.com/dan6erbond/jamboree-api/pkg/models"
	"gorm.io/gorm"
)

// AddSongPlaylist is the resolver for the addSongPlaylist field.
func (r *mutationResolver) AddSongPlaylist(ctx context.Context, payload model.AddSongPlaylistPayload) (*models.SongPlaylist, error) {
	var party models.Party
	tx := r.db.Where("name = ?", payload.PartyName).First(&party)
	if tx.Error != nil {
		return nil, tx.Error
	}
	songPlaylist := models.SongPlaylist{
		Link:      payload.Link,
		PartyName: payload.PartyName,
	}
	tx = r.db.Create(&songPlaylist)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &songPlaylist, nil
}

// EditSongPlaylist is the resolver for the editSongPlaylist field.
func (r *mutationResolver) EditSongPlaylist(ctx context.Context, payload model.EditSongPlaylistPayload) (*models.SongPlaylist, error) {
	var songPlaylist models.SongPlaylist
	tx := r.db.Where("id = ?", payload.ID).First(&songPlaylist)
	if tx.Error != nil {
		return nil, tx.Error
	}
	songPlaylistUpdate := map[string]interface{}{
		"link": payload.Link,
	}
	tx = r.db.Model(&songPlaylist).Updates(songPlaylistUpdate)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &songPlaylist, nil
}

// ToggleSongPlaylistVote is the resolver for the toggleSongPlaylistVote field.
func (r *mutationResolver) ToggleSongPlaylistVote(ctx context.Context, songPlaylistID int, username string) (*models.SongPlaylistVote, error) {
	var songPlaylistVote models.SongPlaylistVote
	tx := r.db.Where("song_playlist_id = ? AND username = ?").First(&songPlaylistVote)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			songPlaylistVote = models.SongPlaylistVote{
				SongPlaylistID: songPlaylistID,
				Username:       username,
			}
			tx = r.db.Create(&songPlaylistVote)
			if tx.Error != nil {
				return nil, tx.Error
			}
			return &songPlaylistVote, nil
		}
		return nil, tx.Error
	}
	tx = r.db.Delete(&songPlaylistVote)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return nil, nil
}

// DeleteSongPlaylist is the resolver for the deleteSongPlaylist field.
func (r *mutationResolver) DeleteSongPlaylist(ctx context.Context, songPlaylistID int) (*model.DeleteSongPlaylistResult, error) {
	var songPlaylist models.SongPlaylist
	tx := r.db.Where("id = ?", songPlaylistID).First(&songPlaylist)
	if tx.Error != nil {
		return nil, tx.Error
	}
	tx = r.db.Delete(&songPlaylist)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &model.DeleteSongPlaylistResult{
		Success: true,
	}, nil
}

// ID is the resolver for the id field.
func (r *songPlaylistResolver) ID(ctx context.Context, obj *models.SongPlaylist) (int, error) {
	return int(obj.ID), nil
}

// Votes is the resolver for the votes field.
func (r *songPlaylistResolver) Votes(ctx context.Context, obj *models.SongPlaylist) ([]*models.SongPlaylistVote, error) {
	var votes []*models.SongPlaylistVote
	tx := r.db.Where("song_playlist_id = ?", obj.ID).Find(&votes)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return votes, nil
}

// ID is the resolver for the id field.
func (r *songPlaylistVoteResolver) ID(ctx context.Context, obj *models.SongPlaylistVote) (int, error) {
	return int(obj.ID), nil
}

// SongPlaylist returns generated.SongPlaylistResolver implementation.
func (r *Resolver) SongPlaylist() generated.SongPlaylistResolver { return &songPlaylistResolver{r} }

// SongPlaylistVote returns generated.SongPlaylistVoteResolver implementation.
func (r *Resolver) SongPlaylistVote() generated.SongPlaylistVoteResolver {
	return &songPlaylistVoteResolver{r}
}

type songPlaylistResolver struct{ *Resolver }
type songPlaylistVoteResolver struct{ *Resolver }
