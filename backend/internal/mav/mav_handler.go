package mav

import (
    "context"
    "errors"
    "log"
    "sync"
    "time"
    "fmt"

    "agro-bot/internal/ws"

    gomavlib "github.com/bluenviron/gomavlib/v2"
    "github.com/bluenviron/gomavlib/v2/pkg/dialects/ardupilotmega"
    "github.com/bluenviron/gomavlib/v2/pkg/dialects/common"
)

var posMask = common.POSITION_TARGET_TYPEMASK_VX_IGNORE |
    common.POSITION_TARGET_TYPEMASK_VY_IGNORE |
    common.POSITION_TARGET_TYPEMASK_VZ_IGNORE |
    common.POSITION_TARGET_TYPEMASK_AX_IGNORE |
    common.POSITION_TARGET_TYPEMASK_AY_IGNORE |
    common.POSITION_TARGET_TYPEMASK_AZ_IGNORE |
    common.POSITION_TARGET_TYPEMASK_YAW_IGNORE |
    common.POSITION_TARGET_TYPEMASK_YAW_RATE_IGNORE

type Waypoint struct {
    Lat float64
    Lon float64
}

type Options struct {
    UDPAddr          string
    OutSystemID      byte
    OutComponentID   byte
    TargetSystem     uint8
    TargetComponent  uint8
}

type Client struct {
    mu     sync.RWMutex
    node   *gomavlib.Node
    opts   Options
    ctx    context.Context
    cancel context.CancelFunc
    wg     sync.WaitGroup

    missionCh chan any
    cmdAckCh  chan *common.MessageCommandAck
    OnMissionReached func(seq uint16)
    OnPos func(ws.Pos)
    OnErr func(error)
}

func New(opts Options) (*Client, error) {
    if opts.UDPAddr == "" {
        return nil, errors.New("UDPAddr required")
    }
    if opts.OutSystemID == 0 {
        opts.OutSystemID = 255
    }
    if opts.OutComponentID == 0 {
        opts.OutComponentID = 190
    }
    if opts.TargetSystem == 0 {
        opts.TargetSystem = 1
    }
    if opts.TargetComponent == 0 {
        opts.TargetComponent = 1
    }

    node, err := gomavlib.NewNode(gomavlib.NodeConf{
        Endpoints:      []gomavlib.EndpointConf{gomavlib.EndpointUDPServer{Address: opts.UDPAddr}},
        Dialect:        ardupilotmega.Dialect,
        OutVersion:     gomavlib.V2,
        OutSystemID:    opts.OutSystemID,
        OutComponentID: opts.OutComponentID,
    })
    if err != nil {
        return nil, err
    }


    ctx, cancel := context.WithCancel(context.Background())
    c := &Client{
        node:   node,
        opts:   opts,
        ctx:    ctx,
        cancel: cancel,
    }
    c.wg.Add(1)
    c.missionCh = make(chan any, 64)
    c.cmdAckCh = make(chan *common.MessageCommandAck, 32)

    go c.readLoop()

    return c, nil
}

func (c *Client) readLoop() {
    defer c.wg.Done()
    defer c.node.Close()

    for evt := range c.node.Events() {
        f, ok := evt.(*gomavlib.EventFrame);
        if (!ok) { continue }

        switch m := f.Message().(type) {
            case *ardupilotmega.MessageGlobalPositionInt:
                if c.OnPos != nil {
                    c.OnPos(ws.Pos{ Lat: float64(m.Lat) / 1e7, Lon: float64(m.Lon) / 1e7, })
                }

            case *common.MessageHeartbeat:
                sysID := f.SystemID()
                c.mu.Lock()
                c.opts.TargetSystem = sysID
                c.mu.Unlock()

            case *common.MessageMissionRequestInt,
                 *common.MessageMissionRequest,
                 *common.MessageMissionAck,
                 *common.MessageMissionCount:

                log.Printf("ACK: %T %+v", m, m)
                select {
                    case c.missionCh <- m:
                    default:
                        log.Printf("missionCh full, drop %T", m)
                }

            case *common.MessageMissionItemReached:
                if c.OnMissionReached != nil { c.OnMissionReached(m.Seq) }

            case *common.MessageCommandAck:
                select { case c.cmdAckCh <- m: default: }
            default:
                //log.Printf("Unhandled MAVLink message: %T %+v", m, m)
}
    }
}

func (c *Client) Close() {
    c.cancel()
    done := make(chan struct{})
    go func() {
        c.wg.Wait()
        close(done)
    }()
    select {
    case <-done:
    case <-time.After(2 * time.Second):
        log.Printf("mav: forced shutdown")
    }
}

func (c *Client) SendGoto(lat, lon float64) error {
    c.mu.RLock()
    n := c.node
    opts := c.opts
    c.mu.RUnlock()
    if n == nil {
        return errors.New("node not initialized")
    }

    const typeMaskPositionOnly uint16 = 4088 // ignore v*, a*, yaw, yaw_rate

    msg := &common.MessageSetPositionTargetGlobalInt{
        TimeBootMs:  0,
        TargetSystem:    opts.TargetSystem,
        TargetComponent: opts.TargetComponent,
        CoordinateFrame: common.MAV_FRAME_GLOBAL_RELATIVE_ALT_INT,
        TypeMask:        posMask,
        LatInt:          int32(lat * 1e7),
        LonInt:          int32(lon * 1e7),
        Alt:             0,
    }
    return n.WriteMessageAll(msg)
}

