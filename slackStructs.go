package slack

type (
	// UserInfo - call https://slack.com/api/users.info
	UserInfo struct {
		Ok    bool   `json:"ok"`
		Error string `json:"error"`
		User  User   `json:"user"`
	}

	// ChanelInfo - The short vsrion of info
	// Call - https://slack.com/api/channels.join
	ChanelInfo struct {
		Ok               bool    `json:"ok"`
		Error            bool    `json:"error"`
		Channel          Channel `json:"channel"`
		AlreadyInChannel bool    `json:"already_in_channel"`
	}
	// Channel has all information about a chat channel
	Channel struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		Created    int    `json:"created"`
		Creator    string `json:"creator"`
		IsArchived bool   `json:"is_archived"`
		IsGeneral  bool   `json:"is_general"`
		IsChannel  bool   `json:"is_channel"`
		IsMember   bool   `json:"is_member"`
		HasPins    bool   `json:"has_pins"`
	}
	// User is all info about a user
	User struct {
		ID                string      `json:"id"`
		Name              string      `json:"name"`
		RealName          string      `json:"real_name"`
		Status            interface{} `json:"status"`
		Color             string      `json:"color"`
		Deleted           bool        `json:"deleted"`
		Has2fa            bool        `json:"has_2fa"`
		TwoFactorType     string      `json:"two_factor_type"`
		HasFiles          bool        `json:"has_files"`
		IsAdmin           bool        `json:"is_admin"`
		IsBot             bool        `json:"is_bot"`
		IsOwner           bool        `json:"is_owner"`
		IsPrimaryOwner    bool        `json:"is_primary_owner"`
		IsRestricted      bool        `json:"is_restricted"`
		IsUltraRestricted bool        `json:"is_ultra_restricted"`
		Tz                string      `json:"tz"`
		TzLabel           string      `json:"tz_label"`
		TzOffset          int         `json:"tz_offset"`
		Profile           struct {
			Title              string `json:"title"`
			FirstName          string `json:"first_name"`
			LastName           string `json:"last_name"`
			RealName           string `json:"real_name"`
			RealNameNormalized string `json:"real_name_normalized"`
			Email              string `json:"email"`
			Image24            string `json:"image_24"`
			Image32            string `json:"image_32"`
			Image48            string `json:"image_48"`
			Image72            string `json:"image_72"`
			Image192           string `json:"image_192"`
			Image512           string `json:"image_512"`
		} `json:"profile"`
	}
	// RxMessage - a receive type message
	RxMessage struct {
		Type    string
		Channel Channel
		User    User
		Text    string
		Time    string
	}
)
