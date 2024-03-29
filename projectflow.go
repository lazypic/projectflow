package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
	now = time.Now()
	// db setting
	flagRegion  = flag.String("region", "ap-northeast-2", "AWS region name")
	flagProfile = flag.String("profile", "lazypic", "AWS Credentials profile name")
	flagTable   = flag.String("table", "projectflow", "AWS Dynamodb table name")

	// mode and partition key
	flagAdd = flag.Bool("add", false, "user addition mode")
	flagSet = flag.Bool("set", false, "user update mode")

	// date
	flagHelp = flag.Bool("help", false, "print help")

	// attributes
	flagID            = flag.String("id", "", "project id")
	flagBudget        = flag.Uint64("budget", 0, "project budget")
	flagStartDate     = flag.String("start", now.Format(time.RFC3339), "project start date")
	flagUpdateDate    = flag.String("updatedate", now.Format(time.RFC3339), "project start date")
	flagEndDate       = flag.String("end", "", "project end date")
	flagProjectStatus = flag.String("status", "", "project status")
	flagSearchword    = flag.String("search", "", "search word")
	flagMonetaryUnit  = flag.String("monetaryunit", "KRW", "monetary unit")
)

func main() {
	log.SetPrefix("projectflow: ")
	log.SetFlags(0)
	flag.Parse()

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config:            aws.Config{Region: aws.String(*flagRegion)},
		Profile:           *flagProfile,
	}))
	db := dynamodb.New(sess)

	// 테이블이 존재하는지 점검하고 없다면 테이블을 생성한다.
	if !validTable(*db, *flagTable) {
		_, err := db.CreateTable(tableStruct(*flagTable))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		fmt.Println("Created table:", *flagTable)
		fmt.Println("Please try again in one minute.")
		os.Exit(0)
	}
	if *flagHelp {
		flag.Usage()
	}
	if *flagAdd && *flagID != "" {
		err := AddProject(*db)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	} else if *flagSet && *flagID != "" {
		err := SetProject(*db)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	} else if *flagSearchword != "" {
		err := GetProjects(*db, *flagSearchword)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	} else {
		flag.PrintDefaults()
	}
}
