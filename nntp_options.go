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
	messageID      MessageID
	articleNumber  int
	dotEncodedBody bool
}

// Article number in a newsgroup. The lowest article number is 1. Number 0 is only used for special meanings.
func ArticleNumber(articleNumber int) ArticleOption {
	return func(o *articleOptions) {
		o.articleNumber = articleNumber
	}
}

// The message-id provided can be unwrapped without the "<>" outter pair. The function will rewrap it by calling the
// FullMessageID() function.
func ArticleMessageID(messageID MessageID) ArticleOption {
	return func(o *articleOptions) {
		o.messageID = messageID
	}
}

// Stream raw article body without dot decoding.
func WithDotEncodedBody() ArticleOption {
	return func(o *articleOptions) {
		o.dotEncodedBody = true
	}
}

type OverOption func(*overOptions)

type overOptions struct {
	messageID    MessageID
	articleRange *Range
}

func OverMessageID(messageID MessageID) OverOption {
	return func(o *overOptions) {
		o.messageID = messageID
	}
}

func WithArticleRange(firstNum, lastNum int) OverOption {
	return func(o *overOptions) {
		o.articleRange = &Range{firstNum, lastNum}
	}
}
