package nntp

type GroupOption func(*groupOptions)

type groupOptions struct {
	groupName  string
	groupRange *Range
}

func GroupName(group string) GroupOption {
	return func(o *groupOptions) {
		o.groupName = group
	}
}

func GroupRange(firstNum, lastNum int) GroupOption {
	return func(o *groupOptions) {
		o.groupRange = &Range{firstNum, lastNum}
	}
}

type ArticleOption func(*articleOptions)

type articleOptions struct {
	messageID     string
	articleNumber int
}

// Article number in a newsgroup. The lowest article number is 1. Number 0 is only used for special meanings.
func ArticleNumber(articleNumber int) ArticleOption {
	return func(o *articleOptions) {
		o.articleNumber = articleNumber
	}
}

// The message-id provided can be unwrapped without the "<>" outter pair. The function will rewrap it by calling the
// FullMessageID() function.
func MessageID(messageID string) ArticleOption {
	return func(o *articleOptions) {
		o.messageID = FullMessageID(messageID)
	}
}

type OverOption func(*overOptions)

type overOptions struct {
	messageID    string
	articleRange *Range
}

func OverMessageID(messageID string) OverOption {
	return func(o *overOptions) {
		o.messageID = messageID
	}
}

func OverArticleRange(firstNum, lastNum int) OverOption {
	return func(o *overOptions) {
		o.articleRange = &Range{firstNum, lastNum}
	}
}
