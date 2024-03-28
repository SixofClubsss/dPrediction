package prediction

import (
	dreams "github.com/dReam-dApps/dReams"
)

// dSports and dPredictions account data
type accountData struct {
	Predict []string `json:"predict,omitempty"`
	Sports  []string `json:"sports,omitempty"`
}

// Get account data for dSports and dPredictions
func GetAccount() interface{} {
	return accountData{
		Predict: Predict.Favorites.SCIDs,
		Sports:  Sports.Favorites.SCIDs,
	}
}

// Set stored account data to dSports and dPredictions variables
func SetAccount(ad interface{}) (err error) {
	var account accountData
	err = dreams.SetAccount(ad, &account)
	if err != nil {
		clearAccountData()
		return
	}

	Predict.Favorites.SCIDs = account.Predict
	Sports.Favorites.SCIDs = account.Sports

	return
}

// Clear existing account data for dSports and dPredictions
func clearAccountData() {
	Predict.Favorites.SCIDs = []string{}
	Sports.Favorites.SCIDs = []string{}
}

// Store dSports and dPredictions to datashard
func saveAccount() *dreams.AccountEncrypted {
	return dreams.AddAccountData(GetAccount(), "prediction")
}
