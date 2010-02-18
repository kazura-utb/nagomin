package twitter

import (
	"json"
	"bufio"
	"io"
	"net"
	"bytes"
	"encoding/base64"
	"strings"
	"http"
	"fmt"
	"os"
)

const(
	twitterHost	= "twitter.com"
	dataFormat = "json"
//Timeline Methods
	publicTimeline = "/statuses/public_timeline"
	homeTimeline = "/statuses/home_timeline"
	friendsTimeline = "/statuses/friends_timeline"
	userTimeline = "/statuses/user_timeline"
	mentions = "/statuses/mentions"
	retweetedByMe = "/statuses/retweeted_by_me"
	retweetedToMe = "/statuses/retweeted_to_me"
	retweetedOfMe = "/statuses/retweeted_of_me"
//Status Methods
	statusesShow = "/statuses/show"
	statusesUpdate = "/statuses/update"
	statusesDestroy = "/statuses/destroy"
	statusesRetweet = "/statuses/retweet"
	statusesRetweets = "/statuses/retweets"
//User Methods
	usersShow = "/users/show"
	usersSearch = "/users/search"
	statusesFriends = "/statuses/friends"
	statusesFollowers = "/statuses/followers"
)

// 以下3つはjsonデータのパース用
type User struct {
	Name        string
	Screen_name string
}
type Tweet struct {
	User User
	Text string
	Id uint64
}

// ユーザアカウント情報
type Auth struct {
	Name	string
	Passwd	string
}
// クライアント
// こいつにいくつかインタフェースを実装
type Client struct{
	Auth Auth
}

func ClientBuilder(name string, passwd string) (cc Client){
	var c Client
	c.Auth.Name = name
	c.Auth.Passwd = passwd

	return c
}

// Auth構造体に関連付けたメソッド
// BASIC認証用にName/PasswdをBASE64エンコード
func (a *Auth) Base64enc() (s string){
	authSrc := a.Name + ":" + a.Passwd
	bb := &bytes.Buffer{}
	encoder := base64.NewEncoder(base64.StdEncoding, bb)
	encoder.Write(strings.Bytes(authSrc))
	encoder.Close()
	authEnc := bb.String()

	return authEnc
}

// ツイートをポストするメソッド
func (c *Client)Post(stat string) {
	// 送信本文をURLエンコード
	sStatURL := http.URLEscape(stat)

	if con, err := net.Dial("tcp", "", "twitter.com:80"); err == nil {
		io.WriteString(con, "POST /statuses/update.json?status=" + sStatURL + " HTTP/1.1\r\n")
		io.WriteString(con, "Host: twitter.com\r\n")
		io.WriteString(con, "Authorization: Basic " + c.Auth.Base64enc() + "\r\n")
		io.WriteString(con, "\r\n")

		con.Close()
	}
}

func addParam(params string, name string, value string)(p string){
	var pre int

	if params == ""{
		pre = '?'
	}else{
		pre = '&'
	}

	params += fmt.Sprintf("%c%s=%s", pre, name, value)

	return params
}

func (c *Client)makeAuthURL(request string, suffix string)(s string){
	return fmt.Sprintf("http://%s:%s@%s%s.%s%s", c.Auth.Name, c.Auth.Passwd, twitterHost, request, dataFormat, suffix)
}

func (c *Client)makePublicURL(request string, suffix string)(s string){
	return fmt.Sprintf("http://%s%s.%s%s", twitterHost, request, dataFormat, suffix)
}

func (c *Client)PublicTimeline() (t []Tweet){
	var tweets []Tweet

	url := c.makePublicURL(publicTimeline, "")
	res, _, err := http.Get(url)
	if err != nil{
		return nil
	}

	if res.Status != "200 OK"{
		return nil
	}

    reader := bufio.NewReader(res.Body);
    line, _ := reader.ReadString(0);

	json.Unmarshal(line, &tweets)

	return tweets
}

func (c *Client)HomeTimeline(sinceId uint64, maxId uint64, count uint, page uint) (t []Tweet){
	var params string
	var tweets []Tweet

	if sinceId != 0{
		params = addParam(params, "since_id", fmt.Sprintf("%d", sinceId))
	}
	if maxId != 0{
		params = addParam(params, "max_id", fmt.Sprintf("%d", maxId))
	}
	if count != 0{
		params = addParam(params, "count", fmt.Sprintf("%d", count))
	}
	if page != 0{
		params = addParam(params, "page", fmt.Sprintf("%d", page))
	}

	url := c.makeAuthURL(homeTimeline, params)
	res, _, err := http.Get(url)
	if err != nil{
		return nil
	}

	if res.Status != "200 OK"{
		return nil
	}

    reader := bufio.NewReader(res.Body);
    line, _ := reader.ReadString(0);

	json.Unmarshal(line, &tweets)

	return tweets
}

func (c *Client)FriendsTimeline(sinceId uint64, maxId uint64, count uint, page uint) (t []Tweet){
	var params string
	var tweets []Tweet

	if sinceId != 0{
		params = addParam(params, "since_id", fmt.Sprintf("%d", sinceId))
	}
	if maxId != 0{
		params = addParam(params, "max_id", fmt.Sprintf("%d", maxId))
	}
	if count != 0{
		params = addParam(params, "count", fmt.Sprintf("%d", count))
	}
	if page != 0{
		params = addParam(params, "page", fmt.Sprintf("%d", page))
	}

	url := c.makeAuthURL(friendsTimeline, params)
	res, _, err := http.Get(url)
	if err != nil{
		return nil
	}

	if res.Status != "200 OK"{
		return nil
	}

    reader := bufio.NewReader(res.Body);
    line, _ := reader.ReadString(0);

	json.Unmarshal(line, &tweets)

	return tweets
}

