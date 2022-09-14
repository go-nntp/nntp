package nntp

import (
	"fmt"
	"io"
	"net/textproto"
	"strconv"
	"strings"
	"time"

	"gopkg.in/option.v0"
)

/**
 * Returns servers capabilities
 *
 * @return mixed (array) list of capabilities
 */
func (conn *Conn) CmdCapabilities() (capabilities []string, err error) {
	if err = conn.PrintfLine("CAPABILITIES"); err != nil {
		err = fmt.Errorf("[nntp.CmdCapabilities] failed to send CAPABILITIES command: %w", err)
		return
	}
	code, msg, err := conn.ReadCodeLine(0)
	if err != nil {
		err = fmt.Errorf("[nntp.CmdCapabilities] failed to read CAPABILITIES response: %w", err)
		return
	}
	switch code {
	case ResponseCodeCapabilitiesFollow: // 101
		if capabilities, err = conn.ReadDotLines(); err != nil {
			err = fmt.Errorf("[nntp.CmdCapabilities] failed to read CAPABILITIES response body: %w", err)
			return
		}
	default:
		err = fmt.Errorf("[nntp.CmdCapabilities] unexpected response: %w", &Error{code, msg})
	}
	return
}

/**
 * Tell the news server we want an article.
 *
 * @return true when posting allowed, false when posting disallowed.
 */
func (conn *Conn) CmdModeReader() (postingAllowed bool, err error) {
	if err = conn.PrintfLine("MODE READER"); err != nil {
		err = fmt.Errorf("[nntp.CmdModeReader] failed to send MODE READER command: %w", err)
		return
	}
	code, msg, err := conn.ReadCodeLine(0)
	if err != nil {
		err = fmt.Errorf("[nntp.CmdModeReader] failed to read MODE READER response: %w", err)
		return
	}
	switch code {
	case ResponseCodeReadyPostingAllowed: // 200
		postingAllowed = true
	case ResponseCodeReadyPostingProhibited: // 201
		postingAllowed = false
	default:
		err = fmt.Errorf("[nntp.CmdModeReader] unexpected response: %w", &Error{code, msg})
	}
	return
}

/**
 * Disconnect from the NNTP server.
 */
func (conn *Conn) CmdQuit() (err error) {
	if err = conn.PrintfLine("QUIT"); err != nil {
		err = fmt.Errorf("[nntp.CmdQuit] failed to send QUIT command: %w", err)
		return
	}
	code, msg, err := conn.ReadCodeLine(0)
	if err != nil {
		err = fmt.Errorf("[nntp.CmdQuit] failed to read QUIT response: %w", err)
		return
	}
	switch code {
	case ResponseCodeDisconnectingRequested: // 205
		if err = conn.Close(); err != nil {
			err = fmt.Errorf("[nntp.CmdQuit] failed to close connection: %w", err)
			return
		}
	default:
		err = fmt.Errorf("[nntp.CmdQuit] unexpected response: %w", &Error{code, msg})
	}
	return
}

/**
 * Selects a news group (issue a GROUP command to the server)
 *
 * @param string $newsgroup The newsgroup name
 *
 * @return groupinfo on success
 */
func (conn *Conn) CmdGroup(newsgroup string) (groupinfo *GroupStat, err error) {
	if err = conn.PrintfLine("GROUP %s", newsgroup); err != nil {
		err = fmt.Errorf("[nntp.CmdGroup] failed to send GROUP command: %w", err)
		return
	}
	code, msg, err := conn.ReadCodeLine(0)
	if err != nil {
		err = fmt.Errorf("[nntp.CmdGroup] failed to read GROUP response: %w", err)
		return
	}
	switch code {
	case ResponseCodeGroupSelected: // 211
		info := &GroupStat{}
		if _, err = fmt.Sscanf(msg, "%d %d %d %s", &info.Count, &info.First, &info.Last, &info.Group); err != nil {
			err = fmt.Errorf("[nntp.CmdGroup] failed to parse group info response: %#v: %w", msg, ErrorParsingResponse)
			return
		}
		groupinfo = info
	default:
		err = fmt.Errorf("[nntp.CmdGroup] unexpected response: %w", &Error{code, msg})
	}
	return
}

/**
 * @param optional string $newsgroup
 * @param optional mixed $range
 *
 * @return optional mixed (array) on success or (object) pear_error on failure
 * @access protected
 */
