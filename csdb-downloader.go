// [][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][]
// ------------------------------------------------------------------------------------------------
// csdb downloader by DKT/Samar
// ------------------------------------------------------------------------------------------------
// [][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][]

package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// [][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][]
// ------------------------------------------------------------------------------------------------
// Zmienne globalne
// ------------------------------------------------------------------------------------------------
// [][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][]

var cacheDir string
var sep string
var fullDownload bool

// RssItem - pojednyczy wpis w XML
// ------------------------------------------------------------------------------------------------
type RssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	GUID        string `xml:"guid"`
	PubDate     string `xml:"pubDate"`
}

// XMLRssFeed - tabela XML
// ------------------------------------------------------------------------------------------------
type XMLRssFeed struct {
	Items []RssItem `xml:"channel>item"`
}

// XMLHandle - kto jest autorem wydał
// ------------------------------------------------------------------------------------------------
type XMLHandle struct {
	ID        string `xml:"ID"`
	XMLHandle string `xml:"Handle"`
}

// XMLGroup - kto jest autorem wydał
// ------------------------------------------------------------------------------------------------
type XMLGroup struct {
	ID   string `xml:"ID"`
	Name string `xml:"Name"`
}

// XMLReleasedBy - kto wydał
// ------------------------------------------------------------------------------------------------
type XMLReleasedBy struct {
	XMLHandle []XMLHandle `xml:"Handle"`
	XMLGroup  []XMLGroup  `xml:"Group"`
}

// XMLCredit - XMLCredit za produkcję
// ------------------------------------------------------------------------------------------------
type XMLCredit struct {
	CreditType string    `xml:"CreditType"`
	XMLHandle  XMLHandle `xml:"Handle"`
}

// XMLDownloadLink - download links
// ------------------------------------------------------------------------------------------------
type XMLDownloadLink struct {
	Link string `xml:"Link"`
}

// XMLEvent - kompo
// ------------------------------------------------------------------------------------------------
type XMLEvent struct {
	ID   string `xml:"ID"`
	Name string `xml:"Name"`
}

// XMLReleasedAt - kompa
// ------------------------------------------------------------------------------------------------
type XMLReleasedAt struct {
	XMLEvent XMLEvent `xml:"Event"`
}

// XMLUsedSID - SIDy
// ------------------------------------------------------------------------------------------------
type XMLUsedSID struct {
	ID       string `xml:"ID"`
	HVSCPath string `xml:"HVSCPath"`
	Name     string `xml:"Name"`
	Author   string `xml:"Author"`
}

// XMLRelease - wydanie produkcji na csdb
// ------------------------------------------------------------------------------------------------
type XMLRelease struct {
	ReleaseID         string            `xml:"Release>ID"`
	ReleaseName       string            `xml:"Release>Name"`
	ReleaseType       string            `xml:"Release>Type"`
	ReleaseYear       string            `xml:"Release>ReleaseYear"`
	ReleaseMonth      string            `xml:"Release>ReleaseMonth"`
	ReleaseDay        string            `xml:"Release>ReleaseDay"`
	ReleaseScreenShot string            `xml:"Release>ScreenShot"`
	Rating            float32           `xml:"Release>Rating"`
	XMLReleasedBy     XMLReleasedBy     `xml:"Release>ReleasedBy"`
	XMLReleasedAt     XMLReleasedAt     `xml:"Release>ReleasedAt"`
	Credits           []XMLCredit       `xml:"Release>Credits>Credit"`
	DownloadLinks     []XMLDownloadLink `xml:"Release>DownloadLinks>DownloadLink"`
	UsedSIDs          []XMLUsedSID      `xml:"Release>UsedSIDs>SID"`
}

// LatestRelease - najwyższy numer ID
// ------------------------------------------------------------------------------------------------
type LatestRelease struct {
	ID int `xml:"LatestReleaseId"`
}

// Release - wydanie produkcji na csdb
// ================================================================================================
type Release struct {
	ReleaseID         int
	ReleaseYear       int
	ReleaseMonth      int
	ReleaseDay        int
	ReleaseName       string
	ReleaseType       string
	ReleaseScreenShot string
	ReleasedAt        string
	SIDPath           string
	Rating            float32
	ReleasedBy        []string
	Credits           []string
	DownloadLinks     []string
	// SrcCached         bool
	// WAVCached         bool
	// SrcExt            string
	// Disabled          bool
	// UsedSIDs          []UsedSID
}

