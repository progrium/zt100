package irc

import (
	"fmt"
	"log"
	"os"
	"reflect"

	ircx "github.com/nickvanw/ircx/v2"
	"github.com/progrium/zt100/pkg/manifold"
	"gopkg.in/sorcix/irc.v2"
)

type Handler interface {
	HandleMessage(s ircx.Sender, m *irc.Message)
}

type IRCClient struct {
	Server string
	Nick   string
	User   string
	pass   string

	Handler Handler `com:"singleton"`

	bot *ircx.Bot
}

func (c *IRCClient) Mounted(obj manifold.Object) error {
	c.pass = os.Getenv("TWITCH_IRC_TOKEN")
	c.bot = ircx.WithLogin(c.Server, c.Nick, c.User, c.pass)
	if err := c.bot.Connect(); err != nil {
		return err
	}
	c.bot.HandleFunc(irc.RPL_WELCOME, c.OnRegisterConnect)
	c.bot.HandleFunc(irc.PING, c.OnPingHandler)
	c.bot.HandleFunc(irc.PRIVMSG, c.OnMsgHandler)
	go c.bot.HandleLoop()
	return nil
}

func (c *IRCClient) OnRegisterConnect(s ircx.Sender, m *irc.Message) {
	channel := fmt.Sprintf("#%s", c.Nick)
	log.Print("Connected, joining ", channel, " ...")
	s.Send(&irc.Message{
		Command: irc.JOIN,
		Params:  []string{channel},
	})
}

func (c *IRCClient) OnMsgHandler(s ircx.Sender, m *irc.Message) {
	log.Print(m)
	if c.Handler != nil {
		c.Handler.HandleMessage(s, m)
	}
}

func (c *IRCClient) OnPingHandler(s ircx.Sender, m *irc.Message) {
	s.Send(&irc.Message{
		Command: irc.PONG,
		Params:  m.Params,
	})
}

type BangMux struct {
	obj manifold.Object `hash:"ignore"`
}

func (c *BangMux) InitializeComponent(obj manifold.Object) {
	c.obj = obj
}

func (c *BangMux) HandleMessage(s ircx.Sender, m *irc.Message) {
	if m.Trailing()[0] != '!' {
		return
	}
	cmd := m.Trailing()[1:]
	for _, child := range c.obj.Children() {
		if child.Name() != cmd {
			continue
		}
		var handler Handler
		child.ValueTo(reflect.ValueOf(&handler))
		if handler != nil {
			handler.HandleMessage(s, m)
		}
	}
}
