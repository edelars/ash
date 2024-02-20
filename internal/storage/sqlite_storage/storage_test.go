package sqlite_storage

import (
	"os"
	"testing"

	"ash/internal/commands"
	"ash/internal/configuration"
	"ash/internal/dto"
	"ash/internal/internal_context"
	"ash/internal/storage"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func Test_sqliteStorage_InitStorage(t *testing.T) {
	const filename = "test.sql"
	h := NewSqliteStorage(configuration.StorageSqliteOpts{FileName: filename})
	err := h.initStorage()
	assert.NoError(t, err)
	assert.NoError(t, h.db.Close())

	h = NewSqliteStorage(configuration.StorageSqliteOpts{FileName: filename})
	err = h.initStorage()
	assert.NoError(t, err)
	assert.NoError(t, h.db.Close())

	assert.NoError(t, os.Remove(filename))
}

func Test_convertDataToExecString(t *testing.T) {
	type args struct {
		data storage.DataIface
	}
	tests := []struct {
		name     string
		args     args
		wantDir  string
		wantExec string
	}{
		{
			name: "1",
			args: args{
				data: internal_context.InternalContext{}.WithExecutionList([]dto.CommandIface{commands.NewCommand("1", nil, false).WithArgs([]string{"-l", "-b"}), commands.NewCommand("ls", nil, false).WithArgs([]string{"-l4", "-n"})}),
			},
			wantDir:  internal_context.InternalContext{}.GetCurrentDir(),
			wantExec: "1 -l -b|ls -l4 -n",
		},
		{
			name: "2",
			args: args{
				data: internal_context.InternalContext{}.WithExecutionList([]dto.CommandIface{commands.NewCommand("1", nil, false).WithArgs([]string{"-l", "-b"})}),
			},
			wantDir:  internal_context.InternalContext{}.GetCurrentDir(),
			wantExec: "1 -l -b",
		},
		{
			name: "3",
			args: args{
				data: internal_context.InternalContext{},
			},
			wantDir:  internal_context.InternalContext{}.GetCurrentDir(),
			wantExec: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDir, gotExec := convertDataToExecString(tt.args.data)
			if gotDir != tt.wantDir {
				t.Errorf("convertDataToExecString() gotDir = %v, want %v", gotDir, tt.wantDir)
			}
			if gotExec != tt.wantExec {
				t.Errorf("convertDataToExecString() gotExec = %v, want %v", gotExec, tt.wantExec)
			}
		})
	}
}

func Test_sqliteStorage_putData(t *testing.T) {
	const filename = "test.sql"
	h := NewSqliteStorage(configuration.StorageSqliteOpts{FileName: filename})
	err := h.initStorage()
	assert.NoError(t, err)
	assert.NoError(t, h.putData(internal_context.InternalContext{}.WithExecutionList([]dto.CommandIface{commands.NewCommand("ls", nil, false)}), 1))

	rows, err := h.db.Query("select lastUsedTime,usedCounter  from history")

	assert.NoError(t, err)

	var data []row
	for rows.Next() {
		i := row{}
		err = rows.Scan(&i.LastUsedTime, &i.UsedCounter)
		assert.NoError(t, err)
		data = append(data, i)
	}
	assert.Equal(t, 1, len(data))
	assert.Equal(t, 1, data[0].UsedCounter)

	prevInt := data[0].LastUsedTime

	rows.Close()

	assert.NoError(t, h.putData(internal_context.InternalContext{}.WithExecutionList([]dto.CommandIface{commands.NewCommand("ls", nil, false)}), 22))
	rows, err = h.db.Query("select lastUsedTime,usedCounter  from history")

	assert.NoError(t, err)
	defer rows.Close()

	data = nil
	for rows.Next() {
		i := row{}
		err = rows.Scan(&i.LastUsedTime, &i.UsedCounter)
		assert.NoError(t, err)
		data = append(data, i)
	}
	assert.Equal(t, 1, len(data))
	assert.Equal(t, 2, data[0].UsedCounter)

	rows.Close()

	assert.NoError(t, h.putData(internal_context.InternalContext{}.WithExecutionList([]dto.CommandIface{commands.NewCommand("lss", nil, false)}), 22))
	rows, err = h.db.Query("select lastUsedTime from history")

	assert.NoError(t, err)
	defer rows.Close()

	data = nil
	for rows.Next() {
		i := row{}
		err = rows.Scan(&i.LastUsedTime)
		assert.NoError(t, err)
		data = append(data, i)
	}
	assert.Equal(t, 2, len(data))
	rows.Close()

	assert.Equal(t, true, prevInt < data[0].LastUsedTime)
	assert.NoError(t, h.db.Close())
	assert.NoError(t, os.Remove(filename))
}

