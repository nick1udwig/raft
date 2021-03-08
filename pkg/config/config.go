package config

const (
	// The element in the slice that must be modified
	//  based on user input: the bare metal Urbit dir.
	WindowsDockerCmdArgsUrbitDirSource = 6
	WindowsDockerCmdName               = "docker"
	WindowsDockerUrl                   = "https://desktop.docker.com/win/stable/Docker%20Desktop%20Installer.exe"
)

var (
	// golang does not allow lists as consts,
	//  otherwise this would be `const`.
	WindowsDockerCmdArgs = []string{
		"run",
		"-P",
		"-d",
		"--restart",
		"always",
		"-v",
		"",
		"alpine",
		"ls",
		"/urbit",
		"--name",
		"urbit",
		"tloncorp/urbit",
	}
)
