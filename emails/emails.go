package emails

import (
	"bytes"
	"io"
	"mime/multipart"
	"mime/quotedprintable"
	"net/textproto"
	"time"

	"github.com/gofrs/uuid"
)

type Email struct {
	From     string
	To       string
	Subject  string
	TextPart string
	HTMLPart string
}

func (email *Email) ToBytes() ([]byte, error) {
	var err error
	buffer := bytes.NewBuffer([]byte{})
	write := func(str string) {
		if err != nil {
			return
		}
		_, err = buffer.WriteString(str)
	}
	write("From: ")
	write(email.From)
	write("\r\nTo: ")
	write(email.To)
	write("\r\nSubject: ")
	if err != nil {
		return nil, err
	}
	err = encode(buffer, []byte(email.Subject))
	if err != nil {
		return nil, err
	}
	write("\r\nDate: ")
	write(time.Now().Format(time.RFC1123Z))
	write("\r\nMessage-Id: <")
	if err != nil {
		return nil, err
	}
	messageId, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	write(messageId.String())
	write("@filharmonic.beta.gouv.fr>")
	write("\r\nMIME-version: 1.0")
	if err != nil {
		return nil, err
	}

	w := multipart.NewWriter(buffer)
	write("\r\nContent-Type: multipart/alternative; boundary=" + w.Boundary())
	write("\r\n\r\n")
	if err != nil {
		return nil, err
	}
	ww, err := w.CreatePart(textproto.MIMEHeader{
		"Content-Type":              []string{"text/plain; charset=utf-8"},
		"Content-Transfer-Encoding": []string{"quoted-printable"},
	})
	if err != nil {
		return nil, err
	}
	encoder := quotedprintable.NewWriter(ww)
	_, err = encoder.Write([]byte(email.TextPart))
	if err != nil {
		return nil, err
	}
	err = encoder.Close()
	if err != nil {
		return nil, err
	}

	ww, err = w.CreatePart(textproto.MIMEHeader{
		"Content-Type":              []string{"text/html; charset=utf-8"},
		"Content-Transfer-Encoding": []string{"quoted-printable"},
	})
	if err != nil {
		return nil, err
	}
	encoder = quotedprintable.NewWriter(ww)
	_, err = encoder.Write([]byte(email.HTMLPart))
	if err != nil {
		return nil, err
	}
	err = encoder.Close()
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func encode(writer io.Writer, data []byte) error {
	encoder := quotedprintable.NewWriter(writer)
	_, err := encoder.Write(data)
	if err != nil {
		return err
	}
	err = encoder.Close()
	return err
}
