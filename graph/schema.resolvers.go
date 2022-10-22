package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/dan6erbond/jamboree-api/graph/generated"
	graphqlModel "github.com/dan6erbond/jamboree-api/graph/model"
	randomstring "github.com/dan6erbond/jamboree-api/internal/random_string"
	uniquename "github.com/dan6erbond/jamboree-api/internal/unique_name"
	"github.com/dan6erbond/jamboree-api/pkg/auth"
	"github.com/dan6erbond/jamboree-api/pkg/models"
	"gorm.io/gorm"
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

// EditParty is the resolver for the editParty field.
func (r *mutationResolver) EditParty(ctx context.Context, partyOptions graphqlModel.EditPartyRequest) (*models.Party, error) {
	var party models.Party
	tx := r.db.First(&party).Where("name = ?", partyOptions.PartyName)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if *auth.ForContext(ctx) == party.AdminCode {
		return nil, fmt.Errorf("cannot edit party without an admin code")
	}
	partyUpdate := models.Party{}
	if partyOptions.DateOptionsEnabled != nil {
		partyUpdate.DateOptionsEnabled = *partyOptions.DateOptionsEnabled
	}
	if partyOptions.DateVotingEnabled != nil {
		partyUpdate.DateVotingEnabled = *partyOptions.DateVotingEnabled
	}
	if partyOptions.LocationOptionsEnabled != nil {
		partyUpdate.LocationOptionsEnabled = *partyOptions.LocationOptionsEnabled
	}
	if partyOptions.LocationVotingEnabled != nil {
		partyUpdate.LocationVotingEnabled = *partyOptions.LocationVotingEnabled
	}
	tx = r.db.Model(&party).Updates(partyUpdate)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &party, nil
}

// AddLocation is the resolver for the addLocation field.
func (r *mutationResolver) AddLocation(ctx context.Context, partyName string, location string) (*models.PartyLocation, error) {
	var party models.Party
	tx := r.db.First(&party).Where("name = ?", partyName)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if !party.LocationVotingEnabled {
		return nil, fmt.Errorf("cannot add locations to this party, please ask an admin to enable this")
	}
	if !party.LocationOptionsEnabled && *auth.ForContext(ctx) != party.AdminCode {
		return nil, fmt.Errorf("user location options not enabled for this party, please ask an admin to enable this")
	}
	partyLocation := models.PartyLocation{
		Location:  location,
		PartyName: partyName,
	}
	tx = r.db.Create(&partyLocation)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &partyLocation, nil
}

// AddDate is the resolver for the addDate field.
func (r *mutationResolver) AddDate(ctx context.Context, partyName string, date string) (*models.PartyDate, error) {
	var party models.Party
	tx := r.db.First(&party).Where("name = ?", partyName)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if !party.DateVotingEnabled {
		return nil, fmt.Errorf("cannot add dates to this party, please ask an admin to enable this")
	}
	if !party.DateOptionsEnabled && *auth.ForContext(ctx) != party.AdminCode {
		return nil, fmt.Errorf("user date options not enabled for this party, please ask an admin to enable this")
	}
	i, err := strconv.ParseInt(date, 10, 64)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(i, 0)
	partyDate := models.PartyDate{
		Date:      &tm,
		PartyName: partyName,
	}
	tx = r.db.Create(&partyDate)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &partyDate, nil
}

// AddSupply is the resolver for the addSupply field.
func (r *mutationResolver) AddSupply(ctx context.Context, payload graphqlModel.AddSupplyPayload) (*models.Supply, error) {
	var party models.Party
	tx := r.db.First(&party).Where("name = ?", payload.PartyName)
	if tx.Error != nil {
		return nil, tx.Error
	}
	supply := models.Supply{
		Name:      payload.Name,
		IsUrgent:  payload.IsUrgent,
		Emoji:     payload.Emoji,
		PartyName: party.Name,
	}
	if payload.Assignee != nil {
		supply.Assignee = *payload.Assignee
	}
	if payload.Quantity != nil {
		supply.Quantity = int32(*payload.Quantity)
	} else {
		supply.Quantity = 1
	}
	tx = r.db.Create(&supply)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &supply, nil
}

// EditSupply is the resolver for the editSupply field.
func (r *mutationResolver) EditSupply(ctx context.Context, payload graphqlModel.EditSupplyPayload) (*models.Supply, error) {
	var supply models.Supply
	tx := r.db.First(&supply).Where("id = ?", payload.ID)
	if tx.Error != nil {
		return nil, tx.Error
	}
	supplyUpdate := models.Supply{}
	if payload.Emoji != nil {
		supplyUpdate.Emoji = *payload.Emoji
	}
	if payload.IsUrgent != nil {
		supplyUpdate.IsUrgent = *payload.IsUrgent
	}
	if payload.Name != nil {
		supplyUpdate.Name = *payload.Name
	}
	if payload.Quantity != nil {
		if *payload.Quantity <= 0 {
			return nil, fmt.Errorf("supplyUpdate quantity cannot be below zero")
		}
		supplyUpdate.Quantity = int32(*payload.Quantity)
	}
	tx = r.db.Model(&supply).Updates(supplyUpdate)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &supply, nil
}

