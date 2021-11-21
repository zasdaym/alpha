# Alpha (Server & Client)

AlphaServer and AlphaClient are tools to record SSH login attempts.

## Software design

The system consists of two binaries, AlphaServer and AlphaClient which communicates via HTTP. AlphaServer will serves an endpoint which can be called by multiple AlphaClients to "increment" the SSH log in attempt for each AlphaClient identified by client ID. AlphaClient detects a SSH log in attempt by watching the SSH server log file. AlphaServer use in-memory storage, which means every time it is restared the data will disappear.

## Deployment design

The system is deployed using Docker Compose to spawns multiple containers. The roles of each containers are described below:
- server: Runs the AlphaServer binary, provides web UI to see total SSH log in attempts for all clients.
- node-abc: Runs the AlphaServer client which watch log file shared with `node-abc-ssh` via bind mounted volume and sends data to `server` to be aggregated.
- node-abc-ssh: Runs the SSH server client shares log file with `node-abc`. An extra container is used because docker container are designed to run only single process. This pattern also known as "Sidecar pattern".
- node-xyz: Same like `node-abc`, but with different client ID.
- node-xyz-ssh: Same like `node-abc-ssh`, but shares log file with `node-xyz`.
