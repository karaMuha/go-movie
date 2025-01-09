package core

import (
	"github.com/karaMuha/go-movie/rating/internal/core/commands"
	"github.com/karaMuha/go-movie/rating/internal/core/ports/driven"
	"github.com/karaMuha/go-movie/rating/internal/core/ports/driving"
	"github.com/karaMuha/go-movie/rating/internal/core/queries"
)

type Application struct {
	appCommands
	appQueries
}

var _ driving.IApplication = (*Application)(nil)

type appCommands struct {
	commands.SubmitRatingCommand
}

type appQueries struct {
	queries.GetAggregatedRatingQuery
}

func New(ratingRepo driven.IRatingRepository) Application {
	return Application{
		appCommands: appCommands{
			SubmitRatingCommand: commands.NewSubmitRatingCommand(ratingRepo),
		},
		appQueries: appQueries{
			GetAggregatedRatingQuery: queries.NewGetAggregatedRatingQuery(ratingRepo),
		},
	}
}
