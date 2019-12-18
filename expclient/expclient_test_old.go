package expclient

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"testing"
)

func Test_New(t *testing.T) {

	exec := &http.Client{}
	expected := &Client{
		exec: exec,
	}

	cases := []struct {
		exec     *http.Client
		expected *Client
	}{
		{
			exec:     &http.Client{},
			expected: expected,
		},
		{
			exec:     nil,
			expected: expected,
		},
	}

	for _, tc := range cases {
		got := New(tc.exec)
		equals(t, tc.expected, got)
	}

}

var postsOkResp = &Resp{
	Id: 3,
	User: 0,
	Title: "ea molestias quasi exercitationem repellat qui ipsa sit aut",
	Body:  "et iusto sed quo iure\nvoluptatem occaecati omnis eligendi aut ad\nvoluptatem doloribus vel accusantium quis pariatur\nmolestiae porro eius odio et labore et velit aut",
}

var postsWrongIdType = &Resp{}


func TestClient_GetPosts(t *testing.T) {
	exec := &http.Client{}

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		base := path.Base(r.URL.String())

		if _ , err := strconv.Atoi(base); err != nil {
			_, _ = w.Write([]byte(`{}`))
			return
		}
		_, _ =  w.Write([]byte(`{"user":0,"id":3,"title":"ea molestias quasi exercitationem repellat qui ipsa sit aut","body":"et iusto sed quo iure\nvoluptatem occaecati omnis eligendi aut ad\nvoluptatem doloribus vel accusantium quis pariatur\nmolestiae porro eius odio et labore et velit aut","hash":[222,195,103,28,209,30,27,20,159,76,185,39,167,99,62,204]}`))
	})

	testexec, teardown := testingHttpClient(h, exec)
	client := New(testexec)
	defer teardown()

	cases := []struct {
		id     string
		expected *Resp
	}{
		{
			id:     "3",
			expected: postsOkResp,
		},
		{
			id:     "f",
			expected: postsWrongIdType,
		},
		{
			id:     "",
			expected: postsWrongIdType,
		},
	}

	for _, tc := range cases {


		postsOkResp.ToUpperBody()

		got, err := client.GetPosts(tc.id)
		ok(t, err)



		equals(t, tc.expected, got)
	}
}

func testingHttpClient(handler http.Handler, exec *http.Client) (*http.Client, func()) {
	s := httptest.NewTLSServer(handler)

	url := s.Listener.Addr()

	fmt.Println(url.String())

	exec.Transport = &http.Transport{
		DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
			return net.Dial(network, url.String())
		},
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	return exec, s.Close
}

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