// releases - glówna i globalna tablica z aktualnymi produkcjami
// ================================================================================================
var releases []Release

// Config - info o ostatniej ściągniętej produkcji
// ================================================================================================
type Config struct {
	DownloadDirectory string
	NoCompoDirectory  string
	LastID            int
	// HistoryYear       int
	// HistoryMonth      int
}

// releases - glówna i globalna tablica z aktualnymi produkcjami
// ================================================================================================
var config Config

// [][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][]
// ------------------------------------------------------------------------------------------------
// Funkcje
// ------------------------------------------------------------------------------------------------
// [][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][]

// ErrCheck - obsługa błedów
// ================================================================================================
func ErrCheck(errNr error) bool {
	if errNr != nil {
		fmt.Println(errNr)
		return false
	}
	return true
}

// ErrCheck2 - obsługa błedów bez komunikatu
// ================================================================================================
func ErrCheck2(errNr error) bool {
	if errNr != nil {
		return false
	}
	return true
}

// // ReadDb - Odczyt bazy
// // ================================================================================================
// func ReadDb() {
// 	file, _ := ioutil.ReadFile("releases.json")
// 	_ = json.Unmarshal([]byte(file), &releases)
// }

// // WriteDb - Zapis bazy
// // ================================================================================================
// func WriteDb() {
// 	file, _ := json.MarshalIndent(releases, "", " ")
// 	_ = ioutil.WriteFile("releases.json", file, 0666)
// }

// ReadConfig - Odczyt konfiguracji
// ================================================================================================
func ReadConfig() {
	file, _ := ioutil.ReadFile("config.json")
	_ = json.Unmarshal([]byte(file), &config)
}

// WriteConfig - Zapis konfiguracji
// ================================================================================================
func WriteConfig() {
	file, _ := json.MarshalIndent(config, "", " ")
	_ = ioutil.WriteFile("config.json", file, 0666)
}

// fileExists - sprawdzenie czy plik istnieje
// ================================================================================================
func fileExists(filename string) bool {

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}

	return true
}

// // Sortowanie datami
// // ================================================================================================

// type byDate []Release

// func (s byDate) Len() int {
// 	return len(s)
// }
// func (s byDate) Swap(i, j int) {
// 	s[i], s[j] = s[j], s[i]
// }
// func (s byDate) Less(i, j int) bool {

// 	d1 := time.Date(s[i].ReleaseYear, time.Month(s[i].ReleaseMonth), s[i].ReleaseDay, 0, 0, 0, 0, time.Local)
// 	d2 := time.Date(s[j].ReleaseYear, time.Month(s[j].ReleaseMonth), s[j].ReleaseDay, 0, 0, 0, 0, time.Local)

// 	return d2.Before(d1)
// }

// // Sortowanie byID
// // ================================================================================================

// type byID []Release

// func (s byID) Len() int {
// 	return len(s)
// }
// func (s byID) Swap(i, j int) {
// 	s[i], s[j] = s[j], s[i]
// }
// func (s byID) Less(i, j int) bool {
// 	return s[i].ReleaseID > s[j].ReleaseID
// }

// // Sortowanie datami i ID
// // ================================================================================================

// type byDateAndID []Release

// func (s byDateAndID) Len() int {
// 	return len(s)
// }
// func (s byDateAndID) Swap(i, j int) {
// 	s[i], s[j] = s[j], s[i]
// }
// func (s byDateAndID) Less(i, j int) bool {

// 	d1 := time.Date(s[i].ReleaseYear, time.Month(s[i].ReleaseMonth), s[i].ReleaseDay, 0, 0, 0, 0, time.Local)
// 	d2 := time.Date(s[j].ReleaseYear, time.Month(s[j].ReleaseMonth), s[j].ReleaseDay, 0, 0, 0, 0, time.Local)
// 	id1 := s[i].ReleaseID
// 	id2 := s[j].ReleaseID

// 	return d2.Before(d1) && id1 > id2
// }

