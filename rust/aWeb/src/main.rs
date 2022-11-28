use futures::StreamExt;
use serde::{Serialize, Deserialize};
use reqwest;
use tokio::time::Instant;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>>{
    let ct = Instant::now();
    let init = reqwest::get("https://api.hypixel.net/skyblock/auctions")
        .await?
        .json::<Root>()
        .await?;
    println!("{:?}", init.success);

    let paths = gen_url(init.total_pages).await;

    let fetches = futures::stream::iter(
        paths.into_iter().map(|path| {
            async move {
                match reqwest::get(&path).await {
                    Ok(resp) => {
                        match resp.json::<Root>().await {
                            Ok(_) => {
                                // println!("Got {}", root.page)
                            },
                            Err(_) => {
                                println!("ERROR unable to read")
                            }
                        }
                    },
                    Err(_) => {
                        println!("ERROR unable to download {}", path)
                    }
                }
            }
        })
    ).buffer_unordered(100).collect::<Vec<()>>();
    println!("Waiting...");
    fetches.await;
    println!("{:?}", ct.elapsed());
    Ok(())
}

async fn gen_url(total_pages: i64) -> Vec<String> {
    let path = "https://api.hypixel.net/skyblock/auctions";
    let mut paths:Vec<String> = vec![];
    for i in 1..total_pages {
        paths.push(format!("{}?page={}", path, i))
    }
    paths
}

#[derive(Default, Debug, Clone, PartialEq, Eq, Hash, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Root {
    pub success: bool,
    pub page: i64,
    pub total_pages: i64,
    pub total_auctions: i64,
    pub last_updated: i64,
    pub auctions: Vec<Auction>,
}
#[derive(Default, Debug, Clone, PartialEq, Eq, Hash, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Auction {
    pub uuid: String,
    pub auctioneer: String,
    #[serde(rename = "profile_id")]
    pub profile_id: String,
    pub coop: Vec<String>,
    pub start: i64,
    pub end: i64,
    #[serde(rename = "item_name")]
    pub item_name: String,
    #[serde(rename = "item_lore")]
    pub item_lore: String,
    pub extra: String,
    pub category: String,
    pub tier: String,
    #[serde(rename = "starting_bid")]
    pub starting_bid: i64,
    #[serde(rename = "item_bytes")]
    pub item_bytes: String,
    pub claimed: bool,
    //#[serde(rename = "claimed_bidders")]
    //pub claimed_bidders: Vec<Value>,
    #[serde(rename = "highest_bid_amount")]
    pub highest_bid_amount: i64,
    #[serde(rename = "last_updated")]
    pub last_updated: i64,
    pub bin: bool,
    pub bids: Vec<Bid>,
    #[serde(rename = "item_uuid")]
    pub item_uuid: Option<String>,
}
#[derive(Default, Debug, Clone, PartialEq, Eq, Hash, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Bid {
    #[serde(rename = "auction_id")]
    pub auction_id: String,
    pub bidder: String,
    #[serde(rename = "profile_id")]
    pub profile_id: String,
    pub amount: i64,
    pub timestamp: i64,
}