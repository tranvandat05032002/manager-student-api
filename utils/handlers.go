package utils

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func ConvertDurationToTimeUTC(timeDuration time.Duration) time.Time {
	// Tạo múi giờ UTC+7 (Asia/Ho_Chi_Minh)
	fmt.Println("Time component:             ", timeDuration)
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return time.Time{}
	}
	// Thời gian hiện tại ở múi giờ UTC
	nowUTC := time.Now().UTC()

	// Thêm thời gian duration vào thời gian hiện tại
	expirationTimeUTC := nowUTC.Add(timeDuration)

	// Chuyển đổi thời gian từ UTC sang UTC+7
	expirationTimeUTCPlus7 := expirationTimeUTC.In(loc)

	return expirationTimeUTCPlus7
}

func ConvertStringToObjectId(userId string) primitive.ObjectID {
	userIdObjecId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		fmt.Println("Error converting user id to string:", err)
	}
	return userIdObjecId
}

func GetCurrentTimeInLocal(local string) (time.Time, error) {
	location, err := time.LoadLocation(local)
	if err != nil {
		fmt.Println("Error loading location:", err)
		return time.Time{}, err
	}
	return time.Now().In(location), nil
}

func ConvertStringToISODate(local string, dateString string) (time.Time, error) {
	location, err := time.LoadLocation(local)
	if err != nil {
		return time.Time{}, fmt.Errorf("error loading location: %w", err)
	}

	t, err := time.ParseInLocation("02.01.2006, 15:04:05", dateString, location)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing date string: %w", err)
	}

	return t.UTC(), nil
}
func ConvertISO8601ToDate(isoString string) (time.Time, error) {
	// Phân tích chuỗi ISO 8601 thành time.Time
	t, err := time.Parse(time.RFC3339Nano, isoString)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing ISO 8601 string: %w", err)
	}
	return t, nil
}