func (conn *Conn) CmdListGroup(options ...GroupOption) (groupinfo *GroupStat, articles []int, err error) {
	opts := option.New(options)
	if opts.groupName == "" {
		err = conn.PrintfLine("LISTGROUP")
	} else if opts.groupRange == nil {
		err = conn.PrintfLine("LISTGROUP %s", opts.groupName)
	} else {
		err = conn.PrintfLine("LISTGROUP %s %s", opts.groupName, opts.groupRange.String())
	}
	if err != nil {
		err = fmt.Errorf("[nntp.CmdListGroup] failed to send LISTGROUP command: %w", err)
		return
	}
	code, msg, err := conn.ReadCodeLine(0)
	if err != nil {
		err = fmt.Errorf("[nntp.CmdListGroup] failed to read LISTGROUP response: %w", err)
		return
	}
	switch code {
	case ResponseCodeGroupSelected: // 211
		info := &GroupStat{}
		if _, err = fmt.Sscanf(msg, "%d %d %d %s", &info.Count, &info.First, &info.Last, &info.Group); err != nil {
			err = fmt.Errorf("[nntp.CmdListGroup] failed to parse group info response: %#v: %w", msg, ErrorParsingResponse)
			return
		}
		var lines []string
		if lines, err = conn.ReadDotLines(); err != nil {
			err = fmt.Errorf("[nntp.CmdListGroup] failed to read LISTGROUP response body: %w", err)
			return
		}
		ids := make([]int, len(lines))
		for i, line := range lines {
			if ids[i], err = strconv.Atoi(line); err != nil {
				err = fmt.Errorf("[nntp.CmdListGroup] failed to parse article ID: %#v: %w", line, err)
				return
			}
		}
		groupinfo, articles = info, ids
	default:
		err = fmt.Errorf("[nntp.CmdListGroup] unexpected response: %w", &Error{code, msg})
	}
	return
}

// If the currently selected newsgroup is valid, the current article * number MUST be set to the previous article in
// that newsgroup (that * is, the highest existing article number less than the current article * number).  If
// successful, a response indicating the new current * article number and the message-id of that article MUST be
// returned. * No article text is sent in response to this command.
func (conn *Conn) CmdLast(newsgroup string) (err error) {
	if err = conn.PrintfLine("LAST %s", newsgroup); err != nil {
		err = fmt.Errorf("[nntp.CmdLast] failed to send LAST command: %w", err)
		return
	}
	code, msg, err := conn.ReadCodeLine(0)
	if err != nil {
		err = fmt.Errorf("[nntp.CmdLast] failed to read LAST response: %w", err)
		return
	}
	switch code {
	case ResponseCodeArticleSelected: // 223
	default:
		err = fmt.Errorf("[nntp.CmdLast] unexpected response: %w", &Error{code, msg})
	}
	return
}

// If the currently selected newsgroup is valid, the current article number MUST be set to the next article in that
// newsgroup (that is, the lowest existing article number greater than the current article number).  If successful, a
// response indicating the new current article number and the message-id of that article MUST be returned. No article
// text is sent in response to this command.
func (conn *Conn) CmdNext(newsgroup string) (err error) {
	if err = conn.PrintfLine("NEXT %s", newsgroup); err != nil {
		err = fmt.Errorf("[nntp.CmdNext] failed to send NEXT command: %w", err)
		return
	}
	code, msg, err := conn.ReadCodeLine(0)
	if err != nil {
		err = fmt.Errorf("[nntp.CmdNext] failed to read NEXT response: %w", err)
		return
	}
	switch code {
	case ResponseCodeArticleSelected: // 223
	default:
		err = fmt.Errorf("[nntp.CmdNext] unexpected response: %w", &Error{code, msg})
	}
	return
}

// Selects and presents the entire article.
func (conn *Conn) CmdArticle(options ...ArticleOption) (article *Article, err error) {
	opts := option.New(options)
	if opts.messageID != "" {
		err = conn.PrintfLine("ARTICLE %s", opts.messageID)
	} else if opts.articleNumber != 0 {
		err = conn.PrintfLine("ARTICLE %d", opts.articleNumber)
	} else {
		err = conn.PrintfLine("ARTICLE")
	}
	if err != nil {
		err = fmt.Errorf("[nntp.CmdArticle] failed to send ARTICLE command: %w", err)
		return
	}
	code, msg, err := conn.ReadCodeLine(0)
	if err != nil {
		err = fmt.Errorf("[nntp.CmdArticle] failed to read ARTICLE response: %w", err)
		return
	}
	switch code {
	case ResponseCodeArticleFollows: // 220
		article = new(Article)
		if _, err = fmt.Sscanf(msg, "%d %s", &article.ArticleNumber, &article.MessageID); err != nil {
			err = fmt.Errorf("[nntp.CmdArticle] failed to parse ARTICLE command status line: %#v: %w", msg, ErrorParsingResponse)
			return
		}
		if article.Header, err = conn.ReadMIMEHeader(); err != nil {
			err = fmt.Errorf("[nntp.CmdArticle] failed to parse MIME header: %#v: %w", msg, ErrorParsingResponse)
			return
		}
		article.Body = conn.DotReader()
	default:
		err = fmt.Errorf("[nntp.CmdArticle] unexpected response: %w", &Error{code, msg})
	}
	return
}

