package telegram

type User struct {
	Id                      int    `json:id`
	IsBot                   bool   `json:is_bot`
	FirstName               string `json:first_name`
	LastName                string `json:last_name`
	Username                string `json:username`
	LanguageCode            string `json:language_code`
	CanJoinGroups           bool   `json:can_join_groups`
	CanReadAllGroupMessages bool   `json:can_read_all_group_messages`
	SupportsInlineQueries   bool   `json:supports_inline_queries`
}

// TODO: Not finished, add missing fields
type Chat struct {
	Id          int    `json:id`
	Type        string `json:type` // private, group, supergroup, channel
	Title       string `json:title`
	Description string `json:description`
	Bio         string `json:bio`
}

// TODO: Not finished, only required field were added
type Message struct {
	MessageId  int  `json:message_id`
	From       User `json:from`
	SenderChat Chat `json:sender_chat`
	Date       int  `json:date`
	Chat       Chat `json:chat`
}

// TODO: Not finished, only required field were added
type Update struct {
	UpdateId      int     `json:update_id`
	Message       Message `json:message`
	EditedMessage Message `json:edited_message`
}
