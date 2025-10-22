go build .
cd ../maelstrom
./maelstrom test -w unique-ids --bin /Users/andremoreira9/Documents/personal/learning/flyio/unique-id/unique-id --time-limit 10 --rate 100 --node-count 5 --availability total --nemesis partition