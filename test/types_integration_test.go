package test

import (
	"database/sql"
	"fmt"
	"math/rand"
	"testing"

	"github.com/Shopify/ghostferry/testhelpers"
	"github.com/sirupsen/logrus"
)

func addTypesToTable(db *sql.DB, dbName, tableName string) {
	query := "ALTER TABLE %s.%s " +
		"ADD tiny_col TINYINT," +
		"ADD float_col FLOAT," +
		"ADD double_col DOUBLE," +
		"ADD decimal_col DECIMAL(4, 2)," +
		"ADD year_col YEAR," +
		"ADD date_col DATE," +
		"ADD time_col TIME," +
		"ADD dt_col DATETIME," +
		"ADD ts_col TIMESTAMP," + // TODO broken on master
		"ADD varchar_col VARCHAR(128)," +
		"ADD enum_col ENUM('foo', 'bar')," +
		"ADD set_col SET('foo', 'bar', 'baz')," +
		"ADD utfmb4_col TEXT CHARSET utf8mb4," +
		"ADD utf32_col TEXT CHARSET utf32," +
		"ADD latin1_col TEXT CHARSET latin1 COLLATE latin1_swedish_ci," +
		"ADD blob_col BLOB," +
		"ADD uint64_col BIGINT UNSIGNED," +
		"ADD uint32_col INT UNSIGNED," +
		"ADD uint16_col SMALLINT UNSIGNED," +
		"ADD uint8_col TINYINT UNSIGNED"

	query = fmt.Sprintf(query, dbName, tableName)
	_, err := db.Exec(query)
	testhelpers.PanicIfError(err)
}

func setupMultiTypeTable(f *testhelpers.TestFerry) {
	testhelpers.SeedInitialData(f.SourceDB, "gftest", "table1", 0)
	testhelpers.SeedInitialData(f.TargetDB, "gftest", "table1", 0)

	addTypesToTable(f.SourceDB, "gftest", "table1")
	addTypesToTable(f.TargetDB, "gftest", "table1")

	tx, err := f.SourceDB.Begin()
	testhelpers.PanicIfError(err)

	for i := 0; i < 100; i++ {
		query := "INSERT INTO gftest.table1 " +
			"(id, data, tiny_col, float_col, double_col, decimal_col, year_col, date_col, time_col, dt_col, varchar_col, enum_col, set_col, utfmb4_col, utf32_col, latin1_col, blob_col, uint64_col, uint32_col, uint16_col, uint8_col)" +
			"VALUES (NULL, ?, ?, 3.14, 2.72, 42.42, NOW(), NOW(), NOW(), NOW(), ?, ?, 'foo,baz', ?, ?, ?, ?, 18446744073709551615, 3221225472, 49152, 192)"

		enumVal := "foo"
		if i%2 == 0 {
			enumVal = "bar"
		}
		randStr := testhelpers.RandData()
		randUStr := testhelpers.RandUTF8MB4Data()
		randLStr := testhelpers.RandLatin1Data()
		randBytes := testhelpers.RandByteData()

		_, err = tx.Exec(query, randStr, rand.Intn(2), randStr, enumVal, randUStr, randUStr, randLStr, randBytes)
		testhelpers.PanicIfError(err)
	}

	testhelpers.PanicIfError(tx.Commit())
}

func TestCopyDataWithManyTypes(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	testcase := &testhelpers.IntegrationTestCase{
		T:           t,
		SetupAction: setupMultiTypeTable,
		DataWriter: &testhelpers.MixedActionDataWriter{ // TODO: there's no guarantee that this data writer will update the data of the existing rows.
			ProbabilityOfInsert: 1.0 / 3.0,
			ProbabilityOfUpdate: 1.0 / 3.0,
			ProbabilityOfDelete: 1.0 / 3.0,
			NumberOfWriters:     3,
			Tables:              []string{"gftest.table1"},
		},
		Ferry: testhelpers.NewTestFerry(),
	}

	testcase.Run()
}