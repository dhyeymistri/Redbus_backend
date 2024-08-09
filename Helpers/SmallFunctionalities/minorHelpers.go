package minorhelpers

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FindCommonElements(arr1, arr2 []primitive.ObjectID) []primitive.ObjectID {
	elementSet := make(map[primitive.ObjectID]struct{})
	for _, elem := range arr1 {
		elementSet[elem] = struct{}{}
	}
	var commonElements []primitive.ObjectID

	for _, elem := range arr2 {
		if _, exists := elementSet[elem]; exists {
			commonElements = append(commonElements, elem)
			// Remove the element from the set to avoid duplicates in the result
			delete(elementSet, elem)
		}
	}

	return commonElements
}

func HasDateChanged(startTimeStr, endTimeStr string) bool {
	// Define the time layout (format)
	layout := "15:04"

	// Parse the start and end times
	startTime, err := time.Parse(layout, startTimeStr)
	if err != nil {
		fmt.Println("Error parsing start time:", err)
		return false
	}

	endTime, err := time.Parse(layout, endTimeStr)
	if err != nil {
		fmt.Println("Error parsing end time:", err)
		return false
	}

	// Assume today’s date for the start and end times
	now := time.Now().Local()
	startTime = time.Date(now.Year(), now.Month(), now.Day(), startTime.Hour(), startTime.Minute(), startTime.Second(), startTime.Local().Hour(), time.UTC)
	endTime = time.Date(now.Year(), now.Month(), now.Day(), endTime.Hour(), endTime.Minute(), endTime.Second(), endTime.Local().Hour(), time.UTC)

	// If end time is earlier than start time, adjust end time to the next day
	if endTime.Before(startTime) {
		endTime = endTime.Add(24 * time.Hour)
	}

	// Compare the dates
	return endTime.YearDay() != startTime.YearDay()
}

func TimeDifference(t1, t2 string) (time.Duration, error) {
	// Define the time layout for parsing
	const layout = "15:04"

	// Parse the input times
	time1, err := time.Parse(layout, t1)
	if err != nil {
		return 0, err
	}
	time2, err := time.Parse(layout, t2)
	if err != nil {
		return 0, err
	}

	diff := time2.Sub(time1)

	// If the difference is negative, adjust by adding 24 hours
	if diff < 0 {
		diff += 24 * time.Hour
	}

	return diff, nil
}