// AssignSupply is the resolver for the assignSupply field.
func (r *mutationResolver) AssignSupply(ctx context.Context, supplyID int, username string) (*models.Supply, error) {
	var supply models.Supply
	tx := r.db.First(&supply).Where("id = ?", supplyID)
	if tx.Error != nil {
		return nil, tx.Error
	}
	supplyUpdate := models.Supply{}
	supplyUpdate.Assignee = username
	tx = r.db.Model(&supply).Updates(supplyUpdate)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &supply, nil
}

// DeleteSupply is the resolver for the deleteSupply field.
func (r *mutationResolver) DeleteSupply(ctx context.Context, supplyID int) (*graphqlModel.DeleteSupplyResult, error) {
	var supply models.Supply
	tx := r.db.First(&supply).Where("id = ?", supplyID)
	if tx.Error != nil {
		return nil, tx.Error
	}
	tx = r.db.Where("id = ?", supplyID).Delete(&supply)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &graphqlModel.DeleteSupplyResult{
		Success: true,
	}, nil
}

// ToggleDateVote is the resolver for the toggleDateVote field.
func (r *mutationResolver) ToggleDateVote(ctx context.Context, partyDateID int, username string) (*models.PartyDateVote, error) {
	var partyDateVote models.PartyDateVote
	tx := r.db.First(&partyDateVote).Where("username = ?", username)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			partyDateVote = models.PartyDateVote{
				PartyDateID: partyDateID,
				Username:    username,
			}
			tx = r.db.Create(&partyDateVote)
			if tx.Error != nil {
				return nil, tx.Error
			}
			return &partyDateVote, nil
		}
		return nil, tx.Error
	}
	tx = r.db.Delete(&partyDateVote)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return nil, nil
}

// ToggleLocationVote is the resolver for the toggleLocationVote field.
func (r *mutationResolver) ToggleLocationVote(ctx context.Context, partyLocationID int, username string) (*models.PartyLocationVote, error) {
	var partyLocationVote models.PartyLocationVote
	tx := r.db.First(&partyLocationVote).Where("username = ?", username)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			partyLocationVote = models.PartyLocationVote{
				PartyLocationID: partyLocationID,
				Username:        username,
			}
			tx = r.db.Create(&partyLocationVote)
			if tx.Error != nil {
				return nil, tx.Error
			}
			return &partyLocationVote, nil
		}
		return nil, tx.Error
	}
	tx = r.db.Delete(&partyLocationVote)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return nil, nil
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

// Dates is the resolver for the dates field.
func (r *partyResolver) Dates(ctx context.Context, obj *models.Party) ([]*models.PartyDate, error) {
	var dates []*models.PartyDate
	tx := r.db.Find(&dates).Where("party_name = ?", obj.Name)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return dates, nil
}

// Locations is the resolver for the locations field.
func (r *partyResolver) Locations(ctx context.Context, obj *models.Party) ([]*models.PartyLocation, error) {
	var locations []*models.PartyLocation
	tx := r.db.Find(&locations).Where("party_name = ?", obj.Name)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return locations, nil
}

// ID is the resolver for the id field.
func (r *partyDateResolver) ID(ctx context.Context, obj *models.PartyDate) (int, error) {
	return int(obj.ID), nil
}

// Date is the resolver for the date field.
func (r *partyDateResolver) Date(ctx context.Context, obj *models.PartyDate) (string, error) {
	return obj.Date.String(), nil
}

// Votes is the resolver for the votes field.
func (r *partyDateResolver) Votes(ctx context.Context, obj *models.PartyDate) ([]*models.PartyDateVote, error) {
	panic(fmt.Errorf("not implemented: Votes - votes"))
}

// ID is the resolver for the id field.
func (r *partyDateVoteResolver) ID(ctx context.Context, obj *models.PartyDateVote) (int, error) {
	return int(obj.ID), nil
}

// ID is the resolver for the id field.
func (r *partyLocationResolver) ID(ctx context.Context, obj *models.PartyLocation) (int, error) {
	return int(obj.ID), nil
}

// Votes is the resolver for the votes field.
func (r *partyLocationResolver) Votes(ctx context.Context, obj *models.PartyLocation) ([]*models.PartyLocationVote, error) {
	panic(fmt.Errorf("not implemented: Votes - votes"))
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

// Votes is the resolver for the votes field.
func (r *songPlaylistResolver) Votes(ctx context.Context, obj *models.SongPlaylist) ([]*models.SongPlaylistVote, error) {
	panic(fmt.Errorf("not implemented: Votes - votes"))
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
