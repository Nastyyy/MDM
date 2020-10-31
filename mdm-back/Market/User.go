package market

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	satori "github.com/satori/go.uuid"
)

var LIMIT float32 = 1000000000.0

type User struct {
	Name      string       `json:"username,omitempty"`
	Timestamp *time.Time   `json:"timestamp,omitempty"`
	UUID      *satori.UUID `json:"uuid,omitempty"`
	Portfolio *[]Holding
	balance   float32
	Conn      *websocket.Conn `json:"-"`
	Cookie    *http.Cookie    `json:"-"`
}

type Holding struct {
	Asset  *Stock
	Volume int
}

func (user *User) Withdraw(amount float32) error {
	if user.balance-amount >= 0.0 {
		user.balance -= amount
		return nil
	}
	return fmt.Errorf("Erorr withdrawing from %v: %v", user, amount)
}

func (user *User) Deposit(amount float32) error {
	if user.balance+amount <= LIMIT {
		user.balance += amount
		return nil
	}

	return fmt.Errorf("Error depositing into%v: %v", user, amount)
}

func (user *User) UpdateHolding(stock *Stock, volume int) error {
	for _, holding := range *user.Portfolio {
		if holding.Asset == stock && holding.Volume-volume >= 0 {
			holding.Volume += volume
			holding.Asset.Volume -= volume
			return nil
		}
	}

	// Holding not found, so need to add
	holding := Holding{Asset: stock, Volume: volume}
	newPortfolio := append(*user.Portfolio, holding)
	user.Portfolio = &newPortfolio
	return nil
}

func (user *User) SetUUID(newUUID string) error {
	gen, err := satori.FromString(newUUID)
	if err != nil {
		return err
	}

	user.UUID = &gen
	return nil
}

func (user *User) GetWorth() float32 {
	return user.balance
}

func (user *User) GetPortfolioVolume(ticker string) int {
	for _, holding := range *user.Portfolio {
		if holding.Asset.Ticker == ticker {
			return holding.Volume
		}
	}
	return 0
}

func (user *User) GetBalance() float32 {
	return user.balance
}

func (user User) String() string {
	return fmt.Sprintf("%s -- %v | Worth:$%v | Balance:$%v |", user.Name, user.UUID, user.GetWorth(), user.GetBalance())
}

// NewUser creats and initializes a new user instance
func NewUser(name string, uuid string) (*User, error) {
	currentTime := time.Now()
	playerUUID, err := satori.FromString(uuid)
	if err != nil {
		return nil, err
	}

	player := User{Name: name, Timestamp: &currentTime, UUID: &playerUUID}
	return &player, nil
}
