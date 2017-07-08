# Telegram Hooks for [Logrus](https://github.com/sirupsen/logrus) <img src="http://i.imgur.com/hTeVwmJ.png" width="40" height="40" alt=":walrus:" class="emoji" title=":walrus:"/>

## Install

```shell
$ go get github.com/Hu13er/telegrus
```

## Usage

```go
package main

import (
	"github.com/sirupsen/logrus"
	"github.com/Hu13er/telegrus"
)

func main() {
	log := logrus.New()
    
    var (
        botToken = ""      // Your Bot token
        chatID = uint32(0) // Your chatID
    )

    hooker := telegrus.NewHooker(botToken, chatID)
    hooker.SetMention(map[logrus.Level][]string{
			logrus.WarnLevel:  []string{"Hu13er", "foobar"},
			logrus.ErrorLevel: []string{"Hu13er"},
			logrus.PanicLevel: []string{"Hu13er"},
		})
    log.Hooks.Add(hooker)

	log.WithFields(logrus.Fields{
		"name": "huber",
		"age":  20,
	}).Error("Hello world!")
}
```

## License
*MIT*