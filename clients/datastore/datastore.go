package datastore

import (
	"errors"
	"time"

	"github.com/jcherianucla/bloggo/.gen/idl/proto"

	"github.com/go-pg/pg"
	"github.com/jcherianucla/bloggo/clients/instrumenter"
	"github.com/jcherianucla/bloggo/config"
	"github.com/jcherianucla/bloggo/utils"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	INSERT       = "insert"
	UPDATE       = "update"
	SELECT       = "select"
	CREATE       = "create"
	DELETE       = "delete"
	_modelLogTag = "model"
	_opLogTag    = "operation"
	_maxRetries  = 4
	_minBackoff  = 250 * time.Millisecond
)

// Module provides the Datastore through Fx
var Module = fx.Provide(New)

// Params defines the input dependencies for the Datastore
type Params struct {
	fx.In
	config.AppConfig
	instrumenter.Instrument
}

// Result defines the output dependency that is the Datastore
type Result struct {
	fx.Out
	Store
}

// New creates a connection to the running data store server in the service
// and returns a new data store
func New(p Params) Result {
	_db := pg.Connect(&pg.Options{
		User:                  p.Config().DBConfig.User,
		Password:              p.Config().DBConfig.Password,
		Database:              p.Config().DBConfig.Name,
		Addr:                  p.Config().DBConfig.HostPort,
		RetryStatementTimeout: true,
		MaxRetries:            _maxRetries,
		MinRetryBackoff:       _minBackoff,
	})
	post := &proto.Post{}
	p.Logger(utils.DebugLogType).Info("Creating database",
		commonLogTags(CREATE, post)...,
	)
	if err := _db.CreateTable(post, nil); err != nil {
		p.Logger(utils.DebugLogType).Fatal("Failed to create database",
			commonLogTags(CREATE, post, zap.Error(err))...,
		)
		utils.HandleErr(err)
	}
	return Result{Store: &store{db: _db, Instrument: p.Instrument}}
}

// Store defines a simple wrapper on top of the underlying data store
type Store interface {
	// DB exposes the singleton reference to the data store
	DB() *pg.DB
	// Operate provides a simple factory-like ability to perform ORM operations
	Operate(op string, obj interface{}) error
}

type store struct {
	db *pg.DB
	instrumenter.Instrument
}

func (s *store) DB() *pg.DB {
	s.Logger(utils.DebugLogType).Info("Get access to DB")
	return s.db
}

func (s *store) Operate(op string, obj interface{}) error {
	s.Logger(utils.DebugLogType).Info(utils.JoinStrings(op, " object"),
		commonLogTags(op, obj)...,
	)
	var err error
	switch op {
	case INSERT:
		err = s.db.Insert(obj)
	case UPDATE:
		err = s.db.Update(obj)
	case SELECT:
		err = s.db.Select(obj)
	case DELETE:
		err = s.db.Delete(obj)
	default:
		err = errors.New("invalid operation")
	}
	if err != nil {
		s.Logger(utils.DebugLogType).Error(utils.JoinStrings("Failed to ", op),
			commonLogTags(op, obj, zap.Error(err))...,
		)
		return err
	}
	return nil
}

// commonLogTags includes the basic operation type and model log tags, along
// with any extra arguments that might want to be logged
func commonLogTags(op string, obj interface{}, args ...zap.Field) []zap.Field {
	baseTags := []zap.Field{
		zap.String(_opLogTag, op),
		zap.Any(_modelLogTag, obj),
	}
	baseTags = append(baseTags, args...)
	return baseTags
}
