package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	header := `
░░░░░░░░░▄░░░░░░░░░░░░░░▄░░░░
░░░░░░░░▌▒█░░░░░░░░░░░▄▀▒▌░░░
░░░░░░░░▌▒▒█░░░░░░░░▄▀▒▒▒▐░░░
░░░░░░░▐▄▀▒▒▀▀▀▀▄▄▄▀▒▒▒▒▒▐░░░
░░░░░▄▄▀▒░▒▒▒▒▒▒▒▒▒█▒▒▄█▒▐░░░
░░░▄▀▒▒▒░░░▒▒▒░░░▒▒▒▀██▀▒▌░░░ 
░░▐▒▒▒▄▄▒▒▒▒░░░▒▒▒▒▒▒▒▀▄▒▒▌░░
░░▌░░▌█▀▒▒▒▒▒▄▀█▄▒▒▒▒▒▒▒█▒▐░░ Wow!
░▐░░░▒▒▒▒▒▒▒▒▌██▀▒▒░░░▒▒▒▀▄▌░ Such updates!
░▌░▒▄██▄▒▒▒▒▒▒▒▒▒░░░░░░▒▒▒▒▌░ 
▐▒▀▐▄█▄█▌▄░▀▒▒░░░░░░░░░░▒▒▒▐░
▐▒▒▐▀▐▀▒░▄▄▒▄▒▒▒▒▒▒░▒░▒░▒▒▒▒▌
▐▒▒▒▀▀▄▄▒▒▒▄▒▒▒▒▒▒▒▒░▒░▒░▒▒▐░
░▌▒▒▒▒▒▒▀▀▀▒▒▒▒▒▒░▒░▒░▒░▒▒▒▌░
░▐▒▒▒▒▒▒▒▒▒▒▒▒▒▒░▒░▒░▒▒▄▒▒▐░░
░░▀▄▒▒▒▒▒▒▒▒▒▒▒░▒░▒░▒▄▒▒▒▒▌░░
░░░░▀▄▒▒▒▒▒▒▒▒▒▒▄▄▄▀▒▒▒▒▄▀░░░
░░░░░░▀▄▄▄▄▄▄▀▀▀▒▒▒▒▒▄▄▀░░░░░
░░░░░░░░░▒▒▒▒▒▒▒▒▒▒▀▀░░░░░░░░

`

	fmt.Print(header)
	config, err := NewConfig()
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		fmt.Print("Please check the config file! Press any key to close the program...")
		fmt.Scanf("h")
		return
	}

	enabledAddonModIDs := getEnabledAddonModIDs(config)
	fmt.Printf("Downloading updates to: %s\n\n", config.AddonPath)

	if len(enabledAddonModIDs) == 0 {
		fmt.Println("You configured no AddOns in the config.yaml. Please set your needed AddOns to True.")
		fmt.Print("Press any key to close the program...")
		fmt.Scanf("h")
		return
	}

	cacheMap, err := getCache()
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	for addon, modID := range enabledAddonModIDs {
		apiURL := fmt.Sprintf(config.API + strconv.Itoa(modID) + "/" + "files")

		fmt.Printf("===== AddOn: %s =====\n", addon)
		jsonData, err := queryAPI(apiURL)
		if err != nil {
			fmt.Printf("Error querying API for addon %s: %v\n", addon, err)
			continue
		}

		newestReleaseID, err := getNewestReleaseID(jsonData)
		if err != nil {
			fmt.Printf("Error getting newest release ID for addon %s: %v\n", addon, err)
			continue
		}

		if lastReleaseID, ok := cacheMap[addon]; ok {
			if newestReleaseID == lastReleaseID {
				fmt.Printf("Wowies! You are already up-to-date. Skipping...\n\n")
				continue
			}
		}

		downloadURL := fmt.Sprintf(config.API+"%d/files/%d/download", modID, newestReleaseID)
		err = downloadFile(downloadURL, config.AddonPath, addon)
		if err != nil {
			fmt.Printf("Error downloading file for addon %s: %v\n", addon, err)
			continue
		}
		fmt.Println("Downloaded file.")

		zipFilePath := filepath.Join(config.AddonPath, addon+".zip")
		err = unzip(zipFilePath, config.AddonPath)
		if err != nil {
			fmt.Printf("Error unzipping file for addon %s: %v\n", addon, err)
			continue
		}

		err = os.Remove(zipFilePath)
		if err != nil {
			fmt.Printf("Error deleting downloaded .zip file for addon %s: %v\n", addon, err)
			continue
		}

		fmt.Printf("Unzipped file and deleted remains.\n\n")

		// Update the cacheMap with the newestReleaseID after successful update
		cacheMap[addon] = newestReleaseID
	}

	dumpCache(cacheMap)
	fmt.Print("Wow! Much speed! Updates are done. Press any key to close the program...")
	fmt.Scanf("h")
}
