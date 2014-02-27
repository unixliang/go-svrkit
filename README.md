go_udp_svrkit
=============

to create a service communicated with others, generally focus on 3 points belows:<br>
1. the data struct, like input/output
2. protocal defined by remote services
3. the main logic of the service to create, generally procedure-oriented

go_udp_svrkit is a udp framework separate what mentioned above into corresponded 3 part memtioned above.

see detail in example/svr
or try a bench test: 
go test -test.bench=".*"
