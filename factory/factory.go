package factory

import (
	"context"
	"github.com/pangpanglabs/goutils/behaviorlog"

	"github.com/go-xorm/xorm"
	"github.com/pangpanglabs/goutils/echomiddleware"
	"github.com/sirupsen/logrus"
)

func DB(ctx context.Context) xorm.Interface {
	v := ctx.Value(echomiddleware.ContextDBName)
	if v == nil {
		panic("DB is not exist")
	}
	if db, ok := v.(xorm.Interface); ok {
		return db
	}
	if db, ok := v.(*xorm.Engine); ok {
		return db.NewSession()
	}
	panic("DB is not exist")
}

func Logger(ctx context.Context) *logrus.Entry {
	v := ctx.Value(behaviorlog.LogContextName)
	if v == nil {
		return logrus.WithFields(logrus.Fields{})
	}
	if logger, ok := v.(*logrus.Entry); ok {
		return logger
	}
	return logrus.WithFields(logrus.Fields{})
}
