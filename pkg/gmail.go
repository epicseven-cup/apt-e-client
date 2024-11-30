package pkg

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/api/googleapi"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func watch(user string, srv *gmail.Service) {
	watch := gmail.WatchRequest{
		LabelFilterAction:   "",
		LabelFilterBehavior: "",
		LabelIds:            nil,
		TopicName:           "",
		ForceSendFields:     nil,
		NullFields:          nil,
	}
	webhook := srv.Users.Watch(user, &watch)

	r, err := webhook.Do()
	if err != nil {
		log.Fatal("Unable to retrieve user's Gmail webhook %v", err)
	}
	fmt.Println(r.MarshalJSON())
}

type Inbox struct {
	messages      []*gmail.Message
	nextPageToken string
}

func populateFakeMessages(size int64) []*gmail.Message {
	messages := make([]*gmail.Message, size)
	for i := range size {
		messages[i] = &gmail.Message{
			HistoryId:       0,
			Id:              "",
			InternalDate:    0,
			LabelIds:        nil,
			Payload:         nil,
			Raw:             "RAW CONTENT\n",
			SizeEstimate:    0,
			Snippet:         "Fake Message Snippet",
			ThreadId:        "",
			ServerResponse:  googleapi.ServerResponse{},
			ForceSendFields: nil,
			NullFields:      nil,
		}
	}
	return messages
}

func InitFakeInbox(maxDisplay int64) *Inbox {
	return &Inbox{
		messages:      populateFakeMessages(maxDisplay),
		nextPageToken: "",
	}
}

func InitInbox(maxDisplay int64) *Inbox {
	ctx := context.Background()
	b, err := os.ReadFile("secret/credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}
	user := "me"

	messageListRequest := srv.Users.Messages.List(user)
	messageListRequest.MaxResults(maxDisplay)
	// Getting the message list request now
	messageListRespond, err := messageListRequest.Do()
	if err != nil {
		log.Fatalf("Fail to get respond from message list %v", err)
	}
	token := messageListRespond.NextPageToken
	messages := make([]*gmail.Message, maxDisplay)
	// Going through the list of message that we got, since the GetMessageList only return the array of Message Ids we need to grab it again with out GetMessage request
	for index, message := range messageListRespond.Messages {
		messageRequest := srv.Users.Messages.Get(user, message.Id).Format("full")
		messageRespond, err := messageRequest.Do()
		if err != nil {
			log.Fatalf("Fail to get respond from message resond: %v", err)
		}
		messages[index] = messageRespond

	}

	return &Inbox{
		messages:      messages,
		nextPageToken: token,
	}
}
