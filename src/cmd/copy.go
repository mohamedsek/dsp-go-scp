package cmd

import (
	"fmt"
	"math"
	"os"
	"runtime"
	"strings"

	"github.com/mohamedsek/dsp-go-scp.git/src/fs"
	"github.com/mohamedsek/dsp-go-scp.git/src/multithread"
	"github.com/mohamedsek/dsp-go-scp.git/src/scp"
	"github.com/pkg/sftp"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var RootCommand = &cobra.Command{
	Use:   "copy <source> <target> <user> <remote-host> [flags]",
	Short: "Allows you to synchronize 2 directories recursively",
	Args:  cobra.MinimumNArgs(2),

	Run: func(cmd *cobra.Command, args []string) {
		origin := args[0]
		target := args[1]

		var target_user string
		var target_host string
		if len(args) == 4 {
			target_user = args[2]
			target_host = args[3]
		}

		fi, err := os.Stat(origin)
		if os.IsNotExist(err) {
			fmt.Println(err)
			os.Exit(1)
		}

		if fs.REMOTE {
			scp.SSH_CONNEXION_TO_TARGET, err = scp.ConnectToHost(target_user, target_host)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		// source is a file
		if !fi.IsDir() {
			if fs.PROGRESS {
				fs.PROGRESS_BAR = progressbar.Default(int64(1))
			}
			// clean target
			if fs.REMOTE {
				SFTP_CLIENT_TO_TARGET, err := sftp.NewClient(scp.SSH_CONNEXION_TO_TARGET)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				// defer SFTP_CLIENT_TO_TARGET.Close()
				err = SFTP_CLIENT_TO_TARGET.Remove(target)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			} else {
				if err := os.RemoveAll(target); err != nil {
					os.Exit(1)
				}
			}
			// copy file origin -> target
			if fs.REMOTE {
				// new SFTP client to target
				SFTP_CLIENT_TO_TARGET, err := sftp.NewClient(scp.SSH_CONNEXION_TO_TARGET)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				// defer SFTP_CLIENT_TO_TARGET.Close()
				fs.CopyFileRemoteTarget(origin, target, SFTP_CLIENT_TO_TARGET)
			} else {
				fs.CopyFileLocal(origin, target)
			}

		} else {
			// source is a folder

			// scan files and folders
			folders, files, err := fs.Scan(origin, target)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			// initialize progress bar
			if fs.PROGRESS {
				fs.PROGRESS_BAR = progressbar.Default(int64(len(files) + len(folders)))
			}
			// clean target
			if fs.REMOTE {
				// new SFTP client to target
				SFTP_CLIENT_TO_TARGET, err := sftp.NewClient(scp.SSH_CONNEXION_TO_TARGET)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				// defer SFTP_CLIENT_TO_TARGET.Close()
				if fi, _ := SFTP_CLIENT_TO_TARGET.Stat(target); fi != nil {
					maxNbrRemoteFiles, _ := fs.ScanRemote(target, SFTP_CLIENT_TO_TARGET)
					for i := 0; i <= maxNbrRemoteFiles; i++ {
						if fi, _ := SFTP_CLIENT_TO_TARGET.Stat(target); fi != nil {
							err := fs.RemoveFolder(target, SFTP_CLIENT_TO_TARGET)
							if err != nil {
								fmt.Println(err)
								os.Exit(1)
							} else {
								break
							}

						}
					}
				}
			} else {
				if fi, _ := os.Stat(target); fi != nil {
					if err := os.RemoveAll(target); err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
				}
			}

			// copy folders and files origin -> target

			if fs.MULTI_THREADS {

				total := float64(len(files) + len(folders))
				cpus := float64(runtime.NumCPU())
				chunck := int(math.Ceil((total / cpus)))
				tasks := make(chan bool)

				if fs.REMOTE {
					// limit number of threads for a remote copy
					if cpus > float64(6) {
						chunck = int(math.Ceil((total / float64(fs.SSH_CONNEXIONS_LIMIT))))
					}
				}

				// FOLDER
				// chunck = int(math.Ceil((float64(len(folders)) / cpus)))
				missions := multithread.SplitFiles(folders, chunck)
				for _, f := range missions {
					go multithread.CopyFolderThread(tasks, f)
				}
				var missionsComplete int = 0
				for range tasks {
					missionsComplete++
					if missionsComplete == len(missions) {
						break
					}
				}

				// FILE
				// chunck = int(math.Ceil((float64(len(files)) / cpus)))
				missions = multithread.SplitFiles(files, chunck)
				for _, f := range missions {
					go multithread.CopyFileThread(tasks, f)
				}
				missionsComplete = 0
				for range tasks {
					missionsComplete++
					if missionsComplete == len(missions) {
						break
					}
				}

			} else {

				// copy folders
				if fs.REMOTE {
					// new SFTP client to target
					SFTP_CLIENT_TO_TARGET, err := sftp.NewClient(scp.SSH_CONNEXION_TO_TARGET)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					defer SFTP_CLIENT_TO_TARGET.Close()

					for _, folder := range folders {
						s := strings.Split(folder, ":")
						err := fs.CopyFolderRemoteTarget(s[0], s[1], SFTP_CLIENT_TO_TARGET)
						if err != nil {
							fmt.Println(err)
							os.Exit(1)
							break
						}
					}

				} else {
					for _, folder := range folders {
						s := strings.Split(folder, ":")
						err := fs.CopyFolderLocal(s[0], s[1])
						if err != nil {
							fmt.Println(err)
							os.Exit(1)
							break
						}
					}

				}

				// copy files
				if fs.REMOTE {
					// new SFTP client to target
					SFTP_CLIENT_TO_TARGET, err := sftp.NewClient(scp.SSH_CONNEXION_TO_TARGET)
					if err != nil {
						fmt.Println(err)
					}
					defer SFTP_CLIENT_TO_TARGET.Close()

					for _, file := range files {
						s := strings.Split(file, ":")
						err := fs.CopyFileRemoteTarget(s[0], s[1], SFTP_CLIENT_TO_TARGET)
						if err != nil {
							fmt.Println(err)
							os.Exit(1)
							break
						}
					}
				} else {
					for _, file := range files {
						s := strings.Split(file, ":")
						err := fs.CopyFileLocal(s[0], s[1])
						if err != nil {
							fmt.Println(err)
							os.Exit(1)
							break
						}
					}
				}
			}
		}
	},
}
