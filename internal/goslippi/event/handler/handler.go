package handler

import (
	"github.com/ethangamma24/slippi-go/pkg/slippi/types"
	"github.com/ethangamma24/slippi-go/internal/goslippi/event"
)

// EventHandler defines the behaviour for parsing Slippi events.
type EventHandler interface {
	Parse(dec *event.Decoder, data *slippi.Data) error
}