func (conn *Conn) CmdHead(options ...ArticleOption) (article *Article, err error) {
	opts := option.New(options)
	if opts.messageID != "" {
		err = conn.PrintfLine("HEAD %s", opts.messageID)
	} else {
		err = conn.PrintfLine("HEAD %d", opts.articleNumber)
	}
	if err != nil {
		err = fmt.Errorf("[nntp.CmdHead] failed to send HEAD command: %w", err)
		return
	}
	code, msg, err := conn.ReadCodeLine(0)
	if err != nil {
		err = fmt.Errorf("[nntp.CmdHead] failed to read HEAD response: %w", err)
		return
	}
	switch code {
	case ResponseCodeHeadFollows: // 221
		article = new(Article)
		if _, err = fmt.Sscanf(msg, "%d %s", &article.ArticleNumber, &article.MessageID); err != nil {
			err = fmt.Errorf("[nntp.CmdHead] failed to parse HEAD response: %#v: %w", msg, ErrorParsingResponse)
			return
		}
		if article.Header, err = conn.ReadMIMEHeader(); err != nil {
			err = fmt.Errorf("[nntp.CmdArticle] failed to parse MIME header: %#v: %w", msg, ErrorParsingResponse)
			return
		}
	default:
		err = fmt.Errorf("[nntp.CmdHead] unexpected response: %w", &Error{code, msg})
	}
	return
}

func (conn *Conn) CmdBody(options ...ArticleOption) (article *Article, err error) {
	opts := option.New(options)
	if opts.messageID != "" {
		err = conn.PrintfLine("BODY %s", opts.messageID)
	} else if opts.articleNumber != 0 {
		err = conn.PrintfLine("BODY %d", opts.articleNumber)
	} else {
		err = conn.PrintfLine("BODY")
	}
	if err != nil {
		err = fmt.Errorf("[nntp.CmdBody] failed to send BODY command: %w", err)
		return
	}
	code, msg, err := conn.ReadCodeLine(0)
	if err != nil {
		err = fmt.Errorf("[nntp.CmdBody] failed to read BODY response: %w", err)
		return
	}
	switch code {
	case ResponseCodeBodyFollows: // 222
		article = new(Article)
		if _, err = fmt.Sscanf(msg, "%d %s", &article.ArticleNumber, &article.MessageID); err != nil {
			err = fmt.Errorf("[nntp.CmdBody] failed to parse BODY command status line: %#v: %w", msg, ErrorParsingResponse)
			return
		}
		article.Body = conn.DotReader()
	default:
		err = fmt.Errorf("[nntp.CmdBody] unexpected response: %w", &Error{code, msg})
	}
	return
}

func (conn *Conn) CmdStat(options ...ArticleOption) (article *Article, err error) {
	opts := option.New(options)
	if opts.messageID != "" {
		err = conn.PrintfLine("STAT %s", opts.messageID)
	} else if opts.articleNumber != 0 {
		err = conn.PrintfLine("STAT %d", opts.articleNumber)
	} else {
		err = conn.PrintfLine("STAT")
	}
	if err != nil {
		err = fmt.Errorf("[nntp.CmdStat] failed to send STAT command: %w", err)
		return
	}
	code, msg, err := conn.ReadCodeLine(0)
	if err != nil {
		err = fmt.Errorf("[nntp.CmdStat] failed to read STAT response: %w", err)
		return
	}
	switch code {
	case ResponseCodeArticleSelected: // 223
		article = new(Article)
		if _, err = fmt.Sscanf(msg, "%d %s", &article.ArticleNumber, &article.MessageID); err != nil {
			err = fmt.Errorf("[nntp.CmdStat] failed to parse STAT command status line: %#v: %w", msg, ErrorParsingResponse)
			return
		}
	default:
		err = fmt.Errorf("[nntp.CmdStat] unexpected response: %w", &Error{code, msg})
	}
	return
}

