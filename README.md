go_udp_svrkit
=============

to create a service communicated with others, generally focus on 3 points belows:<br>
1. the data struct, like input/output<br>
2. protocal defined by remote services<br>
3. the main logic of the service to create, generally procedure-oriented<br>
<br>
go_udp_svrkit is a udp framework separate what mentioned above into corresponded 3 part memtioned above.<br>
<br>
see detail in example/svr<br>
or try a bench test: <br>
go test -test.bench=".*"<br>
