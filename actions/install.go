package actions

import (
	"fmt"
	"github.com/aemengo/concourse-worker-manager/config"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

func (a *Action) Install(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	a.logger.Printf("Received [%s] %s ...\n", r.Method, r.URL.Path)

	var (
		version         = ps.ByName("version")
		downloadLink    = downloadLink(version)
		fileName        = downloadName(version)
		destinationPath = filepath.Join(config.Homedir(), fileName)
	)

	fmt.Fprintf(w, "Downloading %q to %q...\n", downloadLink, destinationPath)
	err := downloadFile(downloadLink, destinationPath)
	if err != nil {
		fmt.Fprintf(w, "Error: %s\n", err)
		a.logger.Printf("Error: %s\n", err)
		return
	}

	fmt.Fprintln(w, "Success")
	a.logger.Println("Success")
}

func downloadName(version string) string {
	switch runtime.GOOS {
	case "linux":
		return fmt.Sprintf("concourse-%s-linux-amd64.tgz", version)
	case "darwin":
		return fmt.Sprintf("concourse-%s-darwin-amd64.tgz", version)
	case "windows":
		return fmt.Sprintf("concourse-%s-windows-amd64.zip", version)
	default:
		return ""
	}
}

func downloadLink(version string) string {
	switch runtime.GOOS {
	case "linux":
		return fmt.Sprintf("https://github.com/concourse/concourse/releases/download/v%[1]s/concourse-%[1]s-linux-amd64.tgz", version)
	case "darwin":
		return fmt.Sprintf("https://github.com/concourse/concourse/releases/download/v%[1]s/concourse-%[1]s-darwin-amd64.tgz", version)
	case "windows":
		return fmt.Sprintf("https://github.com/concourse/concourse/releases/download/v%[1]s/concourse-%[1]s-windows-amd64.zip", version)
	default:
		return ""
	}
}

func downloadFile(url string, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