// fileExists - sprawdzenie czy plik istnieje
// ================================================================================================
func fileSize(filename string) (int64, error) {
	// Sprawdzamy rozmiar pliku
	fileStat, err := os.Stat(filename)
	if ErrCheck(err) {
		return fileStat.Size(), err
	}
	return fileStat.Size(), err
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
// ================================================================================================
func DownloadFile(path string, filename string, url string) error {

	var err error

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, 0666)
	}
	if err != nil {
		return err
	}

	filepath := path + sep + filename

	fmt.Println("Downloading new file " + url)
	fmt.Println("Writing to " + filepath)

	httpClient := http.Client{
		Timeout: time.Second * 5, // Timeout after 5 seconds
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	out.Close()

	if strings.Contains(strings.ToLower(filename), ".zip") {
		zipReader, err := zip.OpenReader(filepath)
		if ErrCheck(err) {
			defer zipReader.Close()
			for _, file := range zipReader.File {

				fmt.Println("Extracting from ZIP: " + file.Name)
				if !file.FileInfo().IsDir() {

					outputFile, err := os.OpenFile(
						path+sep+file.Name,
						os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
						file.Mode(),
					)
					if ErrCheck(err) {
						defer outputFile.Close()

						// fmt.Println("Opening: " + file.Name)
						// fmt.Println("Output: " + path + sep + file.Name)

						zippedFile, err := file.Open()
						if ErrCheck(err) {
							defer zippedFile.Close()
							fmt.Println("Writing extracted file " + path + sep + file.Name)
							_, err = io.Copy(outputFile, zippedFile)
							ErrCheck(err)
						}
					}
				}
			}
		}
	}

	// Sprawdzzamy rozmiar pliku
	// fi, err := os.Stat(filepath)
	// if err != nil {
	// 	return err
	// }
	// get the size
	// size := fi.Size()
	// fmt.Println("Downloaded the file with size of " + strconv.Itoa(int(size)) + " bytes.")

	return err
}

// makeCharsetReader - decode reader
// ================================================================================================
func makeCharsetReader(charset string, input io.Reader) (io.Reader, error) {
	return input, nil

	// if charset == "ISO-8859-1" {
	// 	// Windows-1252 is a superset of ISO-8859-1, so should do here
	// 	return charmap.Windows1252.NewDecoder().Reader(input), nil
	// }
	// return nil, fmt.Errorf("Unknown charset: %s", charset)
}