func (conn *Conn) CmdPost(article *Article) (err error) {
	if err = conn.PrintfLine("POST"); err != nil {
		err = fmt.Errorf("[nntp.CmdPost] failed to send POST command: %w", err)
		return
	}
	code, msg, err := conn.ReadCodeLine(0)
	if err != nil {
		err = fmt.Errorf("[nntp.CmdPost] failed to read POST response: %w", err)
		return
	}
	switch code {
	case ResponseCodePostingSend: // 340
		err = nil
	default:
		err = fmt.Errorf("[nntp.CmdPost] unexpected response: %w", &Error{code, msg})
		return
	}
	if article.MessageID != "" {
		article.Header.Set("Message-Id", article.MessageID)
	}
	for key, values := range article.Header {
		for _, value := range values {
			if err = conn.PrintfLine("%s: %s", textproto.CanonicalMIMEHeaderKey(key), value); err != nil {
				err = fmt.Errorf("[nntp.CmdPost] failed to send article header: %w", err)
				return
			}
		}
	}
	if err = conn.PrintfLine(""); err != nil {
		err = fmt.Errorf("[nntp.CmdPost] failed to send article header termination line: %w", err)
		return
	}
	writer := conn.DotWriter()
	if _, err = io.Copy(writer, article.Body); err != nil {
		err = fmt.Errorf("[nntp.CmdPost] failed to send article body: %w", err)
		return
	}
	if err = writer.Close(); err != nil {
		err = fmt.Errorf("[nntp.CmdPost] failed to send article body termination sequence: %w", err)
		return
	}
	code, msg, err = conn.ReadCodeLine(0)
	if err != nil {
		err = fmt.Errorf("[nntp.CmdPost] failed to read POST result: %w", err)
		return
	}
	switch code {
	case ResponseCodePostingSuccess: // 240
		err = nil
	default:
		err = fmt.Errorf("[nntp.CmdPost] unexpected response: %w", &Error{code, msg})
	}
	return
}

func (conn *Conn) CmdIHave(article *Article) (err error) {
	if err = conn.PrintfLine("IHAVE"); err != nil {
		err = fmt.Errorf("[nntp.CmdIHave] failed to send IHAVE command: %w", err)
		return
	}
	code, msg, err := conn.ReadCodeLine(0)
	if err != nil {
		err = fmt.Errorf("[nntp.CmdIHave] failed to read IHAVE response: %w", err)
		return
	}
	switch code {
	case ResponseCodeTransferSend: // 335
		err = nil
	default:
		err = fmt.Errorf("[nntp.CmdIHave] unexpected response: %w", &Error{code, msg})
		return
	}
	if article.MessageID != "" {
		article.Header.Set("Message-Id", article.MessageID)
	}
	for key, values := range article.Header {
		for _, value := range values {
			if err = conn.PrintfLine("%s: %s", textproto.CanonicalMIMEHeaderKey(key), value); err != nil {
				err = fmt.Errorf("[nntp.CmdIHave] failed to send article header: %w", err)
				return
			}
		}
	}
	if err = conn.PrintfLine(""); err != nil {
		err = fmt.Errorf("[nntp.CmdIHave] failed to send article header termination line: %w", err)
		return
	}
	writer := conn.DotWriter()
	if _, err = io.Copy(writer, article.Body); err != nil {
		err = fmt.Errorf("[nntp.CmdIHave] failed to send article body: %w", err)
		return
	}
	if err = writer.Close(); err != nil {
		err = fmt.Errorf("[nntp.CmdIHave] failed to send article body termination sequence: %w", err)
		return
	}
	code, msg, err = conn.ReadCodeLine(0)
	if err != nil {
		err = fmt.Errorf("[nntp.CmdIHave] failed to read IHAVE result: %w", err)
		return
	}
	switch code {
	case ResponseCodePostingSuccess: // 240
		err = nil
	default:
		err = fmt.Errorf("[nntp.CmdIHave] unexpected response: %w", &Error{code, msg})
	}
	return
}

// Get the date from the newsserver format of returned date
func (conn *Conn) CmdDate() (date time.Time, err error) {
	if err = conn.PrintfLine("DATE"); err != nil {
		err = fmt.Errorf("[nntp.CmdDate] failed to send DATE command: %w", err)
		return
	}
	code, msg, err := conn.ReadCodeLine(0)
	if err != nil {
		err = fmt.Errorf("[nntp.CmdDate] failed to read DATE response: %w", err)
		return
	}
	switch code {
	case ResponseCodeServerDate: // 223
		if date, err = time.Parse("20060102150405", msg); err != nil {
			err = fmt.Errorf("[nntp.CmdDate] failed to parse date string %#v: %w", msg, err)
			return
		}
	default:
		err = fmt.Errorf("[nntp.CmdDate] unexpected response: %w", &Error{code, msg})
	}
	return
}

// Returns the server's help text
func (conn *Conn) CmdHelp() (help io.Reader, err error) {
	if err = conn.PrintfLine("HELP"); err != nil {
		err = fmt.Errorf("[nntp.CmdHelp] failed to send HELP command: %w", err)
		return
	}
	code, msg, err := conn.ReadCodeLine(0)
	if err != nil {
		err = fmt.Errorf("[nntp.CmdHelp] failed to read HELP response: %w", err)
		return
	}
	switch code {
	case ResponseCodeHelpFollows: // 100
		help = conn.DotReader()
	default:
		err = fmt.Errorf("[nntp.CmdHelp] unexpected response: %w", &Error{code, msg})
	}
	return
}

