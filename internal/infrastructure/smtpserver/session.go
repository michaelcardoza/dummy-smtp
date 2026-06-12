package smtpserver

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/michaelcardoza/dummy-smtp/internal/core/mail"
)

type Capturer interface {
	Capture(ctx context.Context, params mail.CaptureParams) (*mail.Message, error)
}

type state int

const (
	stateGreeting state = iota
	stateReady
	stateMail
	stateRcpt
)

const idleTimeout = 5 * time.Minute

type session struct {
	conn     net.Conn
	reader   *bufio.Reader
	capturer Capturer

	state state
	from  string
	to    []string
}

func newSession(conn net.Conn, capturer Capturer) *session {
	return &session{
		conn:     conn,
		reader:   bufio.NewReader(conn),
		capturer: capturer,
		state:    stateGreeting,
	}
}

func (s *session) serve(ctx context.Context) {
	defer s.conn.Close()

	s.reply(220, "dummy-smtp ready")

	for {
		s.conn.SetDeadline(time.Now().Add(idleTimeout))

		line, err := s.reader.ReadString('\n')
		if err != nil {
			return
		}

		cmd, args := splitCommand(line)
		if quit := s.handle(ctx, cmd, args); quit {
			return
		}
	}
}

func (s *session) handle(ctx context.Context, cmd, args string) (quit bool) {
	switch cmd {
	case "HELO", "EHLO":
		s.resetEnvelope()
		s.state = stateReady
		s.reply(250, "hello "+args)

	case "MAIL":
		if s.state != stateReady {
			s.reply(503, "bad sequence of commands")
			return false
		}
		s.from = extractArgs(args)
		s.state = stateMail
		s.reply(250, "ok")

	case "RCPT":
		if s.state != stateMail && s.state != stateRcpt {
			s.reply(503, "need MAIL before RCPT")
			return false
		}
		s.to = append(s.to, extractArgs(args))
		s.state = stateRcpt
		s.reply(250, "ok")

	case "DATA":
		if s.state != stateRcpt {
			s.reply(503, "need RCPT before DATA")
			return false
		}
		s.handleData(ctx)

	case "RSET":
		s.resetEnvelope()
		s.state = stateReady
		s.reply(250, "ok")

	case "NOOP":
		s.reply(250, "ok")

	case "QUIT":
		s.reply(221, "bye")
		return true

	default:
		s.reply(502, "command not implemented")
	}
	return false
}

func (s *session) handleData(ctx context.Context) {
	s.reply(354, "end data with <CRLF>.<CRLF>")

	raw, err := s.readData()
	if err != nil {
		s.reply(451, "error reading data")
		return
	}

	parsed := parse(raw)

	_, err = s.capturer.Capture(ctx, mail.CaptureParams{
		From:        s.from,
		To:          s.to,
		Subject:     parsed.Subject,
		TextBody:    parsed.TextBody,
		HTMLBody:    parsed.HTMLBody,
		Raw:         raw,
		Attachments: parsed.Attachments,
	})
	if err != nil {
		s.reply(451, "error processing message")
		return
	}

	s.reply(250, "message accepted")
	s.resetEnvelope()
	s.state = stateReady
}

func (s *session) readData() (string, error) {
	var sb strings.Builder
	for {
		line, err := s.reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		trimmed := strings.TrimRight(line, "\r\n")
		if trimmed == "." {
			break
		}
		if strings.HasPrefix(trimmed, "..") {
			line = line[1:]
		}
		sb.WriteString(line)
	}
	return sb.String(), nil
}

func (s *session) reply(code int, msg string) {
	fmt.Fprintf(s.conn, "%d %s\r\n", code, msg)
}

func (s *session) resetEnvelope() {
	s.from = ""
	s.to = nil
}

func splitCommand(line string) (cmd, args string) {
	line = strings.TrimRight(line, "\r\n")
	parts := strings.SplitN(line, " ", 2)
	cmd = strings.ToUpper(parts[0])
	if len(parts) > 1 {
		args = strings.TrimSpace(parts[1])
	}
	return cmd, args
}

func extractArgs(args string) string {
	start := strings.IndexByte(args, '<')
	end := strings.IndexByte(args, '>')
	if start >= 0 && end > start {
		return args[start+1 : end]
	}
	if i := strings.IndexByte(args, ':'); i >= 0 {
		return strings.TrimSpace(args[i+1:])
	}
	return strings.TrimSpace(args)
}