// CSDBPrepareData - Wątek odczygtujący wszystkie releasy z csdb
// ================================================================================================
func CSDBPrepareData(full bool) {

	// lastDate := time.Now().AddDate(0, -historyMaxMonths, 0)
	// lastDate := time.Date(config.HistoryYear, time.Month(config.HistoryMonth), 1, 0, 0, 0, 0, time.Local)

	// pobranie ostatniego release'u
	netClient := &http.Client{Timeout: time.Second * 10}
	resp, err := netClient.Get("https://csdb.dk/webservice/?type=release&id=0")

	if ErrCheck(err) {

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		ErrCheck(err)
		// fmt.Println(string(body))
		resp.Body.Close()

		// Przerobienie na strukturę

		var entry LatestRelease
		reader := bytes.NewReader(body)
		decoder := xml.NewDecoder(reader)
		decoder.CharsetReader = makeCharsetReader
		err = decoder.Decode(&entry)
		ErrCheck(err)

		newestCSDbID := entry.ID

		if config.LastID == 0 {
			config.LastID = newestCSDbID - 64
			fmt.Println("Running for a first time, downloading 64 last releases. Change your config.json file to adjust the number.")
		}
		lastDownloadedID := config.LastID

		fmt.Println("Checking...")
		fmt.Println("Newest ID on CSDb is " + strconv.Itoa(newestCSDbID))
		fmt.Println("Latest downloaded ID is " + strconv.Itoa(lastDownloadedID))

		// zaczynamy od ostatniego zawsze, nawet jeżeli robimy tylko update bo może ktoś update'ował dane
		checkingID := lastDownloadedID

		searching := true
		for searching {

			resp, err := netClient.Get("https://csdb.dk/webservice/?type=release&id=" + strconv.Itoa(checkingID))

			fmt.Println("ID " + strconv.Itoa(checkingID))

			if ErrCheck(err) {
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				defer resp.Body.Close()

				if ErrCheck(err) {
					resp.Body.Close()
					// fmt.Println(string(body))

					// Przerobienie na strukturę

					var entry XMLRelease
					reader := bytes.NewReader(body)
					decoder := xml.NewDecoder(reader)
					decoder.CharsetReader = makeCharsetReader
					err = decoder.Decode(&entry)
					if ErrCheck2(err) {

						// Szukamy takiego release w naszej bazie
						//

						var relTypesAllowed = [...]string{"C64 Music", "C64 Graphics", "C64 Demo", "C64 One-File Demo", "C64 Intro", "C64 4K Intro", "C64 Crack intro", "C64 Music Collection", "C64 Graphics Collection", "C64 Diskmag", "C64 Charts", "C64 Invitation", "C64 1K Intro", "C64 Fake Demo"}
						typeOK := false
						for _, relType := range relTypesAllowed {
							if relType == entry.ReleaseType {
								typeOK = true
								break
							}
						}

						// prodYear, _ := strconv.Atoi(entry.ReleaseYear)
						// prodMonth, _ := strconv.Atoi(entry.ReleaseMonth)
						// prodDay, _ := strconv.Atoi(entry.ReleaseDay)
						// prodTime := time.Date(prodYear, time.Month(prodMonth), prodDay, 0, 0, 0, 0, time.Local)

						// Sprawdzamy czy doszliśmy do max daty
						// if !prodTime.After(lastDate) {
						// 	fmt.Println("Update finished (date reached)")
						// 	searching = false
						// } else {
						// sprawdzamy typ
						if typeOK {

							// Tworzymy nowy obiekt release który dodamy do slice
							//
							var newRelease Release
							id, _ := strconv.Atoi(entry.ReleaseID)
							newRelease.ReleaseID = id
							newRelease.ReleaseName = entry.ReleaseName
							newRelease.ReleaseScreenShot = entry.ReleaseScreenShot
							newRelease.Rating = entry.Rating
							newRelease.ReleaseYear, _ = strconv.Atoi(entry.ReleaseYear)
							newRelease.ReleaseMonth, _ = strconv.Atoi(entry.ReleaseMonth)
							newRelease.ReleaseDay, _ = strconv.Atoi(entry.ReleaseDay)
							newRelease.ReleaseType = entry.ReleaseType
							newRelease.ReleasedAt = entry.XMLReleasedAt.XMLEvent.Name

							if len(entry.UsedSIDs) == 1 {
								newRelease.SIDPath = entry.UsedSIDs[0].HVSCPath
							}

							fmt.Println("Entry name: " + entry.ReleaseName)
							// fmt.Println("ID:     ", entry.ReleaseID)
							fmt.Println("Type: " + entry.ReleaseType)
							// fmt.Println("Event:  ", entry.XMLReleasedAt.XMLEvent.Name)

							for _, group := range entry.XMLReleasedBy.XMLGroup {
								fmt.Println("Released by: " + group.Name)
								newRelease.ReleasedBy = append(newRelease.ReleasedBy, group.Name)
							}
							for _, handle := range entry.XMLReleasedBy.XMLHandle {
								// fmt.Println("XMLHandle: ", handle.XMLHandle)
								newRelease.ReleasedBy = append(newRelease.ReleasedBy, handle.XMLHandle)
							}

							// Linki dościągnięcia
							// Najpierw SIDy

							for _, link := range entry.DownloadLinks {
								newLink, _ := url.PathUnescape(link.Link)
								// fmt.Println("Download link: " + newLink)
								newRelease.DownloadLinks = append(newRelease.DownloadLinks, newLink)
							}

							//
							// Dodajemy
							//
							if len(newRelease.DownloadLinks) > 0 {
								// releases = append(releases, newRelease)
								DownloadRelease(newRelease)
								config.LastID = checkingID
								WriteConfig()

								// sprawdzamy czy przerwać ściąganie
								if !full && checkingID == newestCSDbID {
									fmt.Println("Update finished")
									searching = false
								}
							}
							// }
						}
					}
				} else {
					fmt.Println("Błąd komunikacji z csdb.dk")
				}
			} else {
				fmt.Println("Błąd komunikacji z csdb.dk")
				break
			}

			if checkingID < newestCSDbID {
				checkingID++
			} else {
				searching = false
			}
		}

	} else {
		fmt.Println("Błąd komunikacji z csdb.dk")
	}

	WriteConfig()
}