type row struct {
	LastUsedTime int
	UsedCounter  int
	Dir          string
	Count        int
}

func Test_sqliteStorage_SaveData(t *testing.T) {
	const filename = "test.sql"
	h := NewSqliteStorage(configuration.StorageSqliteOpts{FileName: filename, CleanupInterval: 999})
	err := h.initStorage()
	assert.NoError(t, err)
	go h.Run()
	defer h.Stop()

	h.SaveData(internal_context.InternalContext{}.WithExecutionList([]dto.CommandIface{commands.NewCommand("ls", nil, false)}))

	rows, err := h.db.Query("select lastUsedTime from history")

	assert.NoError(t, err)

	var data []row
	for rows.Next() {
		i := row{}
		err = rows.Scan(&i.LastUsedTime)
		assert.NoError(t, err)
		data = append(data, i)
	}
	assert.Equal(t, 1, len(data))

	rows.Close()
	assert.NoError(t, os.Remove(filename))
}

func Test_sqliteStorage_cleanupOldDirData(t *testing.T) {
	const filename = "test.sql"
	h := NewSqliteStorage(configuration.StorageSqliteOpts{FileName: filename, MaxHistoryPerDir: 3000})
	err := h.initStorage()
	assert.NoError(t, err)
	assert.NoError(t, h.putData(&icontextImpl{[]dto.CommandIface{commands.NewCommand("ls", nil, false)}, "dir1"}, 1))
	assert.NoError(t, h.putData(&icontextImpl{[]dto.CommandIface{commands.NewCommand("ls1", nil, false)}, "dir1"}, 2))
	assert.NoError(t, h.putData(&icontextImpl{[]dto.CommandIface{commands.NewCommand("ls2", nil, false)}, "dir1"}, 3))
	assert.NoError(t, h.putData(&icontextImpl{[]dto.CommandIface{commands.NewCommand("ls3", nil, false)}, "dir1"}, 4))

	assert.NoError(t, h.putData(&icontextImpl{[]dto.CommandIface{commands.NewCommand("ls", nil, false)}, "dir2"}, 1))
	assert.NoError(t, h.putData(&icontextImpl{[]dto.CommandIface{commands.NewCommand("ls1", nil, false)}, "dir2"}, 2))
	assert.NoError(t, h.putData(&icontextImpl{[]dto.CommandIface{commands.NewCommand("ls2", nil, false)}, "dir2"}, 3))
	assert.NoError(t, h.putData(&icontextImpl{[]dto.CommandIface{commands.NewCommand("ls3", nil, false)}, "dir2"}, 4))

	assert.NoError(t, h.putData(&icontextImpl{[]dto.CommandIface{commands.NewCommand("ls", nil, false)}, "dir3"}, 1))
	assert.NoError(t, h.putData(&icontextImpl{[]dto.CommandIface{commands.NewCommand("ls1", nil, false)}, "dir3"}, 5))

	assert.NoError(t, h.cleanupOldDirData())

	rows, err := h.db.Query("select count() as count from history")

	assert.NoError(t, err)

	i := row{}
	for rows.Next() {
		err = rows.Scan(&i.Count)
		assert.NoError(t, err)
		break
	}

	assert.Equal(t, 10, i.Count)
	rows.Close()

	// 2 test
	h.maxHistoryPerDir = 3
	assert.NoError(t, h.cleanupOldDirData())

	rows, err = h.db.Query("select min(lastUsedTime) as lastUsedTime,dir from history where dir in (select dir from history group by dir HAVING count(dir) = 3) group by dir order by dir")

	assert.NoError(t, err)

	var data []row
	for rows.Next() {
		i := row{}
		err = rows.Scan(&i.LastUsedTime, &i.Dir)
		assert.NoError(t, err)
		data = append(data, i)
	}
	assert.Equal(t, 2, len(data))
	assert.Equal(t, "dir1", data[0].Dir)
	assert.Equal(t, 2, data[0].LastUsedTime)

	assert.Equal(t, "dir2", data[1].Dir)
	assert.Equal(t, 2, data[1].LastUsedTime)

	rows.Close()

	rows, err = h.db.Query(`select lastUsedTime,dir from history where dir="dir3"`)

	assert.NoError(t, err)

	data = nil
	for rows.Next() {
		i := row{}
		err = rows.Scan(&i.LastUsedTime, &i.Dir)
		assert.NoError(t, err)
		data = append(data, i)
	}
	assert.Equal(t, 2, len(data))
	assert.Equal(t, "dir3", data[0].Dir)

	rows.Close()
	assert.NoError(t, os.Remove(filename))
}

