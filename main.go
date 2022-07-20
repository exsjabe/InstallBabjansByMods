package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	spinner "github.com/alecrabbit/go-cli-spinner"
	"github.com/fatih/color"
	"github.com/zs5460/art"
)

func main() {
	mcDir := os.Getenv("APPDATA") + "\\.minecraft"
	appDataDir := os.Getenv("APPDATA") + "\\BabjansByTemp"
	zipsDownload := "https://github.com/exsjabe/BabjansByZips/archive/refs/heads/master.zip"

	toExtract := ExtractAddresses{
		ExtractAddress{".ResourcePacks.zip", "resourcepacks"},
		ExtractAddress{".Shaders.zip", "shaderpacks"},
	}

	if !commandExists("java") {
		log.Fatal("Please install Java before initializing mods.")
	}

	if !pathExists(mcDir) {
		log.Fatal("Please install Minecraft before initializing mods.")
	}

	print(art.String("BabjansBy"))
	fmt.Println("")

	if !pathExists(mcDir + "\\versions\\1.16.5-forge-36.2.23") {
		installForge(appDataDir, mcDir)
	}

	s, _ := spinner.New()
	s.Message("Downloading mods...")
	s.Start()
	fileName := downloadFile(zipsDownload, appDataDir)
	s.Stop()

	Unzip(appDataDir+"\\"+fileName, appDataDir)

	folderName := appDataDir + "\\" + strings.Split(fileName, ".zip")[0]

	toExtract.extract(folderName, appDataDir, mcDir)

	zips, err := os.ReadDir(folderName)

	checkError(err)

	os.MkdirAll(mcDir+"\\mods", 0755)

	for _, zip := range zips {
		if toExtract.includesZipfile(zip.Name()) {
			continue
		}

		err := Unzip(folderName+"\\"+zip.Name(), appDataDir)

		checkError(err)

		fullDir := appDataDir + "\\" + strings.Split(zip.Name(), ".zip")[0]

		files, err := os.ReadDir(fullDir)

		checkError(err)

		for _, file := range files {
			if !pathExists(mcDir + "\\mods\\" + file.Name()) {
				err := CopyFile(fullDir+"\\"+file.Name(), mcDir+"\\mods\\"+file.Name())
				checkError(err)
				color.Yellow("Successfully added " + file.Name() + " to " + "mods folder.")
			}
		}
	}

	s.Message("Cleaning up...")
	s.Start()
	os.RemoveAll(appDataDir)
	s.Stop()
	color.Magenta("Removed all temp files.")
	color.Green("Successfully installed all BabjansBy mods!.")

}

func installForge(workdir string, dest string) {
	s, _ := spinner.New()
	s.Message("Downloading forge...")
	s.Start()

	forgeDownload := "https://github.com/exsjabe/1.16.5-forge-36.2.23/archive/refs/heads/master.zip"
	forgeZip := downloadFile(forgeDownload, workdir)
	s.Stop()
	s.Message("Extracting forge...")
	s.Start()
	err := Unzip(workdir+"\\"+forgeZip, workdir)
	checkError(err)
	forgeFolder := workdir + "\\" + strings.Split(forgeZip, ".zip")[0]

	fullDir := forgeFolder + "\\forge-1.16.5-36.2.23-installer.jar"
	s.Stop()

	err = exec.Command("java", "-jar", fullDir).Run()
	os.RemoveAll(".//forge-1.16.5-36.2.23-installer.jar.log")
	checkError(err)
	color.Green("Successfully installed Forge!")
}