// DownloadFiles - Ściągnięcie plików
// ================================================================================================
func DownloadFiles() {
	for i := len(releases) - 1; i >= 0; i-- {
		// for _, release := range releases {
		release := releases[i]
		// fmt.Println(release)
		// if release.ReleaseID > config.LastDownloadedID {
		// news
		for _, downloadLink := range release.DownloadLinks {
			filename := filepath.Base(downloadLink)
			filename = filepath.Clean(filename)
			filename = strings.ReplaceAll(filename, "...", "")
			if release.ReleasedAt == "" {
				release.ReleasedAt = config.NoCompoDirectory
			}
			dir := cacheDir + sep + release.ReleasedAt + sep + release.ReleasedBy[0] + sep + release.ReleaseName
			dir = filepath.Clean(dir)
			dir = strings.ReplaceAll(dir, "...", "")

			if !fileExists(dir + sep + filename) {
				err := DownloadFile(dir, filename, downloadLink)
				if err == nil {
					fmt.Println("New release: " + release.ReleaseName + " by " + release.ReleasedBy[0])
					// fmt.Println("File " + filename + " downloaded for ID " + strconv.Itoa(release.ReleaseID))
					// config.LastDownloadedID = release.ReleaseID
					// WriteConfig()
				}
			} else {
				// fmt.Println("File " + filename + " already exists for ID " + strconv.Itoa(release.ReleaseID))
			}
		}
		// }
	}
}

// DownloadRelease - Ściągnięcie pojedynczego release'u i zapisanie
// ================================================================================================
func DownloadRelease(release Release) {
	for _, downloadLink := range release.DownloadLinks {
		filename := filepath.Base(downloadLink)
		filename = filepath.Clean(filename)
		filename = strings.ReplaceAll(filename, "...", "")
		if release.ReleasedAt == "" {
			release.ReleasedAt = config.NoCompoDirectory
		}
		var dir string
		if len(release.ReleasedBy) > 0 {
			dir = cacheDir + sep + release.ReleasedAt + sep + release.ReleasedBy[0] + sep + release.ReleaseName
		} else if len(release.Credits) > 0 {
			dir = cacheDir + sep + release.ReleasedAt + sep + release.Credits[0] + sep + release.ReleaseName
		} else {
			dir = cacheDir + sep + release.ReleasedAt + sep + "unknown" + sep + release.ReleaseName
		}

		dir = filepath.Clean(dir)
		dir = strings.ReplaceAll(dir, "...", "")

		if !fileExists(dir + sep + filename) {
			err := DownloadFile(dir, filename, downloadLink)
			if err == nil {
				// fmt.Println("New release: " + release.ReleaseName + " by " + release.ReleasedBy[0])
				// fmt.Println("File " + filename + " downloaded for ID " + strconv.Itoa(release.ReleaseID))
				// config.LastDownloadedID = release.ReleaseID
				// WriteConfig()
			}
		} else {
			// fmt.Println("File " + filename + " already exists for ID " + strconv.Itoa(release.ReleaseID))
		}
	}
}

// [][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][]
// ------------------------------------------------------------------------------------------------
// MAIN()
// ------------------------------------------------------------------------------------------------
// [][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][][]

func main() {

	// Info powitalne
	//
	fmt.Println("==========================================")
	fmt.Println("=======          APP START        ========")
	fmt.Println("==========================================")

	if len(os.Args) > 1 {
		sFullDownload := os.Args[1]
		if sFullDownload == "full" {
			fullDownload = true
		}
	}

	sep = string(os.PathSeparator)

	//
	// Odczyt Configa
	//
	if fileExists("config.json") {
		ReadConfig()
	} else {
		config.DownloadDirectory = "csdb_news"
		config.NoCompoDirectory = "!out_of_compo"
		config.LastID = 0
		WriteConfig()
	}

	cacheDir = config.DownloadDirectory

	fmt.Println("Your download directory is " + cacheDir)

	for {
		CSDBPrepareData(fullDownload)
		fmt.Println("Sleeping for 1 minute...")
		time.Sleep(time.Minute)
	}
}
