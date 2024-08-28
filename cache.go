package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const CACHEFILE = "cache.txt"

// Reads the cached addons and returns a map of addonName: addonID
func getCache() (map[string]int, error) {
	cacheMap := make(map[string]int)

	file, err := os.Open(CACHEFILE)
	if err != nil {
		os.Create(CACHEFILE)
		return cacheMap, fmt.Errorf("%v not found. Wow! Created it instead", CACHEFILE)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")

		if len(parts) == 2 {
			releaseID, err := strconv.Atoi(parts[1])
			if err != nil {
				log.Println("Warning: Skipping line due to error converting releaseID to integer:", err)
				continue
			}

			cacheMap[parts[0]] = releaseID
		} else {
			log.Println("Warning: Skipping line with invalid format:", line)
		}
	}
	if err := scanner.Err(); err != nil {
		return cacheMap, fmt.Errorf("error scanning cache file %v", err)
	}
	return cacheMap, nil
}

// Dumps the cache to the cache.txt file
func dumpCache(cacheMap map[string]int) error {
	file, err := os.Create(CACHEFILE)
	if err != nil {
		return fmt.Errorf("failed to create cache file %s: %v", CACHEFILE, err)
	}
	defer file.Close()

	for addonName, releaseID := range cacheMap {
		line := fmt.Sprintf("%s %d\n", addonName, releaseID)
		_, err := file.WriteString(line)
		if err != nil {
			return fmt.Errorf("failed to write to cache file %s: %v", CACHEFILE, err)
		}
	}
	return nil
}
