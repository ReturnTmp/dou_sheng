package crud

import (
	"fmt"

	"gitee.com/Whitroom/imitate-tiktok/sql"
	"gitee.com/Whitroom/imitate-tiktok/sql/models"
)

func UserLikeVideo(UserID uint, VideoID uint) error {
	var user *models.User
	var video *models.Video

	sql.DB.First(&user, UserID)
	sql.DB.First(&video, VideoID)

	if user == nil || video == nil {
		return fmt.Errorf("找不到用户或视频")
	}

	sql.DB.Model(&user).Association("FavoriteVideos").Append(video)
	sql.DB.Commit()

	return nil
}

func UserDislikeVideo(UserID uint, VideoID uint) error {
	var user *models.User
	var video *models.Video

	sql.DB.First(&user, UserID)
	sql.DB.First(&video, VideoID)

	if user == nil || video == nil {
		return fmt.Errorf("找不到用户或视频")
	}

	err := sql.DB.Model(&user).Association("FavoriteVideos").Delete(video)
	if err != nil {
		return fmt.Errorf("找不到点赞的视频")
	}
	sql.DB.Commit()

	return nil
}

func GetUserLikeVideosByUserID(UserID uint) []models.Video {
	var user *models.User
	sql.DB.Preload("FavoriteVideos").Find(&user, UserID)
	return user.FavoriteVideos
}

func GetVideoLikesCount(VideoID uint) int64 {
	var count int64
	sql.DB.Raw("select count(user_id) from user_favorite_videos where video_id = ?", VideoID).Scan(&count)
	return count
}

func IsUserFavoriteVideo(userID, videoID uint) bool {
	var video *models.Video
	if userID == 0 {
		return false
	}
	sql.DB.Raw("select * from videos where id in "+
		"(select video_id from user_favorite_videos where user_id = ? and video_id = ?)",
		userID, videoID).Scan(&video)
	return video != nil
}
