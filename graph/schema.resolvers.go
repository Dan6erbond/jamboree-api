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
	adminCode := randomstring.GenerateRandomString(64)
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
	tx := r.db.Where("name = ?", partyOptions.PartyName).First(&party)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if *auth.ForContext(ctx) != party.AdminCode {
		return nil, fmt.Errorf("cannot edit party without an admin code")
	}
	updates := map[string]interface{}{}
	if partyOptions.DateOptionsEnabled != nil {
		updates["date_options_enabled"] = *partyOptions.DateOptionsEnabled
	}
	if partyOptions.DateVotingEnabled != nil {
		updates["date_voting_enabled"] = *partyOptions.DateVotingEnabled
	}
	if partyOptions.LocationOptionsEnabled != nil {
		updates["location_options_enabled"] = *partyOptions.LocationOptionsEnabled
	}
	if partyOptions.LocationVotingEnabled != nil {
		updates["location_voting_enabled"] = *partyOptions.LocationVotingEnabled
	}
	tx = r.db.Model(&party).Updates(updates)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &party, nil
}

// AddLocation is the resolver for the addLocation field.
func (r *mutationResolver) AddLocation(ctx context.Context, partyName string, location string) (*models.PartyLocation, error) {
	var party models.Party
	tx := r.db.Where("name = ?", partyName).First(&party)
	if tx.Error != nil {
		return nil, tx.Error
	}
	var partyLocations []*models.PartyLocation
	tx = r.db.Where("party_name = ?", party.Name).Find(&partyLocations)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if !party.LocationVotingEnabled && len(partyLocations) >= 1 {
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
	tx := r.db.Where("name = ?", partyName).First(&party)
	if tx.Error != nil {
		return nil, tx.Error
	}
	var partyDates []*models.PartyDate
	tx = r.db.Where("party_name = ?", party.Name).Find(&partyDates)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if !party.DateVotingEnabled && len(partyDates) >= 1 {
		return nil, fmt.Errorf("cannot add dates to this party, please ask an admin to enable this")
	}
	if !party.DateOptionsEnabled && *auth.ForContext(ctx) != party.AdminCode {
		return nil, fmt.Errorf("user date options not enabled for this party, please ask an admin to enable this")
	}
	i, err := strconv.ParseInt(date, 10, 64)
	if err != nil {
		return nil, err
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
	tx := r.db.Where("name = ?", payload.PartyName).First(&party)
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
	tx := r.db.Where("id = ?", payload.ID).First(&supply)
	if tx.Error != nil {
		return nil, tx.Error
	}
	supplyUpdate := map[string]interface{}{}
	if payload.Emoji != nil {
		supplyUpdate["emoji"] = *payload.Emoji
	}
	if payload.IsUrgent != nil {
		supplyUpdate["is_urgent"] = *payload.IsUrgent
	}
	if payload.Name != nil {
		supplyUpdate["name"] = *payload.Name
	}
	if payload.Assignee != nil {
		supplyUpdate["assignee"] = *payload.Assignee
	}
	if payload.Quantity != nil {
		if *payload.Quantity <= 0 {
			return nil, fmt.Errorf("supplyUpdate quantity cannot be below zero")
		}
		supplyUpdate["quantity"] = int32(*payload.Quantity)
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
	tx := r.db.Where("id = ?", supplyID).First(&supply)
	if tx.Error != nil {
		return nil, tx.Error
	}
	tx = r.db.Model(&supply).Updates(map[string]interface{}{"assignee": username})
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &supply, nil
}

// DeleteSupply is the resolver for the deleteSupply field.
func (r *mutationResolver) DeleteSupply(ctx context.Context, supplyID int) (*graphqlModel.DeleteSupplyResult, error) {
	var supply models.Supply
	tx := r.db.Where("id = ?", supplyID).First(&supply)
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
	tx := r.db.Where("party_date_id = ? AND username = ?", partyDateID, username).First(&partyDateVote)
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
	tx := r.db.Where("party_location_id = ? AND username = ?", partyLocationID, username).First(&partyLocationVote)
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

// EditDate is the resolver for the editDate field.
func (r *mutationResolver) EditDate(ctx context.Context, payload graphqlModel.EditDatePayload) (*models.PartyDate, error) {
	var partyDate models.PartyDate
	tx := r.db.Where("id = ?", payload.ID).First(&partyDate)
	if tx.Error != nil {
		return nil, tx.Error
	}
	i, err := strconv.ParseInt(payload.Date, 10, 64)
	if err != nil {
		return nil, err
	}
	tm := time.Unix(i, 0)
	tx = r.db.Model(&partyDate).Updates(map[string]interface{}{"date": &tm})
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &partyDate, nil
}

// EditLocation is the resolver for the editLocation field.
func (r *mutationResolver) EditLocation(ctx context.Context, payload graphqlModel.EditLocationPayload) (*models.PartyLocation, error) {
	var partyLocation models.PartyLocation
	tx := r.db.Where("id = ?", payload.ID).First(&partyLocation)
	if tx.Error != nil {
		return nil, tx.Error
	}
	tx = r.db.Model(&partyLocation).Updates(map[string]interface{}{"location": payload.Location})
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &partyLocation, nil
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
	tx := r.db.Where("party_name = ?", obj.Name).Order("id").Find(&dates)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return dates, nil
}

// Locations is the resolver for the locations field.
func (r *partyResolver) Locations(ctx context.Context, obj *models.Party) ([]*models.PartyLocation, error) {
	var locations []*models.PartyLocation
	tx := r.db.Where("party_name = ?", obj.Name).Order("id").Find(&locations)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return locations, nil
}

// Supplies is the resolver for the supplies field.
func (r *partyResolver) Supplies(ctx context.Context, obj *models.Party) ([]*models.Supply, error) {
	var supplies []*models.Supply
	tx := r.db.Where("party_name = ?", obj.Name).Order("id").Find(&supplies)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return supplies, nil
}

// SongPlaylists is the resolver for the songPlaylists field.
func (r *partyResolver) SongPlaylists(ctx context.Context, obj *models.Party) ([]*models.SongPlaylist, error) {
	var songPlaylists []*models.SongPlaylist
	tx := r.db.Where("party_name = ?", obj.Name).Order("id").Find(&songPlaylists)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return songPlaylists, nil
}

// ID is the resolver for the id field.
func (r *partyDateResolver) ID(ctx context.Context, obj *models.PartyDate) (int, error) {
	return int(obj.ID), nil
}

// Votes is the resolver for the votes field.
func (r *partyDateResolver) Votes(ctx context.Context, obj *models.PartyDate) ([]*models.PartyDateVote, error) {
	var votes []*models.PartyDateVote
	tx := r.db.Where("party_date_id = ?", obj.ID).Find(&votes)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return votes, nil
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
	var votes []*models.PartyLocationVote
	tx := r.db.Where("party_location_id = ?", obj.ID).Find(&votes)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return votes, nil
}

// ID is the resolver for the id field.
func (r *partyLocationVoteResolver) ID(ctx context.Context, obj *models.PartyLocationVote) (int, error) {
	return int(obj.ID), nil
}

// Party is the resolver for the party field.
func (r *queryResolver) Party(ctx context.Context, name *string, adminCode *string) (*models.Party, error) {
	var party *models.Party
	if name != nil {
		tx := r.db.First(&party, "name = ?", *name)
		if tx.Error != nil {
			return nil, tx.Error
		}
	} else if adminCode != nil {
		tx := r.db.First(&party, "admin_code = ?", *adminCode)
		if tx.Error != nil {
			return nil, tx.Error
		}
	} else {
		return nil, fmt.Errorf("must provide either one of name or admin code")
	}
	return party, nil
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
