package lookup

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	apiURL string = "https://od-api.oxforddictionaries.com/api/v1"
	//apiId  string = //Enter ID
	//apiKey string = //Enter Key
)

type Definition struct {
	Word         string
	PartOfSpeech string
	Definition   string
}

//Allows Definition to be printed as formatted string
func (def Definition) String() string {
	return fmt.Sprintf("%s (%s):\n%s", def.Word, def.PartOfSpeech, def.Definition)
}

//Uses Oxford Dictionaries Lemmatron API to find the root of a word
func getRoot(word string) (string, error) {
	URL := fmt.Sprintf("%s/%s/en/%s", apiURL, "inflections", word)
	client := &http.Client{}
	request, err := http.NewRequest("GET", URL, nil)
	request.Header.Add("app_id", apiID)
	request.Header.Add("app_key", apiKey)

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	root := gjson.GetBytes(responseBody, "results.0.lexicalEntries.0.inflectionOf.0.text")
	return fmt.Sprint(root), nil
}

//Retrieves dictionary entry of word using Oxford Dictionaries API
func getDictionaryEntry(word string) (definition string, partOfSpeech string, err error) {
	URL := fmt.Sprintf("%s/%s/en/%s", apiURL, "entries", word)
	client := &http.Client{}
	request, err := http.NewRequest("GET", URL, nil)
	request.Header.Add("app_id", apiID)
	request.Header.Add("app_key", apiKey)

	response, err := client.Do(request)
	if err != nil {
		return "", "", err
	}
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", "", err
	}

	def := gjson.GetBytes(responseBody, "results.0.lexicalEntries.0.entries.0.senses.0.definitions.0")
	part := gjson.GetBytes(responseBody, "results.0.lexicalEntries.0.lexicalCategory")
	return fmt.Sprint(def), fmt.Sprint(part), nil
}

//Retrieves the English definition of a word
func GetDefinition(word string) (Definition, error) {
	root, err := getRoot(word)
	if err != nil {
		return Definition{}, err
	}
	definition, partOfSpeech, err := getDictionaryEntry(root)
	if err != nil {
		return Definition{}, err
	}
	return Definition{Word: strings.Title(root), Definition: definition, PartOfSpeech: partOfSpeech}, nil
}
