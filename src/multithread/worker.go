package multithread

import (
	"fmt"
	"strings"

	"github.com/mohamedsek/dsp-go-scp.git/src/fs"
	"github.com/mohamedsek/dsp-go-scp.git/src/scp"
	"github.com/pkg/sftp"
)

func Worker(channel chan error, callback func() error) {
	channel <- callback()
}

func SplitFiles(slice []string, chunkSize int) [][]string {
	var chunks [][]string
	for i := 0; i < len(slice); i = i + chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])
	}

	return chunks
}

func CopyFolderThread(ending chan bool, folders []string) {
	if fs.REMOTE {
		// new SFTP client to target
		SFTP_CLIENT_TO_TARGET, err := sftp.NewClient(scp.SSH_CONNEXION_TO_TARGET)
		if err != nil {
			fmt.Println(err)
		}
		// defer SFTP_CLIENT_TO_TARGET.Close()

		for _, folder := range folders {
			s := strings.Split(folder, ":")
			err := fs.CopyFolderRemoteTarget(s[0], s[1], SFTP_CLIENT_TO_TARGET)
			if err != nil {
				fmt.Println(err)
				ending <- true
				break
			}
		}

		ending <- true
	} else {
		for _, folder := range folders {
			s := strings.Split(folder, ":")
			err := fs.CopyFolderLocal(s[0], s[1])
			if err != nil {
				fmt.Println(err)
				ending <- true
				break
			}
		}

		ending <- true
	}
}

func CopyFileThread(ending chan bool, files []string) {

	if fs.REMOTE {
		// new SFTP client to target
		SFTP_CLIENT_TO_TARGET, err := sftp.NewClient(scp.SSH_CONNEXION_TO_TARGET)
		if err != nil {
			fmt.Println(err)
		}
		// defer SFTP_CLIENT_TO_TARGET.Close()

		for _, file := range files {
			s := strings.Split(file, ":")
			err := fs.CopyFileRemoteTarget(s[0], s[1], SFTP_CLIENT_TO_TARGET)
			if err != nil {
				fmt.Println(err)
				ending <- true
				break
			}
		}

		ending <- true
	} else {
		for _, file := range files {
			s := strings.Split(file, ":")
			err := fs.CopyFileLocal(s[0], s[1])
			if err != nil {
				fmt.Println(err)
				ending <- true
				break
			}
		}

		ending <- true
	}

}
