package slack

import "fmt"

func (c Channel) String() string {
	return fmt.Sprintf("#%s [%s]", c.Name, c.ID)
}

func (u User) String() string {
	return fmt.Sprintf("@%s '%s' [%s]", u.Name, u.RealName, u.ID)
}

func (m RxMessage) String() string {
	return fmt.Sprintf("[%s] %v, %v\ntxt='%s'\n", m.Type, m.Channel, m.User, m.Text)
}

func (m Message) String() string {
	return fmt.Sprintf("%v %v : '%s'", m.User, m.Channel, m.Text)
}

func (i BotInfo) String() string {
	return fmt.Sprintf("@%s", i.ID)
}
