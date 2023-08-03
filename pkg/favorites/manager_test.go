//nolint:paralleltest,funlen
package favorites_test

import (
	"context"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"

	favorites2 "github.com/gerladeno/favorites-mechanics/pkg/favorites"
)

type TestManagerSuite struct {
	suite.Suite
	manager *favorites2.Manager
}

func TestManager(t *testing.T) {
	suite.Run(t, new(TestManagerSuite))
}

func (s *TestManagerSuite) SetupSuite() {
	var err error

	s.manager, err = favorites2.NewManager(context.Background(), logrus.New(),
		favorites2.NewOptions(
			true,
			"rubbish",
			time.Minute,
			40,
		))
	s.Require().NoError(err)
}

func (s *TestManagerSuite) TearDownSuite() {
	time.Sleep(100 * time.Millisecond)
}

func (s *TestManagerSuite) TestManager() {
	var list []favorites2.Entry

	name := "first command"
	exec := "sudo do nothing"

	s.Run("add command", func() {
		s.manager.AddCommand(name, exec, 0, 0)
		list = s.manager.ListDirectory(0)
		s.Require().Len(list, 1)
		s.Require().Equal(name, list[0].Name)
		s.Require().Equal(exec, list[0].Exec)
		s.Require().Equal(name, s.manager.DisplayEntry(&list[0]))
	})

	s.Run("add empty command without name", func() {
		s.manager.AddCommand("", "", 0, 0)
		s.Require().Len(s.manager.ListDirectory(0), 1)
	})

	s.Run("add command without name", func() {
		s.manager.AddCommand("", exec, 0, list[0].ID)
		list = s.manager.ListDirectory(0)
		s.Require().Len(list, 2)
		s.Require().Equal("", list[0].Name)
		s.Require().Equal(exec, s.manager.DisplayEntry(&list[0]))
	})

	s.Run("move command within a dir", func() {
		s.manager.MoveEntry(list[0].ID, 0, 0)
		s.Require().Len(list, 2)
		list = s.manager.ListDirectory(0)
		s.Require().Equal(name, list[0].Name)
		s.Require().Equal("", list[1].Name)
	})

	s.Run("change exec", func() {
		exec = "sudo do something veeeeery long bla-bla-bla"
		s.manager.ModifyExec(list[1].ID, exec)
		s.Require().Equal(exec, s.manager.ListDirectory(0)[1].Exec)
		s.Require().Equal("sudo do something veeeeery long bla-b...", s.manager.DisplayEntry(&s.manager.ListDirectory(0)[1]))
	})

	s.Run("rename command", func() {
		name = "second command"
		s.manager.RenameEntry(list[1].ID, name)
		s.Require().Equal(name, s.manager.ListDirectory(0)[1].Name)
	})

	s.Run("add dir", func() {
		s.manager.AddDir("first dir", 0, list[0].ID)
		list = s.manager.ListDirectory(0)
		s.Require().Len(list, 3)
		s.Require().Equal("first dir", list[0].Name)
		s.Require().True(list[0].IsDir)
		s.Require().Equal(exec, list[2].Exec)
	})

	s.Run("add dir without a name", func() {
		s.manager.AddDir("", 0, list[0].ID)
		s.Require().Len(s.manager.ListDirectory(0), 3)
	})

	var list2 []favorites2.Entry

	s.Run("add dir to dir", func() {
		s.manager.AddDir("second dir", list[0].ID, 0)
		list = s.manager.ListDirectory(0)
		s.Require().Len(list, 3)
		list2 = s.manager.ListDirectory(list[0].ID)
		s.Require().Len(list2, 1)
		s.Require().Equal("second dir", list2[0].Name)
		s.Require().True(list2[0].IsDir)
	})

	s.Run("move command to another dir", func() {
		s.manager.MoveEntry(list[2].ID, list[0].ID, 0)
		list = s.manager.ListDirectory(0)
		s.Require().Len(list, 2)
		list2 = s.manager.ListDirectory(list[0].ID)
		s.Require().Len(list2, 2)
		s.Require().Equal(exec, list2[1].Exec)
	})

	s.Run("move another command to the dir", func() {
		s.manager.MoveEntry(list[1].ID, list[0].ID, list2[0].ID)
		list = s.manager.ListDirectory(0)
		s.Require().Len(list, 1)
		list2 = s.manager.ListDirectory(list[0].ID)
		s.Require().Len(list2, 3)
		s.Require().Equal(exec, list2[2].Exec)
		s.Require().True(list2[1].IsDir)
		s.Require().Equal("first command", list2[0].Name)
	})

	s.Run("sync out sync in", func() {
		s.manager.SyncOut()
		s.manager.SyncIn()
		list = s.manager.ListDirectory(0)
		s.Require().Len(list, 1)
		list2 = s.manager.ListDirectory(list[0].ID)
		s.Require().Len(list2, 3)
		s.Require().Equal(exec, list2[2].Exec)
		s.Require().True(list2[1].IsDir)
		s.Require().Equal("first command", list2[0].Name)
	})

	s.Run("delete command", func() {
		s.manager.DeleteCommand(list2[2].ID)
		list2 = s.manager.ListDirectory(list[0].ID)
		s.Require().Len(list2, 2)
	})

	s.Run("add tmp command", func() {
		s.manager.AddCommand("tmp cmd", "", list[0].ID, 0)
		list2 = s.manager.ListDirectory(list[0].ID)
		s.Require().Len(list2, 3)
	})

	s.Run("delete first command", func() {
		s.manager.DeleteCommand(list2[0].ID)
		list2 = s.manager.ListDirectory(list[0].ID)
		s.Require().Len(list2, 2)
	})

	s.Run("delete parent dir", func() {
		s.manager.DeleteDir(list[0].ID)
		list = s.manager.ListDirectory(0)
		s.Require().Len(list, 0)
		s.Require().Len(s.manager.EntryIDs, 0)
	})
}
