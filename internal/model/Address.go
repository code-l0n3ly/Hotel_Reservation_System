package model

type Address struct {
	AddressID        string `json:"addressID"`
	Country          string `json:"Country"`
	City             string `json:"city"`
	State            string `json:"state"`
	Street           string `json:"street"`
	PostalCode       string `json:"PostalCode"`
	AdditionalNumber string `json:"additionalNumber"`
	MapLocation      string `json:"mapLocation"`
	Latitude         string `json:"Latitude"`
	Longitude        string `json:"Longitude"`
}
