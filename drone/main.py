import argparse
import uuid
import paho.mqtt.client as mqtt
import os
import random
import requests
from pathlib import Path
import mimetypes
from picamera import make_photo

BASE_URL   = os.getenv("BASE_URL", "http://10.120.154.245:8080")
IMAGES_DIR = Path(os.getenv("IMAGES_DIR", "./"))
URL = f"{BASE_URL}/image/1"

def pick_image(p: Path) -> Path:
    imgs = [x for x in p.iterdir() if x.suffix.lower() in (".jpg", ".jpeg", ".png", ".webp", ".gif")]
    if not imgs:
        raise RuntimeError(f"No images in {p.resolve()}")
    return random.choice(imgs)

def send_image():
    img = pick_image(IMAGES_DIR)
    ctype = mimetypes.guess_type(img.name)[0] or "application/octet-stream"
    print(f"POST {URL} <- {img.name} ({ctype})")

    with open(img, "rb") as fh:
        files = {"image": (img.name, fh, ctype)}
        resp = requests.post(URL, files=files, timeout=60)

    print("Status:", resp.status_code)
    print("Response:", resp.text[:500])


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
        send_image()

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

main()

