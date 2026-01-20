
# November 28, 2025 (v0.2)

### What’s New?

This update focuses on improving how missions are organized, viewed, and managed inside the Agricultural Rover Platform. With the introduction of mission statuses and custom mission metadata, planning and reviewing rover activities becomes more streamlined and user-friendly.


### Mission Status System

We’re introducing a new **mission status model** to provide clearer visibility into mission progress and lifecycle. Every mission is now automatically assigned one of three states:

* **Planned**: mission route is configured but not yet executed
* **In Progress**: the rover is currently performing the mission
* **Completed**: the rover has finished the route and all images have been captured

These statuses appear throughout the interface: in the mission list, mission details page, and mission history. This change helps users quickly understand the state of each mission and improves filtering, sorting, and general workflow clarity.


### Custom Mission Metadata

Missions now support **user-defined metadata**, allowing operators to better describe and catalog rover activities. When saving a mission, users can now specify:

* **Mission Name**: human-readable names such as *“North Field Morning Sweep”*
* **Description**: optional notes for context, such as *“Test for spectral camera focus”*
* **Date & Time**: automatically tracked timestamps for start and completion

This enriches how mission data is displayed across the platform and makes it easier to identify past missions without relying on raw IDs.



### Enhanced Mission List & Filters

With the addition of statuses and metadata, the mission list has been upgraded to support:

* Search by mission name
* Filtering by status (Planned, In Progress, Completed)
* Date-range filtering (From–To)
* Sorting by date, mission duration, and number of captured images

These tools give operators full control over large mission histories and significantly improve navigation within the platform.


### Why the Change?

As rover operations scale, mission history becomes more complex. Clear statuses, descriptive mission names, and rich metadata improve usability, reduce operator error, and enable more advanced mission planning and analytics in future releases.

### What’s Next?

Upcoming versions will continue to build on these enhancements with:

- Expanded mission editing tools
- Status-based automation and validation rules
- Deep mission analytics and performance metrics
- Smarter search and advanced filtering options

---

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

