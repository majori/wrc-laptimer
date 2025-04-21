# Changelog

## 0.0.1 (2025-04-21)


### Features

* add location ID to routes ([f881707](https://github.com/majori/wrc-laptimer/commit/f881707cbe476c127c548fd90f4e1ab0bbe132cf))
* allow querying with sql ([591309a](https://github.com/majori/wrc-laptimer/commit/591309a36a293c7d1c2cd760aaefe0395517804f))
* class+stage combo instead of vehicle+stage ([#18](https://github.com/majori/wrc-laptimer/issues/18)) ([51c4448](https://github.com/majori/wrc-laptimer/commit/51c4448494c3b606a7e3ec4a202a41dcafd73f08))
* create generic query endpoint ([5b53a7b](https://github.com/majori/wrc-laptimer/commit/5b53a7b7502a4c41252b8038e9e0d6557ae1b4b6))
* default driver name generator ([#20](https://github.com/majori/wrc-laptimer/issues/20)) ([cec38a8](https://github.com/majori/wrc-laptimer/commit/cec38a89699359aa68b0a3afa74c6b7664ba1d35))
* enhance leaderboard UI ([#15](https://github.com/majori/wrc-laptimer/issues/15)) ([3680778](https://github.com/majori/wrc-laptimer/commit/3680778a9b4ee26d1f6bd910f93ce897bfea79ae))
* event and series support ([#16](https://github.com/majori/wrc-laptimer/issues/16)) ([bfd9f12](https://github.com/majori/wrc-laptimer/commit/bfd9f122e779ff0e3b554aa6ade91be5b6dea747))
* handle system signals properly (ctrl-c) ([93d1f15](https://github.com/majori/wrc-laptimer/commit/93d1f15521d31bfc36129b12ea347e28cbbab357))
* implement point system ([772c350](https://github.com/majori/wrc-laptimer/commit/772c35001a580dcdcac24a24526d108a34040d66))
* link telemetry to a session ([a3e33ae](https://github.com/majori/wrc-laptimer/commit/a3e33ae08a1d2e59bf12327063546b8444ee1190))
* logout user if inactive ([0204fe6](https://github.com/majori/wrc-laptimer/commit/0204fe6b76dfe0c33ea729c8a9181408473e312b))
* save default telemetry packet to DuckDB ([0bbd7d9](https://github.com/majori/wrc-laptimer/commit/0bbd7d98d8f87051a72107d888182a22842085ed))
* serve static web files ([37612da](https://github.com/majori/wrc-laptimer/commit/37612da3b74c9d51dc96221b83d4cef93e15daac))
* split telemetry to separate channels ([3ae0582](https://github.com/majori/wrc-laptimer/commit/3ae05825cab53ab44858b6d085220f995ea4cdbb))
* store sessions to database ([a936bcd](https://github.com/majori/wrc-laptimer/commit/a936bcd155017993e4304d1a7ad82e0bfae91c88))
* **web:** initialize alpine project ([#2](https://github.com/majori/wrc-laptimer/issues/2)) ([ce59fd8](https://github.com/majori/wrc-laptimer/commit/ce59fd84f77eee9d8a5acf642eff17303fd162f9))


### Bug Fixes

* close db properly ([b85e7d0](https://github.com/majori/wrc-laptimer/commit/b85e7d02569bac83cd62c1a84ce526ddd2c43087))
* ignore shakedowns in sessions query ([e1510cf](https://github.com/majori/wrc-laptimer/commit/e1510cf85e6b4ba8a30fa480c92e749660324ae5))
* ignore telemetry if no active session ([c53a460](https://github.com/majori/wrc-laptimer/commit/c53a46023f337974c6fc9fb65bd3844a6a4c2cea))
* insert login if user exists ([8d0aa0f](https://github.com/majori/wrc-laptimer/commit/8d0aa0ffd594ef9a6527a7db5be9ca06f0e37825))
* LastInsertId() does not work on DuckDB ([b797596](https://github.com/majori/wrc-laptimer/commit/b7975969a62fc6fa1e5a8f0f9b03d1ee1b4ef086))
* logout previous users on login ([6f95edf](https://github.com/majori/wrc-laptimer/commit/6f95edfcb2f36197d03a17714e4b622facc8a113))
* remove foreign key constraint ([5561f50](https://github.com/majori/wrc-laptimer/commit/5561f50804b1c9090c38a1ff699ac57a16949aa8))
* set correct types for telemetry fields ([33163be](https://github.com/majori/wrc-laptimer/commit/33163be9a50fb38435c2110a5333dfd61e34d294))
* set session ID correctly ([e13316f](https://github.com/majori/wrc-laptimer/commit/e13316f0ae8fb4efd126bc26b181673ab09f4103))


### Continuous Integration

* try out building to different architectures ([4546cee](https://github.com/majori/wrc-laptimer/commit/4546ceef331c18a3b86588c495c86331d9f209e7))
