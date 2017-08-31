package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"net/mail"
	"net/smtp"
	"os"
	"os/exec"
)

func sendMail(addr, from, rcpt string, msg []byte) error {
	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()
	defer c.Quit()

	c.Mail(from)
	c.Rcpt(rcpt)
	wc, err := c.Data()
	if err != nil {
		return err
	}
	defer wc.Close()

	_, err = wc.Write(msg)
	return err
}

func encryptMail(certdir, from, rcpt, subject string, buf bytes.Buffer) ([]byte, error) {
	// Openssl for smime encryption.
	c := exec.Command("openssl", "smime", "-encrypt", "-des3", "-from", from, "-to", rcpt, "-subject", subject, certdir+"/"+rcpt+".crt")

	// Pipe msg to openssl.
	stdin, err := c.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	// Write to stdin.
	go func() {
		defer stdin.Close()
		stdin.Write(buf.Bytes())
	}()

	// exec openssl smime
	return c.CombinedOutput()
}

func main() {
	var (
		addr, rcpt, from, subject, certdir string
	)
	flag.StringVar(&addr, "addr", "localhost:10025", "smtp address to send to")
	flag.StringVar(&certdir, "dir", "/etc/postfix/smime", "smime cert directory")
	flag.StringVar(&from, "from", "sender@example.com", "mail sender")
	flag.StringVar(&rcpt, "to", "recipient@example.com", "mail recipient")
	flag.Parse()

	// Read Mail from stdin
	var buf bytes.Buffer
	tee := io.TeeReader(os.Stdin, &buf)

	// Parse E-Mail
	m, err := mail.ReadMessage(tee)
	if err != io.EOF && err != nil {
		log.Fatal("error parsing mail: ", err)
	}

	header := m.Header
	subject = header.Get("Subject")

	emsg, err := encryptMail(certdir, from, rcpt, subject, buf)
	if err != nil {
		// Send mail unencrypted
		log.Println("not encrypting")
		err := sendMail(addr, from, rcpt, buf.Bytes())
		if err != nil {
			log.Fatal("error sending unencrypted mail: ", err)
		}
	}
	err = sendMail(addr, from, rcpt, emsg)
	if err != nil {
		log.Fatal("error sending encrypted mail: ", err)
	}
}
