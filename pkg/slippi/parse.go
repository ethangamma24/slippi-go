package slippi

import (
	goslippi "github.com/ethangamma24/slippi-go/internal/goslippi"
	types "github.com/ethangamma24/slippi-go/pkg/slippi/types"
)

// ParseMetaFromBytes parses just the UBJSON metadata block from a .slp byte slice.
// This is significantly cheaper than full parse + Game.MetadataTyped when only
// metadata is needed (e.g. dashboard listings).
func ParseMetaFromBytes(name string, data []byte) (types.Metadata, error) {
	return goslippi.ParseMetaFromBytes(name, data)
}
