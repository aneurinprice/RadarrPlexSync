# RadarrPlexSync

## Description

This project aims to allow a user to automatically transfer radarr data to plex. The main use case I have for this is to sync the `Edition` field on my movies so that Radarr and Plex are always up to date. Currently the `Edition` field is the only supported datapoint, but adding more should be easy enough now that the basics work.

## How to use

- Set the following environmental variables:
  ```
    PlexServerUrl
    PlexServerKey
    RadarrServerUrl
    RadarrServerKey
    LogLevel #(Optional)
  ```

- `go run cmd/radarrplexsync.go`
    Depending on your log level, Your output should resemble:
    ```
        National Treasure: Collectors Edition
        Now You See Me: Extended Edition
        Robin Hood: Special Edition
        RoboCop: Directors Cut
        Rock of Ages: Extended Edition
        Seven Brides for Seven Brothers: Widescreen Edition
        Short Circuit: Special Edition
    ```

    Any movie displayed in the above output has been successfully sync'd

## Upcoming
  - Proper releases
  - Build binaries
  - Docker container
  - Attempt to suport multiple editions of the same movie