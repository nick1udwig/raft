package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hosted-fornet/raft/pkg/config"

	"github.com/skratchdot/open-golang/open"
	"go.uber.org/zap"
)

var (
	sugar *zap.SugaredLogger
)

func downloadFile(url, filePath string) (err error) {
	response, err := http.Get(url)
	if err != nil {
		return
	}
	defer response.Body.Close()
	sugar.Debugw(
		"Got response code for url",
		"url", url,
		"responseCode", response.StatusCode,
	)

	file, err := os.Create(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	// TODO: Add progress reader;
	// perhaps https://github.com/machinebox/progress

	_, err = io.Copy(file, response.Body)

	return
}

func getDockerWindows(dockerPath string) (err error) {
	_, err = os.Stat(dockerPath)
	if err == nil {
		// File already exists.
		return
	}
	dockerUrl := config.WindowsDockerUrl
	err = downloadFile(dockerUrl, dockerPath)
	return
}

func getUserInput(prompt string) (userInput string, err error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s", prompt)
	userInput, err = reader.ReadString('\n')
	if err == nil {
		userInput = strings.TrimRight(userInput, "\n\r")
	}
	return
}

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	sugar = logger.Sugar()
	sugar.Debugw("Initialized logger.")

	// 1. Get configuration information from user.
	// Make Urbit dir if it does not exist.
	urbitDirPath, err := getUserInput("Enter path to store Urbit directory: ")
	if err != nil {
		sugar.Errorw(
			"Failed to get Urbit dir path.",
			"err", err,
		)
		return
	}
	urbitDirPath, err = filepath.Abs(urbitDirPath)
	if err != nil {
		sugar.Errorw(
			"Failed to set Urbit abs dir path.",
			"urbitDirPath", urbitDirPath,
			"err", err,
		)
		return
	}
	err = os.MkdirAll(urbitDirPath, os.ModePerm)
	if err != nil {
		sugar.Errorw(
			"Failed to make Urbit dir.",
			"urbitDirPath", urbitDirPath,
			"err", err,
		)
		return
	}

	// Determine how to set up Urbit for user.
	//  Allowed inputs:
	//  a. I'm just exploring, get me to Mars!
	//  b. I already have a ship and pier.
	//  c. I have a `.key` file to set up a ship.
	//  If a, set up a comet for user (and provide link to learn more about ship types).
	//  Else, if b, ask user to point to the pier.
	//  Else, if c, ask user to point to the `.key` file.
	//  Else, make a comet for user (and provide link to learn more about ship types).
	// TODO: Implement options; for dev purposes, just mine a coment (i.e. option a).
	userOptionString := "a"
	userOptionString = strings.ToLower(userOptionString)
	if userOptionString == "a" {
		// Comet pier will have name with prefix `myCometPrefix`
		//  and suffix an int, either 0 or the lowest number that
		//  does not appear in the Urbit dir (to avoid overwriting
		//  in case the user has instantiated multiple comets using
		//  this method).
		myCometPrefix := "myComet"
		urbitDirFiles, err := ioutil.ReadDir(urbitDirPath)
		if err != nil {
			sugar.Errorw(
				"Could not read Urbit dir. Bailing out.",
				"urbitDirPath", urbitDirPath,
				"err", err,
			)
			return
		}
		myCometNumber := 0
		for _, file := range urbitDirFiles {
			if strings.Contains(file.Name(), myCometPrefix) {
				fileCometNumber, err := strconv.Atoi(file.Name()[len(myCometPrefix):])
				if err != nil {
					sugar.Errorw(
						"Could not get comet number from existing comet pier.",
						"myCometPrefix", myCometPrefix,
						"fileName", file.Name(),
						"err", err,
					)
					return
				}
				myCometNumber = fileCometNumber + 1
			}
		}
		myCometName := fmt.Sprintf("%v%v.comet", myCometPrefix, myCometNumber)
		sugar.Infow(
			"User chose to mine a comet.",
			"myCometName", myCometName,
		)
	} else if userOptionString == "b" {
		sugar.Errorw(
			"Not yet implemented.",
			"userOptionString", userOptionString,
		)
		return
	} else if userOptionString == "c" {
		sugar.Errorw(
			"Not yet implemented.",
			"userOptionString", userOptionString,
		)
		return
	}

	// 2. Set up Docker.
	// TODO: Only get, install Docker if it is not already installed.
	setupDocker := true
	// GET Docker.
	if setupDocker {
		sugar.Debugw(
			"Got input from user. Creating dir and GETing Docker...",
			"urbitDirPath", urbitDirPath,
		)
		windowsDockerPath := filepath.Join(urbitDirPath, "docker.exe")
		err = getDockerWindows(windowsDockerPath)
		if err != nil {
			sugar.Errorw(
				"Failed to GET Docker for Windows.",
				"err", err,
			)
			return
		}

		// Install Docker.
		sugar.Infow(
			"Starting Docker Desktop installer. Follow the prompt to install.",
		)
		// err = open.Run(windowsDockerPath)
		installWindowsDockerCmd := exec.Command("cmd", "/C", "start", "/wait", windowsDockerPath)
		err := installWindowsDockerCmd.Run()
		if err != nil {
			sugar.Errorw(
				"Failed to install Docker for Windows.",
				"err", err,
			)
			return
		}

		sugar.Infow(
			"Docker installation successful.",
		)
	}

	// 3. Ensure Docker Desktop is running.
	err = open.Start("docker")
	if err != nil {
		sugar.Errorw(
			"Failed to start Docker.",
			"err", err,
		)
		return
	}

	// 4. Get the Urbit image and run a container.
	args := config.WindowsDockerCmdArgs
	args[config.WindowsDockerCmdArgsUrbitDirSource] = urbitDirPath
	startContainerCmd := exec.Command(
		config.WindowsDockerCmdName,
		args...,
	)
	err = startContainerCmd.Run()
	if err != nil {
		sugar.Errorw(
			"Failed to start the Urbit container.",
			"err", err,
		)
		return
	}

	sugar.Infow(
		"Success!",
	)
}
