# dnsBench

## Summary
dnsBench is a dns server stree test tool

## usage
set your gopath and just run go build  
command options:  
  -n int  
    	run through domains N times (default 1)  
  -p int  
    	server port (default 53)  
  -path string  
    	file path for domain sets (default "./domains")  
  -s string  
    	serverAddress (default "127.0.0.1")  
  -t string  
    	dns type,support A/CNAME (default "A")  


