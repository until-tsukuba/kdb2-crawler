package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

const USER_AGENT = "kdb2-crawler (+https://github.com/until-tsukuba/kdb2-crawler)"

func main() {
	client := new(http.Client)
	client.Jar, _ = cookiejar.New(nil)

	req, _ := http.NewRequest(
		"GET",
		"https://kdb.tsukuba.ac.jp/",
		nil,
	)
	req.Header.Set("User-Agent", USER_AGENT)
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "initial req: %v", err)
		return
	}

	fmt.Printf("url: %s\n", resp.Request.URL)
	courseResp, err := searchCourse(client, resp.Request.URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "searchCourse: %v", err)
		return
	}
	csvResp, err := downloadCSV(client, courseResp.Request.URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "downloadCSV: %v", err)
		return
	}

	sjisReader := transform.NewReader(csvResp.Body, japanese.ShiftJIS.NewDecoder())
	tee := io.TeeReader(sjisReader, os.Stdout)
	s := bufio.NewScanner(tee)
	for s.Scan() {
	}

	if err := s.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		return
	}

}

func searchCourse(client *http.Client, requestURL *url.URL) (*http.Response, error) {
	values := &url.Values{}
	initValues(values, 2022)
	values.Set("_eventId", "searchOpeningCourse")
	//fmt.Println(values.Encode())

	req, err := http.NewRequest(
		"POST",
		requestURL.String(),
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func downloadCSV(client *http.Client, requestURL *url.URL) (*http.Response, error) {
	values := &url.Values{}
	initValues(values, 2022)
	values.Set("_eventId", "output")
	values.Set("outputFormat", "0")
	//fmt.Println(values.Encode())

	csvReq, err := http.NewRequest(
		"POST",
		requestURL.String(),
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		return nil, err
	}
	csvReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(csvReq)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func initValues(values *url.Values, nendo uint) {
	values.Add("index", "")
	values.Add("locale", "")
	values.Add("nendo", fmt.Sprint(nendo))
	values.Add("termCode", "")
	values.Add("dayCode", "")
	values.Add("periodCode", "")
	values.Add("campusCode", "")
	values.Add("hierarchy1", "")
	values.Add("hierarchy2", "")
	values.Add("hierarchy3", "")
	values.Add("hierarchy4", "")
	values.Add("hierarchy5", "")
	values.Add("freeWord", "")
	values.Add("_orFlg", "1")
	values.Add("_andFlg", "1")
	values.Add("_gaiyoFlg", "1")
	values.Add("_risyuFlg", "1")
	values.Add("_excludeFukaikoFlg", "1")
	//values.Add("_eventId", "")
	//values.Add("_outputFormat", "")
}
