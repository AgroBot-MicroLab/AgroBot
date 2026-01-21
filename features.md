# Rover Telemetry (90%)
## Implemented
- Rover position is received and displayed in real time
- Heartbeat messages are handled correctly to prevent UDP connection timeouts

## Missing
- Display of battery level
- Display of temperature readings

# Mission Planning (100%)
## Implemented
- Missions can be created and stored in the database
- Missions are successfully uploaded to the rover via MAVLink
- Missions can be started remotely on the rover

# Image Capture & Display (75%)
## Implemented
- Rover receives MQTT commands
- Rover captures an image on MQTT trigger
- Images are transmitted back to the server
- Images are stored in a gallery and can be viewed

## Missing
- Gimbal control to set camera orientation

# Mission Preview & Replay (100%)
## Implemented
- Missions can be saved with title and description
- Missions can be previewed on the map
- Missions can be replayed

# Integration Layer (50%)
## Implemented
- APIs available for fetching:
  - images
  - missions
  - rover positions

## Missing
- Abstract device interface (swap camera drivers without changing code)
- Generic device abstraction (support devices beyond cameras)

# System Robustness (75%)
## Implemented
- Rover position stored using WAL-style persistence
- Rover continues mission execution during network outages

## Missing
- Encrypted transport (DTLS or QUIC over UDP)
- Automatic reconnection logic when connection drops
