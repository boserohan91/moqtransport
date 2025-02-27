package chat

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/mengelbart/moqtransport"
)

type joinedRooms struct {
	trackID uint64
	st      *moqtransport.SendTrack
	rts     []*moqtransport.ReceiveTrack
}

type Client struct {
	session     *moqtransport.Session
	rooms       map[string]*joinedRooms
	lock        sync.Mutex
	nextTrackID uint64
}

func NewQUICClient(addr string) (*Client, error) {
	p, err := moqtransport.DialQUIC(addr, 3)
	if err != nil {
		return nil, err
	}
	return NewClient(p)
}

func NewWebTransportClient(addr string) (*Client, error) {
	p, err := moqtransport.DialWebTransport(addr, 3)
	if err != nil {
		return nil, err
	}
	return NewClient(p)
}

func NewClient(p *moqtransport.Session) (*Client, error) {
	log.SetOutput(io.Discard)
	c := &Client{
		session:     p,
		rooms:       map[string]*joinedRooms{},
		lock:        sync.Mutex{},
		nextTrackID: 0,
	}
	go func() {
		for {
			var a *moqtransport.Announcement
			a, err := c.session.ReadAnnouncement(context.Background())
			if err != nil {
				panic(err)
			}
			log.Printf("got Announcement: %v", a.Namespace())
			a.Accept()
		}
	}()
	go func() {
		for {
			s, err := c.session.ReadSubscription(context.Background())
			if err != nil {
				panic(err)
			}
			parts := strings.SplitN(s.Namespace(), "/", 2)
			if len(parts) < 2 {
				s.Reject(errors.New("invalid trackname"))
				continue
			}
			moq_chat, id := parts[0], parts[1]
			if moq_chat != "moq-chat" {
				s.Reject(errors.New("invalid moq-chat namespace"))
				continue
			}
			if _, ok := c.rooms[id]; !ok {
				s.Reject(errors.New("invalid subscribe request"))
				continue
			}
			s.SetTrackID(c.rooms[id].trackID)
			c.rooms[id].st = s.Accept()
		}
	}()
	return c, nil
}

func (c *Client) handleCatalogDeltas(roomID, username string, catalogTrack *moqtransport.ReceiveTrack) error {
	buf := make([]byte, 64_000)
	for {
		n, err := catalogTrack.Read(buf)
		if err != nil {
			return err
		}
		delta, err := parseDelta(string(buf[:n]))
		if err != nil {
			return err
		}
		for _, p := range delta.joined {
			if p == username {
				continue
			}
			t, err := c.session.Subscribe(context.Background(), fmt.Sprintf("moq-chat/%v", roomID), p, username)
			if err != nil {
				return err
			}
			go func(room, user string) {
				fmt.Printf("%v joined the chat %v\n> ", user, room)
				for {
					buf := make([]byte, 64_000)
					n, err = t.Read(buf)
					if err != nil {
						log.Fatal(err)
					}
					fmt.Fprintf(os.Stdout, "room %v|user %v: %v\n> ", room, user, string(buf[:n]))
				}
			}(roomID, p)
		}
	}
}

func (c *Client) joinRoom(roomID, username string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.nextTrackID += 1
	c.rooms[roomID] = &joinedRooms{
		trackID: c.nextTrackID,
		st:      nil,
		rts:     []*moqtransport.ReceiveTrack{},
	}
	if err := c.session.Announce(context.Background(), fmt.Sprintf("moq-chat/%v/participant/%v", roomID, username)); err != nil {
		return err
	}
	catalogTrack, err := c.session.Subscribe(context.Background(), fmt.Sprintf("moq-chat/%v", roomID), "/catalog", username)
	if err != nil {
		return err
	}
	buf := make([]byte, 64_000)
	var n int
	n, err = catalogTrack.Read(buf)
	if err != nil {
		return err
	}
	var participants *chatalog
	participants, err = parseChatalog(string(buf[:n]))
	if err != nil {
		return err
	}
	for p := range participants.participants {
		if p == username {
			continue
		}
		t, err := c.session.Subscribe(context.Background(), fmt.Sprintf("moq-chat/%v", roomID), p, username)
		if err != nil {
			log.Fatalf("failed to subscribe to participant track: %v", err)
		}
		go func(room, user string) {
			for {
				buf := make([]byte, 64_000)
				n, err = t.Read(buf)
				if err != nil {
					log.Fatalf("failed to read from participant track: %v", err)
				}
				fmt.Fprintf(os.Stdout, "room %v|user %v: %v\n> ", room, user, string(buf[:n]))
			}
		}(roomID, p)
	}
	go c.handleCatalogDeltas(roomID, username, catalogTrack)
	return nil
}

func (c *Client) Run() error {
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprintf(os.Stdout, "> ")
		cmd, err := r.ReadString('\n')
		if err != nil {
			return err
		}
		if strings.HasPrefix(cmd, "join") {
			fields := strings.Fields(cmd)
			if len(fields) < 3 {
				fmt.Println("invalid join command, usage: 'join <room id> <username>'")
				continue
			}
			if err = c.joinRoom(fields[1], fields[2]); err != nil {
				return err
			}
			continue
		}
		if strings.HasPrefix(cmd, "msg") {
			fields := strings.Fields(cmd)
			if len(fields) < 3 {
				fmt.Println("invalid join command, usage: 'msg <room id> <msg>'")
				continue
			}
			msg, ok := strings.CutPrefix(cmd, fmt.Sprintf("msg %v", fields[1]))
			if !ok {
				fmt.Println("invalid msg command, usage: 'msg <room id> <msg>'")
				continue
			}
			w, err := c.rooms[fields[1]].st.StartReliableObject()
			if err != nil {
				fmt.Printf("failed to send object: %v", err)
				continue
			}
			_, err = w.Write([]byte(strings.TrimSpace(msg)))
			if err != nil {
				return fmt.Errorf("failed to write to room: %v", err)
			}
			continue
		}
		fmt.Println("invalid command, try 'join' or 'msg'")
	}
}
