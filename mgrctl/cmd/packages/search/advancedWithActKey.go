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

type advancedWithActKeyFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	LuceneQuery          string
	ActivationKey          string
}

func advancedWithActKeyCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "advancedWithActKey",
		Short: "Advanced method to search lucene indexes with a passed in query written
 in Lucene Query Parser syntax, additionally this method will limit results to those
 which are associated with a given activation key.
 Lucene Query Parser syntax is defined at
 http://lucene.apache.org/java/3_5_0/queryparsersyntax.html_blank
 lucene.apache.org.
 Fields searchable for Packages:
 name, epoch, version, release, arch, description, summary
 Lucene Query Example: "name:kernel AND version:2.6.18 AND -description:devel"",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags advancedWithActKeyFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, advancedWithActKey)
		},
	}

	cmd.Flags().String("LuceneQuery", "", "a query written in the form of Lucene QueryParser Syntax")
	cmd.Flags().String("ActivationKey", "", "activation key to look for packages in")

	return cmd
}

func advancedWithActKey(globalFlags *types.GlobalFlags, flags *advancedWithActKeyFlags, cmd *cobra.Command, args []string) error {

res, err := search.Search(&flags.ConnectionDetails, flags.LuceneQuery, flags.ActivationKey)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

