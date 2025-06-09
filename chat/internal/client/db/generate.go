package db

//go:generate sh -c "rm -rf mocks && mkdir mocks"
//go:generate minimock -i Transactor -o ./mocks/ -s "_minimock.go"
//go:generate minimock -i TxManager -o ./mocks/ -s "minimock.go"
//go:generate minimock -i Client -o ./mocks/ -s "minimock.go"
