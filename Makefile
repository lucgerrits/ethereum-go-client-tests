PROGRAM_NAME = send-transaction

all: build run

build:
	go build -o $(PROGRAM_NAME)

run:
	./$(PROGRAM_NAME)

profile-cpu:
	go tool pprof http://localhost:6060/debug/pprof/profile

profile-heap:
	go tool pprof http://localhost:6060/debug/pprof/heap

clean:
	rm -f $(PROGRAM_NAME)
