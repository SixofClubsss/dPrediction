# dPrediction
P2P predictions and sports bet built on Dero.

Written in Go and using [Fyne Toolkit](https://fyne.io/), dPrediction is built on Dero's private L1. Powered by [Gnomon](https://github.com/civilware/Gnomon) and [dReams](https://github.com/dReam-dApps/dReams), dPrediction allows users to facilitates P2P predictions and Sports bets through a combination of Dero services and smart contracts. dPrediction was entered into dArch Event 1 under the On-Chain Options and Futures category. Contracts and services are each individually ran and operated, no oracles or bridges are used in its operation. 

![goMod](https://img.shields.io/github/go-mod/go-version/SixofClubsss/dPrediction.svg)![goReport](https://goreportcard.com/badge/github.com/SixofClubsss/dPrediction)[![goDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://pkg.go.dev/github.com/SixofClubsss/dPrediction)

### Features
- Install and run your own dPredictions/dSports
- 8 current markets with API
- 9 current leagues with API
- Contract ranking system
- Cli service

dPrediction contracts can be set up and operated through [dReams](https://dreamdapps.io). 

![windowsOS](https://raw.githubusercontent.com/SixofClubsss/dreamdappsite/main/assets/os-windows-green.svg)![macOS](https://raw.githubusercontent.com/SixofClubsss/dreamdappsite/main/assets/os-macOS-green.svg)![linuxOS](https://raw.githubusercontent.com/SixofClubsss/dreamdappsite/main/assets/os-linux-green.svg)

### Build
Following these build instructions, you can build dPrediction as a *individual* dApp.
- Install latest [Go version](https://go.dev/doc/install)
- Install [Fyne](https://developer.fyne.io/started/) dependencies
- Clone repo and build with:
```
git clone https://github.com/SixofClubsss/dPrediction.git
cd dPrediction/cmd/dPrediction
go build .
./dPrediction
```

Owners can run the cli service for their contracts. Up to nine wallets can be used per contract. Once connected the service will automatically find and operate any associated contracts.

### Using service
- Install latest [Go version](https://go.dev/doc/install)
- Clone repo and build with:

```
git clone https://github.com/SixofClubsss/dPrediction.git
cd dPrediction/cmd/dService
go build .
```
- Options
```
  --daemon=<127.0.0.1:10102>     Set daemon rpc address to connect.
  --wallet=<127.0.0.1:10103>     Set wallet rpc address to connect.
  --login=<user:pass>     	 Wallet rpc user:pass for auth.
  --transfers=<false>            True/false value for enabling processing transfers to integrated address.
  --debug=<true>     		 True/false value for enabling terminal debug.
  --fastsync=<true>	         Gnomon option,  true/false value to define loading at chain height on start up.
  --num-parallel-blocks=<5>      Gnomon option,  defines the number of parallel blocks to index.
```

- On local daemon, with wallet running rpc server start the service using:
```
./dService --login=user:pass
```
- Add `--transfers=true` flag if facilitating through Dero transactions.

### Donations
- **Dero Address**: dero1qyr8yjnu6cl2c5yqkls0hmxe6rry77kn24nmc5fje6hm9jltyvdd5qq4hn5pn

![DeroDonations](https://raw.githubusercontent.com/SixofClubsss/dreamdappsite/main/assets/DeroDonations.jpg)

---

#### Licensing

dPrediction is free and open source.   
The source code is published under the [MIT](https://github.com/SixofClubsss/dPrediction/blob/main/LICENSE) License.   
Copyright © 2023 SixofClubs   