# November 1, 2025 (v0.1)

### What’s New?

We’re introducing updates to improve the Agricultural Rover Mission Planner experience and prepare for the next major iteration.

### Mission Planner Enhancements

We’re redesigning the mission planner interface to make it more intuitive and responsive. This includes clearer mission visualization, streamlined control layouts, and improved feedback when interacting with the rover.

### Mission Persistence and Replay

You’ll soon be able to save, load, and replay rover missions directly from the database. This feature will make it easier to test, analyze, and repeat mission runs without manual reconfiguration.

### WebSocket Performance Improvements

The current WebSocket implementation uses a global lock on a shared HashMap for managing active connections. We’re refactoring this to use a non-blocking or sharded connection model, greatly improving scalability and reducing latency under multiple simultaneous connections.

### Image Management

A new Image Gallery page is coming soon. All captured images will be displayed chronologically with timestamps and mission context, providing a single place to review and analyze collected data.

### Why the Change?

These updates aim to improve system reliability, scalability, and overall user experience while laying the groundwork for more advanced mission planning features in upcoming releases.

### What’s Next?

In upcoming versions, expect:

- A fully optimized WebSocket backend
- Persistent mission history and playback controls
- Modernized, responsive UI components
- Enhanced media management and filtering options

