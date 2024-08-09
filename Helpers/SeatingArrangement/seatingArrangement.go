package seating

import (
	models "Redbus_backend/Models"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ArrangingSeats(arrangementType string) []models.Seat {
	var theArrangement []models.Seat

	if arrangementType == "45SE" {
		groupNo := 1
		for col := 0; col < 10; col++ {
			for row := 0; row < 5; row++ {
				seat := models.Seat{
					SeatID:           primitive.NewObjectID(),
					SeatType:         "SE",
					OnlyFemale:       false,
					OnlyMale:         false,
					PassengerGender:  "",
					SeatAvailibility: true,
					PassengerName:    "",
					PassengerAge:     0,
					Row:              row,
					Column:           col,
					Group:            "G" + strconv.Itoa(groupNo),
					IsUpperDeck:      false,
					Class:            "C",
					IsSelected:       false,
				}
				if row == 2 {
					groupNo++
					seat.Group = ""
					seat.SeatType = "N"
				}
				if row == 0 || row == 4 {
					if col < 7 {
						seat.Class = "A"
					} else {
						seat.Class = "B"
					}
				} else {
					if col < 3 {
						seat.Class = "A"
					} else if col < 7 {
						seat.Class = "B"
					}
				}
				theArrangement = append(theArrangement, seat)
			}
			groupNo++
		}
		for row := 0; row < 5; row++ {
			seat := models.Seat{
				SeatID:           primitive.NewObjectID(),
				SeatType:         "SE",
				OnlyFemale:       false,
				OnlyMale:         false,
				PassengerGender:  "",
				SeatAvailibility: true,
				PassengerName:    "",
				PassengerAge:     0,
				Row:              row,
				Column:           10,
				Group:            "G" + strconv.Itoa(groupNo),
				IsUpperDeck:      false,
				Class:            "C",
				IsSelected:       false,
			}
			if row == 0 || row == 4 {
				seat.Class = "B"
			}
			theArrangement = append(theArrangement, seat)
		}
	} else if arrangementType == "42SL" {
		groupNo := 1
		for deckCount := 0; deckCount < 2; deckCount++ {
			for col := 0; col < 7; col++ {
				for row := 0; row < 4; row++ {
					seat := models.Seat{
						SeatID:           primitive.NewObjectID(),
						SeatType:         "SL",
						OnlyFemale:       false,
						OnlyMale:         false,
						PassengerGender:  "",
						SeatAvailibility: true,
						PassengerName:    "",
						PassengerAge:     0,
						Row:              row,
						Column:           col,
						Group:            "G" + strconv.Itoa(groupNo),
						IsUpperDeck:      false,
						Class:            "C",
						IsSelected:       false,
					}
					if deckCount != 0 {
						seat.IsUpperDeck = true
					}
					if row == 2 {
						groupNo++
						seat.Group = ""
						seat.SeatType = "N"
					} else if row == 3 || col == 0 {
						seat.Class = "A"
					} else if col < 4 {
						seat.Class = "B"
					}
					theArrangement = append(theArrangement, seat)
				}
				groupNo++
			}
		}
	} else if arrangementType == "30SL" {
		groupNo := 1
		for deckCount := 0; deckCount < 2; deckCount++ {
			for col := 0; col < 5; col++ {
				for row := 0; row < 4; row++ {
					seat := models.Seat{
						SeatID:           primitive.NewObjectID(),
						SeatType:         "SL",
						OnlyFemale:       false,
						OnlyMale:         false,
						PassengerGender:  "",
						SeatAvailibility: true,
						PassengerName:    "",
						PassengerAge:     0,
						Row:              row,
						Column:           col,
						Group:            "G" + strconv.Itoa(groupNo),
						IsUpperDeck:      false,
						Class:            "C",
						IsSelected:       false,
					}
					if deckCount != 0 {
						seat.IsUpperDeck = true
					}
					if row == 2 {
						groupNo++
						seat.Group = ""
						seat.SeatType = "N"
					} else if row == 3 || col == 0 {
						seat.Class = "A"
					} else if col < 3 {
						seat.Class = "B"
					}
					theArrangement = append(theArrangement, seat)
				}
				groupNo++
			}
		}
	} else if arrangementType == "38SL" {
		groupNo := 1
		for deckCount := 0; deckCount < 2; deckCount++ {
			for col := 0; col < 5; col++ {
				for row := 0; row < 4; row++ {
					seat := models.Seat{
						SeatID:           primitive.NewObjectID(),
						SeatType:         "SL",
						OnlyFemale:       false,
						OnlyMale:         false,
						PassengerGender:  "",
						SeatAvailibility: true,
						PassengerName:    "",
						PassengerAge:     0,
						Row:              row,
						Column:           col,
						Group:            "G" + strconv.Itoa(groupNo),
						IsUpperDeck:      false,
						Class:            "C",
						IsSelected:       false,
					}
					if deckCount != 0 {
						seat.IsUpperDeck = true
					}
					if row == 2 {
						groupNo++
						seat.Group = ""
						seat.SeatType = "N"
					} else if row == 3 || col < 2 {
						seat.Class = "A"
					} else if col < 4 {
						seat.Class = "B"
					}
					theArrangement = append(theArrangement, seat)
				}
				groupNo++
			}
			for row := 0; row < 4; row++ {
				seat := models.Seat{
					SeatID:           primitive.NewObjectID(),
					SeatType:         "SL",
					OnlyFemale:       false,
					OnlyMale:         false,
					PassengerGender:  "",
					SeatAvailibility: true,
					PassengerName:    "",
					PassengerAge:     0,
					Row:              row,
					Column:           5,
					Group:            "G" + strconv.Itoa(groupNo),
					IsUpperDeck:      false,
					Class:            "C",
					IsSelected:       false,
				}
				if deckCount != 0 {
					seat.IsUpperDeck = true
				}
				theArrangement = append(theArrangement, seat)
			}
			groupNo++
		}
	} else if arrangementType == "8SE33SL" {
		groupNo := 1
		for col := 0; col < 4; col++ {
			for row := 0; row < 4; row++ {
				seat := models.Seat{
					SeatID:           primitive.NewObjectID(),
					SeatType:         "SE",
					OnlyFemale:       false,
					OnlyMale:         false,
					PassengerGender:  "",
					SeatAvailibility: true,
					PassengerName:    "",
					PassengerAge:     0,
					Row:              row,
					Column:           col,
					Group:            "G" + strconv.Itoa(groupNo),
					IsUpperDeck:      false,
					Class:            "B",
					IsSelected:       false,
				}
				if row == 2 {
					seat.SeatType = "N"
					seat.Group = ""
					groupNo++
				} else if row == 0 {
					seat.Class = "A"
				} else if row == 3 {
					if col%2 == 0 {
						seat.SeatType = "SL"
						seat.Class = "A"
						theArrangement = append(theArrangement, seat)
						groupNo++
					}
					continue
				}
				theArrangement = append(theArrangement, seat)
			}
			groupNo++
		}
		for col := 4; col < 10; col += 2 {
			for row := 0; row < 4; row++ {
				seat := models.Seat{
					SeatID:           primitive.NewObjectID(),
					SeatType:         "SL",
					OnlyFemale:       false,
					OnlyMale:         false,
					PassengerGender:  "",
					SeatAvailibility: true,
					PassengerName:    "",
					PassengerAge:     0,
					Row:              row,
					Column:           col,
					Group:            "G" + strconv.Itoa(groupNo),
					IsUpperDeck:      false,
					Class:            "C",
					IsSelected:       false,
				}
				if row == 2 {
					groupNo++
					seat.Group = ""
					seat.SeatType = "N"
				}
				theArrangement = append(theArrangement, seat)
			}
			groupNo++
		}
		for row := 0; row < 4; row++ {
			seat := models.Seat{
				SeatID:           primitive.NewObjectID(),
				SeatType:         "SL",
				OnlyFemale:       false,
				OnlyMale:         false,
				PassengerGender:  "",
				SeatAvailibility: true,
				PassengerName:    "",
				PassengerAge:     0,
				Row:              row,
				Column:           12,
				Group:            "G" + strconv.Itoa(groupNo),
				IsUpperDeck:      false,
				Class:            "C",
				IsSelected:       false,
			}
			theArrangement = append(theArrangement, seat)
		}
		groupNo++
		for col := 0; col < 6; col++ {
			for row := 0; row < 4; row++ {
				seat := models.Seat{
					SeatID:           primitive.NewObjectID(),
					SeatType:         "SL",
					OnlyFemale:       false,
					OnlyMale:         false,
					PassengerGender:  "",
					SeatAvailibility: true,
					PassengerName:    "",
					PassengerAge:     0,
					Row:              row,
					Column:           col,
					Group:            "G" + strconv.Itoa(groupNo),
					IsUpperDeck:      true,
					Class:            "C",
					IsSelected:       false,
				}
				if row == 2 {
					groupNo++
					seat.Group = ""
					seat.SeatType = "N"
				} else if row == 3 || col == 1 {
					seat.Class = "A"
				} else if col < 3 {
					seat.Class = "B"
				}
				theArrangement = append(theArrangement, seat)
			}
			groupNo++
		}
	}

	var toBeReturnedSlice []models.Seat

	for _, item := range theArrangement {
		if item.Group != "" {
			toBeReturnedSlice = append(toBeReturnedSlice, item)
		}
	}

	return toBeReturnedSlice
}
