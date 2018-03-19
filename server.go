package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/groggygopher/oyster/handlers"
	"github.com/groggygopher/oyster/session"
)

var (
	// csvFile  = flag.String("csv_file", "", "The CSV file to import")
	// ruleFile = flag.String("rule_file", "", "The JSON encoded rule file")

	port    = flag.Int("port", 8080, "The port to serve HTTP on")
	saveDir = flag.String("save_dir", filepath.Join(os.TempDir(), "oyster"), "The directory to save user data")
)

func main() {
	sessMgr := session.NewManager(*saveDir)
	defer func() {
		if err := sessMgr.Close(); err != nil {
			log.Printf("error: sessionManager.Close: %v", err)
		}
	}()
	http.Handle("/session", handlers.NewSessionHandler(sessMgr))
	http.Handle("/upload", handlers.NewUploadHandler(sessMgr))
	http.Handle("/transactions", handlers.NewTransactionsHandler(sessMgr))

	http.Handle("/", http.FileServer(http.Dir("html")))
	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))

	// var err error
	// defer func() {
	// 	if err != nil {
	// 		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	// 		os.Exit(1)
	// 	}
	// }()

	// flag.Parse()

	// if *csvFile == "" {
	// 	err = errors.New("csv_file must be given")
	// 	return
	// }

	// f, err := os.Open(*csvFile)
	// if err != nil {
	// 	err = fmt.Errorf("Open(%s): %v", *csvFile, err)
	// }
	// defer f.Close()
	// trans, err := register.ReadAllTransactions(f)
	// if err != nil {
	// 	err = fmt.Errorf("ReadAllTransactions: %v", err)
	// 	return
	// }
	// fmt.Printf("Imported %d transactions\n", len(trans))

	// mngr := rule.NewManager()
	// if *ruleFile != "" {
	// 	var f *os.File
	// 	f, err = os.Open(*ruleFile)
	// 	if err != nil {
	// 		err = fmt.Errorf("Open(%s): %v", *ruleFile, err)
	// 		return
	// 	}
	// 	defer f.Close()
	// 	err = mngr.LoadRules(f)
	// 	if err != nil {
	// 		err = fmt.Errorf("LoadRules: %v", err)
	// 		return
	// 	}
	// }

	// var matched []*register.Transaction
	// for _, t := range trans {
	// 	var match bool
	// 	match, err = mngr.Evaluate(t)
	// 	if err != nil {
	// 		return
	// 	}
	// 	if match {
	// 		matched = append(matched, t)
	// 	}
	// }
	// fmt.Printf("Matched %d transactions\n", len(matched))
	// for _, t := range matched {
	// 	fmt.Println(t)
	// }

	// out, err := os.Create("rules.json")
	// if err != nil {
	// 	return
	// }
	// defer out.Close()
	// err = mngr.DumpRules(out)
	// if err != nil {
	// 	return
	// }
}
