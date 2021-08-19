package monitor

import (
	"bytes"
	"encoding/gob"

	"github.com/MadBase/MadNet/blockchain/objects"
	"github.com/MadBase/MadNet/consensus/db"
	"github.com/MadBase/MadNet/logging"
	"github.com/MadBase/MadNet/utils"
	"github.com/dgraph-io/badger/v2"
	"github.com/sirupsen/logrus"
)

// import (
// 	"bytes"
// 	"context"
// 	"encoding/gob"

// 	"github.com/MadBase/MadNet/blockchain/objects"
// 	"github.com/MadBase/MadNet/logging"
// 	"github.com/MadBase/MadNet/utils"
// 	"github.com/dgraph-io/badger/v2"
// 	"github.com/sirupsen/logrus"
// )

var stateKey = []byte("monitorStateKey")

// Database describes required functionality for monitor persistence
type Database interface {
	FindState() (*objects.MonitorState, error)
	UpdateState(state *objects.MonitorState) error
}

type monitorDB struct {
	database *db.Database
	logger   *logrus.Entry
}

// NewDatabase initializes a new monitor database
func NewDatabase(db *db.Database) Database {
	logger := logging.GetLogger("monitor").WithField("Component", "database")
	return &monitorDB{
		logger:   logger,
		database: db}
}

func (mon *monitorDB) FindState() (*objects.MonitorState, error) {

	state := &objects.MonitorState{}

	fn := func(txn *badger.Txn) error {
		data, err := utils.GetValue(txn, stateKey)
		if err != nil {
			return err
		}

		buf := bytes.NewBuffer(data)
		dec := gob.NewDecoder(buf)
		err = dec.Decode(state)
		if err != nil {
			return err
		}
		return nil
	}

	err := mon.database.View(fn)
	if err != nil {
		return nil, err
	}

	return state, nil
}

func (mon *monitorDB) UpdateState(state *objects.MonitorState) error {

	buf := &bytes.Buffer{}

	enc := gob.NewEncoder(buf)
	err := enc.Encode(state)
	if err != nil {
		return err
	}

	fn := func(txn *badger.Txn) error {
		return utils.SetValue(txn, stateKey, buf.Bytes())
	}

	return mon.database.Update(fn)
}