type icontextImpl struct {
	ExecList   []dto.CommandIface
	CurrentDir string
}

func (icontextimpl *icontextImpl) GetExecutionList() []dto.CommandIface {
	return icontextimpl.ExecList
}

func (icontextimpl *icontextImpl) GetCurrentDir() string {
	return icontextimpl.CurrentDir
}

func Test_sqliteStorage_cleanupOldAllData(t *testing.T) {
	const filename = "test.sql"
	h := NewSqliteStorage(configuration.StorageSqliteOpts{FileName: filename, MaxHistoryTotal: 1000})
	err := h.initStorage()
	assert.NoError(t, err)
	assert.NoError(t, h.putData(&icontextImpl{[]dto.CommandIface{commands.NewCommand("ls", nil, false)}, "dir1"}, 1))
	assert.NoError(t, h.putData(&icontextImpl{[]dto.CommandIface{commands.NewCommand("ls1", nil, false)}, "dir1"}, 2))
	assert.NoError(t, h.putData(&icontextImpl{[]dto.CommandIface{commands.NewCommand("ls2", nil, false)}, "dir1"}, 3))
	assert.NoError(t, h.putData(&icontextImpl{[]dto.CommandIface{commands.NewCommand("ls3", nil, false)}, "dir1"}, 4))

	assert.NoError(t, h.putData(&icontextImpl{[]dto.CommandIface{commands.NewCommand("ls", nil, false)}, "dir2"}, 1))
	assert.NoError(t, h.putData(&icontextImpl{[]dto.CommandIface{commands.NewCommand("ls1", nil, false)}, "dir2"}, 2))
	assert.NoError(t, h.putData(&icontextImpl{[]dto.CommandIface{commands.NewCommand("ls2", nil, false)}, "dir2"}, 3))
	assert.NoError(t, h.putData(&icontextImpl{[]dto.CommandIface{commands.NewCommand("ls3", nil, false)}, "dir2"}, 4))

	assert.NoError(t, h.putData(&icontextImpl{[]dto.CommandIface{commands.NewCommand("ls", nil, false)}, "dir3"}, 1))
	assert.NoError(t, h.putData(&icontextImpl{[]dto.CommandIface{commands.NewCommand("ls1", nil, false)}, "dir3"}, 5))

	assert.NoError(t, h.cleanupOldAllData())
	rows, err := h.db.Query("select count() as count from history")

	assert.NoError(t, err)

	i := row{}
	for rows.Next() {
		err = rows.Scan(&i.Count)
		assert.NoError(t, err)
		break
	}

	assert.Equal(t, 10, i.Count)
	rows.Close()
	///
	h.maxHistoryTotal = 7
	assert.NoError(t, h.cleanupOldAllData())
	rows, err = h.db.Query("select min(lastUsedTime) as lastUsedTime,dir from history group by dir order by dir")

	assert.NoError(t, err)

	var data []row
	for rows.Next() {
		i := row{}
		err = rows.Scan(&i.LastUsedTime, &i.Dir)
		assert.NoError(t, err)
		data = append(data, i)
	}

	assert.Equal(t, 3, len(data))
	assert.Equal(t, "dir1", data[0].Dir)
	assert.Equal(t, 2, data[0].LastUsedTime)

	assert.Equal(t, "dir2", data[1].Dir)
	assert.Equal(t, 2, data[1].LastUsedTime)

	assert.Equal(t, "dir3", data[2].Dir)
	assert.Equal(t, 5, data[2].LastUsedTime)

	rows.Close()

	h.maxHistoryTotal = 5
	assert.NoError(t, h.cleanupOldAllData())
	rows, err = h.db.Query("select min(lastUsedTime) as lastUsedTime,dir from history group by dir order by dir")

	assert.NoError(t, err)

	data = nil
	for rows.Next() {
		i := row{}
		err = rows.Scan(&i.LastUsedTime, &i.Dir)
		assert.NoError(t, err)
		data = append(data, i)
	}
	assert.Equal(t, 3, len(data))
	assert.Equal(t, "dir1", data[0].Dir)
	assert.Equal(t, 3, data[0].LastUsedTime)

	assert.Equal(t, "dir2", data[1].Dir)
	assert.Equal(t, 3, data[1].LastUsedTime)

	assert.Equal(t, "dir3", data[2].Dir)
	assert.Equal(t, 5, data[2].LastUsedTime)

	rows.Close()
	assert.NoError(t, os.Remove(filename))
}

