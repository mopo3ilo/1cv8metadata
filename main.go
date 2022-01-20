package main

import (
	"database/sql"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gen2brain/dlgs"
	"github.com/mopo3ilo/sql1cv8"
	"gopkg.in/yaml.v2"

	_ "github.com/denisenkom/go-mssqldb"
)

func openConfig() map[string]string {
	log.Println("open config")
	f, err := os.Open("bases.yml")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	cfg := map[string]string{}
	if err = yaml.Unmarshal(b, cfg); err != nil {
		log.Fatal(err)
	}
	return cfg
}

func openScript(pth string) (string, string) {
	log.Println("open script")
	sel, err := true, error(nil)
	if pth == "" {
		pth, sel, err = dlgs.File("Выберите файл скрипта", "*.sql", false)
		if err != nil {
			log.Fatal(err)
		}
		if !sel {
			log.Fatal("canceled")
		}
	}
	f, err := os.Open(pth)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	return pth, string(b)
}

func main() {
	var (
		fDatabase string
		fFileName string
		fSave     bool
		fExec     bool
	)
	flag.StringVar(&fDatabase, "database", "", "база данных из конфигурационного файла")
	flag.StringVar(&fFileName, "filename", "", "файл скрипта, который нужно обработать")
	flag.BoolVar(&fSave, "save", false, "сохранить скрипт после обработки")
	flag.BoolVar(&fExec, "exec", false, "выполнить скрипт после обработки")
	flag.Parse()

	log.Println("START")
	defer log.Println("STOP")

	var (
		sel bool
		err error
		txt string
		exc string
	)

	cfg := openConfig()
	log.Println("open database")
	if fDatabase == "" {
		lst := make([]string, 0, len(cfg))
		for k := range cfg {
			lst = append(lst, k)
		}
		sort.Strings(lst)
		fDatabase, sel, err = dlgs.List("Базы", "Выберите базу для обработки:", lst)
		if err != nil {
			log.Fatal(err)
		}
		if !sel {
			log.Fatal("canceled")
		}
	}
	con, ok := cfg[fDatabase]
	if !ok {
		log.Fatal("database not exist")
	}
	fDatabase += ".json"
	log.Println("load metadata")
	met, err := sql1cv8.LoadNewer(con, fDatabase)
	if err != nil {
		log.Fatal(err)
	}
	fFileName, txt = openScript(fFileName)
	scr, err := met.Parse(txt)
	if err != nil {
		log.Fatal(err)
	}
	if !fSave && !fExec {
		exc, sel, err = dlgs.List("Действие", "Что делать с выбранным скриптом:", []string{"Сохранить", "Выполнить"})
		if err != nil {
			log.Fatal(err)
		}
		if !sel {
			log.Fatal("canceled")
		}
	}
	if fSave || exc == "Сохранить" {
		dir, fnm := filepath.Split(fFileName)
		ext := filepath.Ext(fnm)
		fnm = strings.TrimSuffix(fnm, ext)
		fFileName = dir + fnm + "_new" + ext
		f, err := os.Create(fFileName)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		if _, err := f.WriteString(scr); err != nil {
			log.Fatal(err)
		}
	}
	if fExec || exc == "Выполнить" {
		db, err := sql.Open("sqlserver", con)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		scr = strings.ReplaceAll(scr, "\ngo\n", "\nGO\n")
		scr = strings.ReplaceAll(scr, "\nGo\n", "\nGO\n")
		scr = strings.ReplaceAll(scr, "\ngO\n", "\nGO\n")
		for _, v := range strings.Split(scr, "\nGO\n") {
			if _, err := db.Exec(v); err != nil {
				log.Fatal(err)
			}
		}
	}
}
