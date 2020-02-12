```
  ____ ___  ____    
 / / /|_  |/ __/____
/_  _/ __// _ \/ __/
 /_//____/\___/\__/ 
                    
```
  
## Overview

426c is an end-to-end encrypted messenger written in Go.

## Why

An opportunity to learn and manufacture a tool that respects people's privacy. 

## Features

* E2E Encryption
* Encryption at rest
* No message logging or storage server side
* Re-attempt message which failed to send
* "Block" issuing and cost per message.

## Prerequisites

* Go v1.13 or later

## Build & Install

First you will need to clone this repository into `$GOPATH/src/github.com/syleron/426c` and execute the following command(s):


```
$ sudo make
...
```

## License

426c source code is available under the GPLv3 License which can be found in the LICENSE file.

Copyright (c) 2020 Andrew Zak <andrew@linux.com>