Simple golang script to dump contents of an OracleDB instance.
Accepts a single argument which should be your connection string.
Will dump colour coded output to stdout and errors to stderr.
Your terminal buffer might(will) not be enough for the contennts so I suggest
redirecting output to a file.

TODO:
1. Fix file stream issues for the file itself
2. Find a neater and more presentable way to display contents
3. Resolve current errors you might get in a smarter way.