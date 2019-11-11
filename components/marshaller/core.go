package marshaller

import (
	"errors"
	"fmt"

	"github.com/bluest-eel/state/components/config"
	log "github.com/sirupsen/logrus"
)

// db-wide constants
const (
	// Connection types
	GOB         string = "gob"
	GOTINY      string = "gotiny"
	GOGO        string = "gogoprotobuf"
	PROTO       string = "protobuf"
	PASSTHROUGH string = "passthrough"
)

// StateMetadata ...
type StateMetadata struct {
	Name string
}

// Marsh ...
type Marsh interface {
	Marshal(data *StateMetadata) (uint64, error)
	Unmarshal(data uint64) (StateMetadata, error)
}

// New dispatches the creation of a specific marshaller type
func New(cfg *config.Config) (Marsh, error) {
	switch cfg.Marshaller.Type {
	case GOB:
		return NewGob(), nil
	// case GOTINY:
	// 	return NewGoTiny(), nil
	default:
		errMsg := fmt.Sprintf("database connection type %s not supported", cfg.DB.Type)
		log.Error(errMsg)
		err := errors.New(errMsg)
		return NewPassthrough(), err
	}
}
