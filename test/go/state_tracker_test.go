package test

import (
	"testing"

	"github.com/Shopify/ghostferry"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/stretchr/testify/suite"
)

type StateTrackerTestSuite struct {
	suite.Suite
}

func (s *StateTrackerTestSuite) TestMinBinlogPosition() {
	serializedState := &ghostferry.SerializableState{
		LastWrittenBinlogPosition: mysql.Position{
			Name: "mysql-bin.00003",
			Pos:  4,
		},

		LastStoredSourceBinlogPositionForInlineVerifier: mysql.Position{
			Name: "mysql-bin.00003",
			Pos:  10,
		},
		LastStoredTargetBinlogPositionForInlineVerifier: mysql.Position{
			Name: "mysql-bin.00003",
			Pos:  10,
		},
	}
	s.Require().Equal(serializedState.MinSourceBinlogPosition(), mysql.Position{"mysql-bin.00003", 4})
	s.Require().Equal(serializedState.LastStoredTargetBinlogPositionForInlineVerifier, mysql.Position{"mysql-bin.00003", 10})

	serializedState = &ghostferry.SerializableState{
		LastWrittenBinlogPosition: mysql.Position{
			Name: "mysql-bin.00003",
			Pos:  4,
		},

		LastStoredSourceBinlogPositionForInlineVerifier: mysql.Position{
			Name: "mysql-bin.00002",
			Pos:  10,
		},

		LastStoredTargetBinlogPositionForInlineVerifier: mysql.Position{
			Name: "mysql-bin.00002",
			Pos:  10,
		},
	}
	s.Require().Equal(serializedState.MinSourceBinlogPosition(), mysql.Position{"mysql-bin.00002", 10})
	s.Require().Equal(serializedState.LastStoredTargetBinlogPositionForInlineVerifier, mysql.Position{"mysql-bin.00002", 10})

	serializedState = &ghostferry.SerializableState{
		LastWrittenBinlogPosition: mysql.Position{
			Name: "",
			Pos:  0,
		},

		LastStoredSourceBinlogPositionForInlineVerifier: mysql.Position{
			Name: "mysql-bin.00002",
			Pos:  10,
		},

		LastStoredTargetBinlogPositionForInlineVerifier: mysql.Position{
			Name: "mysql-bin.00002",
			Pos:  10,
		},
	}
	s.Require().Equal(serializedState.MinSourceBinlogPosition(), mysql.Position{"mysql-bin.00002", 10})
	s.Require().Equal(serializedState.LastStoredTargetBinlogPositionForInlineVerifier, mysql.Position{"mysql-bin.00002", 10})

	serializedState = &ghostferry.SerializableState{
		LastStoredSourceBinlogPositionForInlineVerifier: mysql.Position{
			Name: "",
			Pos:  0,
		},

		LastStoredTargetBinlogPositionForInlineVerifier: mysql.Position{
			Name: "mysql-bin.00002",
			Pos:  10,
		},

		LastWrittenBinlogPosition: mysql.Position{
			Name: "mysql-bin.00002",
			Pos:  10,
		},
	}
	s.Require().Equal(serializedState.MinSourceBinlogPosition(), mysql.Position{"mysql-bin.00002", 10})
	s.Require().Equal(serializedState.LastStoredTargetBinlogPositionForInlineVerifier, mysql.Position{"mysql-bin.00002", 10})
}

func TestStateTrackerTestSuite(t *testing.T) {
	suite.Run(t, new(StateTrackerTestSuite))
}
