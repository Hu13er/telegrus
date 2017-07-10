package telegrus

import (
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	TextFormatter = &logrus.TextFormatter{DisableColors: true}
	JSONFormatter = &logrus.JSONFormatter{}
)

type hooker struct {
	bot       *telegramBot
	MinLevel  logrus.Level
	mention   map[logrus.Level][]string
	formatter logrus.Formatter

	mutex sync.Mutex
}

func NewHooker(botToken string, chatID int64) *hooker {
	return &hooker{
		bot:       newTelegramBot(botToken, chatID, 128),
		MinLevel:  logrus.DebugLevel,
		mention:   make(map[logrus.Level][]string),
		formatter: TextFormatter,
	}
}

func (h *hooker) SetLevel(level logrus.Level) *hooker {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.MinLevel = level
	return h
}

func (h *hooker) SetMention(m map[logrus.Level][]string) *hooker {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.mention = m
	return h
}

func (h *hooker) MentionOn(level logrus.Level, users ...string) *hooker {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, ok := h.mention[level]; !ok {
		h.mention[level] = make([]string, 0)
	}
	h.mention[level] = append(h.mention[level], users...)
	return h
}

func (h *hooker) SetFormatter(formatter logrus.Formatter) *hooker {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.formatter = formatter
	return h
}

func (h *hooker) Fire(entry *logrus.Entry) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	users := []string{}
	for _, user := range h.mention[entry.Level] {
		users = append(users, "@"+user)
	}

	if len(users) > 0 {
		time, level, msg := entry.Time, entry.Level, entry.Message
		entry = entry.WithField("mention", strings.Join(users, ", "))
		entry.Time, entry.Level, entry.Message = time, level, msg
	}

	buf, err := h.formatter.Format(entry)
	if err != nil {
		return err
	}

	h.bot.SendMsg(string(buf))
	return nil
}

func (h *hooker) Levels() []logrus.Level {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	outp := []logrus.Level{}
	for _, level := range logrus.AllLevels {
		if level <= h.MinLevel {
			outp = append(outp, level)
		}
	}

	return outp
}
