package main

import (
	"log"

	"github.com/typical-go/typical-go/pkg/typapp"
	"github.com/typical-go/typical-go/pkg/typast"
	"github.com/typical-go/typical-go/pkg/typgo"
	"github.com/typical-go/typical-go/pkg/typmock"
	"github.com/typical-go/typical-go/pkg/typrls"
	"github.com/typical-go/typical-rest-server/internal/generated/config"
	"github.com/typical-go/typical-rest-server/pkg/dbtool"
	"github.com/typical-go/typical-rest-server/pkg/dbtool/mysqltool"
	"github.com/typical-go/typical-rest-server/pkg/dbtool/pgtool"
	"github.com/typical-go/typical-rest-server/pkg/typcfg"
	"github.com/typical-go/typical-rest-server/pkg/typdocker"
	"github.com/typical-go/typical-rest-server/pkg/typrepo"
)

var descriptor = typgo.Descriptor{
	ProjectName:    "typical-rest-server",
	ProjectVersion: "0.9.9",

	Tasks: []typgo.Tasker{
		// annotate
		&typast.AnnotateMe{
			Sources: []string{"internal"},
			Annotators: []typast.Annotator{
				&typapp.CtorAnnotation{},
				&typrepo.EntityAnnotation{},
				&typcfg.EnvconfigAnnotation{DotEnv: ".env", UsageDoc: "USAGE.md"},
			},
		},
		// test
		&typgo.GoTest{
			Includes: []string{"internal/app/*/**", "pkg/**"},
			// Excludes: []string{"internal/app/*"},
		},
		// compile
		&typgo.GoBuild{},
		// run
		&typgo.RunBinary{
			Before: typgo.TaskNames{"annotate", "build"},
		},
		// mock
		&typmock.GenerateMock{},
		// docker
		&typdocker.DockerTool{},
		// pg
		&pgtool.PgTool{
			Name: "pg",
			ConfigFn: func() dbtool.Configurer {
				cfg, err := config.LoadPgDatabaseCfg()
				if err != nil {
					log.Fatal(err)
				}
				return cfg
			},
			DockerName:   "typical-rest-server_pg01_1",
			MigrationSrc: "file://databases/postgresdb/migration",
			SeedSrc:      "databases/postgresdb/seed",
		},
		// mysql
		&mysqltool.MySQLTool{
			Name: "mysql",
			ConfigFn: func() dbtool.Configurer {
				cfg, err := config.LoadMysqlDatabaseCfg()
				if err != nil {
					log.Fatal(err)
				}
				return cfg
			},
			DockerName:   "typical-rest-server_mysql01_1",
			MigrationSrc: "file://databases/mysqldb/migration",
			SeedSrc:      "databases/mysqldb/seed",
		},
		// reset
		&typgo.Task{
			Name:  "reset",
			Usage: "reset the project locally (postgres/etc)",
			Action: typgo.TaskNames{
				"pg.drop", "pg.create", "pg.migrate", "pg.seed",
				"mysql.drop", "mysql.create", "mysql.migrate", "mysql.seed",
			},
		},
		// release
		&typrls.ReleaseTool{
			Before: typgo.TaskNames{"test", "build"},
			// Releaser:  &typrls.CrossCompiler{Targets: []typrls.Target{"darwin/amd64", "linux/amd64"}},
			Publisher: &typrls.Github{Owner: "typical-go", Repo: "typical-rest-server"},
		},
	},
}

func main() {
	typgo.Start(&descriptor)
}
