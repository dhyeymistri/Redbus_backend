package filtersearchedbuses

import (
	getbusdetail "Redbus_backend/Helpers/GetBusDetail"
	models "Redbus_backend/Models"
)

func Filtering(filters map[string]interface{}, searchedBusResult []models.Booking) []models.Booking {
	var returnedBookingArray []models.Booking
	for key, value := range filters {
		switch key {
		case "ac":
			if value.(bool) {
				for idx := range searchedBusResult {
					bus := getbusdetail.GetBusDetail(searchedBusResult[idx].BusID)
					if bus.IsAcAvailable {
						returnedBookingArray = append(returnedBookingArray, searchedBusResult[idx])
						// searchedBusResult = append(searchedBusResult[:idx], searchedBusResult[idx+1:]...)
					}
				}
			}
		case "nonac":
			if value.(bool) {
				for idx := range searchedBusResult {
					bus := getbusdetail.GetBusDetail(searchedBusResult[idx].BusID)
					if !bus.IsAcAvailable {
						returnedBookingArray = append(returnedBookingArray, searchedBusResult[idx])
					}
				}
			}
		case "seater":
			if value.(bool) {
				for idx := range searchedBusResult {
					bus := getbusdetail.GetBusDetail(searchedBusResult[idx].BusID)
					if bus.BusType == "45SE" || bus.BusType != "8SE33SL" {
						returnedBookingArray = append(returnedBookingArray, searchedBusResult[idx])
					}
				}
			}
		case "sleeper":
			if value.(bool) {
				for idx := range searchedBusResult {
					bus := getbusdetail.GetBusDetail(searchedBusResult[idx].BusID)
					if bus.BusType != "45SE" && bus.BusType != "8SE33SL" {
						returnedBookingArray = append(returnedBookingArray, searchedBusResult[idx])
					}
				}
			}
		case "arr-b-6":
			if value.(bool) {
				for idx := range searchedBusResult {
					if searchedBusResult[idx].TravelEndTime <= "6:00" {
						returnedBookingArray = append(returnedBookingArray, searchedBusResult[idx])
					}
				}
			}
		case "arr-a-18":
			if value.(bool) {
				for idx := range searchedBusResult {
					if searchedBusResult[idx].TravelEndTime >= "18:00" {
						returnedBookingArray = append(returnedBookingArray, searchedBusResult[idx])
					}
				}
			}
		case "arr-a-6-b-12":
			if value.(bool) {
				for idx := range searchedBusResult {
					if searchedBusResult[idx].TravelEndTime >= "6:00" && searchedBusResult[idx].TravelEndTime <= "12:00" {
						returnedBookingArray = append(returnedBookingArray, searchedBusResult[idx])
					}
				}
			}
		case "arr-a-12-b-18":
			if value.(bool) {
				for idx := range searchedBusResult {
					if searchedBusResult[idx].TravelEndTime >= "12:00" && searchedBusResult[idx].TravelEndTime <= "18:00" {
						returnedBookingArray = append(returnedBookingArray, searchedBusResult[idx])
					}
				}
			}
		case "dep-b-6":
			if value.(bool) {
				for idx := range searchedBusResult {
					if searchedBusResult[idx].TravelStartTime <= "6:00" {
						returnedBookingArray = append(returnedBookingArray, searchedBusResult[idx])
					}
				}
			}
		case "dep-a-18":
			if value.(bool) {
				for idx := range searchedBusResult {
					if searchedBusResult[idx].TravelStartTime >= "18:00" {
						returnedBookingArray = append(returnedBookingArray, searchedBusResult[idx])
					}
				}
			}
		case "dep-a-6-b-12":
			if value.(bool) {
				for idx := range searchedBusResult {
					if searchedBusResult[idx].TravelStartTime >= "6:00" && searchedBusResult[idx].TravelStartTime <= "12:00" {
						returnedBookingArray = append(returnedBookingArray, searchedBusResult[idx])
					}
				}
			}
		case "dep-a-12-b-18":
			if value.(bool) {
				for idx := range searchedBusResult {
					if searchedBusResult[idx].TravelStartTime >= "12:00" && searchedBusResult[idx].TravelStartTime <= "18:00" {
						returnedBookingArray = append(returnedBookingArray, searchedBusResult[idx])
					}
				}
			}
		case "wifi":
			if value.(bool) {
				for idx := range searchedBusResult {
					flag := false
					bus := getbusdetail.GetBusDetail(searchedBusResult[idx].BusID)
					for _, item := range bus.Amenities {
						if item == "wifi" {
							flag = true
						}
					}
					if flag {
						returnedBookingArray = append(returnedBookingArray, searchedBusResult[idx])
					}
				}
			}
		case "toilet":
			if value.(bool) {
				for idx := range searchedBusResult {
					flag := false
					bus := getbusdetail.GetBusDetail(searchedBusResult[idx].BusID)
					for _, item := range bus.Amenities {
						if item == "toilet" {
							flag = true
						}
					}
					if flag {
						returnedBookingArray = append(returnedBookingArray, searchedBusResult[idx])
					}
				}
			}
		case "waterBottle":
			if value.(bool) {
				for idx := range searchedBusResult {
					flag := false
					bus := getbusdetail.GetBusDetail(searchedBusResult[idx].BusID)
					for _, item := range bus.Amenities {
						if item == "waterBottle" {
							flag = true
						}
					}
					if flag {
						returnedBookingArray = append(returnedBookingArray, searchedBusResult[idx])
					}
				}
			}
		case "chargingPoint":
			if value.(bool) {
				for idx := range searchedBusResult {
					flag := false
					bus := getbusdetail.GetBusDetail(searchedBusResult[idx].BusID)
					for _, item := range bus.Amenities {
						if item == "chargingPoint" {
							flag = true
						}
					}
					if flag {
						returnedBookingArray = append(returnedBookingArray, searchedBusResult[idx])
					}
				}
			}
		case "blankets":
			if value.(bool) {
				for idx := range searchedBusResult {
					flag := false
					bus := getbusdetail.GetBusDetail(searchedBusResult[idx].BusID)
					for _, item := range bus.Amenities {
						if item == "blankets" {
							flag = true
						}
					}
					if flag {
						returnedBookingArray = append(returnedBookingArray, searchedBusResult[idx])
					}
				}
			}
		case "liveTracking":
			if value.(bool) {
				for idx := range searchedBusResult {
					bus := getbusdetail.GetBusDetail(searchedBusResult[idx].BusID)
					if bus.LiveTracking {
						returnedBookingArray = append(returnedBookingArray, searchedBusResult[idx])
					}
				}
			}
		}
	}
	return returnedBookingArray
}
