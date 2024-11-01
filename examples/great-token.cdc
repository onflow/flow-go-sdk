access(all) contract GreatToken {

  access(all) resource interface NFT {
    access(all) fun id(): Int {
      post {
        result > 0
      }
    }
  }

  access(all) resource GreatNFT: NFT {
    access(self) let _id: Int
    access(self) let _special: Bool

    access(all) fun id(): Int {
      return self._id
    }

    access(all) fun isSpecial(): Bool {
      return self._special
    }

    init(id: Int, isSpecial: Bool) {
      pre {
        id > 0
      }
      self._id = id
      self._special = isSpecial
    }
  }

  access(all) resource GreatNFTMinter {
    access(all) var nextID: Int
    access(all) let specialMod: Int

    access(all) fun mint(): @GreatNFT {
      var isSpecial = self.nextID % self.specialMod == 0
      let nft <- create GreatNFT(id: self.nextID, isSpecial: isSpecial)
      self.nextID = self.nextID + 1
      return <-nft
    }

    init(firstID: Int, specialMod: Int) {
      pre {
        firstID > 0
        specialMod > 1
      }
      self.nextID = firstID
      self.specialMod = specialMod
    }
  }

  access(all) fun createGreatNFTMinter(firstID: Int, specialMod: Int): @GreatNFTMinter {
    return <-create GreatNFTMinter(firstID: firstID, specialMod: specialMod)
  }
}
