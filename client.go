package main

import (
	"AI/models"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ByteArena/box2d"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

const (
	WAITING       = 0
	IN_BATTLE     = 1
	IN_SETTLEMENT = 2
	IN_DISMISSAL  = 3
)

const (
	// You can equivalently use the `GroupIndex` approach, but the more complicated and general purpose approach is used deliberately here. Reference http://www.aurelienribon.com/post/2011-07-box2d-tutorial-collision-filtering.
	COLLISION_CATEGORY_CONTROLLED_PLAYER = (1 << 1)
	COLLISION_CATEGORY_TREASURE          = (1 << 2)
	COLLISION_CATEGORY_TRAP              = (1 << 3)
	COLLISION_CATEGORY_TRAP_BULLET       = (1 << 4)
	COLLISION_CATEGORY_BARRIER           = (1 << 5)
	COLLISION_CATEGORY_PUMPKIN           = (1 << 6)
	COLLISION_CATEGORY_SPEED_SHOES       = (1 << 7)

	COLLISION_MASK_FOR_CONTROLLED_PLAYER = (COLLISION_CATEGORY_TREASURE | COLLISION_CATEGORY_TRAP | COLLISION_CATEGORY_TRAP_BULLET | COLLISION_CATEGORY_SPEED_SHOES)
	COLLISION_MASK_FOR_TREASURE          = (COLLISION_CATEGORY_CONTROLLED_PLAYER)
	COLLISION_MASK_FOR_TRAP              = (COLLISION_CATEGORY_CONTROLLED_PLAYER)
	COLLISION_MASK_FOR_TRAP_BULLET       = (COLLISION_CATEGORY_CONTROLLED_PLAYER)
	COLLISION_MASK_FOR_BARRIER           = (COLLISION_CATEGORY_PUMPKIN)
	COLLISION_MASK_FOR_PUMPKIN           = (COLLISION_CATEGORY_BARRIER)
	COLLISION_MASK_FOR_SPEED_SHOES       = (COLLISION_CATEGORY_CONTROLLED_PLAYER)
)

type wsReq struct {
	MsgId int             `json:"msgId"`
	Act   string          `json:"act"`
	Data  json.RawMessage `json:"data"`
}

type wsResp struct {
	Ret         int32           `json:"ret,omitempty"`
	EchoedMsgId int32           `json:"echoedMsgId,omitempty"`
	Act         string          `json:"act,omitempty"`
	Data        json.RawMessage `json:"data,omitempty"`
}

type wsRespPb struct {
	Ret         int32  `json:"ret,omitempty"`
	EchoedMsgId int32  `json:"echoedMsgId,omitempty"`
	Act         string `json:"act,omitempty"`
	Data        []byte `json:"data,omitempty"`
}

type Client struct {
	Id                    int //roomId
	LastRoomDownsyncFrame *models.RoomDownsyncFrame
	BattleState           int
	c                     *websocket.Conn
	Player                *models.Player
	CollidableWorld       *box2d.B2World
	Barrier               map[int32]*models.Barrier
	PlayerCollidableBody  *box2d.B2Body `json:"-"`
}

var addr = flag.String("addr", "localhost:9992", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/tsrht"}
	q := u.Query()
	q.Set("intAuthToken", "e3595bb3b2e44c651033aa88f3399b71")
	u.RawQuery = q.Encode()
	//ref to the NewClient and DefaultDialer.Dial https://github.com/gorilla/websocket/issues/54
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		client := &Client{
			LastRoomDownsyncFrame: nil,
			BattleState:           -1,
			c:                     c,
			Player:                &models.Player{Id: 10},
			Barrier:               make(map[int32]*models.Barrier),
		}
		client.initMapStaticResource()
		for {
			var resp *wsResp
			resp = new(wsResp)
			err := c.ReadJSON(resp)
			if err != nil {
				//log.Println("marshal wsResp:", err)
			}
			if resp.Act == "RoomDownsyncFrame" {
				var respPb *wsRespPb
				respPb = new(wsRespPb)
				err := c.ReadJSON(respPb)
				if err != nil {
					//log.Println("marshal respPb:", err)
				}
				client.decodeProtoBuf(respPb.Data)
				client.controller()
				client.upsyncFrameData()
			} else {
				//handleHbRequirements(resp)
			}
			time.Sleep(time.Duration(int64(20)))
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")

			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func (client *Client) controller() {
	if client.Player.Speed == 0 {
		return
	}
	if client.LastRoomDownsyncFrame.Id == 1 {
		client.BattleState = IN_BATTLE
		log.Println("Game Start")
	}
	for edge := client.PlayerCollidableBody.GetContactList(); edge != nil; edge = edge.Next {
		log.Println("player conteact")
		if edge.Contact.IsTouching() {
			log.Println("player conteact")
			if _, ok := edge.Other.GetUserData().(*models.Barrier); ok {
				log.Println("player conteact to the barrier")
			}
		}
	}
	client.Player.Y = client.Player.Y - 10
	//client.Player.Y++
	newB2Vec2Pos := box2d.MakeB2Vec2(client.Player.X, client.Player.Y)
	models.MoveDynamicBody(client.PlayerCollidableBody, &newB2Vec2Pos, 0)
	//models.PrettyPrintBody(client.PlayerCollidableBody)
	for edge := client.PlayerCollidableBody.GetContactList(); edge != nil; edge = edge.Next {
		log.Println("player conteact")
		if edge.Contact.IsTouching() {
			log.Println("player conteact")
			if _, ok := edge.Other.GetUserData().(*models.Barrier); ok {
				log.Println("player conteact to the barrier")
			}
		}
	}
}

func (client *Client) upsyncFrameData() {
	if client.BattleState == IN_BATTLE {
		newFrame := &struct {
			Id            int64             `json:"id"`
			X             float64           `json:"x"`
			Y             float64           `json:"y"`
			Dir           *models.Direction `json:"dir"`
			AckingFrameId int32             `json:"AckingFrameId"`
		}{client.Player.Id, client.Player.X, client.Player.Y, client.Player.Dir, client.LastRoomDownsyncFrame.Id}
		newFrameByte, err := json.Marshal(newFrame)
		if err != nil {
			log.Println("json Marshal:", err)
			return
		}
		req := &wsReq{
			MsgId: 1,
			Act:   "PlayerUpsyncCmd",
			Data:  newFrameByte,
		}
		reqByte, err := json.Marshal(req)
		err = client.c.WriteMessage(websocket.TextMessage, reqByte)
		if err != nil {
			log.Println("write:", err)
			return
		}
	}
}

func (client *Client) decodeProtoBuf(message []byte) {
	room_downsync_frame := models.RoomDownsyncFrame{}
	err := proto.Unmarshal(message, &room_downsync_frame)
	if err != nil {
		log.Fatal(err)
	}
	client.LastRoomDownsyncFrame = &room_downsync_frame
	client.Player.Speed = room_downsync_frame.Players[int32(client.Player.Id)].Speed
	client.Player.Dir = room_downsync_frame.Players[int32(client.Player.Id)].Dir
	client.Player.X = room_downsync_frame.Players[int32(client.Player.Id)].X
	client.Player.Y = room_downsync_frame.Players[int32(client.Player.Id)].Y
}

func (client *Client) initMapStaticResource() {

	relativePath := "./map/map/treasurehunter.tmx"
	execPath, err := os.Executable()
	ErrFatal(err)

	pwd, err := os.Getwd()
	ErrFatal(err)

	fmt.Printf("execPath = %v, pwd = %s, returning...\n", execPath, pwd)

	tmxMapIns := models.TmxMap{}
	pTmxMapIns := &tmxMapIns
	fp := filepath.Join(pwd, relativePath)
	fmt.Printf("fp == %v\n", fp)
	if !filepath.IsAbs(fp) {
		panic("Tmx filepath must be absolute!")
	}

	byteArr, err := ioutil.ReadFile(fp)
	ErrFatal(err)
	models.DeserializeToTmxMapIns(byteArr, pTmxMapIns)

	tsxIns := models.Tsx{}
	pTsxIns := &tsxIns
	relativePath = "./map/map/tile_1.tsx"
	fp = filepath.Join(pwd, relativePath)
	fmt.Printf("fp == %v\n", fp)
	if !filepath.IsAbs(fp) {
		panic("Filepath must be absolute!")
	}

	byteArr, err = ioutil.ReadFile(fp)
	ErrFatal(err)
	models.DeserializeToTsxIns(byteArr, pTsxIns)

	client.InitColliders()
	client.InitBarrier(pTmxMapIns, pTsxIns)
}

func (client *Client) InitColliders() {
	gravity := box2d.MakeB2Vec2(0.0, 0.0)
	world := box2d.MakeB2World(gravity)
	world.SetContactFilter(&box2d.B2ContactFilter{})
	client.CollidableWorld = &world

	log.Println("InitColliders for client.Players:", zap.Any("roomId", client.Id))
	player := client.Player
	var bdDef box2d.B2BodyDef
	colliderOffset := box2d.MakeB2Vec2(0, 0) // Matching that of client-side setting.
	bdDef = box2d.MakeB2BodyDef()
	bdDef.Type = box2d.B2BodyType.B2_dynamicBody
	bdDef.Position.Set(player.X+colliderOffset.X, player.Y+colliderOffset.Y)

	b2PlayerBody := client.CollidableWorld.CreateBody(&bdDef)

	b2CircleShape := box2d.MakeB2CircleShape()
	b2CircleShape.M_radius = 32 // Matching that of client-side setting.

	fd := box2d.MakeB2FixtureDef()
	fd.Shape = &b2CircleShape
	fd.Filter.CategoryBits = COLLISION_CATEGORY_PUMPKIN
	fd.Filter.MaskBits = COLLISION_MASK_FOR_PUMPKIN
	fd.Density = 0.0
	b2PlayerBody.CreateFixtureFromDef(&fd)

	client.PlayerCollidableBody = b2PlayerBody
	b2PlayerBody.SetUserData(player)
	//models.PrettyPrintBody(client.PlayerCollidableBody)
}

func (client *Client) InitBarrier(pTmxMapIns *models.TmxMap, pTsxIns *models.Tsx) {
	for _, lay := range pTmxMapIns.Layers {
		if lay.Name != "tile_1 human skeleton" && lay.Name != "tile_1 board" && lay.Name != "tile_1 stone" {
			continue
		}
		for index, tile := range lay.Tile {
			if tile == nil || tile.Tileset == nil {
				continue
			}
			if tile.Tileset.Source != "tile_1.tsx" {
				continue
			}

			barrier := &models.Barrier{}
			barrier.X, barrier.Y = pTmxMapIns.GetCoordByGid(index)
			barrier.Type = tile.Id
			if v, ok := pTsxIns.BarrierPolyLineList[int(tile.Id)]; ok {
				thePoints := make([]*models.Vec2D, 0)
				for _, p := range v.Points {
					thePoints = append(thePoints, &models.Vec2D{
						X: p.X,
						Y: p.Y,
					})
				}
				barrier.Boundary = &models.Polygon2D{Points: thePoints}
			}

			var bdDef box2d.B2BodyDef
			bdDef = box2d.MakeB2BodyDef()
			bdDef.Type = box2d.B2BodyType.B2_staticBody
			bdDef.Position.Set(barrier.X, barrier.Y) // todo ？？？？？
			b2EmelementBody := client.CollidableWorld.CreateBody(&bdDef)

			fd := box2d.MakeB2FixtureDef()
			if barrier.Boundary != nil {
				b2Vertices := make([]box2d.B2Vec2, len(barrier.Boundary.Points))
				for vIndex, v2 := range barrier.Boundary.Points {
					b2Vertices[vIndex] = v2.ToB2Vec2()
				}
				b2PolygonShape := box2d.MakeB2PolygonShape()
				b2PolygonShape.Set(b2Vertices, len(barrier.Boundary.Points))
				fd.Shape = &b2PolygonShape
			} else {
				b2CircleShape := box2d.MakeB2CircleShape()
				b2CircleShape.M_radius = 32
				fd.Shape = &b2CircleShape
			}

			fd.Filter.CategoryBits = COLLISION_CATEGORY_BARRIER
			fd.Filter.MaskBits = COLLISION_MASK_FOR_BARRIER
			fd.Density = 0.0
			b2EmelementBody.CreateFixtureFromDef(&fd)

			barrier.CollidableBody = b2EmelementBody
			b2EmelementBody.SetUserData(barrier)
			client.Barrier[int32(index)] = barrier
		}
	}
}

func ErrFatal(err error) {
	if err != nil {
		log.Fatal("ErrFatal", zap.NamedError("err", err))
	}
}
