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

type advancedWithChannelFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	LuceneQuery          string
	ChannelLabel          string
}

func advancedWithChannelCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "advancedWithChannel",
		Short: "Advanced method to search lucene indexes with a passed in query written
 in Lucene Query Parser syntax, additionally this method will limit results to those
 which are in the passed in channel label.
 Lucene Query Parser syntax is defined at
 http://lucene.apache.org/java/3_5_0/queryparsersyntax.html_blank
 lucene.apache.org.
 Fields searchable for Packages:
 name, epoch, version, release, arch, description, summary
 Lucene Query Example: "name:kernel AND version:2.6.18 AND -description:devel"",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags advancedWithChannelFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, advancedWithChannel)
		},
	}

	cmd.Flags().String("LuceneQuery", "", "a query written in the form of Lucene QueryParser Syntax")
	cmd.Flags().String("ChannelLabel", "", "the channel Label")

	return cmd
}

func advancedWithChannel(globalFlags *types.GlobalFlags, flags *advancedWithChannelFlags, cmd *cobra.Command, args []string) error {

res, err := search.Search(&flags.ConnectionDetails, flags.LuceneQuery, flags.ChannelLabel)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

