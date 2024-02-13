package sqlite_storage

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"ash/internal/configuration"
	"ash/internal/storage"

	_ "github.com/mattn/go-sqlite3"
)

type sqliteStorage struct {
	maxHistoryPerDir int
	maxHistoryTotal  int
	cleanupInterval  int

	db        *sql.DB
	stopChan  chan struct{}
	inputData chan storage.DataIface
	dbFile    string
}

func NewSqliteStorage(opts configuration.StorageSqliteOpts) sqliteStorage {
	return sqliteStorage{
		stopChan:         make(chan struct{}),
		inputData:        make(chan storage.DataIface, opts.WriteBuffer),
		dbFile:           opts.FileName,
		maxHistoryPerDir: opts.MaxHistoryPerDir,
		maxHistoryTotal:  opts.MaxHistoryTotal,
		cleanupInterval:  opts.CleanupInterval,
	}
}

func (s *sqliteStorage) Run() error {
	if err := s.initStorage(); err != nil {
		return err
	}

	var cleanNeeded bool

	timer := time.NewTicker(time.Duration(s.cleanupInterval) * time.Second)

	defer func() {
		close(s.stopChan)
		close(s.inputData)
		timer.Stop()
		s.db.Close()
	}()
	for {
		select {
		case i := <-s.inputData:
			err := s.putData(i, 0)
			if err != nil {
				return err
			}
			if cleanNeeded {
				if err = s.cleanupOldDirData(); err != nil {
					return err
				}
				if err = s.cleanupOldAllData(); err != nil {
					return err
				}
				cleanNeeded = !cleanNeeded
				timer.Reset(time.Duration(s.cleanupInterval) * time.Second)
			}
		case <-timer.C:
			cleanNeeded = true
			timer.Stop()
		case <-s.stopChan:
			return nil
		}
	}
}

func (s *sqliteStorage) Stop() {
	s.stopChan <- struct{}{}
}

func (s *sqliteStorage) initStorage() error {
	if _, err := os.Stat(s.dbFile); err != nil {
		if _, err := os.Create(s.dbFile); err != nil {
			return err
		}
	}

	var err error
	if s.db, err = sql.Open("sqlite3", s.dbFile); err != nil {
		return err
	}

	if err = s.db.Ping(); err != nil {
		return err
	}

	if err = s.initTables(); err != nil {
		return err
	}

	return nil
}

func (s *sqliteStorage) SaveData(data storage.DataIface) {
	s.inputData <- data
}

// By default dont needed to setup unixtime. Leave 0
func (s *sqliteStorage) putData(data storage.DataIface, unixtime int64) error {
	dir, exec := convertDataToExecString(data)
	if exec == "" {
		return nil
	}

	if unixtime == 0 {
		unixtime = time.Now().Unix()
	}

	_, err := s.db.Exec("INSERT INTO history (lastUsedTime,usedCounter,dir,execWithArgs) VALUES(@time,1,@dir,@exec) ON CONFLICT(dir,execWithArgs) DO UPDATE SET usedCounter = usedCounter + 1, lastUsedTime = @time;",
		sql.Named("time", unixtime),
		sql.Named("dir", dir),
		sql.Named("exec", exec),
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *sqliteStorage) initTables() error {
	const historyTable string = `
	  CREATE TABLE IF NOT EXISTS history (
	  lastUsedTime INTEGER NOT NULL,
	  usedCounter INTEGER NOT NULL,
	  dir TEXT NOT NULL,
	  execWithArgs TEXT NOT NULL
	  );

	  CREATE INDEX IF NOT EXISTS dirIndex ON history (dir);
	  CREATE INDEX IF NOT EXISTS lastUsedTimeIndex ON history (lastUsedTime);
	  CREATE INDEX IF NOT EXISTS usedCounterIndex ON history (usedCounter);
	  CREATE UNIQUE INDEX IF NOT EXISTS uniqDirExeci on history (dir,execWithArgs);
	  `

	if _, err := s.db.Exec(historyTable); err != nil {
		return err
	}
	return nil
}

func convertDataToExecString(data storage.DataIface) (dir, exec string) {
	var str []string
	for _, v := range data.GetExecutionList() {
		str = append(str, fmt.Sprintf("%s %s", v.GetName(), strings.Join(v.GetArgs(), " ")))
	}
	return data.GetCurrentDir(), strings.Join(str, "|")
}

func (s *sqliteStorage) cleanupOldDirData() error {
	_, err := s.db.Exec(`delete from 
	  history 
	where 
	  (
		lastUsedTime, usedCounter, dir, execWithArgs
	  ) in (
		select 
		  lastUsedTime, 
		  usedCounter, 
		  dir, 
		  execWithArgs 
		from 
		  (
			select 
			  *, 
			  ROW_NUMBER() OVER (
				PARTITION BY dir 
				ORDER BY 
				  lastUsedTime desc
			  ) count 
			from 
			  history 
			where 
			  dir in (
				select 
				  dir 
				from 
				  history 
				GROUP BY 
				  dir 
				HAVING 
				  count(dir) > @c
			  ) 
			order by 
			  dir
		  ) 
		where 
		  count > @c
	  )`, sql.Named("c", s.maxHistoryPerDir))
	if err != nil {
		return err
	}
	return nil
}

func (s *sqliteStorage) cleanupOldAllData() error {
	_, err := s.db.Exec(`delete from 
	  history 
	where 
	  (
		lastUsedTime, usedCounter, dir, execWithArgs
	  ) in (
		select 
		  lastUsedTime, 
		  usedCounter, 
		  dir, 
		  execWithArgs 
		from 
		  history 
		order by 
		  lastUsedTime 
		limit 
		  (
			select 
			  count() 
			from 
			  history
		  )-@c
	  )
	`, sql.Named("c", s.maxHistoryTotal))
	if err != nil {
		return err
	}
	return nil
}

func (s *sqliteStorage) GetTopHistoryForCurrentDirAndAll(currentDir string, limit int) (res []storage.StorageResult) {
	rows, err := s.db.Query("select min(lastUsedTime) as lastUsedTime,dir from history group by dir order by dir")
	if err != nil {
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var item storageResult

		if err = rows.Scan(&item); err == nil {
			res = append(res, &item)
		}

	}
	return
}

func (s *sqliteStorage) GetHistoryMathPrefix(prefix string, limit int) []storage.StorageResult {
	panic("not implemented") // TODO: Implement
}
