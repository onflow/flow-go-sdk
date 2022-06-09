## RLP Encoding

This is a git subtree to the ethereum module found here: https://github.com/ethereum/go-ethereum/

Contents have been modified to only include the needed RLP encoding 
functionality according to the licence found here: https://github.com/ethereum/go-ethereum/blob/master/COPYING.LESSER

Copyright 2014 The go-ethereum Authors
This file is part of the go-ethereum library.

The go-ethereum library is free software: you can redistribute it and/or modify
it under the terms of the GNU Lesser General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

The go-ethereum library is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public License
along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

## Updating

You can use git subtree command to update from the upstream branch.

```
git subtree pull --prefix rlp https://github.com/ethereum/go-ethereum main --squash
```
