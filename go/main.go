package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func main() {
	ct := time.Now()
	c := http.DefaultClient
	init := hr(c, 0)
	var wg sync.WaitGroup
	for i := 1; i < init.TotalPages; i++ {
		wg.Add(1)
		println("calling ", i)
		go wgRun(&wg, c, i)
	}
	wg.Wait()
	fmt.Println(time.Since(ct))
	fmt.Println("Done!!!")
}
func wgRun(wg *sync.WaitGroup, c *http.Client, p int) {
	defer wg.Done()
	var tgw sync.WaitGroup
	tgw.Add(1)
	go hrwg(c, p, &tgw)
	tgw.Wait()
}

type AuctionData struct {
	Success       bool      `json:"success" bson:"success"`
	Page          int       `json:"page" bson:"page"`
	TotalPages    int       `json:"totalPages" bson:"totalPages"`
	TotalAuctions int       `json:"totalAuctions" bson:"totalAuctions"`
	LastUpdated   int64     `json:"lastUpdated" bson:"lastUpdated"`
	Auctions      []Auction `json:"auctions" bson:"auctions"`
}
type Auction struct {
	Uuid             string   `json:"uuid" bson:"uuid"`
	Auctioneer       string   `json:"auctioneer" bson:"auctioneer"`
	ProfileId        string   `json:"profile_id" bson:"profileId"`
	Coop             []string `json:"coop" bson:"coop"`
	CoopUser         []string `json:"coopUser" bson:"coopUser"`
	Start            int64    `json:"start" bson:"start"`
	End              int64    `json:"end" bson:"end"`
	ItemName         string   `json:"item_name" bson:"itemName"`
	ItemLore         string   `json:"item_lore" bson:"itemLore"`
	Extra            string   `json:"extra" bson:"extra"`
	Category         string   `json:"category" bson:"category"`
	Tier             string   `json:"tier" bson:"tier"`
	StartingBid      int64    `json:"starting_bid" bson:"startingBid"`
	Claimed          bool     `json:"claimed" bson:"claimed"`
	HighestBidAmount int64    `json:"highest_bid_amount" bson:"highestBidAmount"`
	Bin              bool     `json:"bin,omitempty" bson:"bin,omitempty"`
	LowestPrice      int64    `bson:"lowestPrice"`
	HighestPrice     int64    `bson:"highestPrice"`
	Reforge          string   `bson:"reforge"`
	Recombobulated   bool     `bson:"recombobulated"`
	Dungeoned        bool     `bson:"dungeoned"`
	DungeonedLvl     int64    `json:"dungeoned_lvl" bson:"dungeonedLvl"`
	LimitedUsage     bool     `json:"limited_usage" bson:"limitedUsage"`
}

func hr(c *http.Client, p int) AuctionData {
	println("getting ", p)
	req, _ := http.NewRequest(http.MethodGet, "https://api.hypixel.net/skyblock/auctions", nil)
	req.URL.RawQuery = "page=" + strconv.Itoa(p)

	res, err := c.Do(req)
	if err != nil {
		println("Error ", err.Error())
		return AuctionData{}
	}
	defer func(Body io.ReadCloser) {

		err = Body.Close()
		if err != nil {
			log.Panicf("unable to close body: %v\n", err)
			println("Error ", err.Error())
		}

	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		println("Error ", err.Error())
		return AuctionData{}
	}
	var data AuctionData
	err = json.Unmarshal(body, &data)
	if err != nil {
		println("Error ", err.Error())
		return AuctionData{}
	}
	println("done")
	return data
}
func hrwg(c *http.Client, p int, wg *sync.WaitGroup) AuctionData {
	defer wg.Done()
	println("getting ", p)
	req, _ := http.NewRequest(http.MethodGet, "https://api.hypixel.net/skyblock/auctions", nil)
	req.URL.RawQuery = "page=" + strconv.Itoa(p)
	res, _ := c.Do(req)


	body, _ := io.ReadAll(res.Body)

	var data AuctionData
	json.Unmarshal(body, &data)

	println("done")
	return data
}
