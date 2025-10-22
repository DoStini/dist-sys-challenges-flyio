cd cmd/broadcast
go build .

cd ../../maelstrom
./maelstrom test -w broadcast --bin /Users/andremoreira9/Documents/personal/learning/flyio/broadcast/broadcast --node-count 1 --time-limit 20 --rate 10
