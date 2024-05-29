package function_util

import (
	"io"
	"minik8s/pkg/api"
	"minik8s/pkg/config"
	"minik8s/util/log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// We have a docker registry running at localhost:5050
// docker run -d -p 5050:5000 --name minik8s_registry registry:2
const RegistryPort = "5050"

func CreateImage(function *api.Function) (string, error) {
	// create image
	// BUG: docker build can only use relative path
	// Fix: copy source code to a temp dir and build image
	curPath := "/root/minik8s/pkg/serverless/function/function_util"
	tempPath := filepath.Join(curPath, function.Metadata.UUID)
	err := os.Mkdir(tempPath, os.ModePerm)
	if err != nil {
		return "", err
	}
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			log.Error("error remove temp dir: %s", err)
		}
	}(tempPath)

	// copy source code to temp dir
	err = copyDirContents(function.FilePath, tempPath)
	// wait for copy
	time.Sleep(1 * time.Second)

	// Step 1: Build Image
	// docker build --build-arg SOURCE_DIR=/path/to/source -t my-python-app .
	cmd := exec.Command("docker", "build", "--build-arg", "SOURCE_DIR="+function.Metadata.UUID, "-t",
		GetImageName(function.Metadata.Name, function.Metadata.NameSpace), curPath)
	output, err := cmd.CombinedOutput()
	//log.Info("output: %s", string(output))
	if err != nil {
		return "", err
	}

	// Step 2: Tag Image
	// docker tag myimage:latest localhost:5000/myimage:latest
	cmd = exec.Command("docker", "tag", GetImageName(function.Metadata.Name, function.Metadata.NameSpace),
		config.Remotehost+":"+RegistryPort+"/"+GetImageName(function.Metadata.Name, function.Metadata.NameSpace))
	log.Info("cmd: %s", cmd.String())
	output, err = cmd.CombinedOutput()
	//log.Info("output: %s", string(output))
	if err != nil {
		return "", err
	}

	// Step 3: Push Image
	// docker push localhost:5000/myimage:latest
	cmd = exec.Command("docker", "push",
		config.Remotehost+":"+RegistryPort+"/"+GetImageName(function.Metadata.Name, function.Metadata.NameSpace))
	//log.Info("cmd: %s", cmd.String())
	output, err = cmd.CombinedOutput()
	log.Info("output: %s", string(output))
	if err != nil {
		return "", err
	}
	return config.Remotehost + ":" + RegistryPort + "/" + GetImageName(function.Metadata.Name, function.Metadata.NameSpace), nil
}

func DeleteFunctionImage(name string, namespace string) error {
	// Step 0: Delete Image in namespace
	cmd := exec.Command("nerdctl", "-n", namespace, "rm",
		config.Remotehost+":"+RegistryPort+"/"+GetImageName(name, namespace))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	// Step 1: Delete Image
	cmd = exec.Command("docker", "rmi", GetImageName(name, namespace))
	output, err = cmd.CombinedOutput()
	//log.Info("output: %s", string(output))
	if err != nil {
		return err
	}

	// Step 2: Delete Image in local registry
	cmd = exec.Command("docker", "rmi",
		config.Remotehost+":"+RegistryPort+"/"+GetImageName(name, namespace))
	log.Info("cmd: %s", cmd.String())
	output, err = cmd.CombinedOutput()
	log.Info("output: %s", string(output))
	if err != nil {
		return err
	}
	return nil
}
func copyDirContents(srcDir, dstDir string) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Construct the destination path
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		destPath := filepath.Join(dstDir, relPath)

		if info.IsDir() {
			// Create destination directory
			err := os.MkdirAll(destPath, info.Mode())
			if err != nil {
				return err
			}
		} else {
			// Copy file
			err := copyFile(path, destPath)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	// Ensure the copied file has the same permissions as the source file
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, srcInfo.Mode())
}
