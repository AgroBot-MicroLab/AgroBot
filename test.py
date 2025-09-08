import time
from pymavlink import mavutil
from pymavlink.dialects.v20 import common as mavlink

m = mavutil.mavlink_connection('udp:0.0.0.0:14550', source_system=255, source_component=190)
m.wait_heartbeat()
tgt_sys, tgt_comp = m.target_system, m.target_component

def set_guided():
    mode_id = m.mode_mapping().get('GUIDED', 15)  # rover GUIDED = 15
    m.set_mode(mode_id)

def arm():
    m.mav.command_long_send(tgt_sys, tgt_comp,
                            mavutil.mavlink.MAV_CMD_COMPONENT_ARM_DISARM,
                            0, 1, 0, 0, 0, 0, 0, 0)

def goto(lat, lon):
    type_mask = 0b0000111111111000
    m.mav.set_position_target_global_int_send(
        0, tgt_sys, tgt_comp,
        mavutil.mavlink.MAV_FRAME_GLOBAL_RELATIVE_ALT_INT,
        type_mask,
        int(lat * 1e7), int(lon * 1e7), 0,
        0, 0, 0, 0, 0, 0, 0, 0)

set_guided()
arm()
goto(-35.36221177, 149.16507572)

last_hb = time.time()
while True:
    msg = m.recv_match(blocking=False)
    if msg and msg.get_type() == 'GLOBAL_POSITION_INT':
        gp = msg
        print(f"{gp.lat/1e7:.7f}, {gp.lon/1e7:.7f}")

    if time.time() - last_hb > 1.0:
        m.mav.heartbeat_send(mavlink.MAV_TYPE_GCS, mavlink.MAV_AUTOPILOT_INVALID, 0, 0, 0)
        last_hb = time.time()

    time.sleep(0.02)

