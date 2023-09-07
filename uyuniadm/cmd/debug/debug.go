package debug

import (
	"log"
	"github.com/spf13/cobra"
	"github.com/uyuni-project/uyuni-tools/shared/types"
	"github.com/uyuni-project/uyuni-tools/shared/utils"
)

type flagpole struct {
	Enable			bool
	Restart			bool
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
	if flags.Enable {
		log.Printf("Enabling remote debug")
		utils.Exec(globalFlags, "", false, false, []string{}, "grep -q 'JAVA_OPTS=\" \\$JAVA_OPTS -Xdebug -Xrunjdwp:transport=dt_socket,address=\\*:8000,server=y,suspend=n\"' /etc/tomcat/conf.d/remote_debug.conf || echo 'JAVA_OPTS=\" $JAVA_OPTS -Xdebug -Xrunjdwp:transport=dt_socket,address=*:8000,server=y,suspend=n\" ' >> /etc/tomcat/conf.d/remote_debug.conf")
		utils.Exec(globalFlags, "", false, false, []string{}, "grep -q 'JAVA_OPTS=\" \\$JAVA_OPTS -Xdebug -Xrunjdwp:transport=dt_socket,address=\\*:8001,server=y,suspend=n\"' /etc/rhn/taskomatic.conf || echo 'JAVA_OPTS=\" $JAVA_OPTS -Xdebug -Xrunjdwp:transport=dt_socket,address=*:8001,server=y,suspend=n\" ' >> /etc/rhn/taskomatic.conf")
	} else {
		log.Printf("Disabling remote debug")
		utils.Exec(globalFlags, "", false, false, []string{}, "sed -i '/-Xdebug -Xrunjdwp:transport=dt_socket,address=/d' /etc/tomcat/conf.d/remote_debug.conf")
		utils.Exec(globalFlags, "", false, false, []string{}, "sed -i '/-Xdebug -Xrunjdwp:transport=dt_socket,address=/d' /etc/rhn/taskomatic.conf")
	}
	
	if flags.Restart {
		log.Printf("Running spacewalk-service restart. Changes will be applied after restart ")
		utils.Exec(globalFlags, "", false, false, []string{}, "spacewalk-service restart")
	} else {
		log.Printf("Configuration files changes but not applied. Please restart spacewalk-service to apply them")
	}
}
