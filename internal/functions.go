package internal

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	plexFunctions "github.com/jrudio/go-plex-client"
	log "github.com/sirupsen/logrus"
	"golift.io/starr/radarr"
)

func CheckIfError (err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func CheckEnvVars() () {
	envVars := []string {
		"PlexServerUrl",
		"PlexServerKey",
		"RadarrServerUrl",
		"RadarrServerKey",
	}

	for _, envVar := range envVars {
		_, status := os.LookupEnv(envVar)
		if !status {
			log.Fatal(envVar + " IS UNSET!")
		}
	}
}

func GenerateTitleList (movie *radarr.Movie) ([]string) {
	titleList := []string {
		movie.Title,
	}
	
	for	_, altTitle := range movie.AlternateTitles {
		titleList = append(titleList, altTitle.Title)
	}
	log.Info("Found " + strconv.Itoa(len(titleList)) + " Titles For " + movie.Title)
    log.WithFields(
        log.Fields{
			"TitleList:": titleList,
        },
    ).Debug("Generated Titlelist")
	return titleList
}

func VerifyPlexMatch (potentialMatch string, movie string)(bool) {
	match := strings.Contains(potentialMatch, movie)
	return match
}

func SearchPlexForMatches (plex plexFunctions.Plex, titleList []string, movie radarr.Movie) (plexFunctions.Metadata, error) {
	var (
		results plexFunctions.SearchResults
		err	error
		match bool
		plexMovie plexFunctions.Metadata
	)
	
	for _, title := range titleList {
		results, err = plex.Search(title)
		match, plexMovie = CheckPlexMatches(results, movie)
		if (match) {
			log.Info("Match Found For: " + movie.Title)
			log.Debug("Matched:\n" + "Plex: " + plexMovie.Title + " -> " + "Radarr: " + movie.Title)
			break
		}
	}
	return plexMovie, err
}

func CheckPlexMatches(results plexFunctions.SearchResults, movie radarr.Movie) (bool, plexFunctions.Metadata) {
	var (
		match bool
		potentialMatch plexFunctions.Metadata
	)
	for _, potentialMatch = range results.MediaContainer.MediaContainer.Metadata {
		log.Debug("Checking Match: " + potentialMatch.Media[0].Part[0].File + " -> " + movie.MovieFile.RelativePath)
		match = VerifyPlexMatch(potentialMatch.Media[0].Part[0].File, movie.MovieFile.RelativePath)
		if (match) {
			break
		}
		}	
	return match, potentialMatch
}

func UpdateEdition(plexMovie plexFunctions.Metadata, movie radarr.Movie, PlexServerUrl string, PlexServerKey string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("PUT", PlexServerUrl + "/library/sections/" + string(plexMovie.LibrarySectionID) + "/all", nil)
	CheckIfError(err)
	q := req.URL.Query()
	q.Add("type", "1")
	q.Add("id", string(plexMovie.RatingKey))
	q.Add("includeExternalMedia", "1")
	q.Add("editionTitle.value", movie.MovieFile.Edition)
	q.Add("X-Plex-Token", PlexServerKey)
	q.Add("editionTitle.locked", "0")
	req.URL.RawQuery = q.Encode()
	httpResponse, err := client.Do(req)

	log.Debug("Query: ", strings.Replace(string(req.URL.RawQuery) ,PlexServerKey,"<REDACTED>", -1))
	log.Debug(httpResponse.StatusCode)
	return httpResponse, err
}