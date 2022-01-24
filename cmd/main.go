package main

import (
	"log"

	"github.com/keisuke-m123/gomoduler/cmd/domain"
	"github.com/keisuke-m123/gomoduler/cmd/root"
)

func main() {
	rootCmd := root.NewRootCmd()
	rootCmd.AddCommand(domain.NewDomainAnalysisCmd())
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("execute error: %v", err)
	}
}
