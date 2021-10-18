package pkg

import "os"

//FileExist file exist?
func FileExist(pid string) bool {
	_, err := os.Stat(pid)
	exist := true
	if err != nil && os.IsNotExist(err) {
		exist = false
	}

	return exist
}
