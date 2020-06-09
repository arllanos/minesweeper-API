package logs

import "go.uber.org/zap"

var sugar *zap.SugaredLogger

func InitLogger() error {
	l, err := zap.NewDevelopment()

	if err != nil {
		return err
	}

	sugar = l.Sugar()

	return nil
}

func Log() *zap.SugaredLogger {
	_ = sugar.Sync()
	return sugar
}