// Fetches a list of all newsgroups created since a specified date
func (conn *Conn) CmdNewGroups(date time.Time, useGMT bool) (groups []GroupListItem, err error) {
	datestring := date.Format("20060102 150405")
	if useGMT {
		err = conn.PrintfLine("NEWGROUPS %s GMT", datestring)
	} else {
		err = conn.PrintfLine("NEWGROUPS %s", datestring)
	}
	if err != nil {
		err = fmt.Errorf("[nntp.CmdNewGroups] failed to send NEWGROUPS command: %w", err)
		return
	}
	code, msg, err := conn.ReadCodeLine(0)
	if err != nil {
		err = fmt.Errorf("[nntp.CmdNewGroups] failed to read NEWGROUPS response: %w", err)
		return
	}
	switch code {
	case ResponseCodeNewGroupsFollow: // 231
		var lines []string
		if lines, err = conn.ReadDotLines(); err != nil {
			err = fmt.Errorf("[nntp.CmdNewGroups] failed to read NEWGROUPS response body: %w", err)
			return
		}
		groups = make([]GroupListItem, len(lines))
		for i, line := range lines {
			group := &groups[i]
			if _, err = fmt.Sscanf(line, "%s %d %d %c", &group.Group, &group.Last, &group.First, &group.Permission); err != nil {
				err = fmt.Errorf("[nntp.CmdNewGroups] failed to parse group list item %#v: %w", line, err)
				return
			}
		}
	default:
		err = fmt.Errorf("[nntp.CmdNewGroups] unexpected response: %w", &Error{code, msg})
	}
	return
}

// Fetches a list of message-ids of articles posted or received on the server, in the newsgroups whose names match the
// wildmat, since the specified date and time.
func (conn *Conn) CmdNewNews(wildmat string, date time.Time, useGMT bool) (messageIds []string, err error) {
	if wildmat == "" {
		err = fmt.Errorf("[nntp.CmdNewNews] empty wildmat: %w", ErrorInvalidParams)
		return
	}
	datestring := date.Format("20060102 150405")
	if useGMT {
		err = conn.PrintfLine("NEWNEWS %s %s GMT", wildmat, datestring)
	} else {
		err = conn.PrintfLine("NEWNEWS %s %s", wildmat, datestring)
	}
	if err != nil {
		err = fmt.Errorf("[nntp.CmdNewNews] failed to send NEWNEWS command: %w", err)
		return
	}
	code, msg, err := conn.ReadCodeLine(0)
	if err != nil {
		err = fmt.Errorf("[nntp.CmdNewNews] failed to read NEWNEWS response: %w", err)
		return
	}
	switch code {
	case ResponseCodeNewArticlesFollow: // 230
		if messageIds, err = conn.ReadDotLines(); err != nil {
			err = fmt.Errorf("[nntp.CmdNewNews] failed to read NEWNEWS response body: %w", err)
			return
		}
	default:
		err = fmt.Errorf("[nntp.CmdNewNews] unexpected response: %w", &Error{code, msg})
	}
	return
}

// Fetches a list of all avaible newsgroups.
func (conn *Conn) CmdList() (groups []GroupListItem, err error) {
	if err = conn.PrintfLine("LIST"); err != nil {
		err = fmt.Errorf("[nntp.CmdList] failed to send LIST command: %w", err)
		return
	}
	code, msg, err := conn.ReadCodeLine(0)
	if err != nil {
		err = fmt.Errorf("[nntp.CmdList] failed to read LIST response: %w", err)
		return
	}
	switch code {
	case ResponseCodeGroupsFollow: // 215
		var lines []string
		if lines, err = conn.ReadDotLines(); err != nil {
			err = fmt.Errorf("[nntp.CmdList] failed to read LIST response body: %w", err)
			return
		}
		groups = make([]GroupListItem, len(lines))
		for i, line := range lines {
			group := &groups[i]
			if _, err = fmt.Sscanf(line, "%s %d %d %c", &group.Group, &group.Last, &group.First, &group.Permission); err != nil {
				err = fmt.Errorf("[nntp.CmdList] failed to parse group list item %#v: %w", line, err)
				return
			}
		}
	default:
		err = fmt.Errorf("[nntp.CmdList] unexpected response: %w", &Error{code, msg})
	}
	return
}

