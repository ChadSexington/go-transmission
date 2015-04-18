package transmission

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/auth"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	mux                *http.ServeMux
	transmissionClient TransmissionClient
	server             *httptest.Server
)

func setup(output string) {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	m := martini.New()
	r := martini.NewRouter()
	r.Post("/transmission/rpc", func() string {
		return output
	})
	m.Action(r.Handle)
	m.Use(auth.Basic("test", "test"))
	mux.Handle("/", m)

	transmissionClient = New(server.URL, "test", "test")
}

func teardown() {
	server.Close()
}

func TestGetTorrents(t *testing.T) {
	setup(`{"arguments":{"torrents":[{"eta":-1,"id":5,
  "leftUntilDone":0,"name":"Test", "files":[],
  "rateDownload":0,"rateUpload":0,"status":6,"uploadRatio":0.3114}]},
  "result":"success"}`)
	defer teardown()

	Convey("Test get list torrents", t, func() {
		torrents, err := transmissionClient.GetTorrents()
		So(err, ShouldBeNil)
		So(len(torrents), ShouldEqual, 1)
	})
}

func TestGetDownloadedTorrents(t *testing.T) {
	setup(`{"arguments":{"torrents":[{"eta":-1,"id":5,
  "leftUntilDone":0,"name":"Test", "files":[],
  "rateDownload":0,"rateUpload":0,"status":6,"uploadRatio":0.3114},
	{"eta":-1,"id":5,
  "leftUntilDone":10,"name":"Test2", "files":[],
  "rateDownload":1,"rateUpload":0,"status":6,"uploadRatio":0.3114}]},
  "result":"success"}`)
	defer teardown()

	Convey("Test get list torrents", t, func() {
		torrents, err := transmissionClient.GetDownloadedTorrents()
		So(err, ShouldBeNil)
		So(len(torrents), ShouldEqual, 1)
	})
}

func TestRemoveTorrent(t *testing.T) {
	setup(`{"arguments":{},"result":"success"}`)
	defer teardown()

	Convey("Test removing torrent", t, func() {
		result, err := transmissionClient.RemoveTorrent(1, true)
		So(err, ShouldBeNil)
		So(result, ShouldEqual, "success")
	})
}

func TestAddTorrentByFile(t *testing.T) {
	setup(`{"arguments":{"torrent-added":
  {"hashString":"875a2d90068c32b4ce7992eaf56cd03f5be0d193",
  "id":23,"name":"Test Name"}}
  ,"result":"success"}`)
	defer teardown()

	Convey("Test adding torrent", t, func() {
		addedTorrent, err := transmissionClient.AddTorrentByFile("/tmp/file", "/home/lnguyen")
		So(err, ShouldBeNil)
		So(addedTorrent.Name, ShouldEqual, "Test Name")
		So(addedTorrent.ID, ShouldEqual, 23)
	})
}
