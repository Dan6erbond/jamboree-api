package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/dan6erbond/jamboree-api/graph/generated"
	graphqlModel "github.com/dan6erbond/jamboree-api/graph/model"
	randomstring "github.com/dan6erbond/jamboree-api/internal/random_string"
	uniquename "github.com/dan6erbond/jamboree-api/internal/unique_name"
	"github.com/dan6erbond/jamboree-api/pkg/models"
)

// CreateParty is the resolver for the createParty field.
func (r *mutationResolver) CreateParty(ctx context.Context, username string) (*graphqlModel.CreatePartyResult, error) {
	name := uniquename.GenerateUniqueName()
	for {
		var count int64
		r.db.Model(&models.Party{}).Where("name = ?", name).Count(&count)
		if count == 0 {
			break
		}
		name = uniquename.GenerateUniqueName()
	}
	adminCode := randomstring.GenerateRandomString(32)
	party := models.Party{
		Name:      name,
		AdminCode: adminCode,
		Creator:   username,
	}
	tx := r.db.Create(&party)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &graphqlModel.CreatePartyResult{
		Name:      name,
		AdminCode: adminCode,
	}, nil
}

// AddLocation is the resolver for the addLocation field.
func (r *mutationResolver) AddLocation(ctx context.Context, partyName string, location string) (*models.PartyLocation, error) {
	r.db.Create(&models.PartyLocation{})
}

// Settings is the resolver for the settings field.
func (r *partyResolver) Settings(ctx context.Context, obj *models.Party) (*graphqlModel.PartySettings, error) {
	return &graphqlModel.PartySettings{
		Dates: &graphqlModel.DatePartySettings{
			VotingEnabled:  obj.DateVotingEnabled,
			OptionsEnabled: obj.DateOptionsEnabled,
		},
		Locations: &graphqlModel.LocationPartySettings{
			VotingEnabled:  obj.LocationVotingEnabled,
			OptionsEnabled: obj.LocationOptionsEnabled,
		},
	}, nil
}

// ID is the resolver for the id field.
func (r *partyDateResolver) ID(ctx context.Context, obj *models.PartyDate) (int, error) {
	return int(obj.ID), nil
}

// Date is the resolver for the date field.
func (r *partyDateResolver) Date(ctx context.Context, obj *models.PartyDate) (string, error) {
	return obj.Date.String(), nil
}

// ID is the resolver for the id field.
func (r *partyDateVoteResolver) ID(ctx context.Context, obj *models.PartyDateVote) (int, error) {
	return int(obj.ID), nil
}

// ID is the resolver for the id field.
func (r *partyLocationResolver) ID(ctx context.Context, obj *models.PartyLocation) (int, error) {
	return int(obj.ID), nil
}

// ID is the resolver for the id field.
func (r *partyLocationVoteResolver) ID(ctx context.Context, obj *models.PartyLocationVote) (int, error) {
	return int(obj.ID), nil
}

// Party is the resolver for the party field.
func (r *queryResolver) Party(ctx context.Context, name *string, adminCode *string) (*models.Party, error) {
	var party *models.Party
	if name != nil {
		tx := r.db.First(&party, "name = ?", name)
		if tx.Error != nil {
			return nil, tx.Error
		}
	} else if adminCode != nil {
		tx := r.db.First(&party, "admin_code = ?", adminCode)
		if tx.Error != nil {
			return nil, tx.Error
		}
	}
	return party, nil
}

// ID is the resolver for the id field.
func (r *songPlaylistResolver) ID(ctx context.Context, obj *models.SongPlaylist) (int, error) {
	return int(obj.ID), nil
}

// ID is the resolver for the id field.
func (r *songPlaylistVoteResolver) ID(ctx context.Context, obj *models.SongPlaylistVote) (int, error) {
	return int(obj.ID), nil
}

// ID is the resolver for the id field.
func (r *supplyResolver) ID(ctx context.Context, obj *models.Supply) (int, error) {
	return int(obj.ID), nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Party returns generated.PartyResolver implementation.
func (r *Resolver) Party() generated.PartyResolver { return &partyResolver{r} }

// PartyDate returns generated.PartyDateResolver implementation.
func (r *Resolver) PartyDate() generated.PartyDateResolver { return &partyDateResolver{r} }

// PartyDateVote returns generated.PartyDateVoteResolver implementation.
func (r *Resolver) PartyDateVote() generated.PartyDateVoteResolver { return &partyDateVoteResolver{r} }

// PartyLocation returns generated.PartyLocationResolver implementation.
func (r *Resolver) PartyLocation() generated.PartyLocationResolver { return &partyLocationResolver{r} }

// PartyLocationVote returns generated.PartyLocationVoteResolver implementation.
func (r *Resolver) PartyLocationVote() generated.PartyLocationVoteResolver {
	return &partyLocationVoteResolver{r}
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// SongPlaylist returns generated.SongPlaylistResolver implementation.
func (r *Resolver) SongPlaylist() generated.SongPlaylistResolver { return &songPlaylistResolver{r} }

// SongPlaylistVote returns generated.SongPlaylistVoteResolver implementation.
func (r *Resolver) SongPlaylistVote() generated.SongPlaylistVoteResolver {
	return &songPlaylistVoteResolver{r}
}

// Supply returns generated.SupplyResolver implementation.
func (r *Resolver) Supply() generated.SupplyResolver { return &supplyResolver{r} }

type mutationResolver struct{ *Resolver }
type partyResolver struct{ *Resolver }
type partyDateResolver struct{ *Resolver }
type partyDateVoteResolver struct{ *Resolver }
type partyLocationResolver struct{ *Resolver }
type partyLocationVoteResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type songPlaylistResolver struct{ *Resolver }
type songPlaylistVoteResolver struct{ *Resolver }
type supplyResolver struct{ *Resolver }
