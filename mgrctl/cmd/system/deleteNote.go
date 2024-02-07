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

type deleteNoteFlags struct {
	api.ConnectionDetails `mapstructure:"api"`
	Sid          int
	NoteId          int
}

func deleteNoteCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deleteNote",
		Short: "Deletes the given note from the server.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var flags deleteNoteFlags
			return utils.CommandHelper(globalFlags, cmd, args, &flags, deleteNote)
		},
	}

	cmd.Flags().String("Sid", "", "")
	cmd.Flags().String("NoteId", "", "")

	return cmd
}

func deleteNote(globalFlags *types.GlobalFlags, flags *deleteNoteFlags, cmd *cobra.Command, args []string) error {

res, err := system.System(&flags.ConnectionDetails, flags.Sid, flags.NoteId)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %v", res)

	return nil
}

