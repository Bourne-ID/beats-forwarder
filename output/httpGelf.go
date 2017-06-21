package output

import (
	cfg "github.com/Bourne-ID/beats-forwarder/config"
	"errors"
	"bytes"
	"net/http"
	"github.com/Graylog2/go-gelf/gelf"
	"time"
	"encoding/json"
	"log"
	"reflect"
	"github.com/Sirupsen/logrus"
)

type HTTPGelfClient struct {
	endpoint string
}

func (c *HTTPGelfClient) Init(config *cfg.Config) error {

	if (config.Output.HTTPGelf.Endpoint == nil || *config.Output.HTTPGelf.Endpoint == "" ) {
		return errors.New("No endpoint URL provided.")
	}
	c.endpoint = *config.Output.HTTPGelf.Endpoint

	return nil
}

func (c *HTTPGelfClient) WriteAndRetry(payload []byte) (error) {

	//payload should be uncompressed BEATS JSON.
	//I'm sure there's a better way of doing this, but I'm going to use the method being used in Graylog BEATS importing.
	var f interface{}
	err := json.Unmarshal(payload, &f)
	if err != nil {
		log.Println(err)
		//return error code
	}
	jsonMap := f.(map[string]interface{})
	metadata := jsonMap["@metadata"].(map[string]interface{})
	beat := metadata["beat"]
	var gelfMessage gelf.Message
	if beat == "filebeat" {
		gelfMessage = parseFilebeat(jsonMap)
	} else if beat == "winlogbeat" {
		gelfMessage = parseWinlogbeat(jsonMap)
	} else {
		gelfMessage = parseGenericBeat(jsonMap)
	}

	var buf bytes.Buffer
	buf.Grow(1024)

	gelfMessage.MarshalJSONBuf(&buf)

	req, err := http.NewRequest("POST", c.endpoint, &buf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err = client.Do(req)

	if err != nil {
		return err
	}

	return nil
}

func parseFilebeat(jsonMap map[string]interface{}) (gelf.Message) {
	messageString := jsonMap["message"].(string)
	gelfMessage := CreateMessage(messageString, jsonMap);


	gelfMessage.Extra["_facility"] = "filebeat"
	gelfMessage.Extra["_file"] = jsonMap["source"]
	gelfMessage.Extra["_input_type"] = jsonMap["input_type"]
	gelfMessage.Extra["_count"] = jsonMap["count"]
	gelfMessage.Extra["_offset"] = jsonMap["offset"]
	gelfMessage.Extra["_fields"] =jsonMap["fields"]

	return gelfMessage
}

func parseWinlogbeat(jsonMap map[string]interface{}) (gelf.Message) {
	message := jsonMap["message"].(string)
	gelfMessage := CreateMessage(message, jsonMap)

	gelfMessage.Facility = "winlogbeat"

	flatGelf := FlattenPrefixed(jsonMap, "winlogbeat")

	//merge maps
	for k, v := range flatGelf {
		gelfMessage.Extra["_"+k] = v
	}

	return gelfMessage;
}

func parseGenericBeat(jsonMap map[string]interface{}) (gelf.Message) {
	message := jsonMap["message"].(string)
	gelfMessage := CreateMessage(message, jsonMap)
	gelfMessage.Facility = "genericBeat"

	flatGelf := FlattenPrefixed(jsonMap, "beat")
	//merge maps
	for k, v := range flatGelf {
		gelfMessage.Extra["_"+k] = v
	}
	logrus.Info(gelfMessage)

	return gelfMessage;
}

func CreateMessage(message string, jsonMap map[string]interface{}) (gelf.Message){
	beat := jsonMap["beat"]
	hostname := "unknown"
	name := "unknown"
	if beat != nil {
		hostname = beat.(map[string]interface{})["hostname"].(string)
		name = beat.(map[string]interface{})["name"].(string)
	}

	var timestamp float64 = 0
	if jsonMap["@timestamp"] != nil {
		timestampString := jsonMap["@timestamp"].(string)
		t, _ := time.Parse("2006-01-02T15:04:05.000Z", timestampString)
		timestamp = float64(t.Unix())
	}

	typeString := jsonMap["type"].(string)

	extra := map[string]interface{}{
		"_tags": jsonMap["tags"],
		"_name": name,
		"_type": typeString,
	}

	m := gelf.Message{
		Version:  "1.1",
		Host:     hostname,
		Short:    string(message),
		TimeUnix: timestamp,
		Level:    6, // info
		Extra:    extra,
	}

return m
}

func (c *HTTPGelfClient) Connect() (error) {
	return nil
}

func (c *HTTPGelfClient) Close() {

}

// https://github.com/doublerebel/bellows/blob/master/main.go edited to use _ instead of .
func FlattenPrefixed(value interface{}, prefix string) map[string]interface{} {
	m := make(map[string]interface{})
	FlattenPrefixedToResult(value, prefix, m)
	return m
}

func FlattenPrefixedToResult(value interface{}, prefix string, m map[string]interface{}) {
	base := ""
	if prefix != "" {
		base = prefix+"_"
	}

	original := reflect.ValueOf(value)
	kind := original.Kind()
	if kind == reflect.Ptr || kind == reflect.Interface {
		original = reflect.Indirect(original)
		kind = original.Kind()
	}
	t := original.Type()

	switch kind {
	case reflect.Map:
		if t.Key().Kind() != reflect.String {
			break
		}
		for _, childKey := range original.MapKeys() {
			childValue := original.MapIndex(childKey)
			FlattenPrefixedToResult(childValue.Interface(), base+childKey.String(), m)
		}
	case reflect.Struct:
		for i := 0; i < original.NumField(); i += 1 {
			childValue := original.Field(i)
			childKey := t.Field(i).Name
			FlattenPrefixedToResult(childValue.Interface(), base+childKey, m)
		}
	default:
		if prefix != "" {
			m[prefix] = value
		}
	}
}
