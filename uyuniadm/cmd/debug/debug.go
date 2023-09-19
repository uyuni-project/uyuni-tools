package debug

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type flagpole struct {
	Enable  bool
	Restart bool
}

func NewCommand(globalFlags *types.GlobalFlags) *cobra.Command {
	flags := &flagpole{}

	debugCmd := &cobra.Command{
		Use:   "debug ",
		Short: "Enable/Disable remote debugging",
		Long: `Enable or Disable remote debugging. 
Remote debugging is available using:
- port 8000, for rhn
- port 8001, for taskomatic
`,
		Run: func(cmd *cobra.Command, args []string) {
			run(globalFlags, flags, cmd, args)
		},
	}

	debugCmd.Flags().BoolVar(&flags.Enable, "enable", true, "Enable or disable remote debugging")
	debugCmd.Flags().BoolVar(&flags.Restart, "restart", false, "Apply changes restarting spacewalk-service")

	return debugCmd
}

func run(globalFlags *types.GlobalFlags, flags *flagpole, cmd *cobra.Command, args []string) {
	const tomcatFile = "/etc/tomcat/conf.d/remote_debug.conf"
	const taskoFile = "/etc/rhn/taskomatic.conf"

	if flags.Enable {
		enableTemplate := "grep -q 'JAVA_OPTS=\" \\$JAVA_OPTS -Xdebug -Xrunjdwp:transport=dt_socket,address=\\*:%d,server=y,suspend=n\"' %s ||" +
			" echo 'JAVA_OPTS=\" $JAVA_OPTS -Xdebug -Xrunjdwp:transport=dt_socket,address=*:%d,server=y,suspend=n\" ' >> %s"

		log.Printf("Enabling remote debug")
		utils.Exec(globalFlags, "", false, false, true, []string{}, fmt.Sprintf(enableTemplate, 8000, tomcatFile, 8000, tomcatFile))
		utils.Exec(globalFlags, "", false, false, true, []string{}, fmt.Sprintf(enableTemplate, 8001, taskoFile, 8001, taskoFile))
	} else {
		log.Printf("Disabling remote debug")
		utils.Exec(globalFlags, "", false, false, true, []string{}, "sed -i '/-Xdebug -Xrunjdwp:transport=dt_socket,address=/d' "+tomcatFile)
		utils.Exec(globalFlags, "", false, false, true, []string{}, "sed -i '/-Xdebug -Xrunjdwp:transport=dt_socket,address=/d' "+taskoFile)
	}

	if flags.Restart {
		log.Printf("Running spacewalk-service restart. Changes will be applied after restart ")
		utils.Exec(globalFlags, "", false, false, true, []string{}, "spacewalk-service restart")
	} else {
		log.Printf("Configuration files changed but not applied. Please restart spacewalk-service to apply them")
	}
}
