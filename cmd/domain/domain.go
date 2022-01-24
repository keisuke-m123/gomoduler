package domain

import (
	"fmt"
	"log"

	"github.com/keisuke-m123/gomoduler/analyzer/configuration"
	"github.com/keisuke-m123/gomoduler/analyzer/domain"
	"github.com/keisuke-m123/gomoduler/cmd/root"
	"github.com/spf13/cobra"
	"golang.org/x/tools/go/analysis/multichecker"
)

func NewDomainAnalysisCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "domain",
		Short: "domain analysis",
		Run:   run,
	}
	return cmd
}

func run(cmd *cobra.Command, args []string) {
	fmt.Println(root.ConfigFilePath())
	conf, err := configuration.LoadConfig(root.ConfigFilePath())
	if err != nil {
		log.Fatalf("failed to load config: %s", err)
	}
	multichecker.Main(domain.NewDomainAnalyzer(conf.DomainConfig.Paths))
}
