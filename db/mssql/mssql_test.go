package mssql

import (
	"github.com/kinwyb/go/db"
	"github.com/kinwyb/go/err1"
	"testing"
)

func Test_mstx(t *testing.T) {
	sql, err := Connect("", "", "", "")
	if err != nil {
		t.Fatal(err)
	}
	err = sql.Transaction(func(tx db.TxSQL) err1.Error {
		err := one(tx)
		if err != nil {
			return err1.NewError(-1, err.Error())
		}
		err = three(tx)
		if err != nil {
			return err1.NewError(-1, err.Error())
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func one(tx db.Query) error {
	return tx.Transaction(func(tx db.TxSQL) err1.Error {
		err := two(tx)
		if err != nil {
			return err1.NewError(-1, err.Error())
		}
		return nil
	})
}

func two(tx db.Query) error {
	return nil
}

func three(tx db.Query) error {
	return nil
}
