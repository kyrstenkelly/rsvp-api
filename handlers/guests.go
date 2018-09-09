package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

// Guest type
type Guest struct {
	ID        string   `json:"id,omitempty"`
	Firstname string   `json:"firstname,omitempty"`
	Lastname  string   `json:"lastname,omitempty"`
	Email     string   `json:"email,omitempty"`
	Address   *Address `json:"address,omitempty"`
}

// Address type
type Address struct {
	Line1 string `json:"line1,omitempty"`
	Line2 string `json:"line2,omitempty"`
	Zip   string `json:"zip,omitempty"`
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

var guests []Guest

// GetGuests gets guests
func GetGuests(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(guests)
}

// GetGuest gets a guest
func GetGuest(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var guest Guest
	for _, item := range guests {
		if item.ID == params["id"] {
			guest = item
		}
	}
	json.NewEncoder(w).Encode(guest)
}

// CreateGuest creates a guest
func CreateGuest(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var guest Guest
	_ = json.NewDecoder(r.Body).Decode(&guest)
	guest.ID = params["id"]
	guests = append(guests, guest)
	json.NewEncoder(w).Encode(guests)
}

// DeleteGuest deletes a guest
func DeleteGuest(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range guests {
		if item.ID == params["id"] {
			guests = append(guests[:index], guests[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(guests)
}

func main() {
	var ourAddress = &Address{
		Line1: "7009 Almeda Rd #633",
		Zip:   "77054",
		City:  "Houston",
		State: "TX",
	}
	guests = append(guests, Guest{ID: "1", Firstname: "Kyrsten", Lastname: "Kelly", Email: "kyrsten.kelly@gmail.com", Address: ourAddress})
	guests = append(guests, Guest{ID: "2", Firstname: "James", Lastname: "Custer", Email: "jamescuster0121@gmail.com", Address: ourAddress})
}
