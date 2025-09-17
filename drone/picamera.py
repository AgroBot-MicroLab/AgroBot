from picamera2 import Picamera2
from libcamera import controls

cam_index = 0
picam2 = Picamera2(cam_index)
config = picam2.create_preview_configuration(
    main={"format": "BGR888", "size": (1296, 972)}
)

picam2.configure(config)
picam2.options["capture_timeout"] = 5000

def make_photo():
    picam2.start(show_preview=False)

    try:
        picam2.set_controls({"AeEnable": True})
    except Exception:
        pass

    picam2.capture_file("photo.jpg")
    picam2.stop()
    print("done")
