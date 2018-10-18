package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

func (cm *CommandManager) commandKonesyntees(
	args string,
	message discordgo.Message,
	contextual bool,
) (
	context bool,
	err error,
) {
	text, speed, voice, err := parseVoiceParams(message.Content)
	if err != nil {
		return
	}
	if len(text) > 100 {
		return false, errors.New("text too long")
	}

	url := fmt.Sprintf("http://teenus.eki.ee/konesyntees?haal=%d&kiirus=%d&tekst=%s", voice, speed, strings.Replace(text, " ", "%20", -1))
	response, err := http.Get(url)
	if err != nil {
		return
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	url, err = getVoiceMp3Url(string(data))
	if err != nil {
		return
	}

	response, err = http.Get(url)
	if err != nil {
		return
	}

	cm.Discord.ChannelFileSend(message.ChannelID, "konesyntees.mp3", response.Body)
	return
}

func parseVoiceParams(text string) (string, int, int, error) {
	text = text[13:]
	params := strings.Split(text, " ")
	speed := 0
	voice := 0

	for i := 0; i < len(params); i++ {
		if strings.HasPrefix(params[i], "--") == false {
			break
		}

		text = text[len(params[i])+1:]
		pos := strings.IndexByte(params[i], '=')

		if pos == -1 || pos == len(params[i])-1 {
			continue
		}

		param := params[i][strings.IndexByte(params[i], '=')+1:]
		value, err := strconv.Atoi(param)
		param = params[i][strings.Index(params[i], "--")+2 : strings.IndexByte(params[i], '=')]

		if err != nil {
			return "", 0, 0, errors.New("failed to parse param \"" + param + "\"")
		}

		if strings.Compare(param, "speed") == 0 {
			if value < -9 || value > 6 {
				return "", 0, 0, errors.New("speed must be in the range of -9 .. 6")
			}

			speed = value
		} else if strings.Compare(param, "voice") == 0 {
			if value < 0 || value > 3 {
				return "", 0, 0, errors.New("voice must be in the range of 0 .. 3")
			}

			voice = value
		}
	}

	text = strings.TrimSpace(text)
	if len(text) == 0 {
		return "", 0, 0, errors.New("the text can't be empty")
	}

	return text, speed, voice, nil
}

func getVoiceMp3Url(data string) (url string, err error) {
	dec := json.NewDecoder(strings.NewReader(data))
	var urls map[string]string

	err = dec.Decode(&urls)
	if err != nil {
		return
	}

	url, exists := urls["mp3url"]
	if exists == false {
		err = errors.New("the correct URL doesn't exist for some reason")
	}

	return
}
