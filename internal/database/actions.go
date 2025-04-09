package database

import (
	"errors"
	"ikoyhn/podcast-sponsorblock/internal/models"
	"os"
	"time"

	log "github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

func UpdateEpisodePlaybackHistory(youtubeVideoId string, totalTimeSkipped float64) {
	log.Info("[DB] Updating episode playback history...")
	db.Model(&models.EpisodePlaybackHistory{}).
		Where("youtube_video_id = ?", youtubeVideoId).
		FirstOrCreate(&models.EpisodePlaybackHistory{
			YoutubeVideoId:   youtubeVideoId,
			LastAccessDate:   time.Now().Unix(),
			TotalTimeSkipped: totalTimeSkipped,
		})
}

func DeleteEpisodePlaybackHistory(youtubeVideoId string) {
	db.Where("youtube_video_id = ?", youtubeVideoId).Delete(&models.EpisodePlaybackHistory{})
}

func DeletePodcastCronJob() {
	oneWeekAgo := time.Now().Add(-7 * 24 * time.Hour).Unix()

	var histories []models.EpisodePlaybackHistory
	db.Where("last_access_date < ?", oneWeekAgo).Find(&histories)

	for _, history := range histories {
		os.Remove("/config/audio/" + history.YoutubeVideoId + ".m4a")
		db.Delete(&history)
		log.Info("[DB] Deleted old episode playback history... " + history.YoutubeVideoId)
	}
}

func GetEpisodePlaybackHistory(youtubeVideoId string) *models.EpisodePlaybackHistory {
	var history models.EpisodePlaybackHistory
	db.Where("youtube_video_id = ?", youtubeVideoId).First(&history)
	return &history
}

func EpisodeExists(youtubeVideoId string, episodeType string) (bool, error) {
	var episode models.PodcastEpisode
	err := db.Where("youtube_video_id = ? AND type = ?", youtubeVideoId, episodeType).First(&episode).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func GetLatestEpisode(podcastId string) (*models.PodcastEpisode, error) {
	var episode models.PodcastEpisode
	err := db.Where("podcast_id = ?", podcastId).Order("published_date DESC").First(&episode).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &episode, nil
}

func GetAllPlaylistVideosByPlaylistId(playlistId string) []models.PodcastEpisode {
	var episodes []models.PodcastEpisode
	db.Where("playlist_id != ?", playlistId).Find(&episodes)
	return episodes
}
func GetPodcastEpisodesByPodcastId(podcastId string) ([]models.PodcastEpisode, error) {
	var episodes []models.PodcastEpisode
	err := db.Where("podcast_id = ?", podcastId).Find(&episodes).Error
	if err != nil {
		return nil, err
	}
	return episodes, nil
}

func SavePlaylistEpisodes(playlistEpisodes []models.PodcastEpisode) {
	db.CreateInBatches(playlistEpisodes, 100)
}

func GetPodcast(id string) *models.Podcast {
	var podcastDb models.Podcast
	err := db.Where("id = ?", id).Find(&podcastDb).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
	}
	if podcastDb.Id == "" {
		return nil
	}
	return &podcastDb
}

func SavePodcast(podcast *models.Podcast) {
	db.Create(&podcast)
}

func GetAllPodcasts() []models.Podcast {
	var podcasts []models.Podcast
	db.Find(&podcasts)
	return podcasts
}
