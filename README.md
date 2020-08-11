# OracleDump script
Simple Golang script that can be easily compiled for any operating system
The following article provides some information on cross compiling for
Windows should the client be using windows --> https://github.com/golang/go/wiki/WindowsCrossCompiling

Instructions are simple. The script needs to be ran with with the connection string provided as an argument.
Otherwise the script will attempt to connect to the database listener using the following credentials on the localhost
username: SYS AS SYSDBA
password: oracle
