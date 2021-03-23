package reference_test

import (
	"context"
	"log"
	"os"
	"testing"

	"0chain.net/blobbercore/config"
	"0chain.net/blobbercore/datastore"
	"0chain.net/blobbercore/reference"
	cconfig "0chain.net/core/config"
	"0chain.net/core/logging"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	connCtx = context.Background()
	db      *gorm.DB
)

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

	logging.Logger = zap.New(nil)

	code := m.Run()
	// teardown
	datastore.GetStore().Close()
	os.Exit(code)
}

func migrate(db *gorm.DB) error {
	tables := []interface{}{
		&reference.Collaborator{},
		&reference.CommitMetaTxn{},
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

	refs := []*reference.Ref{
		&reference.Ref{
			AllocationID: "1",
			Path:         "/",
			PathLevel:    1,
			NumBlocks:    int64(10),
			Type:         reference.DIRECTORY,
			Children: []*reference.Ref{
				&reference.Ref{
					AllocationID: "1",
					Path:         "/1",
					PathLevel:    2,
					ParentPath:   "/",
					NumBlocks:    int64(10),
					Type:         reference.FILE,
				},
				&reference.Ref{
					AllocationID: "1",
					Path:         "/2",
					PathLevel:    2,
					ParentPath:   "/",
					NumBlocks:    int64(10),
					Type:         reference.DIRECTORY,
					Children: []*reference.Ref{
						&reference.Ref{
							AllocationID: "1",
							Path:         "/2/1",
							PathLevel:    3,
							ParentPath:   "/2",
							Type:         reference.FILE,
						},
					},
				},
			},
		},
		&reference.Ref{
			AllocationID: "2",
			Path:         "/",
			PathLevel:    1,
			NumBlocks:    int64(10),
			Type:         reference.DIRECTORY,
			Children: []*reference.Ref{
				&reference.Ref{
					AllocationID: "2",
					Path:         "/1",
					PathLevel:    2,
					ParentPath:   "/",
					NumBlocks:    int64(10),
					Type:         reference.DIRECTORY,
					Children: []*reference.Ref{
						&reference.Ref{
							AllocationID: "2",
							Path:         "/1/1",
							PathLevel:    3,
							ParentPath:   "/1/1",
							NumBlocks:    int64(10),
							Type:         reference.DIRECTORY,
							Children: []*reference.Ref{
								&reference.Ref{
									AllocationID: "2",
									Path:         "/1/1/1",
									PathLevel:    4,
									NumBlocks:    int64(10),
									Type:         reference.FILE,
								},
							},
						},
					},
				},
			},
		},
		&reference.Ref{
			AllocationID: "3",
			Path:         "/1",
			PathLevel:    1,
			NumBlocks:    int64(0),
			Type:         reference.FILE,
		},
		&reference.Ref{
			AllocationID: "4",
			Path:         "/",
			PathLevel:    1,
			NumBlocks:    int64(10),
			Type:         reference.FILE,
			Children: []*reference.Ref{
				&reference.Ref{
					AllocationID: "4",
					Path:         "/1/2",
					ParentPath:   "/",
					PathLevel:    2,
					NumBlocks:    int64(10),
					Type:         reference.FILE,
				},
			},
		},
		&reference.Ref{
			AllocationID: "5",
			Path:         "/1",
			ParentPath:   "/1",
			PathLevel:    3,
			NumBlocks:    int64(10),
			Type:         reference.DIRECTORY,
			Children: []*reference.Ref{
				&reference.Ref{
					AllocationID: "5",
					Path:         "/1/2",
					ParentPath:   "/1",
					PathLevel:    1,
					NumBlocks:    int64(10),
					Type:         reference.DIRECTORY,
				},
			},
		},
		&reference.Ref{
			AllocationID: "6",
			Path:         "/",
			PathLevel:    1,
			NumBlocks:    int64(10),
			Type:         reference.DIRECTORY,
			Children: []*reference.Ref{
				&reference.Ref{
					AllocationID: "6",
					Path:         "/1",
					ParentPath:   "/",
					PathLevel:    2,
					NumBlocks:    int64(10),
					Type:         reference.DIRECTORY,
					Children: []*reference.Ref{
						&reference.Ref{
							AllocationID: "6",
							Path:         "/1/2",
							ParentPath:   "/",
							PathLevel:    3,
							NumBlocks:    int64(10),
							Type:         reference.FILE,
						},
					},
				},
			},
		},
		&reference.Ref{
			AllocationID: "7",
			Path:         "/",
			PathLevel:    1,
			NumBlocks:    int64(10),
			Type:         reference.DIRECTORY,
			Children: []*reference.Ref{
				&reference.Ref{
					AllocationID: "7",
					Path:         "/1",
					ParentPath:   "/",
					PathLevel:    2,
					NumBlocks:    int64(10),
					Type:         reference.FILE,
				},
			},
		},
	}
	for _, ref := range refs {
		if err := saveRefRecurse(ref); err != nil {
			return err
		}
	}
	return nil
}

func saveRefRecurse(ref *reference.Ref) error {
	ref.LookupHash = reference.GetReferenceLookup(ref.AllocationID, ref.Path)
	err := db.Save(ref).Error
	if err != nil {
		return err
	}
	if ref.Children != nil {
		for _, child := range ref.Children {
			err = saveRefRecurse(child)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
