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
}

func Test_sqliteStorage_SaveData(t *testing.T) {
	const filename = "test.sql"
	h := NewSqliteStorage(configuration.StorageSqliteOpts{FileName: filename})
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
	h := NewSqliteStorage(configuration.StorageSqliteOpts{FileName: filename, MaxHistoryPerDir: 3})
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
	rows, err := h.db.Query("select min(lastUsedTime) as lastUsedTime,dir from history where dir in (select dir from history group by dir HAVING count(dir) = 3) group by dir order by dir")

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
	h := NewSqliteStorage(configuration.StorageSqliteOpts{FileName: filename, MaxHistoryTotal: 7})
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
	rows, err := h.db.Query("select min(lastUsedTime) as lastUsedTime,dir from history group by dir order by dir")

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
