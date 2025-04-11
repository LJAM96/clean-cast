package database

import (
	"ikoyhn/podcast-sponsorblock/internal/common"
	"ikoyhn/podcast-sponsorblock/internal/models"
	"os"
	"time"

	"github.com/glebarez/sqlite"
	log "github.com/labstack/gommon/log"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func SetupDatabase() {
	var err error
	// Create the database file if it doesn't exist
	if _, err := os.Stat("C:/Users/jared/Documents/code/config/sqlite.db"); os.IsNotExist(err) {
		err := os.MkdirAll("/config", os.ModePerm)
		if err != nil {
			panic(err)
		}
		f, err := os.Create("C:/Users/jared/Documents/code/config/sqlite.db")
		if err != nil {
			panic(err)
		}
		f.Close()
	}

	db, err = gorm.Open(sqlite.Open("C:/Users/jared/Documents/code/config/sqlite.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&models.EpisodePlaybackHistory{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&models.PodcastEpisode{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&models.Podcast{})
	if err != nil {
		panic(err)
	}
}

func TrackEpisodeFiles() {
	log.Info("[DB] Tracking existing episode files...")
	audioDir := "C:/Users/jared/Documents/code/config/audio"
	if _, err := os.Stat(audioDir); os.IsNotExist(err) {
		os.MkdirAll(audioDir, 0755)
	}
	if _, err := os.Stat("C:/Users/jared/Documents/code/config"); os.IsNotExist(err) {
		os.MkdirAll("C:/Users/jared/Documents/code/config", 0755)
	}
	files, err := os.ReadDir("C:/Users/jared/Documents/code/config/audio/")
	if err != nil {
		log.Fatal(err)
	}

	dbFiles := make([]string, 0)
	db.Model(&models.EpisodePlaybackHistory{}).Pluck("YoutubeVideoId", &dbFiles)

	missingFiles := make([]string, 0)
	nonExistentDbFiles := make([]string, 0)
	for _, file := range files {
		filename := file.Name()
		if !common.IsValidFilename(filename) {
			continue
		}
		found := false
		for _, dbFile := range dbFiles {
			if dbFile == filename[:len(filename)-4] {
				found = true
				break
			}
		}
		if !found {
			missingFiles = append(missingFiles, filename)
		}
	}

	for _, dbFile := range dbFiles {
		found := false
		for _, file := range files {
			if dbFile == file.Name()[:len(file.Name())-4] {
				found = true
				break
			}
		}
		if !found {
			nonExistentDbFiles = append(nonExistentDbFiles, dbFile)
		}
	}

	for _, filename := range missingFiles {
		id := filename[:len(filename)-4]
		if !common.IsValidID(id) {
			continue
		}
		db.Create(&models.EpisodePlaybackHistory{YoutubeVideoId: id, LastAccessDate: time.Now().Unix(), TotalTimeSkipped: 0})
	}

	for _, dbFile := range nonExistentDbFiles {
		if !common.IsValidID(dbFile) {
			continue
		}
		db.Where("youtube_video_id = ?", dbFile).Delete(&models.EpisodePlaybackHistory{})
		log.Info("[DB] Deleted non-existent episode playback history... " + dbFile)
	}
}