const fillBase string = `INSERT INTO history VALUES(unixepoch(),1,'dir1','cmd1') ON CONFLICT(dir,execWithArgs) DO UPDATE SET usedCounter = usedCounter + 1, lastUsedTime = unixepoch();
INSERT INTO history VALUES(unixepoch()+1,1,'dir1','cmd2') ON CONFLICT(dir,execWithArgs) DO UPDATE SET usedCounter = usedCounter + 1, lastUsedTime = unixepoch()+1;
INSERT INTO history VALUES(unixepoch()+2,1,'dir1','cmd3') ON CONFLICT(dir,execWithArgs) DO UPDATE SET usedCounter = usedCounter + 1, lastUsedTime = unixepoch()+2;
INSERT INTO history VALUES(unixepoch()+3,1,'dir1','cmd4') ON CONFLICT(dir,execWithArgs) DO UPDATE SET usedCounter = usedCounter + 1, lastUsedTime = unixepoch()+3;
INSERT INTO history VALUES(unixepoch(),1,'dir2','cmd1') ON CONFLICT(dir,execWithArgs) DO UPDATE SET usedCounter = usedCounter + 1, lastUsedTime = unixepoch();
INSERT INTO history VALUES(unixepoch()+1,1,'dir2','cmd2') ON CONFLICT(dir,execWithArgs) DO UPDATE SET usedCounter = usedCounter + 1, lastUsedTime = unixepoch()+1;
INSERT INTO history VALUES(unixepoch()+2,1,'dir2','cmd3') ON CONFLICT(dir,execWithArgs) DO UPDATE SET usedCounter = usedCounter + 1, lastUsedTime = unixepoch()+2;
INSERT INTO history VALUES(unixepoch()+3,1,'dir2','cmd4') ON CONFLICT(dir,execWithArgs) DO UPDATE SET usedCounter = usedCounter + 1, lastUsedTime = unixepoch()+3;
INSERT INTO history VALUES(unixepoch(),1,'dir3','cmd1') ON CONFLICT(dir,execWithArgs) DO UPDATE SET usedCounter = usedCounter + 1, lastUsedTime = unixepoch();
INSERT INTO history VALUES(unixepoch()+1,1,'dir3','cmd2') ON CONFLICT(dir,execWithArgs) DO UPDATE SET usedCounter = usedCounter + 1, lastUsedTime = unixepoch()+1;

INSERT INTO history VALUES(unixepoch(),1,'dir1','cmd1') ON CONFLICT(dir,execWithArgs) DO UPDATE SET usedCounter = usedCounter + 1, lastUsedTime = unixepoch()-1;
INSERT INTO history VALUES(unixepoch()+1,1,'dir1','cmd2') ON CONFLICT(dir,execWithArgs) DO UPDATE SET usedCounter = usedCounter + 1, lastUsedTime = unixepoch()-2;
INSERT INTO history VALUES(unixepoch()+2,1,'dir1','cmd3') ON CONFLICT(dir,execWithArgs) DO UPDATE SET usedCounter = usedCounter + 1, lastUsedTime = unixepoch()-3;
INSERT INTO history VALUES(unixepoch()+3,1,'dir1','cmd4') ON CONFLICT(dir,execWithArgs) DO UPDATE SET usedCounter = usedCounter + 1, lastUsedTime = unixepoch()-4;
INSERT INTO history VALUES(unixepoch(),1,'dir2','cmd1') ON CONFLICT(dir,execWithArgs) DO UPDATE SET usedCounter = usedCounter + 1, lastUsedTime = unixepoch()-5;
INSERT INTO history VALUES(unixepoch()+1,1,'dir2','cmd2') ON CONFLICT(dir,execWithArgs) DO UPDATE SET usedCounter = usedCounter + 1, lastUsedTime = unixepoch()-6;
INSERT INTO history VALUES(unixepoch()+2,1,'dir2','cmd3') ON CONFLICT(dir,execWithArgs) DO UPDATE SET usedCounter = usedCounter + 1, lastUsedTime = unixepoch()-7;
INSERT INTO history VALUES(unixepoch()+3,1,'dir2','cmd4') ON CONFLICT(dir,execWithArgs) DO UPDATE SET usedCounter = usedCounter + 1, lastUsedTime = unixepoch()-8;
INSERT INTO history VALUES(unixepoch(),1,'dir3','cmd1') ON CONFLICT(dir,execWithArgs) DO UPDATE SET usedCounter = usedCounter + 1, lastUsedTime = unixepoch()-9;
INSERT INTO history VALUES(unixepoch()+1,1,'dir3','cmd2') ON CONFLICT(dir,execWithArgs) DO UPDATE SET usedCounter = usedCounter + 1, lastUsedTime = unixepoch()-10;

INSERT INTO history VALUES(unixepoch(),1,'dir1','cmd155') ON CONFLICT(dir,execWithArgs) DO UPDATE SET usedCounter = usedCounter + 1, lastUsedTime = unixepoch()-1;
INSERT INTO history VALUES(unixepoch()+1,1,'dir1','cmd2345') ON CONFLICT(dir,execWithArgs) DO UPDATE SET usedCounter = usedCounter + 1, lastUsedTime = unixepoch()-2;
INSERT INTO history VALUES(unixepoch()+2,1,'dir1','cmd3343') ON CONFLICT(dir,execWithArgs) DO UPDATE SET usedCounter = usedCounter + 1, lastUsedTime = unixepoch()-3;
INSERT INTO history VALUES(unixepoch()+3,1,'dir1','cmd488') ON CONFLICT(dir,execWithArgs) DO UPDATE SET usedCounter = usedCounter + 1, lastUsedTime = unixepoch()-4;
INSERT INTO history VALUES(unixepoch(),1,'dir2','cmd14322') ON CONFLICT(dir,execWithArgs) DO UPDATE SET usedCounter = usedCounter + 1, lastUsedTime = unixepoch()-5;
INSERT INTO history VALUES(unixepoch()+1,1,'dir2','cmd22342') ON CONFLICT(dir,execWithArgs) DO UPDATE SET usedCounter = usedCounter + 1, lastUsedTime = unixepoch()-6;
INSERT INTO history VALUES(unixepoch()+2,1,'dir2','cmd3223') ON CONFLICT(dir,execWithArgs) DO UPDATE SET usedCounter = usedCounter + 1, lastUsedTime = unixepoch()-7;
INSERT INTO history VALUES(unixepoch()+3,1,'dir2','cmd433') ON CONFLICT(dir,execWithArgs) DO UPDATE SET usedCounter = usedCounter + 1, lastUsedTime = unixepoch()-8;
INSERT INTO history VALUES(unixepoch(),1,'dir3','cmd132111') ON CONFLICT(dir,execWithArgs) DO UPDATE SET usedCounter = usedCounter + 1, lastUsedTime = unixepoch()-9;
INSERT INTO history VALUES(unixepoch()+1,1,'dir3','cmd2876') ON CONFLICT(dir,execWithArgs) DO UPDATE SET usedCounter = usedCounter + 1, lastUsedTime = unixepoch()-10;
`

