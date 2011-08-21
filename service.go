package zeroconf

import (
	"log"
	"strings"
	"net"

	dns "github.com/miekg/godns"
)

type Service struct {
        Type      string
        Name      string
        Domain    string
        Interface *net.Interface
        Host   	Host
	Additional	[]string
}

type Host struct {
	Name	string
	Port	uint16
}

// s.unmarshal may not be complete, return false if so
func (s *Service) valid() bool {
	return len(s.Name) > 0
}

func (s *Service) unmarshal(msg *dns.Msg) {
	for i := range msg.Answer {
		switch rr := msg.Answer[i].(type) {
		case *dns.RR_SRV:
			x := strings.Split(rr.Hdr.Name, ".")
			s.Name = x[0]
			s.Type = strings.Join(x[1:3], ".")
			s.Domain = strings.Join(x[3:], ".")
			s.Host.Name = rr.Target
			s.Host.Port = rr.Port

		case *dns.RR_TXT:
			s.Additional = append(s.Additional, strings.Split(rr.Txt, ",")...)
		default:
			log.Printf("%#v", rr)
		}
	}
}
