package commands

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/kristoisberg/gonesyntees"
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

	response, err := gonesyntees.Request(text, gonesyntees.Voice(voice), speed)

	if err != nil {
		return
	}

	audio, err := http.Get(response.MP3Url)
	if err != nil {
		return
	}

	cm.Discord.ChannelFileSend(message.ChannelID, "konesyntees.mp3", audio.Body)
	return
}

func parseVoiceParams(text string) (string, int, int, error) {
	if len(text) < 14 {
		return "", 0, 0, errors.New("the command must have some sort of params")
	}

	text = text[13:]
	params := strings.Split(text, " ")
	speed := 0
	voice := 0

	for i := 0; i < len(params); i++ {
		if strings.HasPrefix(params[i], "--") == false {
			break
		}

		if len(text) < len(params[i])+2 {
			return "", 0, 0, errors.New("the text can't be empty")
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
			if value < -9 || value > 9 {
				return "", 0, 0, errors.New("speed must be in the range of -9 .. 9")
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