func (c *Client) RunHardcodedMission(ctx context.Context) (error) {
    c.mu.RLock()
    n := c.node
    sys := c.opts.TargetSystem
    c.mu.RUnlock()
    if n == nil { return errors.New("node not initialized") }

    wps := []Waypoint{
        {Lat: -35.36214764686344, Lon: 149.1651090448245},
        {Lat: -35.36214764686344, Lon: 149.1661090448245},
        {Lat: -35.36264000000000, Lon: 149.1666000000000},
        {Lat: -35.36264000000000, Lon: 149.1670000000000},
    }
    if len(wps) == 0 { return errors.New("no waypoints") }

    items := make([]*common.MessageMissionItemInt, len(wps))
    for i, p := range wps {
        items[i] = &common.MessageMissionItemInt{
            TargetSystem:    sys,
            TargetComponent: 1,
            Seq:             uint16(i),
            Frame:           common.MAV_FRAME_GLOBAL_RELATIVE_ALT_INT,
            Command:         common.MAV_CMD_NAV_WAYPOINT,
            Current:         0,
            Autocontinue:    1,
            X:               int32(p.Lat * 1e7),
            Y:               int32(p.Lon * 1e7),
            Z:               0,
            MissionType:     common.MAV_MISSION_TYPE_MISSION,
        }
    }

    items[0].Current = 1

    if err := n.WriteMessageAll(&common.MessageMissionClearAll{
        TargetSystem: sys, TargetComponent: 1,
        MissionType:  common.MAV_MISSION_TYPE_MISSION,
    }); err != nil { return err }

    if err := n.WriteMessageAll(&common.MessageMissionCount{
        TargetSystem: sys,
        TargetComponent: 1,
        Count:        uint16(len(items)),
        MissionType:  common.MAV_MISSION_TYPE_MISSION,
    }); err != nil { return err }

    deadline := time.NewTimer(10 * time.Second)
    defer deadline.Stop()

    sent := make([]bool, len(items))
    count := 0

    for ev := range c.missionCh {
        switch m := ev.(type) {

        case *common.MessageMissionRequest:
            s := int(m.Seq)

            log.Println("Sent point: ", s);
            if s < 0 || s >= len(items) {
                return fmt.Errorf("bad seq %d", s)
            }
            if err := n.WriteMessageAll(items[s]); err != nil {
                return err
            }

            if !sent[s] {
                sent[s] = true
                count++
                if count == len(items) {
                    return nil
                }
            }
        }
    }

    return nil
}

func (c *Client) StartMission(ctx context.Context) error {
    c.mu.RLock()
    n := c.node
    sys := c.opts.TargetSystem
    c.mu.RUnlock()
    if n == nil {
        return errors.New("node not initialized")
    }

    // 1) Set current mission index = 0
    if err := n.WriteMessageAll(&common.MessageMissionSetCurrent{
        TargetSystem:    sys,
        TargetComponent: 1, // MAV_COMP_ID_AUTOPILOT1
        Seq:             0,
    }); err != nil {
        return err
    }

    // helper to wait for COMMAND_ACK for a specific command
    waitAck := func(cmd common.MAV_CMD, d time.Duration) error {
        t := time.NewTimer(d)
        defer t.Stop()
        for {
            select {
            case ack := <-c.cmdAckCh:
                if ack.Command != cmd {
                    continue
                }
                if ack.Result == common.MAV_RESULT_ACCEPTED {
                    return nil
                }
                return fmt.Errorf("ack %v: %v", cmd, ack.Result)
            case <-t.C:
                return fmt.Errorf("ack timeout: %v", cmd)
            case <-ctx.Done():
                return ctx.Err()
            }
        }
    }

    // 2) Arm
    if err := n.WriteMessageAll(&common.MessageCommandLong{
        TargetSystem:    sys,
        TargetComponent: 1,
        Command:         common.MAV_CMD_COMPONENT_ARM_DISARM,
        Confirmation:    0,
        Param1:          1, // arm
    }); err != nil {
        return err
    }
    if err := waitAck(common.MAV_CMD_COMPONENT_ARM_DISARM, 5*time.Second); err != nil {
        return err
    }

    // 3) Switch mode to AUTO (ArduPilot Rover AUTO = custom mode 10)
    if err := n.WriteMessageAll(&common.MessageSetMode{
        TargetSystem: sys,
        BaseMode:     common.MAV_MODE(common.MAV_MODE_FLAG_CUSTOM_MODE_ENABLED | common.MAV_MODE_FLAG_AUTO_ENABLED),
        CustomMode:   10, // Rover: AUTO
    }); err != nil {
        return err
    }

    // Optional: explicitly start mission (many stacks start automatically in AUTO)
    // Use COMMAND_LONG for an ACKable start signal.
    if err := n.WriteMessageAll(&common.MessageCommandLong{
        TargetSystem:    sys,
        TargetComponent: 1,
        Command:         common.MAV_CMD_MISSION_START,
        Confirmation:    0,
        Param1:          0, // first item index
        Param2:          0, // last item (0 means "use planner default"; Rover ignores it)
    }); err != nil {
        return err
    }
    if err := waitAck(common.MAV_CMD_MISSION_START, 5*time.Second); err != nil {
        return err
    }

    return nil
}



