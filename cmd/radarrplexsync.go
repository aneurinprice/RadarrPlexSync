package main

import (
	"fmt"
	"os"
	radarrPlexSync "radarrPlexSync/internal"
	"strings"
	plexFunctions "github.com/jrudio/go-plex-client"
	log "github.com/sirupsen/logrus"
	"golift.io/starr"
	"golift.io/starr/radarr"
)

func main() {
	radarrPlexSync.CheckEnvVars()
	
	LogLevel, exists := os.LookupEnv("LogLevel")
	if !exists {
		LogLevel = "Warn"
	}
	log.Warn("Log Level: " + LogLevel)
	parsedLogLevel, err := log.ParseLevel(LogLevel)
	radarrPlexSync.CheckIfError(err)
	log.SetLevel(parsedLogLevel)

	log.Debug("Radarr Server: " + os.Getenv("RadarrServerUrl"))
	log.Debug("Radarr API Key (Obscured): " + strings.Repeat("#", len(os.Getenv("RadarrServerKey"))))
	log.Debug("Plex Server: " + os.Getenv("PlexServerUrl"))
	log.Debug("Plex API Key (Obscured): " + strings.Repeat("#", len(os.Getenv("PlexServerKey"))))

	radarrConnectionConnection := starr.New(os.Getenv("RadarrServerKey"), os.Getenv("RadarrServerUrl"), 0)
	radarr := radarr.New(radarrConnectionConnection)
	output, err := radarr.GetMovie(0)
	radarrPlexSync.CheckIfError(err)

	for _, movie := range output {
		if movie.HasFile && (len(movie.MovieFile.Edition) > 0 ) {
			plex, err := plexFunctions.New(os.Getenv("PlexServerUrl"), os.Getenv("PlexServerKey"))
			radarrPlexSync.CheckIfError(err)
			titleList := radarrPlexSync.GenerateTitleList(movie)
			plexMovie,err := radarrPlexSync.SearchPlexForMatches(*plex, titleList, *movie)
			radarrPlexSync.CheckIfError(err)
			_, err = radarrPlexSync.UpdateEdition(plexMovie, *movie, os.Getenv("PlexServerUrl"), os.Getenv("PlexServerKey"))
			radarrPlexSync.CheckIfError(err)
			fmt.Println(plexMovie.Title + ": " + movie.MovieFile.Edition)

		}
	}
}
