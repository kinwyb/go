package mssql

import (
	"github.com/kinwyb/go/db"
	"testing"
)

func Test_mstx(t *testing.T) {
	sql, err := Connect("", "", "", "")
	if err != nil {
		t.Fatal(err)
	}
	err = sql.Transaction(func(tx db.TxSQL) error {
		err := one(tx)
		if err != nil {
			return err
		}
		err = three(tx)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func one(tx db.Query) error {
	return tx.Transaction(func(tx db.TxSQL) error {
		return two(tx)
	})
}

func two(tx db.Query) error {
	return nil
}

func three(tx db.Query) error {
	return nil
}
