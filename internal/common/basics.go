package common

type Customer struct {
	XID      string `json:"xid"` // xxh64 of Email+salt
	Email    string `json:"email"`
	Password string `json:"password"` // hashed password
}

type Item struct {
	XID  string `json:"xid"` // xxh64
	Name string `json:"name"`
	// Add other fields as needed
}

type Parcel struct {
	ID          string   `json:"xid"`
	TrackingNr  string   `json:"tracking_nr"`
	WeightKg    float32  `json:"weight_kg"`
	ListOfItems []string `json:"list_of_items"`
}
