package transmitter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	drive "google.golang.org/api/drive/v3"
)

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile := tokenCacheFile()
	// if err != nil {
	// 	log.Fatalf("Unable to get path to cached credential file. %v", err)
	// }
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() string {
	tokenCacheDir := filepath.Join(".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("client.json"))
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.Create(file)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

var client *http.Client
var srv *drive.Service

func init() {
	ctx := context.Background()

	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/drive-go-quickstart.json
	config, err := google.ConfigFromJSON(b, drive.DriveFileScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client = getClient(ctx, config)

	srv, err = drive.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve drive Client %v", err)
	}

	// f, err := os.Open(filename)
	// if err != nil {
	// 	log.Fatalf("error opening %q: %v", filename, err)
	// }
	// defer f.Close()

	//driveFile, err := srv.Files.Create(&drive.File{Name: "Realtime Data"}).Media(f).Do()
	// var temp []string
	// temp = append(temp, "1p4J2b8cDhAA9NjKGd0xNfx91BA_nFpvc")
	// dbfile := drive.File{Name: "RealtimeData.db"}

	// fileID := "1AqR6URNNnWPxdFxkd-tBrP2Uv1CAJAjy"
	//driveFile, err := srv.Files.Create(&dfile).Media(f).Do()

	// driveFile, err := srv.Files.Update(fileID, &dbfile).AddParents("1p4J2b8cDhAA9NjKGd0xNfx91BA_nFpvc").Media(f).Do()

	// if err != nil {
	// 	log.Fatalf("Unable to create file: %v", err)
	// }
	// log.Println(driveFile.Id)
	// log.Println("done")
}
