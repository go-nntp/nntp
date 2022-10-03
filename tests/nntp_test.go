package nntp_test

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"gopkg.in/nntp.v0"
	"gopkg.in/rx.v0"
)

type message struct {
	data []byte
	send bool
}

func send(data string) *message {
	return &message{[]byte(data), true}
}

func recv(data string) *message {
	return &message{[]byte(data), false}
}

type mockServerT struct {
	*io.PipeReader
	*io.PipeWriter
}

func (s *mockServerT) Close() error {
	return nil
}

func mockServer(messages ...*message) io.ReadWriteCloser {
	sr, cw := io.Pipe()
	cr, sw := io.Pipe()
	go func() {
		for _, message := range messages {
			if message.send {
				// from client
				buf := make([]byte, len(message.data))
				_, err := io.ReadFull(sr, buf)
				if err != nil {
					panic(fmt.Errorf("server expects %#v but failed: %w", string(message.data), err))
				}
				if !bytes.Equal(buf, message.data) {
					panic(fmt.Errorf("server expects %#v but got %#v", string(message.data), string(buf)))
				}
			} else {
				_, err := sw.Write(message.data)
				if err != nil {
					panic(fmt.Errorf("server failed to send %#v: %w", string(message.data), err))
				}
			}
		}
	}()
	return &mockServerT{cr, cw}
}

func ArticleOverviewsEquals(a, b *nntp.ArticleOverview) bool {
	if len(a.ExtraFields) != len(b.ExtraFields) {
		return false
	}
	for i := 0; i < len(a.ExtraFields); i++ {
		if a.ExtraFields[i] != b.ExtraFields[i] {
			return false
		}
	}
	return a.ArticleNumber == b.ArticleNumber &&
		a.Subject == b.Subject &&
		a.From == b.From &&
		a.Date == b.Date &&
		a.MessageID == b.MessageID &&
		a.References == b.References &&
		a.Bytes == b.Bytes &&
		a.Lines == b.Lines
}

