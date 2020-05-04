package cmd

import (
	"fmt"
	"github.com/lf-edge/eden/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"runtime"
)

var (
	adamDist         string
	adamPort         string
	adamForce        bool
	eserverImageDist string
	eserverPort      string
	eserverPidFile   string
	eserverLogFile   string
	evePidFile       string
	eveLogFile       string
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start harness",
	Long:  `Start harness.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		assingCobraToViper(cmd)
		viperLoaded, err := utils.LoadConfigFile(config)
		if err != nil {
			return fmt.Errorf("error reading config: %s", err.Error())
		}
		if viperLoaded {
			eveImageFile = utils.ResolveAbsPath(viper.GetString("eve.image-file"))
			adamPort = viper.GetString("adam.port")
			adamDist = utils.ResolveAbsPath(viper.GetString("adam.dist"))
			adamForce = viper.GetBool("adam.force")
			eserverImageDist = utils.ResolveAbsPath(viper.GetString("eden.images.dist"))
			eserverPort = viper.GetString("eden.eserver.port")
			eserverPidFile = utils.ResolveAbsPath(viper.GetString("eden.eserver.pid"))
			eserverLogFile = utils.ResolveAbsPath(viper.GetString("eden.eserver.log"))
			qemuARCH = viper.GetString("eve.arch")
			qemuOS = viper.GetString("eve.os")
			qemuAccel = viper.GetBool("eve.accel")
			qemuSMBIOSSerial = viper.GetString("eve.serial")
			qemuConfigFile = utils.ResolveAbsPath(viper.GetString("eve.qemu-config"))
			evePidFile = utils.ResolveAbsPath(viper.GetString("eve.pid"))
			eveLogFile = utils.ResolveAbsPath(viper.GetString("eve.log"))
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		adamPath, err := filepath.Abs(adamDist)
		if err != nil {
			log.Fatalf("adam-dist problems: %s", err)
		}
		command, err := os.Executable()
		if err != nil {
			log.Fatalf("cannot obtain executable path: %s", err)
		}
		if err := utils.StartAdam(adamPort, adamPath, adamForce, adamTag); err != nil {
			log.Errorf("cannot start adam: %s", err)
		} else {
			log.Infof("Adam is running and accesible on port %s", adamPort)
		}
		if err := utils.StartEServer(command, eserverPort, eserverImageDist, eserverLogFile, eserverPidFile); err != nil {
			log.Errorf("cannot start eserver: %s", err)
		} else {
			log.Infof("Eserver is running and accesible on port %s", eserverPort)
		}
		if err := utils.StartEVEQemu(command, qemuARCH, qemuOS, eveImageFile, qemuSMBIOSSerial, qemuAccel, qemuConfigFile, eveLogFile, evePidFile); err != nil {
			log.Errorf("cannot start eve: %s", err)
		} else {
			log.Infof("EVE is starting")
		}
	},
}

func startInit() {
	currentPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	startCmd.Flags().StringVarP(&adamDist, "adam-dist", "", filepath.Join(currentPath, "dist", "adam"), "adam dist to start (required)")
	startCmd.Flags().StringVarP(&adamPort, "adam-port", "", "3333", "adam dist to start")
	startCmd.Flags().BoolVarP(&adamForce, "adam-force", "", false, "adam force rebuild")
	startCmd.Flags().StringVarP(&eserverImageDist, "image-dist", "", filepath.Join(currentPath, "dist", "images"), "image dist for eserver")
	startCmd.Flags().StringVarP(&eserverPort, "eserver-port", "", "8888", "eserver port")
	startCmd.Flags().StringVarP(&eserverPidFile, "eserver-pid", "", filepath.Join(currentPath, "dist", "eserver.pid"), "file for save eserver pid")
	startCmd.Flags().StringVarP(&eserverLogFile, "eserver-log", "", filepath.Join(currentPath, "dist", "eserver.log"), "file for save eserver log")
	startCmd.Flags().StringVarP(&qemuARCH, "eve-arch", "", runtime.GOARCH, "arch of system")
	startCmd.Flags().StringVarP(&qemuOS, "eve-os", "", runtime.GOOS, "os to run on")
	startCmd.Flags().BoolVarP(&qemuAccel, "eve-accel", "", true, "use acceleration")
	startCmd.Flags().StringVarP(&qemuSMBIOSSerial, "eve-serial", "", "", "SMBIOS serial")
	startCmd.Flags().StringVarP(&qemuConfigFile, "qemu-config", "", filepath.Join(currentPath, "dist", "qemu.conf"), "config file to use")
	startCmd.Flags().StringVarP(&evePidFile, "eve-pid", "", filepath.Join(currentPath, "dist", "eve.pid"), "file for save EVE pid")
	startCmd.Flags().StringVarP(&eveLogFile, "eve-log", "", filepath.Join(currentPath, "dist", "eve.log"), "file for save EVE log")
	startCmd.Flags().StringVar(&config, "config", "", "path to config file")
}