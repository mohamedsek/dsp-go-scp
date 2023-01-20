package fs

import (
	"fmt"
	"os"

	"github.com/pkg/sftp"
)

const PERMS_DIRECTORY = 755

func CopyFolderLocal(origin string, target string) error {
	if DEBUG {
		fmt.Println("COPY FOLDER:", origin, "->", target)
	}

	if PROGRESS {
		PROGRESS_BAR.Add(1)
	}
	return os.MkdirAll(target, PERMS_DIRECTORY)
}

func CopyFileLocal(origin string, target string) error {
	if DEBUG {
		fmt.Println("COPY FILE:", origin, "->", target)
	}
	if PROGRESS {
		PROGRESS_BAR.Add(1)
	}
	return copyFileLocal(origin, target)
}

func copyFileLocal(origin string, target string) error {
	fi, err := os.Open(origin)
	if err != nil {
		return err
	}
	defer fi.Close()

	fo, err := os.Create(target)
	if err != nil {
		return err
	}
	defer fo.Close()

	buffer := make([]byte, 1024)
	for {
		nb, err := fi.Read(buffer)
		if nb == 0 {
			break
		} else if err != nil {
			return err
		}

		if _, err := fo.Write(buffer[:nb]); err != nil {
			return err
		}
	}

	return nil
}

func RemoveFolder(dirToRemove string, SFTP_CLIENT *sftp.Client) error {
	entries, err := SFTP_CLIENT.ReadDir(dirToRemove)
	if err != nil {
		fmt.Println("error: reading folder entries")
		return err
	}
	for _, entry := range entries {

		if entry.IsDir() {
			RemoveFolder(dirToRemove+"/"+entry.Name(), SFTP_CLIENT)
		} else {
			err := SFTP_CLIENT.Remove(dirToRemove + "/" + entry.Name())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func CopyFolderRemoteTarget(origin string, target string, SFTP_CLIENT *sftp.Client) error {
	if DEBUG {
		fmt.Println("COPY FOLDER:", origin, "->", target)
	}

	if PROGRESS {
		PROGRESS_BAR.Add(1)
	}
	return SFTP_CLIENT.MkdirAll(target)
}

func CopyFileRemoteTarget(origin string, target string, SFTP_CLIENT *sftp.Client) error {
	if DEBUG {
		fmt.Println("COPY FILE:", origin, "->", target)
	}
	if PROGRESS {
		PROGRESS_BAR.Add(1)
	}
	return copyFileRemoteTarget(origin, target, SFTP_CLIENT)
}

func copyFileRemoteTarget(origin string, target string, SFTP_CLIENT *sftp.Client) error {
	fi, err := os.Open(origin)
	if err != nil {
		return err
	}
	defer fi.Close()

	fo, err := SFTP_CLIENT.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	if err != nil {
		return err
	}
	defer fo.Close()

	buffer := make([]byte, 1024)
	for {
		nb, err := fi.Read(buffer)
		if nb == 0 {
			break
		} else if err != nil {
			return err
		}

		if _, err := fo.Write(buffer[:nb]); err != nil {
			return err
		}
	}

	return nil
}
