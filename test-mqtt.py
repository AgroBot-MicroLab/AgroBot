import argparse
import time
import uuid
import paho.mqtt.client as mqtt

def run_sub(host, port, topic, qos):
    cid = f"py-sub-{uuid.uuid4().hex[:8]}"
    c = mqtt.Client(client_id=cid, clean_session=True)

    def on_connect(client, userdata, flags, rc):
        print(f"[sub] connected rc={rc}")
        client.subscribe(topic, qos)

    def on_message(client, userdata, msg):
        try:
            payload = msg.payload.decode("utf-8")
        except Exception:
            payload = msg.payload
        print(f"[sub] {msg.topic} qos={msg.qos}: {payload}")

    c.on_connect = on_connect
    c.on_message = on_message
    c.connect(host, port, keepalive=30)
    c.loop_forever()

def run_pub(host, port, topic, qos, payload, retain):
    cid = f"py-pub-{uuid.uuid4().hex[:8]}"
    c = mqtt.Client(client_id=cid, clean_session=True)
    c.connect(host, port, keepalive=30)
    c.loop_start()
    info = c.publish(topic, payload, qos=qos, retain=retain)
    info.wait_for_publish()
    print(f"[pub] sent mid={info.mid} topic={topic} qos={qos} retain={retain}")
    c.loop_stop()
    c.disconnect()

def main():
    p = argparse.ArgumentParser()
    p.add_argument("mode", choices=["sub", "pub"])
    p.add_argument("--host", default="broker.emqx.io")
    p.add_argument("--port", type=int, default=1883)
    p.add_argument("--topic", default="drone/cmd")
    p.add_argument("--qos", type=int, default=1, choices=[0, 1, 2])
    p.add_argument("--payload", default="make_photo")
    p.add_argument("--retain", action="store_true")
    args = p.parse_args()

    if args.mode == "sub":
        run_sub(args.host, args.port, args.topic, args.qos)
    else:
        run_pub(args.host, args.port, args.topic, args.qos, args.payload, args.retain)

if __name__ == "__main__":
    main()

