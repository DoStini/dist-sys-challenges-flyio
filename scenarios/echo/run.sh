cd cmd/echo
go build .

cd ../../maelstrom
./maelstrom test -w echo --bin /Users/andremoreira9/Documents/personal/learning/flyio/cmd/echo/echo --time-limit 10 --rate 100 --node-count 5 --availability total --nemesis partition