// Fetches a list of (all) avaible newsgroups.
func (conn *Conn) CmdListActive(wildmat string) (groups []GroupListItem, err error) {
	if wildmat != "" {
		err = conn.PrintfLine("LIST ACTIVE %s", wildmat)
	} else {
		err = conn.PrintfLine("LIST ACTIVE")
	}
	if err != nil {
		err = fmt.Errorf("[nntp.CmdListActive] failed to send LIST ACTIVE command: %w", err)
		return
	}
	code, msg, err := conn.ReadCodeLine(0)
	if err != nil {
		err = fmt.Errorf("[nntp.CmdListActive] failed to read LIST ACTIVE response: %w", err)
		return
	}
	switch code {
	case ResponseCodeGroupsFollow: // 215
		var lines []string
		if lines, err = conn.ReadDotLines(); err != nil {
			err = fmt.Errorf("[nntp.CmdListActive] failed to read LIST ACTIVE response body: %w", err)
			return
		}
		groups = make([]GroupListItem, len(lines))
		for i, line := range lines {
			group := &groups[i]
			if _, err = fmt.Sscanf(line, "%s %d %d %c", &group.Group, &group.Last, &group.First, &group.Permission); err != nil {
				err = fmt.Errorf("[nntp.CmdListActive] failed to parse group list item %#v: %w", line, err)
				return
			}
		}
	default:
		err = fmt.Errorf("[nntp.CmdListActive] unexpected response: %w", &Error{code, msg})
	}
	return
}

// Fetches a list of (all) avaible newsgroup descriptions.
func (conn *Conn) CmdListNewsgroups(wildmat string) (groups []GroupDescriptionListItem, err error) {
	if wildmat != "" {
		err = conn.PrintfLine("LIST NEWSGROUPS %s", wildmat)
	} else {
		err = conn.PrintfLine("LIST NEWSGROUPS")
	}
	if err != nil {
		err = fmt.Errorf("[nntp.CmdListNewsgroups] failed to send LIST NEWSGROUPS command: %w", err)
		return
	}
	code, msg, err := conn.ReadCodeLine(0)
	if err != nil {
		err = fmt.Errorf("[nntp.CmdListNewsgroups] failed to read LIST NEWSGROUPS response: %w", err)
		return
	}
	switch code {
	case ResponseCodeGroupsFollow: // 215
		var lines []string
		if lines, err = conn.ReadDotLines(); err != nil {
			err = fmt.Errorf("[nntp.CmdListNewsgroups] failed to read LIST NEWSGROUPS response body: %w", err)
			return
		}
		groups = make([]GroupDescriptionListItem, len(lines))
		for i, line := range lines {
			group := &groups[i]
			fields := strings.Fields(line)
			if len(fields) < 2 {
				err = fmt.Errorf("[nntp.CmdListNewsgroups] invalid group description line %#v: %w", line, ErrorParsingResponse)
				return
			}
			group.Group = fields[0]
			group.Description = strings.Join(fields[1:], " ")
		}
	default:
		err = fmt.Errorf("[nntp.CmdListNewsgroups] unexpected response: %w", &Error{code, msg})
	}
	return
}

// Fetches message header for specified articles.
func (conn *Conn) CmdOver(options ...OverOption) (articles []ArticleOverview, err error) {
	opts := option.New(options)
	if opts.messageID != "" {
		err = conn.PrintfLine("OVER %s", opts.messageID)
	} else if opts.articleRange != nil {
		err = conn.PrintfLine("OVER %s", opts.articleRange.String())
	} else {
		err = conn.PrintfLine("OVER")
	}
	if err != nil {
		err = fmt.Errorf("[nntp.CmdOver] failed to send OVER command: %w", err)
		return
	}
	code, msg, err := conn.ReadCodeLine(0)
	if err != nil {
		err = fmt.Errorf("[nntp.CmdOver] failed to read OVER response: %w", err)
		return
	}
	switch code {
	case ResponseCodeOverviewFollows: // 224
		var lines []string
		if lines, err = conn.ReadDotLines(); err != nil {
			err = fmt.Errorf("[nntp.CmdOver] failed to read OVER response body: %w", err)
			return
		}
		articles = make([]ArticleOverview, len(lines))
		for i, line := range lines {
			article := &articles[i]
			fields := strings.Split(line, "\t")
			if len(fields) < 8 {
				err = fmt.Errorf("[nntp.CmdOver] invalid group description line %#v: %w", line, ErrorParsingResponse)
				return
			}
			if article.ArticleNumber, err = strconv.Atoi(fields[0]); err != nil {
				err = fmt.Errorf("[nntp.CmdOver] failed to parse article number %#v: %w", fields[0], ErrorParsingResponse)
				return
			}
			article.Subject, article.From, article.MessageID, article.References = fields[1], fields[2], fields[4], fields[5]
			if article.Date, err = time.Parse(DefaultArticleDateLayout, fields[3]); err != nil {
				err = fmt.Errorf("[nntp.CmdOver] failed to parse article date %#v: %w", fields[3], ErrorParsingResponse)
				return
			}
			if article.Bytes, err = strconv.ParseUint(fields[6], 10, 64); err != nil {
				err = fmt.Errorf("[nntp.CmdOver] failed to parse article bytes count %#v: %w", fields[6], ErrorParsingResponse)
				return
			}
			if article.Lines, err = strconv.ParseUint(fields[7], 10, 64); err != nil {
				err = fmt.Errorf("[nntp.CmdOver] failed to parse article lines count %#v: %w", fields[7], ErrorParsingResponse)
				return
			}
			article.ExtraFields = fields[8:]
		}
	default:
		err = fmt.Errorf("[nntp.CmdOver] unexpected response: %w", &Error{code, msg})
	}
	return
}

