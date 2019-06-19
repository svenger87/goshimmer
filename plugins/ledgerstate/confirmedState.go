package ledgerstate

import (
	"github.com/iotaledger/goshimmer/packages/database"
	"github.com/iotaledger/goshimmer/packages/node"
)

// region database /////////////////////////////////////////////////////////////////////////////////////////////////////

var confirmedStateDatabase database.Database

func configureConfirmedStateDatabase(plugin *node.Plugin) {
	if db, err := database.Get("confirmedState"); err != nil {
		panic(err)
	} else {
		confirmedStateDatabase = db
	}
}

// func storeTransactionInDatabase(transaction *value_transaction.ValueTransaction) errors.IdentifiableError {
// 	if transaction.GetModified() {
// 		if err := transactionDatabase.Set(transaction.GetHash().CastToBytes(), transaction.MetaTransaction.GetBytes()); err != nil {
// 			return ErrDatabaseError.Derive(err, "failed to store transaction")
// 		}

// 		transaction.SetModified(false)
// 	}

// 	return nil
// }

// func getTransactionFromDatabase(transactionHash ternary.Trytes) (*value_transaction.ValueTransaction, errors.IdentifiableError) {
// 	txData, err := transactionDatabase.Get(transactionHash.CastToBytes())
// 	if err != nil {
// 		if err == badger.ErrKeyNotFound {
// 			return nil, nil
// 		} else {
// 			return nil, ErrDatabaseError.Derive(err, "failed to retrieve transaction")
// 		}
// 	}

// 	return value_transaction.FromBytes(txData), nil
// }

// func databaseContainsTransaction(transactionHash ternary.Trytes) (bool, errors.IdentifiableError) {
// 	if contains, err := transactionDatabase.Contains(transactionHash.CastToBytes()); err != nil {
// 		return contains, ErrDatabaseError.Derive(err, "failed to check if the transaction exists")
// 	} else {
// 		return contains, nil
// 	}
// }

// endregion ///////////////////////////////////////////////////////////////////////////////////////////////////////////