func (c *Client)UserTimeline(userId uint64, screenName string, sinceId uint64, maxId uint64, count uint, page uint) (t []Tweet){
	var params string
	var tweets []Tweet

	if userId != 0{
		params = addParam(params, "user_id", fmt.Sprintf("%d", userId))
	}
	if screenName != ""{
		params = addParam(params, "screen_name", screenName)
	}
	if sinceId != 0{
		params = addParam(params, "since_id", fmt.Sprintf("%d", sinceId))
	}
	if maxId != 0{
		params = addParam(params, "max_id", fmt.Sprintf("%d", maxId))
	}
	if count != 0{
		params = addParam(params, "count", fmt.Sprintf("%d", count))
	}
	if page != 0{
		params = addParam(params, "page", fmt.Sprintf("%d", page))
	}

	url := c.makeAuthURL(userTimeline, params)
	res, _, err := http.Get(url)
	if err != nil{
		return nil
	}

	if res.Status != "200 OK"{
		return nil
	}

    reader := bufio.NewReader(res.Body);
    line, _ := reader.ReadString(0);

	json.Unmarshal(line, &tweets)

	return tweets
}

func (c *Client)Mentions(sinceId uint64, maxId uint64, count uint, page uint) (t []Tweet){
	var params string
	var tweets []Tweet

	if sinceId != 0{
		params = addParam(params, "since_id", fmt.Sprintf("%d", sinceId))
	}
	if maxId != 0{
		params = addParam(params, "max_id", fmt.Sprintf("%d", maxId))
	}
	if count != 0{
		params = addParam(params, "count", fmt.Sprintf("%d", count))
	}
	if page != 0{
		params = addParam(params, "page", fmt.Sprintf("%d", page))
	}

	url := c.makeAuthURL(mentions, params)
	res, _, err := http.Get(url)
	if err != nil{
		return nil
	}

	if res.Status != "200 OK"{
		return nil
	}

    reader := bufio.NewReader(res.Body);
    line, _ := reader.ReadString(0);

	json.Unmarshal(line, &tweets)

	return tweets
}

func (c *Client)RetweetedByMe(sinceId uint64, maxId uint64, count uint, page uint) (t []Tweet){
	var params string
	var tweets []Tweet

	if sinceId != 0{
		params = addParam(params, "since_id", fmt.Sprintf("%d", sinceId))
	}
	if maxId != 0{
		params = addParam(params, "max_id", fmt.Sprintf("%d", maxId))
	}
	if count != 0{
		params = addParam(params, "count", fmt.Sprintf("%d", count))
	}
	if page != 0{
		params = addParam(params, "page", fmt.Sprintf("%d", page))
	}

	url := c.makeAuthURL(retweetedByMe, params)
	res, _, err := http.Get(url)
	if err != nil{
		return nil
	}

	if res.Status != "200 OK"{
		return nil
	}

    reader := bufio.NewReader(res.Body);
    line, _ := reader.ReadString(0);

	json.Unmarshal(line, &tweets)

	return tweets
}

func (c *Client)RetweetedToMe(sinceId uint64, maxId uint64, count uint, page uint) (t []Tweet){
	var params string
	var tweets []Tweet

	if sinceId != 0{
		params = addParam(params, "since_id", fmt.Sprintf("%d", sinceId))
	}
	if maxId != 0{
		params = addParam(params, "max_id", fmt.Sprintf("%d", maxId))
	}
	if count != 0{
		params = addParam(params, "count", fmt.Sprintf("%d", count))
	}
	if page != 0{
		params = addParam(params, "page", fmt.Sprintf("%d", page))
	}

	url := c.makeAuthURL(retweetedToMe, params)
	res, _, err := http.Get(url)
	if err != nil{
		return nil
	}

	if res.Status != "200 OK"{
		return nil
	}

    reader := bufio.NewReader(res.Body);
    line, _ := reader.ReadString(0);

	json.Unmarshal(line, &tweets)

	return tweets
}

func (c *Client)RetweetedOfMe(sinceId uint64, maxId uint64, count uint, page uint) (t []Tweet){
	var params string
	var tweets []Tweet

	if sinceId != 0{
		params = addParam(params, "since_id", fmt.Sprintf("%d", sinceId))
	}
	if maxId != 0{
		params = addParam(params, "max_id", fmt.Sprintf("%d", maxId))
	}
	if count != 0{
		params = addParam(params, "count", fmt.Sprintf("%d", count))
	}
	if page != 0{
		params = addParam(params, "page", fmt.Sprintf("%d", page))
	}

	url := c.makeAuthURL(retweetedOfMe, params)
	res, _, err := http.Get(url)
	if err != nil{
		return nil
	}

	if res.Status != "200 OK"{
		return nil
	}

    reader := bufio.NewReader(res.Body);
    line, _ := reader.ReadString(0);

	json.Unmarshal(line, &tweets)

	return tweets
}

func (c *Client)Tweet(tweet string) (err os.Error){
	return c.StatusesUpdate(tweet, 0)
}

func (c *Client)ReplyTweet(tweet string, replyId uint64) (err os.Error){
	return c.StatusesUpdate(tweet, replyId)
}

func (c *Client)StatusesUpdate(status string, replyId uint64) (err os.Error){
	var params string

	if status == ""{
		return os.NewError("must need \"status\" parameter.")
	}

	params = addParam(params, "status", http.URLEscape(status))

	if replyId != 0{
		params = addParam(params, "in_reply_to_status_id", fmt.Sprintf("%d", replyId))
	}

	url := c.makeAuthURL(statusesUpdate, params)
	res, err := http.Post(url, "", bytes.NewBufferString(""))
	if err != nil{
		return err
	}

	if res.Status != "200 OK"{
		return os.NewError("Request failed")
	}

	return nil
}

