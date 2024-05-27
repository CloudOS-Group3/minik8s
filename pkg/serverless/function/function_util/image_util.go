package function_util

import (
	"io"
	"minik8s/pkg/api"
	"minik8s/util/log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func CreateImage(function *api.Function) error {
	// create image
	// BUG: docker build can only use relative path
	// Fix: copy source code to a temp dir and build image
	curPath := "/root/minik8s/pkg/serverless/function/function_util"
	tempPath := filepath.Join(curPath, function.Metadata.UUID)
	//temp, err := os.MkdirTemp("", function.Metadata.UUID)
	err := os.Mkdir(tempPath, os.ModePerm)
	if err != nil {
		return err
	}
	//defer func(path string) {
	//	err := os.RemoveAll(path)
	//	if err != nil {
	//		log.Error("error remove temp dir: %s", err)
	//	}
	//}(temp)

	// copy source code to temp dir
	err = copyDirContents(function.FilePath, tempPath)
	// wait for copy
	time.Sleep(1 * time.Second)

	// docker build --build-arg SOURCE_DIR=/path/to/source -t my-python-app .
	cmd := exec.Command("docker", "build", "--build-arg", "SOURCE_DIR="+function.Metadata.UUID, "-t",
		GetImageName(function.Metadata.Name, function.Metadata.NameSpace), curPath)
	log.Info("cmd: %s", cmd.String())
	output, err := cmd.CombinedOutput()
	log.Info("output: %s", string(output))
	if err != nil {
		return err
	}
	log.Info("output: %s", string(output))

	// push image
	// docker push my-python-app
	//cmd = exec.Command("docker", "push", GetImageName(function.Metadata.Name, function.Metadata.NameSpace))
	//log.Info("cmd: %s", cmd.String())
	//output, err = cmd.CombinedOutput()
	//log.Info("output: %s", string(output))
	//if err != nil {
	//	return err
	//}
	//log.Info("output: %s", string(output))
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