func Test_sqliteStorage_getTopHistory(t *testing.T) {
	const filename = "test.sql"
	h := NewSqliteStorage(configuration.StorageSqliteOpts{FileName: filename, MaxHistoryTotal: 7})
	err := h.initStorage()
	assert.NoError(t, err)
	_, err = h.db.Exec(fillBase)
	assert.NoError(t, err)

	res := h.GetTopHistoryByPattern("cmd", 5)
	assert.Equal(t, 5, len(res))
	for _, v := range res {
		assert.Equal(t, 2, v.GetUsedCount())
		assert.Contains(t, v.GetCommand(), "cmd")
	}
	assert.NoError(t, os.Remove(filename))
}

func Test_sqliteStorage_GetTopHistoryByDirs(t *testing.T) {
	const filename = "test.sql"
	h := NewSqliteStorage(configuration.StorageSqliteOpts{FileName: filename, MaxHistoryTotal: 7})
	err := h.initStorage()
	assert.NoError(t, err)
	_, err = h.db.Exec(fillBase)
	assert.NoError(t, err)

	res := h.GetTopHistoryByDirs("dir1", 10)
	assert.Equal(t, 10, len(res))
	for _, v := range res {
		assert.Contains(t, v.GetCommand(), "cmd")
	}
	assert.NoError(t, os.Remove(filename))
}
