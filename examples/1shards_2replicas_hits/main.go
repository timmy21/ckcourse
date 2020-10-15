package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/timmy21/ckcourse/pkg/chutils"
	"github.com/timmy21/ckcourse/pkg/utils"
)

var (
	action string
	debug  bool
)

func init() {
	flag.StringVar(&action, "action", "", "action name, ex: init, loaddata, cleanup")
	flag.BoolVar(&debug, "debug", false, "output debug message")
}

func main() {
	flag.Parse()

	conn138, err := chutils.CreateConnect("ck1.example.com", chutils.WithDebug(debug))
	if err != nil {
		chutils.PrintError(err)
		log.Fatal(err)
	}
	defer conn138.Close()

	conn139, err := chutils.CreateConnect("ck2.example.com", chutils.WithDebug(debug))
	if err != nil {
		chutils.PrintError(err)
		log.Fatal(err)
	}
	defer conn139.Close()

	switch action {
	case "init":
		createTables(conn138, "01", "138")
		createTables(conn139, "01", "139")
	case "loaddata":
		loadDataFromTable(conn138)
	case "cleanup":
		dropTables(conn138)
		dropTables(conn139)
	default:
		sumRequestNumByUrl(conn139)
	}
}

func sumRequestNumByUrl(conn *sqlx.DB) {
	defer utils.Elapsed("avgDurationByUrl")()

	var items []struct {
		URL             string  `db:"URL"`
		TotalRequestNum float64 `db:"TotalRequestNum"`
	}
	query := `
SELECT
    URL,
    sum(RequestNum) AS TotalRequestNum
FROM tutorial.hits_replica
WHERE EventDate BETWEEN '2014-03-23' AND '2014-03-30'
GROUP BY URL
ORDER BY TotalRequestNum DESC
LIMIT 10
	`
	if err := conn.Select(&items, query); err != nil {
		log.Fatal(err)
	}
	for _, item := range items {
		fmt.Printf("URL: %s, TotalRequestNum: %f\n", item.URL, item.TotalRequestNum)
	}
}
func loadDataFromTable(conn *sqlx.DB) {
	conn.MustExec(`INSERT INTO tutorial.hits_replica SELECT * FROM tutorial.hits_v1`)
}

func createTables(conn *sqlx.DB, shard, replica string) {
	conn.MustExec(createDBSchema)
	conn.MustExec(fmt.Sprintf(createTableSchemaTpl, shard, replica))
}

func dropTables(conn *sqlx.DB) {
	conn.MustExec(dropTableSchema)
}

var dropTableSchema = `
DROP TABLE tutorial.hits_replica
`

var createDBSchema = `
CREATE DATABASE IF NOT EXISTS tutorial
`

var createTableSchemaTpl = `
CREATE TABLE IF NOT EXISTS tutorial.hits_replica
(
    WatchID UInt64,
    JavaEnable UInt8,
    Title String,
    GoodEvent Int16,
    EventTime DateTime,
    EventDate Date,
    CounterID UInt32,
    ClientIP UInt32,
    ClientIP6 FixedString(16),
    RegionID UInt32,
    UserID UInt64,
    CounterClass Int8,
    OS UInt8,
    UserAgent UInt8,
    URL String,
    Referer String,
    URLDomain String,
    RefererDomain String,
    Refresh UInt8,
    IsRobot UInt8,
    RefererCategories Array(UInt16),
    URLCategories Array(UInt16),
    URLRegions Array(UInt32),
    RefererRegions Array(UInt32),
    ResolutionWidth UInt16,
    ResolutionHeight UInt16,
    ResolutionDepth UInt8,
    FlashMajor UInt8,
    FlashMinor UInt8,
    FlashMinor2 String,
    NetMajor UInt8,
    NetMinor UInt8,
    UserAgentMajor UInt16,
    UserAgentMinor FixedString(2),
    CookieEnable UInt8,
    JavascriptEnable UInt8,
    IsMobile UInt8,
    MobilePhone UInt8,
    MobilePhoneModel String,
    Params String,
    IPNetworkID UInt32,
    TraficSourceID Int8,
    SearchEngineID UInt16,
    SearchPhrase String,
    AdvEngineID UInt8,
    IsArtifical UInt8,
    WindowClientWidth UInt16,
    WindowClientHeight UInt16,
    ClientTimeZone Int16,
    ClientEventTime DateTime,
    SilverlightVersion1 UInt8,
    SilverlightVersion2 UInt8,
    SilverlightVersion3 UInt32,
    SilverlightVersion4 UInt16,
    PageCharset String,
    CodeVersion UInt32,
    IsLink UInt8,
    IsDownload UInt8,
    IsNotBounce UInt8,
    FUniqID UInt64,
    HID UInt32,
    IsOldCounter UInt8,
    IsEvent UInt8,
    IsParameter UInt8,
    DontCountHits UInt8,
    WithHash UInt8,
    HitColor FixedString(1),
    UTCEventTime DateTime,
    Age UInt8,
    Sex UInt8,
    Income UInt8,
    Interests UInt16,
    Robotness UInt8,
    GeneralInterests Array(UInt16),
    RemoteIP UInt32,
    RemoteIP6 FixedString(16),
    WindowName Int32,
    OpenerName Int32,
    HistoryLength Int16,
    BrowserLanguage FixedString(2),
    BrowserCountry FixedString(2),
    SocialNetwork String,
    SocialAction String,
    HTTPError UInt16,
    SendTiming Int32,
    DNSTiming Int32,
    ConnectTiming Int32,
    ResponseStartTiming Int32,
    ResponseEndTiming Int32,
    FetchTiming Int32,
    RedirectTiming Int32,
    DOMInteractiveTiming Int32,
    DOMContentLoadedTiming Int32,
    DOMCompleteTiming Int32,
    LoadEventStartTiming Int32,
    LoadEventEndTiming Int32,
    NSToDOMContentLoadedTiming Int32,
    FirstPaintTiming Int32,
    RedirectCount Int8,
    SocialSourceNetworkID UInt8,
    SocialSourcePage String,
    ParamPrice Int64,
    ParamOrderID String,
    ParamCurrency FixedString(3),
    ParamCurrencyID UInt16,
    GoalsReached Array(UInt32),
    OpenstatServiceName String,
    OpenstatCampaignID String,
    OpenstatAdID String,
    OpenstatSourceID String,
    UTMSource String,
    UTMMedium String,
    UTMCampaign String,
    UTMContent String,
    UTMTerm String,
    FromTag String,
    HasGCLID UInt8,
    RefererHash UInt64,
    URLHash UInt64,
    CLID UInt32,
    YCLID UInt64,
    ShareService String,
    ShareURL String,
    ShareTitle String,
    ParsedParams Nested(
        Key1 String,
        Key2 String,
        Key3 String,
        Key4 String,
        Key5 String,
        ValueDouble Float64),
    IslandID FixedString(16),
    RequestNum UInt32,
    RequestTry UInt8
)
ENGINE = ReplicatedMergeTree(
    '/clickhouse/tables/%s/hits_replica',
    '%s'
)
PARTITION BY toYYYYMM(EventDate)
ORDER BY (CounterID, EventDate, intHash32(UserID))
SAMPLE BY intHash32(UserID)
SETTINGS index_granularity = 8192
`
