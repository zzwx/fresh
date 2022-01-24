package runner

import (
	"encoding/csv"
	"github.com/bmatcuk/doublestar"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func initFolders() {
	if isDebug() {
		runnerLog("Initializing folders")
	}
	path := tmpPath()
	if _, errDir := os.Stat(path); os.IsNotExist(errDir) {
		runnerLog("Creating %s", path)
		err := os.Mkdir(path, 0755)

		if err != nil {
			runnerLog(err.Error())
		}
	}
}

func isTmpDir(path string) bool {
	absolutePath, _ := filepath.Abs(path)
	absoluteTmpPath, _ := filepath.Abs(tmpPath())

	return absolutePath == absoluteTmpPath
}

func isIgnored(path string) bool {
	path = filepath.Clean(path)
	r := csv.NewReader(strings.NewReader(settings.Ignore))
	r.LazyQuotes = true
	r.Comma = ','
	r.TrimLeadingSpace = true
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		for i := 0; i < len(record); i++ {
			cleanIgnore := filepath.Clean(strings.TrimSpace(record[i])) // trim surrounding spaces
			m, err := doublestar.PathMatch(cleanIgnore, path)
			if err == nil && m {
				return true
			}
		}
	}
	return false
}

func isWatchedExt(path string) bool {
	absolutePath, _ := filepath.Abs(path)         // it auto-calls Clean
	absoluteTmpPath, _ := filepath.Abs(tmpPath()) // it auto-calls Clean

	if strings.HasPrefix(absolutePath, absoluteTmpPath+string(filepath.Separator)) {
		return false
	}

	ext := filepath.Ext(path)

	r := csv.NewReader(strings.NewReader(settings.ValidExt))
	r.LazyQuotes = true
	r.Comma = ','
	r.TrimLeadingSpace = true
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		for i := 0; i < len(record); i++ {
			if strings.TrimSpace(record[i]) == ext {
				return true
			}
		}
	}
	//for _, e := range strings.Split(settings["valid_ext"], ",") {
	//	if strings.TrimSpace(e) == ext {
	//		return true
	//	}
	//}

	return false
}

func shouldRebuild(eventName string) bool {
	//r := csv.NewReader(strings.NewReader(settings["no_rebuild_ext"]))
	//r.LazyQuotes = true
	//r.Comma = ','
	//r.TrimLeadingSpace = true
	//for {
	//  record, err := r.Read()
	//  if err == io.EOF {
	//    break
	//  }
	//  if err != nil {
	//    log.Fatal(err)
	//  }
	//  for i := 0; i < len(record); i++ {
	//    ext := strings.TrimSpace(record[i]) // trim surrounding spaces
	//  }
	//}

	lastColonPos := strings.LastIndex(eventName, ":")
	fileName := eventName
	if lastColonPos >= 0 {
		fileName = filepath.Clean(eventName[0:lastColonPos])
	}
	if fileName[0] == '"' {
		fileName = fileName[1:]
	}
	if fileName[len(fileName)-1] == '"' {
		fileName = fileName[0 : len(fileName)-1]
	}

	if isIgnored(fileName) {
		return false
	}
	ext := filepath.Ext(fileName)

	r := csv.NewReader(strings.NewReader(settings.NoRebuildExt))
	r.LazyQuotes = true
	r.Comma = ','
	r.TrimLeadingSpace = true
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		for i := 0; i < len(record); i++ {
			if strings.TrimSpace(record[i]) == ext {
				return false
			}
		}
	}

	//for _, e := range strings.Split(settings["no_rebuild_ext"], ",") {
	//	e = strings.TrimSpace(e)
	//	//fileName := filepath.Clean(strings.Replace(strings.Split(eventName, ":")[0], `"`, "", -1))
	//	if strings.HasSuffix(fileName, e) {
	//		return false
	//	}
	//	if isIgnored(fileName) {
	//		return false
	//	}
	//}

	return true
}

func createBuildErrorsLog(message string) bool {
	file, err := os.Create(buildErrorsFilePath())
	if err != nil {
		return false
	}
	defer file.Close()

	_, err = file.WriteString(message)
	if err != nil {
		return false
	}

	return true
}

func removeBuildErrorsLog() error {
	err := os.Remove(buildErrorsFilePath())

	return err
}