// Fetches message header for specified articles.
func (conn *Conn) CmdXOver(options ...OverOption) (articles []ArticleOverview, err error) {
	opts := option.New(options)
	if opts.messageID != "" {
		err = conn.PrintfLine("XOVER %s", opts.messageID)
	} else if opts.articleRange != nil {
		err = conn.PrintfLine("XOVER %s", opts.articleRange.String())
	} else {
		err = conn.PrintfLine("XOVER")
	}
	if err != nil {
		err = fmt.Errorf("[nntp.CmdXOver] failed to send XOVER command: %w", err)
		return
	}
	code, msg, err := conn.ReadCodeLine(0)
	if err != nil {
		err = fmt.Errorf("[nntp.CmdXOver] failed to read XOVER response: %w", err)
		return
	}
	switch code {
	case ResponseCodeOverviewFollows: // 224
		var lines []string
		if lines, err = conn.ReadDotLines(); err != nil {
			err = fmt.Errorf("[nntp.CmdXOver] failed to read XOVER response body: %w", err)
			return
		}
		articles = make([]ArticleOverview, len(lines))
		for i, line := range lines {
			article := &articles[i]
			fields := strings.Split(line, "\t")
			if len(fields) < 8 {
				err = fmt.Errorf("[nntp.CmdXOver] invalid group description line %#v: %w", line, ErrorParsingResponse)
				return
			}
			if article.ArticleNumber, err = strconv.Atoi(fields[0]); err != nil {
				err = fmt.Errorf("[nntp.CmdXOver] failed to parse article number %#v: %w", fields[0], ErrorParsingResponse)
				return
			}
			article.Subject, article.From, article.MessageID, article.References = fields[1], fields[2], fields[4], fields[5]
			if article.Date, err = time.Parse(DefaultArticleDateLayout, fields[3]); err != nil {
				err = fmt.Errorf("[nntp.CmdXOver] failed to parse article date %#v: %w", fields[3], ErrorParsingResponse)
				return
			}
			if article.Bytes, err = strconv.ParseUint(fields[6], 10, 64); err != nil {
				err = fmt.Errorf("[nntp.CmdXOver] failed to parse article bytes count %#v: %w", fields[6], ErrorParsingResponse)
				return
			}
			if article.Lines, err = strconv.ParseUint(fields[7], 10, 64); err != nil {
				err = fmt.Errorf("[nntp.CmdXOver] failed to parse article lines count %#v: %w", fields[7], ErrorParsingResponse)
				return
			}
			article.ExtraFields = fields[8:]
		}
	default:
		err = fmt.Errorf("[nntp.CmdXOver] unexpected response: %w", &Error{code, msg})
	}
	return
}

// Fetches description of the fields returned in OVER or XOVER command.
func (conn *Conn) CmdListOverviewFmt() (fields []OverviewFieldFormat, err error) {
	if err = conn.PrintfLine("LIST OVERVIEW.FMT"); err != nil {
		err = fmt.Errorf("[nntp.CmdListOverviewFmt] failed to send LIST OVERVIEW.FMT command: %w", err)
		return
	}
	code, msg, err := conn.ReadCodeLine(0)
	if err != nil {
		err = fmt.Errorf("[nntp.CmdListOverviewFmt] failed to read LIST OVERVIEW.FMT response: %w", err)
		return
	}
	switch code {
	case ResponseCodeInformationFollows: // 215
		var lines []string
		if lines, err = conn.ReadDotLines(); err != nil {
			err = fmt.Errorf("[nntp.CmdListOverviewFmt] failed to read LIST OVERVIEW.FMT response body: %w", err)
			return
		}
		if len(lines) < 7 {
			err = fmt.Errorf("[nntp.CmdListOverviewFmt] overview format must has at least 7 fields: %w", ErrorParsingResponse)
			return
		}
		fields = make([]OverviewFieldFormat, len(lines))
		fixedFields := []string{
			"Subject",
			"From",
			"Date",
			"Message-ID",
			"References",
			"Bytes",
			"Lines",
		}
		for i, fixedField := range fixedFields {
			if strings.EqualFold(lines[i], fixedField+":") {
				fields[i].Name, fields[i].Type = fixedField, ShortHeaderOverviewField
			} else if i == 5 && strings.EqualFold(lines[i], ":bytes") {
				fields[i].Name, fields[i].Type = fixedField, MetadataOverviewField
			} else if i == 6 && strings.EqualFold(lines[i], ":lines") {
				fields[i].Name, fields[i].Type = fixedField, MetadataOverviewField
			} else {
				err = fmt.Errorf("[nntp.CmdListOverviewFmt] expecting fixed %s field instead got %#v: %w", fixedField, lines[i], ErrorParsingResponse)
				return
			}
		}
		for i, line := range lines[7:] {
			field := &fields[i+7]
			parts := strings.SplitN(line, ":", 2)
			if len(parts) != 2 {
				err = fmt.Errorf("[nntp.CmdListOverviewFmt] field with empty name %#v: %w", line, ErrorParsingResponse)
				return
			}
			if parts[0] == "" {
				if parts[1] == "" {
					err = fmt.Errorf("[nntp.CmdListOverviewFmt] metadata field with empty name %#v: %w", line, ErrorParsingResponse)
					return
				}
				field.Name, field.Type = parts[1], MetadataOverviewField
			} else if strings.EqualFold(parts[1], "full") {
				field.Name, field.Type = parts[0], FullHeaderOverviewField
			} else {
				field.Name, field.Type = parts[0], ShortHeaderOverviewField
			}
		}
	default:
		err = fmt.Errorf("[nntp.CmdListOverviewFmt] unexpected response: %w", &Error{code, msg})
	}
	return
}

