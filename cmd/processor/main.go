package main

import (
    "context"
    "encoding/json"
    "flag"
    "time"

    "database/sql"
    _ "github.com/mattn/go-sqlite3"

    "cloud.google.com/go/pubsub"
    "github.com/censys/scan-takehome/pkg/scanning"

)

var (
    services = []string{"HTTP", "SSH", "DNS"}
    dbFile = flag.String("db", "/data/processor.db", "Database File")
)

func main() {
    projectId := flag.String("project", "test-project", "GCP Project ID")
    subscriptionId := flag.String("subscription", "scan-sub", "GCP PubSub Subscription ID")
    flag.Parse()

    ctx := context.Background()

    client, err := pubsub.NewClient(ctx, *projectId)
    if err != nil {
        panic(err)
    }

    subscription := client.Subscription(*subscriptionId)

    db, err := sql.Open("sqlite3", *dbFile)
    if err != nil {
        panic(err)
    }

	createTableSQL := `CREATE TABLE IF NOT EXISTS scans (
		ip TEXT NOT NULL,  
	    port INTEGER NOT NULL,  
	    service TEXT NOT NULL,  
	    response TEXT NOT NULL,  
	    scan_time DATETIME NOT NULL,  
	    PRIMARY KEY (ip, port, service)
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		panic(err)
	}

    err = subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
        var scan scanning.Scan
        if err := json.Unmarshal(msg.Data, &scan); err != nil {
            panic(err)
        }

        var serviceResp string

        switch scan.DataVersion {
        case scanning.V1:
			var respData scanning.V1Data
            dataMap := scan.Data.(map[string]interface{})
			dataBytes, err := json.Marshal(dataMap)
			if err != nil {
				panic(err)
			}
			json.Unmarshal(dataBytes, &respData)
            serviceResp = string(respData.ResponseBytesUtf8)
        case scanning.V2:
			var respData scanning.V2Data
			dataMap := scan.Data.(map[string]interface{})
			dataBytes, err := json.Marshal(dataMap)
			if err != nil {
				panic(err)
			}
			json.Unmarshal(dataBytes, &respData)
            serviceResp = respData.ResponseStr
        default:
            panic("Invalid DataVersion in scan")
        }
        // consider including timestamp in message so retried messages don't overwrite new ones
        _, err = db.Exec(`
            INSERT OR REPLACE INTO scans (ip, port, service, response, scan_time)
            VALUES (?, ?, ?, ?, ?)`, scan.Ip, scan.Port, scan.Service, serviceResp, time.Now()) 
        
        if err != nil {
            panic("Database insert failed") 
        }

        msg.Ack()
    })

    if err != nil {
        panic(err)
    }
}