func TestOverCommand(t *testing.T) {
	const (
		USER = "testuser"
		PASS = "abcd1234!@#"
	)
	var err error

	GROUP := nntp.GroupStat{
		Count: 4582130111,
		First: 1502531334,
		Last:  6084661444,
		Group: "alt.binaries.test",
	}

	ARTICLES := []*nntp.ArticleOverview{
		{ArticleNumber: 6084661434,
			Subject:     "[081/112] \"BFMlDDFLa1SxxromF4cORw.part080.rar\" (2001/2947)",
			From:        "KKsT8FSZc6uWH@ngPost.com",
			Date:        "Sun, 25 Sep 2022 03:03:36 GMT",
			MessageID:   "<90094f167d034fc994736d02757f7d7b@ngPost>",
			References:  "",
			Bytes:       0xb4b09,
			Lines:       0x163b,
			ExtraFields: []string{"Xref: e alt.binaries.test:6084661434"}},
		{ArticleNumber: 6084661435,
			Subject:     "[06/18] - \"0Vcc4lHHfyTlD03O0UyX3X.part05.rar\" yEnc (634/732) 524288000",
			From:        "8oS8FBXSln <8oS8FBXSln@6Y7EjqT.5U>",
			Date:        "Sun, 25 Sep 22 03:03:37 UTC",
			MessageID:   "<AgIxVaMfQgIxAhWcPlBwZvOp-1664075015428@nyuu>",
			References:  "",
			Bytes:       0xb4bf0,
			Lines:       0x0,
			ExtraFields: []string{"Xref: e alt.binaries.test:6084661435"}},
		{ArticleNumber: 6084661436,
			Subject:     "[06/18] - \"0Vcc4lHHfyTlD03O0UyX3X.part05.rar\" yEnc (635/732) 524288000",
			From:        "8oS8FBXSln <8oS8FBXSln@6Y7EjqT.5U>",
			Date:        "Sun, 25 Sep 22 03:03:37 UTC",
			MessageID:   "<FjUqAeQlGkBgBpEaWaKxGrTo-1664075015663@nyuu>",
			References:  "",
			Bytes:       0xb4c54,
			Lines:       0x0,
			ExtraFields: []string{"Xref: e alt.binaries.test:6084661436"}},
		{ArticleNumber: 6084661437,
			Subject:     "[06/18] - \"0Vcc4lHHfyTlD03O0UyX3X.part05.rar\" yEnc (636/732) 524288000",
			From:        "8oS8FBXSln <8oS8FBXSln@6Y7EjqT.5U>",
			Date:        "Sun, 25 Sep 22 03:03:37 UTC",
			MessageID:   "<AqNgDhSmSuDkVzCkAbRzFgLi-1664075015866@nyuu>",
			References:  "",
			Bytes:       0xb4c47,
			Lines:       0x0,
			ExtraFields: []string{"Xref: e alt.binaries.test:6084661437"}},
		{ArticleNumber: 6084661438,
			Subject:     "[06/18] - \"0Vcc4lHHfyTlD03O0UyX3X.part05.rar\" yEnc (637/732) 524288000",
			From:        "8oS8FBXSln <8oS8FBXSln@6Y7EjqT.5U>",
			Date:        "Sun, 25 Sep 22 03:03:37 UTC",
			MessageID:   "<DdQmIoTzVkApQhGoNhIrUnUc-1664075015897@nyuu>",
			References:  "",
			Bytes:       0xb4c8b,
			Lines:       0x0,
			ExtraFields: []string{"Xref: e alt.binaries.test:6084661438"}},
		{ArticleNumber: 6084661439,
			Subject:     "91f7599197524ead8e8172ce5ab9d2fa",
			From:        "p1GculobIZUZv@ngPost.com",
			Date:        "Sun, 25 Sep 2022 03:03:37 GMT",
			MessageID:   "<91f7599197524ead8e8172ce5ab9d2fa@ngPost>",
			References:  "",
			Bytes:       0xb4a70,
			Lines:       0x163a,
			ExtraFields: []string{"Xref: e alt.binaries.test:6084661439 alt.binaries.misc:18216642898 a.b.boneless:192526866"}},
		{ArticleNumber: 6084661440,
			Subject:     "[06/18] - \"0Vcc4lHHfyTlD03O0UyX3X.part05.rar\" yEnc (638/732) 524288000",
			From:        "8oS8FBXSln <8oS8FBXSln@6Y7EjqT.5U>",
			Date:        "Sun, 25 Sep 22 03:03:38 UTC",
			MessageID:   "<XlHsSjSfRvVfBoTjToJuKwSf-1664075016038@nyuu>",
			References:  "",
			Bytes:       0xb4c91,
			Lines:       0x0,
			ExtraFields: []string{"Xref: e alt.binaries.test:6084661440"}},
		{ArticleNumber: 6084661441,
			Subject:     "[06/18] - \"0Vcc4lHHfyTlD03O0UyX3X.part05.rar\" yEnc (639/732) 524288000",
			From:        "8oS8FBXSln <8oS8FBXSln@6Y7EjqT.5U>",
			Date:        "Sun, 25 Sep 22 03:03:38 UTC",
			MessageID:   "<XpDnZfOcOoFpXiVrQfVjRpAr-1664075016335@nyuu>",
			References:  "",
			Bytes:       0xb4be6,
			Lines:       0x0,
			ExtraFields: []string{"Xref: e alt.binaries.test:6084661441"}},
		{ArticleNumber: 6084661442,
			Subject:     "[06/18] - \"0Vcc4lHHfyTlD03O0UyX3X.part05.rar\" yEnc (640/732) 524288000",
			From:        "8oS8FBXSln <8oS8FBXSln@6Y7EjqT.5U>",
			Date:        "Sun, 25 Sep 22 03:03:38 UTC",
			MessageID:   "<UkSwHrXmBvZdRcWuQwOaCvRl-1664075016553@nyuu>",
			References:  "",
			Bytes:       0xb4c79,
			Lines:       0x0,
			ExtraFields: []string{"Xref: e alt.binaries.test:6084661442"}},
		{ArticleNumber: 6084661443,
			Subject:     "[06/18] - \"0Vcc4lHHfyTlD03O0UyX3X.part05.rar\" yEnc (642/732) 524288000",
			From:        "8oS8FBXSln <8oS8FBXSln@6Y7EjqT.5U>",
			Date:        "Sun, 25 Sep 22 03:03:38 UTC",
			MessageID:   "<McTcBwYrEbMtObVrUfYqBfRz-1664075016757@nyuu>",
			References:  "",
			Bytes:       0xb4c4a,
			Lines:       0x0,
			ExtraFields: []string{"Xref: e alt.binaries.test:6084661443"}},
		{ArticleNumber: 6084661444,
			Subject:     "[06/18] - \"0Vcc4lHHfyTlD03O0UyX3X.part05.rar\" yEnc (641/732) 524288000",
			From:        "8oS8FBXSln <8oS8FBXSln@6Y7EjqT.5U>",
			Date:        "Sun, 25 Sep 22 03:03:38 UTC",
			MessageID:   "<JlCkQlSoDmLfUgJmAzQcHoOw-1664075016616@nyuu>",
			References:  "",
			Bytes:       0xb4bdc,
			Lines:       0x0,
			ExtraFields: []string{"Xref: e alt.binaries.test:6084661444"}},
	}

	netconn := mockServer(
		recv("200 Welcome to Usenet\r\n"),
		send("AUTHINFO user "+USER+"\r\n"),
		recv("381 PASS required\r\n"),
		send("AUTHINFO pass "+PASS+"\r\n"),
		recv("281 Welcome to Usenet\r\n"),
		send("GROUP "+GROUP.Group+"\r\n"),
		recv(fmt.Sprintf("211 %d %d %d %s\r\n", GROUP.Count, GROUP.First, GROUP.Last, GROUP.Group)),
		send(fmt.Sprintf("XOVER %d-%d\r\n", GROUP.Last-10, GROUP.Last)),
		recv("224 Overview Information Follows\r\n"+
			"6084661434	[081/112] \"BFMlDDFLa1SxxromF4cORw.part080.rar\" (2001/2947)	KKsT8FSZc6uWH@ngPost.com	Sun, 25 Sep 2022 03:03:36 GMT	<90094f167d034fc994736d02757f7d7b@ngPost>		740105	5691	Xref: e alt.binaries.test:6084661434\r\n"+
			"6084661435	[06/18] - \"0Vcc4lHHfyTlD03O0UyX3X.part05.rar\" yEnc (634/732) 524288000	8oS8FBXSln <8oS8FBXSln@6Y7EjqT.5U>	Sun, 25 Sep 22 03:03:37 UTC	<AgIxVaMfQgIxAhWcPlBwZvOp-1664075015428@nyuu>		740336		Xref: e alt.binaries.test:6084661435\r\n"+
			"6084661436	[06/18] - \"0Vcc4lHHfyTlD03O0UyX3X.part05.rar\" yEnc (635/732) 524288000	8oS8FBXSln <8oS8FBXSln@6Y7EjqT.5U>	Sun, 25 Sep 22 03:03:37 UTC	<FjUqAeQlGkBgBpEaWaKxGrTo-1664075015663@nyuu>		740436		Xref: e alt.binaries.test:6084661436\r\n"+
			"6084661437	[06/18] - \"0Vcc4lHHfyTlD03O0UyX3X.part05.rar\" yEnc (636/732) 524288000	8oS8FBXSln <8oS8FBXSln@6Y7EjqT.5U>	Sun, 25 Sep 22 03:03:37 UTC	<AqNgDhSmSuDkVzCkAbRzFgLi-1664075015866@nyuu>		740423		Xref: e alt.binaries.test:6084661437\r\n"+
			"6084661438	[06/18] - \"0Vcc4lHHfyTlD03O0UyX3X.part05.rar\" yEnc (637/732) 524288000	8oS8FBXSln <8oS8FBXSln@6Y7EjqT.5U>	Sun, 25 Sep 22 03:03:37 UTC	<DdQmIoTzVkApQhGoNhIrUnUc-1664075015897@nyuu>		740491		Xref: e alt.binaries.test:6084661438\r\n"+
			"6084661439	91f7599197524ead8e8172ce5ab9d2fa	p1GculobIZUZv@ngPost.com	Sun, 25 Sep 2022 03:03:37 GMT	<91f7599197524ead8e8172ce5ab9d2fa@ngPost>		739952	5690	Xref: e alt.binaries.test:6084661439 alt.binaries.misc:18216642898 a.b.boneless:192526866\r\n"+
			"6084661440	[06/18] - \"0Vcc4lHHfyTlD03O0UyX3X.part05.rar\" yEnc (638/732) 524288000	8oS8FBXSln <8oS8FBXSln@6Y7EjqT.5U>	Sun, 25 Sep 22 03:03:38 UTC	<XlHsSjSfRvVfBoTjToJuKwSf-1664075016038@nyuu>		740497		Xref: e alt.binaries.test:6084661440\r\n"+
			"6084661441	[06/18] - \"0Vcc4lHHfyTlD03O0UyX3X.part05.rar\" yEnc (639/732) 524288000	8oS8FBXSln <8oS8FBXSln@6Y7EjqT.5U>	Sun, 25 Sep 22 03:03:38 UTC	<XpDnZfOcOoFpXiVrQfVjRpAr-1664075016335@nyuu>		740326		Xref: e alt.binaries.test:6084661441\r\n"+
			"6084661442	[06/18] - \"0Vcc4lHHfyTlD03O0UyX3X.part05.rar\" yEnc (640/732) 524288000	8oS8FBXSln <8oS8FBXSln@6Y7EjqT.5U>	Sun, 25 Sep 22 03:03:38 UTC	<UkSwHrXmBvZdRcWuQwOaCvRl-1664075016553@nyuu>		740473		Xref: e alt.binaries.test:6084661442\r\n"+
			"6084661443	[06/18] - \"0Vcc4lHHfyTlD03O0UyX3X.part05.rar\" yEnc (642/732) 524288000	8oS8FBXSln <8oS8FBXSln@6Y7EjqT.5U>	Sun, 25 Sep 22 03:03:38 UTC	<McTcBwYrEbMtObVrUfYqBfRz-1664075016757@nyuu>		740426		Xref: e alt.binaries.test:6084661443\r\n"+
			"6084661444	[06/18] - \"0Vcc4lHHfyTlD03O0UyX3X.part05.rar\" yEnc (641/732) 524288000	8oS8FBXSln <8oS8FBXSln@6Y7EjqT.5U>	Sun, 25 Sep 22 03:03:38 UTC	<JlCkQlSoDmLfUgJmAzQcHoOw-1664075016616@nyuu>		740316		Xref: e alt.binaries.test:6084661444\r\n"+
			".\r\n"),
	)
	conn := nntp.NewConn(netconn)
	if err = conn.ReadWelcome(); err != nil {
		t.Fatal(err)
	}
	if err = conn.CmdAuthinfo(USER, PASS); err != nil {
		t.Fatal(err)
	}
	group, err := conn.CmdGroup(GROUP.Group)
	if err != nil {
		t.Fatal(err)
	}

	// log.Printf("%#v\n", group)

	if *group != GROUP {
		t.Errorf("client expects %#v but got %#v", GROUP, group)
	}

	writer, reader := rx.Pipe[*nntp.ArticleOverview](nil)
	conn.CmdXOver(nntp.WithArticleRange(group.Last-10, group.Last)).Subscribe(writer)

	i := 0
	for {
		article, ok := reader.Read()
		if !ok {
			break
		}
		// log.Printf("%#v\n", article)
		if !ArticleOverviewsEquals(article, ARTICLES[i]) {
			t.Errorf("client expects %#v but got %#v", ARTICLES[i], article)
		}
		i++
	}

	if err = reader.Err(); err != nil {
		t.Fatal(err)
	}

	if i != len(ARTICLES) {
		t.Errorf("client expects %d articles but only got %d", len(ARTICLES), i)
	}
}
