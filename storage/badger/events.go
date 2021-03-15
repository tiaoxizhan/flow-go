package badger

import (
	"fmt"

	"github.com/dgraph-io/badger/v2"

	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/storage"
	"github.com/onflow/flow-go/storage/badger/operation"
)

type Events struct {
	db *badger.DB
}

func NewEvents(db *badger.DB) *Events {
	return &Events{
		db: db,
	}
}

func (e *Events) BatchStore(blockID flow.Identifier, events []flow.Event, batch storage.BatchStorage) error {
	if writeBatch, ok := batch.(*badger.WriteBatch); ok {
		for _, event := range events {
			err := operation.BatchInsertEvent(blockID, event)(writeBatch)
			if err != nil {
				return fmt.Errorf("cannot batch insert event: %w", err)
			}
		}
		return nil
	}
	return fmt.Errorf("unsupported BatchStore type %T", batch)
}

// ByBlockID returns the events for the given block ID
func (e *Events) ByBlockID(blockID flow.Identifier) ([]flow.Event, error) {

	var events []flow.Event
	err := e.db.View(operation.LookupEventsByBlockID(blockID, &events))
	if err != nil {
		return nil, handleError(err, flow.Event{})
	}

	return events, nil
}

// ByBlockIDTransactionID returns the events for the given block ID and transaction ID
func (e *Events) ByBlockIDTransactionID(blockID flow.Identifier, txID flow.Identifier) ([]flow.Event, error) {

	var events []flow.Event
	err := e.db.View(operation.RetrieveEvents(blockID, txID, &events))
	if err != nil {
		return nil, handleError(err, flow.Event{})
	}

	return events, nil
}

// ByBlockIDEventType returns the events for the given block ID and event type
func (e *Events) ByBlockIDEventType(blockID flow.Identifier, event flow.EventType) ([]flow.Event, error) {

	var events []flow.Event
	err := e.db.View(operation.LookupEventsByBlockIDEventType(blockID, event, &events))
	if err != nil {
		return nil, handleError(err, flow.Event{})
	}

	return events, nil
}

type ServiceEvents struct {
	db *badger.DB
}

func NewServiceEvents(db *badger.DB) *ServiceEvents {
	return &ServiceEvents{
		db: db,
	}
}

func (e *ServiceEvents) BatchStore(blockID flow.Identifier, events []flow.Event, batch storage.BatchStorage) error {
	if writeBatch, ok := batch.(*badger.WriteBatch); ok {
		for _, event := range events {
			err := operation.BatchInsertServiceEvent(blockID, event)(writeBatch)
			if err != nil {
				return fmt.Errorf("cannot batch insert service event: %w", err)
			}
		}
		return nil
	}
	return fmt.Errorf("unsupported BatchStore type %T", batch)
}

// ByBlockID returns the events for the given block ID
func (e *ServiceEvents) ByBlockID(blockID flow.Identifier) ([]flow.Event, error) {

	var events []flow.Event
	err := e.db.View(operation.LookupServiceEventsByBlockID(blockID, &events))
	if err != nil {
		return nil, handleError(err, flow.Event{})
	}

	return events, nil
}