// Fetches list of fields that may be retrieved using the HDR command.
func (conn *Conn) CmdListHeaders() (anyField bool, fields []string, err error) {
	if err = conn.PrintfLine("LIST HEADERS"); err != nil {
		err = fmt.Errorf("[nntp.CmdListHeaders] failed to send LIST HEADERS command: %w", err)
		return
	}
	code, msg, err := conn.ReadCodeLine(0)
	if err != nil {
		err = fmt.Errorf("[nntp.CmdListHeaders] failed to read LIST HEADERS response: %w", err)
		return
	}
	switch code {
	case ResponseCodeInformationFollows: // 215
		var lines []string
		if lines, err = conn.ReadDotLines(); err != nil {
			err = fmt.Errorf("[nntp.CmdListHeaders] failed to read LIST HEADERS response body: %w", err)
			return
		}
		if len(lines) > 0 {
			if lines[0] == ":" {
				anyField = true
				fields = lines[1:]
			} else {
				fields = lines
			}
		}
	default:
		err = fmt.Errorf("[nntp.CmdListHeaders] unexpected response: %w", &Error{code, msg})
	}
	return
}

/**
 * Authenticate using 'original' method
 *
 * @param string $user The username to authenticate as.
 * @param string $pass The password to authenticate with.
 */
func (conn *Conn) CmdAuthinfo(user, pass string) (err error) {
	if user == "" {
		err = fmt.Errorf("[nntp.CmdAuthinfo] empty username: %w", ErrorInvalidParams)
		return
	}

	// send the username
	if err = conn.PrintfLine("AUTHINFO user %s", user); err != nil {
		err = fmt.Errorf("[nntp.CmdAuthinfo] failed to send AUTHINFO user command: %w", err)
		return
	}
	code, msg, err := conn.ReadCodeLine(0)
	if err != nil {
		err = fmt.Errorf("[nntp.CmdAuthinfo] failed to read AUTHINFO user response: %w", err)
		return
	}
	if code == ResponseCodeAuthenticationContinue /* 381 */ {
		if pass == "" {
			err = fmt.Errorf("[nntp.CmdAuthinfo] empty password: %w", ErrorInvalidParams)
			return
		}
		if err = conn.PrintfLine("AUTHINFO pass %s", pass); err != nil {
			err = fmt.Errorf("[nntp.CmdAuthinfo] failed to send AUTHINFO pass command: %w", err)
			return
		}
		if code, msg, err = conn.ReadCodeLine(0); err != nil {
			err = fmt.Errorf("[nntp.CmdAuthinfo] failed to read AUTHINFO pass response: %w", err)
			return
		}
	}
	switch code {
	case ResponseCodeAuthenticationAccepted: // 281
		err = nil
	case ResponseCodeAuthenticationContinue: // 381
		err = fmt.Errorf("[nntp.CmdAuthinfo] authentication uncompleted: %w", &Error{code, msg})
	case ResponseCodeAuthenticationRejected, ResponseCodeNotPermitted: // 482 || 502
		err = fmt.Errorf("[nntp.CmdAuthinfo] authentication rejected: %w", &Error{code, msg})
	default:
		err = fmt.Errorf("[nntp.CmdAuthinfo] unexpected response: %w", &Error{code, msg})
	}
	return
}
