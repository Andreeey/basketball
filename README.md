# Basketball fixture app

## Requirements

You need to have golang installed snd configured to be able to run this app. This app also uses make commands so unix-like systems is preferred.

## How to run application

Cloning repo:

```bash
git clone git@github.com:Andreeey/basketball.git
```

You need to get vendors before running app. Run this command:

```bash
make vendor
```

Now you are ready to run app:

```bash
make run
```

When app is running you should be able to open http://127.0.0.1:8080/ and start creating games

# Application operation

You can create set of games with `Start next week games` button. This will create 7 games at BE which will start immediately. 
FE will subscribe to WS to listen for top players in real time. Every 5 seconds games will be updated with actual scores.
To start new games click `Start next week games` again. It will stop all current games and create new ones, thus `Simulation for upcoming weeks` bonus requirement is covered

Application has 2 REST endpoints to create / list games and one web socket to keep track on top players.
Substitution logic is also implemented, you can see players moving between `Players` and `Bench`. You can track substitution taking place in app logs as well.
All bonus cases in task were covered.

## Simplifications and notes

### WS

App uses web sockets to show real time players. Current implementation works only with 1 client. This issue is known and fixing it is out of scope of this test task

### Storage

App uses no storage. All data is kept in memory. This approach was taken by author because there is no use-cases for storage being used in test task. 
Also adding storage will increase complexity and time.

### API testing

You can find Postman collection under `assets` directory for test purposes

Please feel free to contact me if you have any questions on test task Andrii Bobryshev bobryshev.andrey@gmail.com