/*
 * Flow Go SDK
 *
 * Copyright 2019-2020 Dapper Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package test

import (
	"math/rand"
	"time"
)

var GreetingScript = []byte(`
transaction(greeting: String) {
  execute { 
    log(greeting.concat(", World!")) 
  }
}
`)

const (
	greetingMsgAF = "Hi"
	greetingMsgSQ = "Hi"
	greetingMsgAM = "ሃይ"
	greetingMsgAR = "مرحبا"
	greetingMsgEU = "Hi"
	greetingMsgBN = "হাই"
	greetingMsgBS = "Zdravo"
	greetingMsgCA = "Hola"
	greetingMsgHR = "Bok"
	greetingMsgDA = "Hej"
	greetingMsgNL = "Hoi"
	greetingMsgET = "Tere"
	greetingMsgFI = "Moi"
	greetingMsgEN = "Hi"
	greetingMsgEL = "Γεια"
	greetingMsgHE = "היי"
	greetingMsgIS = "Hæ"
	greetingMsgID = "Hai"
	greetingMsgIT = "Ciao"
	greetingMsgJA = "こんにちは"
	greetingMsgKO = "안녕"
	greetingMsgLT = "Sveiki"
	greetingMsgMK = "Здраво"
	greetingMsgPT = "Oi"
	greetingMsgPA = "ਸਤ ਸ੍ਰੀ ਅਕਾਲ"
	greetingMsgSK = "Ahoj"
	greetingMsgSL = "Zdravo"
	greetingMsgSB = "Hej"
	greetingMsgTH = "สวัสดี"
	greetingMsgTR = "Merhaba"
	greetingMsgUR = "ہیلو"
	greetingMsgVI = "Chào"
	greetingMsgCY = "Hi"
	greetingMsgBG = "здрасти"
	greetingMsgRU = "Здравствуй"
	greetingMsgFR = "Bonjour"
	greetingMsgNB = "Hei"
	greetingMsgNN = "Hei"
	greetingMsgRO = "Bună"
	greetingMsgDE = "Hallo"
	greetingMsgGH = "Haigh"
	greetingMsgZH = "你好"
	greetingMsgPL = "Cześć"
	greetingMsgBE = "прывітанне"
	greetingMsgHI = "नमस्ते"
	greetingMsgHU = "Szia"
	greetingMsgUK = "Привіт"
	greetingMsgES = "Hola"
)

var greetings []string

func init() {
	greetings = []string{
		greetingMsgAF,
		greetingMsgSQ,
		greetingMsgAM,
		greetingMsgAR,
		greetingMsgEU,
		greetingMsgBN,
		greetingMsgBS,
		greetingMsgCA,
		greetingMsgHR,
		greetingMsgDA,
		greetingMsgNL,
		greetingMsgET,
		greetingMsgFI,
		greetingMsgEN,
		greetingMsgEL,
		greetingMsgHE,
		greetingMsgIS,
		greetingMsgID,
		greetingMsgIT,
		greetingMsgJA,
		greetingMsgKO,
		greetingMsgLT,
		greetingMsgMK,
		greetingMsgPT,
		greetingMsgPA,
		greetingMsgSK,
		greetingMsgSL,
		greetingMsgSB,
		greetingMsgTH,
		greetingMsgTR,
		greetingMsgUR,
		greetingMsgVI,
		greetingMsgCY,
		greetingMsgBG,
		greetingMsgRU,
		greetingMsgFR,
		greetingMsgNB,
		greetingMsgNN,
		greetingMsgRO,
		greetingMsgDE,
		greetingMsgGH,
		greetingMsgZH,
		greetingMsgPL,
		greetingMsgBE,
		greetingMsgHI,
		greetingMsgHU,
		greetingMsgUK,
		greetingMsgES,
	}

	rand.Seed(time.Now().Unix())
}

type Greetings struct {
	count int
}

func GreetingGenerator() *Greetings {
	return &Greetings{0}
}

func (g *Greetings) New() string {
	defer func() { g.count++ }()
	return greetings[g.count%len(greetings)]
}

func (g *Greetings) Random() string {
	return greetings[rand.Intn(len(greetings))]
}
