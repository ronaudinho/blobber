package allocation_test

import (
	"context"
	"log"
	"os"
	"testing"

	"0chain.net/blobbercore/allocation"
	"0chain.net/blobbercore/config"
	"0chain.net/blobbercore/datastore"
	"0chain.net/blobbercore/reference"
	"0chain.net/core/chain"
	cconfig "0chain.net/core/config"
	"0chain.net/core/logging"
	"0chain.net/core/node"

	"github.com/0chain/gosdk/core/zcncrypto"
	// "github.com/0chain/gosdk/zcncore"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	connCtx = context.Background()
	db      *gorm.DB
)

// main
func TestMain(m *testing.M) {
	config.Configuration = config.Config{
		Config: &cconfig.Config{
			SignatureScheme: "bls0chain",
			DeploymentMode:  config.DeploymentDevelopment,
		},
		DBHost:     "localhost",
		DBPort:     "5432",
		DBUserName: "postgres",
		DBName:     "postgres",
		DBPassword: "secret",
	}
	config.SetupDefaultConfig()

	sigSch := zcncrypto.NewSignatureScheme("bls0chain")
	wallet, err := sigSch.GenerateKeys()
	if err != nil {
		log.Fatal(err)
	}
	node.Self.SetKeys(wallet.Keys[0].PublicKey, wallet.Keys[0].PrivateKey)

	serverChain := &chain.Chain{
		ID:               "1",
		Version:          "0",
		OwnerID:          "1",
		ParentChainID:    "1",
		BlockWorker:      "http://127.0.0.1:9091",
		GenesisBlockHash: "1",
	}
	serverChain.InitializeCreationDate()
	chain.SetServerChain(serverChain)

	// setup
	store := datastore.GetStore()
	if err := store.Open(); err != nil {
		log.Fatal(err)
	}
	defer store.Close()
	db = store.GetDB()
	connCtx = context.WithValue(context.Background(), datastore.CONNECTION_CONTEXT_KEY, db)
	if err := migrate(db); err != nil {
		log.Fatal(err)
	}

	// err = zcncore.InitZCNSDK(serverChain.BlockWorker, config.Configuration.SignatureScheme)
	// if err != nil {
	//	log.Fatal(err)
	// }
	logging.Logger = zap.New(nil)

	code := m.Run()
	// teardown
	// datastore.GetStore().Close()
	os.Exit(code)
}

func migrate(db *gorm.DB) error {
	tables := []interface{}{
		&allocation.AllocationChangeCollector{},
		&allocation.AllocationChange{},
		&allocation.ReadPool{},
		&allocation.WritePool{},
		&allocation.Pending{},
		&reference.Ref{},
	}
	for _, t := range tables {
		if db.Migrator().HasTable(t) {
			if err := db.Migrator().DropTable(t); err != nil {
				return err
			}
		}
		if err := db.Migrator().CreateTable(t); err != nil {
			return err
		}
	}

	db.Create(&allocation.AllocationChangeCollector{
		ConnectionID: "1",
		AllocationID: "1",
		ClientID:     "1",
		Status:       allocation.CommittedConnection,
	})
	db.Create(&allocation.AllocationChangeCollector{
		ConnectionID: "2",
		AllocationID: "1",
		ClientID:     "1",
		Status:       allocation.DeletedConnection,
	})
	return nil
}
