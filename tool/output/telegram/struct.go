package telegram

// filesInlineKeyboards set Files Inline Key boards
type filesInlineKeyboards struct {
	GidAndName []map[string]string
	Data       string
}

// functionInlineKeyboards set Files Inline Key boards
type functionInlineKeyboards struct {
	Describe string
	Data     string
}
