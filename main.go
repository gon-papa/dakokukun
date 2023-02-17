package main

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sclevine/agouti"
)

func main() {
	// 引数確認
	clock := flag.String("c", "", "in 出勤 out 退勤")
	flag.Parse()
	if *clock == "in" {
		log.Println("出勤打刻します。")
	} else if *clock == "out" {
		log.Println("退勤打刻します。")
	} else {
		log.Fatal("引数が間違っています。" + flag.Arg(0))
	}

	// 環境変数読み込み
	er := godotenv.Load(".env")
	if er != nil {
		log.Fatal("環境変数が読み込めません。")
	}

	// ChromeDriverの立ち上げ
	agoutiDriver := agouti.ChromeDriver()
	agoutiDriver.Start()

	// プログラム終了時にChromeDriverを閉じる
	defer agoutiDriver.Stop()

	// リンク先のページを開く
	page, _ := agoutiDriver.NewPage()
	page.Navigate(os.Getenv("TARGET_URL"))

	// ログイン
	page.Find("#txtID").Fill(os.Getenv("ID"))
	page.Find("#txtPsw").Fill(os.Getenv("PASS"))
	page.Find("#btnLogin").Click()

	// 打刻ページへ
	page.Find("#ctl00_ContentPlaceHolder1_imgBtnSyuugyou").Click()

	if *clock == "in" {
		// 始業打刻
		page.Find("#ctl00_ContentPlaceHolder1_ibtnIn3").Click()
	} else {
		// 就業打刻
		page.Find("#ctl00_ContentPlaceHolder1_imgOut3").Click()
	}

	// 結果のスクショ
	page.Screenshot(os.Getenv("IMAGE_PATH"))

	// 始業打刻チェック
	_, err := page.Find(".tableBdr2").Text()
	if err != nil {
		log.Printf("正常に打刻できていません。err: %+v", err)
	} else {
		log.Println("打刻が完了しました。")
	}
}
