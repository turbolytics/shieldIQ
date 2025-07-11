package engine

import "github.com/turbolytics/sqlsec/internal/db"

type Engine struct {
	db *db.Queries
}

func NewEngine(queries *db.Queries) *Engine {
	return &Engine{
		db: queries,
	}
}

func (e *Engine) Run() error {
	// The purpose of this method is to start the engine in daemon mode.

	// Begin consuming events from the database through polling.
	// In the future this will be replaced with a more efficient event-driven approach.

	// Pass the events to the rule engine for processing.

	// Alert on the results of the rule evaluation.
	return nil
}

/*
The engine is the rules engine, responsible for processing events and evaluating rules.
*/
