### Initializaiton

Make sure there is no legacy docker service running. 

```
docker stop zenrock-mainnet
```




```
docker volume ls
````



Run the init service to generate the folder structure and download the required binaries

```
docker-compose up -d init-mainnet
```

it should complete in about 30-60 seconds, wait for the process to complete before moving to the second compose command, the container should exit afterwards


### Start the zenrockd mainnet service

Once the initialization is completed you can start the zenrockd service:

```
docker-compose up -d zenrock-mainnet
```
