basic example:

first we want to create 2 binaries, 1 that contains no upgrade information or an empty upgrade handler for v1.0.0

```
const UpgradeName = "v1.0.0"

func (app ZenrockApp) RegisterUpgradeHandlers() {
	app.UpgradeKeeper.SetUpgradeHandler(
		UpgradeName,
		func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			return nil, nil
		},
	)
}
```


build the binary:

```
ignite chain build -o ./zenrock-100
```

now we need to init the initial data directory:

```
ignite chain init --home /tmp/node1-data   
```

next we need to build a binary that contains a handler for the upgrade name that will be requested through governance

```
const UpgradeName = "v1.0.1"

func (app ZenrockApp) RegisterUpgradeHandlers() {
	app.UpgradeKeeper.SetUpgradeHandler(
		UpgradeName,
		func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			return app.ModuleManager.RunMigrations(ctx, app.Configurator(), fromVM)
		},
	)
}
```

build the binary:

```
ignite chain build -o ./zenrock-101
```

for testing purposes we want to set the gov voting period to a low value, 
add the following lines to the genesis of the config.yml
```
    gov:
      params:
        voting_period: 0h1m0s
        expedited_voting_period: 0h0m30s
```


once both binaries and the data folder is set up we can continue with installing cosmovisor
more info at https://docs.cosmos.network/main/build/tooling/cosmovisor

```
go install cosmossdk.io/tools/cosmovisor/cmd/cosmovisor@latest
```

we need to configure at least the following env vars,
the reason cosmovisor needs to be in the node's home folder, is that it will check for the upgrade-info.json file,
which will be written in the data folder by the upgrade module, it will use the information to swap to the new binary

```
export DAEMON_HOME=/tmp/node1-data
export DAEMON_NAME=zenrockd
```

now we can init cosmovisor with the initial binary

```
cosmovisor init zenrock-100/zenrockd
```

after this command has run, a cosmovisor folder will be added to the data folder, 
which will contain the genesis binray we used for the init, you can verify with 
```
find /tmp/node1-data/cosmovisor
```

next we need to add the upgrade binary to cosmovisor, note that the upgrade name needs to be the same as specified in the binary

```
cosmovisor add-upgrade v1.0.1 ./zenrock-101/zenrockd
```

a new folder will be created in the cosmovisor folder which will contain the upgrade

once this is in place we can run cosmovisor

```
cosmovisor run start --home /tmp/node1-data
```

now that the chain is running we can propose and approve the software-upgrade through governance, in this case we will use block 100 to swap

```
export TX_INFO=(--from alice --yes --chain-id zenrock --gas-prices 0.0001urock)

zenrockd tx upgrade software-upgrade v1.0.1 \
	--title upgrade \
	--summary upgrade \
	--upgrade-height 100 \
	--upgrade-info "{}" \
	--no-validate \
	--deposit 10000000urock \
	$TX_INFO

zenrockd tx gov deposit 1 10000000urock $TX_INFO

zenrockd tx gov vote 1 yes $TX_INFO

```

to verify the voting went ok, check with the following command, it should show `yes_count` > 0
```
zenrockd q gov tally 1
```

when the voting period is finished, we can see the proposal has state PROPOSAL_STATUS_PASSED

```
zenrockd q gov proposals
```



when the chain reaches block 100, 
the chain will halt with the following error
```
11:15AM ERR UPGRADE "v1.0.1" NEEDED at height: 100: {} module=x/upgrade
11:15AM ERR error in proxyAppConn.FinalizeBlock err="UPGRADE \"v1.0.1\" NEEDED at height: 100: {}" module=state
11:15AM ERR CONSENSUS FAILURE!!!
```

ann an upgrade-info.json file will be create in the data folder, this file will be picked up by cosmovisor so it knows when to swap to the new binary
```
cat /tmp/node1-data/data/upgrade-info.json

{"name":"v1.0.1","time":"0001-01-01T00:00:00Z","height":100,"info":"{}"}
```

after which cosmovisor will terminate the old binary, 
swap the current binary in its folder to the upgrade specified in the json file,
and launch the new binary

upon startup it will take a backup of the current data folder in /tmp/node1-data/data-backup-<date>


```
11:15AM INF daemon shutting down in an attempt to restart module=cosmovisor
11:15AM INF starting to take backup of data directory backup start time=... module=cosmovisor
11:15AM INF backup completed backup completion time=... backup saved at=/tmp/node1-data/data-backup-<date> module=cosmovisor time taken to complete backup=19.209417
11:15AM INF pre-upgrade command does not exist. continuing the upgrade. module=cosmovisor
11:15AM INF upgrade detected, relaunching app=zenrockd module=cosmovisor
...
11:15AM INF applying upgrade "v1.0.1" at height: 100 module=x/upgrade
```

the upgrade process is now finished and the new binary is running
