package search

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/packages/search"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type advancedFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	LuceneQuery          string
}

func advancedCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "advanced",
		Short: "Advanced method to search lucene indexes with a passed in query written
 in Lucene Query Parser syntax.
 Lucene Query Parser syntax is defined at
 http://lucene.apache.org/java/3_5_0/queryparsersyntax.html_blank
 lucene.apache.org.
 Fields searchable for Packages:
 name, epoch, version, release, arch, description, summary
 Lucene Query Example: "name:kernel AND version:2.6.18 AND -description:devel"",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags advancedFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, advanced)
		},
	}

	cmd.Flags().String("LuceneQuery", "", "a query written in the form of Lucene QueryParser Syntax")

	return cmd
}

func advanced(globalFlags *types.GlobalFlags, flags *advancedFlags, cmd *cobra.Command, args []string) error {

res, err := search.Search(&flags.ConnectionDetails, flags.LuceneQuery)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

