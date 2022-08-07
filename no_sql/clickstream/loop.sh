#!/bin/bash
  for i in {1..10000}
 do
   curl -X POST --location "http://127.0.0.1:80/event/send" \
       -H "Content-Type: application/json" \
       -d "{
             \"userID\": 2,
             \"origin\": \"https://avito.ru\",
             \"ipv4\": \"127.0.0.1\",
             \"agent\": 1420460528,
             \"session\": \"sdfdfwdf33451345345\",
             \"referrer\": \"https://google.com/?text=some\",
             \"country\": \"RU\",
             \"lang\": \"en\"
           }"
 done
 --notice here I want var i to be changing with every curl.