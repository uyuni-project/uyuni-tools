package system

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/api"
	"github.com/uyuni-project/uyuni-tools/shared/api/system"
	apiTypes "github.com/uyuni-project/uyuni-tools/shared/api/types"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type addNoteFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	Subject          string
	Body          string
}

func addNoteCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "addNote",
		Short: "Add a new note to the given server.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags addNoteFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, addNote)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("Subject", "", "What the note is about.")
	cmd.Flags().String("Body", "", "Content of the note.")

	return cmd
}

func addNote(globalFlags *types.GlobalFlags, flags *addNoteFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid, flags.Subject, flags.Body)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

