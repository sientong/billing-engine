package helper

import "database/sql"

func CommitOrRollback(tx *sql.Tx) {
	err := recover()
	if err != nil {
		errorRollback := tx.Rollback()
		CheckErrorOrReturn(errorRollback)
	} else {
		errorCommit := tx.Commit()
		CheckErrorOrReturn(errorCommit)
	}
}
