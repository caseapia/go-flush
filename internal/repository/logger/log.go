package loggerrepository

import (
	"time"

	loggermodule "github.com/caseapia/goproject-flush/internal/models/logger"
)

func (l *LoggerRepository) Log(entry *loggermodule.ActionLog) error {
	entry.CreatedAt = time.Now()

	if err := l.db.Create(&entry).Error; err != nil {
		return err
	}
	return l.db.Create(entry).Error
}
