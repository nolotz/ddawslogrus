package ddawslogrus

import (
	"github.com/sirupsen/logrus"
)

func NewFormatter() *logrus.JSONFormatter {
	return &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.999Z07:00",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime: "timestamp",
		},
	}
}
