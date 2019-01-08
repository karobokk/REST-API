package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type TestServCase struct {
	Result     error
	Input      SearchRequest
	FindClient SearchClient
	NameCase   string
}

func TestClientWrongReq(t *testing.T) {

	TestReq := []TestServCase{
		{

			Result: fmt.Errorf("limit must be > 0"),
			Input: SearchRequest{
				Limit:      -25,
				Offset:     0,
				Query:      "do",
				OrderField: "Name",
				OrderBy:    -1,
			},
			FindClient: SearchClient{
				AccessToken: fmt.Sprintf("%x", md5.Sum([]byte("AccessToken"))),
			},
			NameCase: "Case wrong limit",
		},
		{
			Result: fmt.Errorf("offset must be > 0"),
			Input: SearchRequest{
				Limit:      25,
				Offset:     -25,
				Query:      "do",
				OrderField: "Name",
				OrderBy:    -1,
			},
			FindClient: SearchClient{
				AccessToken: fmt.Sprintf("%x", md5.Sum([]byte("AccessToken"))),
			},
			NameCase: "Case wrong offset",
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	fmt.Println("starting Test server")
	TestReq[0].FindClient.URL = ts.URL
	TestReq[1].FindClient.URL = ts.URL
	for _, Req := range TestReq {
		_, err := Req.FindClient.FindUsers(Req.Input)
		time.Sleep(time.Second)
		if err.Error() != Req.Result.Error() {
			log.Println(err)
			t.Errorf("Unexpected error")
		} else {
			log.Println(Req.NameCase, "passed")
		}

	}

	time.Sleep(time.Second)

}

func TestClientErDo(t *testing.T) {

	TestReq := []TestServCase{

		{
			Result: fmt.Errorf("unknown error Get loh/?limit=25&offset=0&order_by=-1&order_field=Name&query=do: unsupported protocol scheme \"\""),
			Input: SearchRequest{
				Limit:      24,
				Offset:     0,
				Query:      "do",
				OrderField: "Name",
				OrderBy:    -1,
			},
			FindClient: SearchClient{
				URL:         "loh/",
				AccessToken: fmt.Sprintf("%x", md5.Sum([]byte("AccessToken"))),
			},
			NameCase: "Case Wrong URL",
		},
		{

			Result: fmt.Errorf("timeout for limit=25&offset=0&order_by=-1&order_field=Name&query=do"),
			Input: SearchRequest{
				Limit:      24,
				Offset:     0,
				Query:      "do",
				OrderField: "Name",
				OrderBy:    -1,
			},
			FindClient: SearchClient{
				AccessToken: fmt.Sprintf("%x", md5.Sum([]byte("AccessToken"))),
			},
			NameCase: "Case Timeout",
		},
	}

	FakeWaitServ := httptest.NewServer(http.HandlerFunc(WaitTimeoutHan))
	fmt.Println("starting server at :8080")
	TestReq[1].FindClient.URL = FakeWaitServ.URL

	for _, Req := range TestReq {
		_, err := Req.FindClient.FindUsers(Req.Input)
		time.Sleep(time.Second)
		if err.Error() != Req.Result.Error() {
			log.Println(err)
			t.Errorf(Req.NameCase, "Failed")
		} else {

			log.Println(Req.NameCase, "passed")
		}

	}

	time.Sleep(time.Second)

}

func TestClientBadAcToken(t *testing.T) {

	TestReq := []TestServCase{

		{
			Result: fmt.Errorf("Bad AccessToken"),
			Input: SearchRequest{
				Limit:      24,
				Offset:     0,
				Query:      "do",
				OrderField: "Name",
				OrderBy:    -1,
			},
			FindClient: SearchClient{
				AccessToken: fmt.Sprintf("%x", md5.Sum([]byte("AccessTok"))),
			},
			NameCase: "Case Wrong Access Token",
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	fmt.Println("starting server at :8080")
	TestReq[0].FindClient.URL = ts.URL
	for _, Req := range TestReq {
		_, err := Req.FindClient.FindUsers(Req.Input)
		time.Sleep(time.Second)
		if err.Error() != Req.Result.Error() {
			log.Println(err)
			t.Errorf(Req.NameCase, "Failed")
		} else {

			log.Println(Req.NameCase, "passed")
		}

	}

	time.Sleep(time.Second)

}

func TestClientFatServer(t *testing.T) {

	TestReq := []TestServCase{

		{
			Result: fmt.Errorf("Bad AccessToken"),
			Input: SearchRequest{
				Limit:      25,
				Offset:     0,
				Query:      "do",
				OrderField: "Name",
				OrderBy:    -1,
			},
			FindClient: SearchClient{
				AccessToken: fmt.Sprintf("%x", md5.Sum([]byte("AccessTok"))),
			},
			NameCase: "Case Fatall Error On server",
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	fmt.Println("starting server at :8080")
	TestReq[0].FindClient.URL = ts.URL
	for _, Req := range TestReq {
		_, err := Req.FindClient.FindUsers(Req.Input)
		time.Sleep(time.Second)
		if err.Error() != Req.Result.Error() {
			log.Println(err)
			t.Errorf(Req.NameCase, "Failed")
		} else {

			log.Println(Req.NameCase, "passed")
		}

	}

	time.Sleep(time.Second)

}

func TestClientBadErrJSON(t *testing.T) {

	TestReq := []TestServCase{

		{
			Result: fmt.Errorf("cant unpack error json: unexpected end of JSON input"),
			Input: SearchRequest{
				Limit:      25,
				Offset:     0,
				Query:      "do",
				OrderField: "Name",
				OrderBy:    -1,
			},
			FindClient: SearchClient{
				AccessToken: fmt.Sprintf("%x", md5.Sum([]byte("AccessTok"))),
			},
			NameCase: "Case Err JSON request",
		},
	}

	ts := httptest.NewServer(http.HandlerFunc(FakeReqToClient))
	fmt.Println("starting server at :8080")
	TestReq[0].FindClient.URL = ts.URL
	for _, Req := range TestReq {
		_, err := Req.FindClient.FindUsers(Req.Input)
		time.Sleep(time.Second)
		if err.Error() != Req.Result.Error() {
			log.Println(err)
			t.Errorf(Req.NameCase, "Failed")
		} else {

			log.Println(Req.NameCase, "passed")
		}

	}

	time.Sleep(time.Second)

}

/*
	FindClient.URL = ts.URL //"http://127.0.0.1:8080/"
	FindClient.AccessToken = fmt.Sprintf("%x", md5.Sum([]byte("AccessToken")))

	FindClient1.URL = "loh/"
	FindClient1.AccessToken = fmt.Sprintf("%x", md5.Sum([]byte("AccessToken")))
*/

/*
	Req := SearchRequest{

		Limit:      25,
		Offset:     0,    // Можно учесть после сортировки
		Query:      "do", // подстрока в 1 из полей
		OrderField: "Name",
		OrderBy:    -1,
	}
*/

/*
	for Response.NextPage {
		Req.Offset += 25
		Response, _ = FindClient.FindUsers(Req)

		time.Sleep(time.Second)

		for _, User := range Response.Users {

			fmt.Println("Id ", "  ", User.Id)

				fmt.Println("Name ", "  ", User.Name)
				fmt.Println("Age ", "  ", User.Age)
				fmt.Println("About ", "  ", User.About)
				fmt.Println("Gender ", "  ", User.Gender)


		}
	}
*/

/*
	Response, err := FindClient.FindUsers(Req)
	if err != nil {
		fmt.Println(err)
		return
	}
*/

/*
	for _, User := range Response.Users {

	}
*/
