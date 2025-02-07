package initializers

import (
	"github.com/even44/JobsearchAPI/pkg/stores"
)

var Store *stores.MariaDBStore

func CreateDbStores() {
	logger.Println("Creating store")
	Store = stores.NewMariaDBStore(Db)
}
