package actions

import (
	"fmt"
	"github.com/aemengo/concourse-worker-manager/config"
	"github.com/julienschmidt/httprouter"
	"github.com/mholt/archiver"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func (a *Action) Install(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	a.logger.Printf("Received [%s] %s ...\n", r.Method, r.URL.Path)

	if a.concourseCmd != nil {
		fmt.Fprintln(w, "Terminating previous concourse process")
		a.concourseCmd.Process.Kill()
		a.concourseCmd.Process.Release()
	}

	var (
		version         = ps.ByName("version")
		downloadLink    = downloadLink(version)
		fileName        = downloadName(version)
		destinationPath = filepath.Join(config.Homedir(), fileName)
	)

	fmt.Fprintf(w, "Downloading %q ...\n", downloadLink)
	err := downloadFile(downloadLink, destinationPath)
	if err != nil {
		fmt.Fprintf(w, "Error: %s\n", err)
		a.logger.Printf("Error: %s\n", err)
		return
	}

	fmt.Fprintf(w, "Prepping %q ...\n", config.Homedir())
	err = os.RemoveAll(filepath.Join(config.Homedir(), "concourse"))
	if err != nil {
		fmt.Fprintf(w, "Error: %s\n", err)
		a.logger.Printf("Error: %s\n", err)
		return
	}

	fmt.Fprintf(w, "Unpacking %q ...\n", destinationPath)
	err = archiver.Unarchive(destinationPath, config.Homedir())
	if err != nil {
		fmt.Fprintf(w, "Error: %s\n", err)
		a.logger.Printf("Error: %s\n", err)
		return
	}

	var (
		executable = filepath.Join(config.Homedir(), "concourse", "bin", "concourse")
		args       = a.cmdArgs()
		logPath    = filepath.Join(config.Homedir(), "concourse.log")
	)

	logFile, err := os.Create(logPath)
	if err != nil {
		fmt.Fprintf(w, "Error: %s\n", err)
		a.logger.Printf("Error: %s\n", err)
		return
	}

	fmt.Fprintf(w, "Executing: %s %s\n", executable, strings.Join(args, " "))
	cmd := exec.Command(executable, args...)
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	err = cmd.Start()
	if err != nil {
		fmt.Fprintf(w, "Error: %s\n", err)
		a.logger.Printf("Error: %s\n", err)
		return
	}

	a.concourseCmd = cmd
	fmt.Fprintln(w, "Success")
	a.logger.Println("Success")
}

func (a *Action) cmdArgs() []string {
	var args []string

	switch runtime.GOOS {
	case "linux":
		args = []string{
			"worker",
			"--name", a.workerTag + "-" + getPublicIP(),
			"--tag", a.workerTag,
			"--work-dir", "/tmp/concourse",
			"--tsa-host", a.tsaHost,
			"--tsa-public-key", config.TsaHostKeyPath(),
			"--tsa-worker-private-key", config.WorkerKeyPath(),
			"--garden-use-houdini",
			"--baggageclaim-disable-user-namespaces",
		}
	case "darwin":
		args = []string{
			"worker",
			"--name", a.workerTag + "-" + getPublicIP(),
			"--tag", a.workerTag,
			"--work-dir", "/tmp/concourse",
			"--tsa-host", a.tsaHost,
			"--tsa-public-key", config.TsaHostKeyPath(),
			"--tsa-worker-private-key", config.WorkerKeyPath(),
		}
	case "windows":
		args = []string{
			"worker",
			"/name", a.workerTag + "-" + getPublicIP(),
			"/tag", a.workerTag,
			"/work-dir", `C:\concourse`,
			"/tsa-host", a.tsaHost,
			"/tsa-public-key", config.TsaHostKeyPath(),
			"/tsa-worker-private-key", config.WorkerKeyPath(),
		}
	}

	return args
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

func getPublicIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}

		if ip != nil && isPublicIP(ip) {
			return ip.String()
		}
	}

	return ""
}

func isPublicIP(IP net.IP) bool {
	if IP.IsLoopback() || IP.IsLinkLocalMulticast() || IP.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := IP.To4(); ip4 != nil {
		switch {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		default:
			return true
		}
	}
	return false
}