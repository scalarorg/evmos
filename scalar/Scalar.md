1. Modify Evmos code base

   Add flowing line into the evmos/server/start.go

   ```
   ![alt text](server.png)
   ```

   evmos/server/flags/flags.go

   ```
   ![alt text](flags.png)
   // Tendermint/cosmos-sdk full-node start flags
   ```

2. Start dockers
   ```
   docker-compose -f docker-cluster-evmos.yaml up -d
   ```
