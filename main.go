package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/fatih/color"
	_ "github.com/godror/godror"
	"github.com/k0kubun/pp"
)

var f *os.File = os.Stdout

//Currently there seems to be an issue with getting the output
//of this script into a file; however, it works well on the terminal
//provided your buffer lets you scroll that far back.
//The ouput is colour coded so that it is easy to read and can redirected
//into a file without the same errors I get when manually using the file stream
//perhaps a file descriptor would prove to be more effective
func main() {
	//create file
	//f, err := os.OpenFile("oracleDB.log", os.O_WRONLY|os.O_CREATE, 0644)
	//gracefulCrash(err, true)
	defer f.Close()
	//println(len(os.Args))
	//TODO let the script accept a connection string provided
	//as an argument
	var conn string
	if len(os.Args) == 2 {
		conn = os.Args[1]
	} else {
		conn = "sys/oracle@orclcdb AS SYSDBA"
	}
	db, err := sql.Open("godror", conn)
	gracefulCrash(err, true)
	defer db.Close()
	//assume that it was opened correctly
	/*userSysPrivsRows, err := db.Query("SELECT * FROM USER_SYS_PRIVS")
	gracefulCrash(err, true)
	fmt.Fprintf(f, "%s", color.HiRedString("User system privileges:\n"))
	fmt.Fprintf(f, "%s", color.GreenString("=====================================================>"))
	prettyPrintQuery(userSysPrivsRows)
	//databases on database
	//commented out for brevity, but could prove to provide valuable information
	/*
		databaseRows, err := db.Query("SELECT * FROM v$database")
		gracefulCrash(err)
		color.HiRed("\nAll databases:\n")
		color.Green("=====================================================>")
		prettyPrintQuery(databaseRows)
	*/
	//tables
	tables, err := db.Query("SELECT table_name FROM all_tables ORDER BY table_name")
	gracefulCrash(err, true)
	fmt.Fprintf(f, "%s", color.HiRedString("\nAll tables in the currently mounted database:\n"))
	fmt.Fprintf(f, "%s", color.GreenString("=====================================================>"))
	tableNames := prettyPrintTableQuery(tables)
	tableNamesStr := pp.Sprint(tableNames)
	fmt.Fprintf(f, "%s", tableNamesStr)
	fmt.Fprintf(f, "%s", color.RedString("\nEach tables contents"))
	fmt.Fprintf(f, "%s", color.GreenString("=====================================================>"))
	for _, name := range tableNames {
		query := fmt.Sprintf("SELECT * FROM %s FETCH NEXT %d ROWS ONLY", name, 5)
		table, err := db.Query(query)
		skip := gracefulCrash(err, false)
		if skip != false {
			continue
		}
		fmt.Fprintf(f, "%s", color.CyanString("=====================================================>"))
		tabName := color.RedString(name)
		fmt.Fprintf(f, "Contents of table %s\n", tabName)
		prettyPrintQuery(table)
	}
	color.Red("Script Completed Execution\n")
}

func prettyPrintQuery(rows *sql.Rows) {
	dbColumns, err := rows.Columns()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return
	}
	for rows.Next() {
		columns := make([]interface{}, len(dbColumns))
		columnPointers := make([]interface{}, len(dbColumns))
		for i, v := range columns {
			_ = v //this is cause omitting the second value from
			//range results in i being a boolean value
			columnPointers[i] = &columns[i]
		}
		//this will create a string -> interface "map" or "dictionary"
		//or the columns from the query

		//Scan result in column pointers "map/dictionary"
		if err := rows.Scan(columnPointers...); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(2)
		}

		//create map and retrieve the value of each column from the
		//columnPointers
		m := make(map[string]interface{}) //the actual key pair
		for i, colName := range dbColumns {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}
		fmt.Fprintf(f, "%s", color.CyanString("=====================================================>"))
		//Dump to output
		str := pp.Sprint(m)
		fmt.Fprintf(f, "%s", str)
	}
}

func prettyPrintTableQuery(rows *sql.Rows) []string {
	dbColumns, err := rows.Columns()
	var tables []string
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return nil
	}
	for rows.Next() {
		columns := make([]interface{}, len(dbColumns))
		columnPointers := make([]interface{}, len(dbColumns))
		for i, v := range columns {
			_ = v //this is cause omitting the second value from
			//range results in i being a boolean value
			columnPointers[i] = &columns[i]
		}
		//this will create a string -> interface "map" or "dictionary"
		//or the columns from the query

		//Scan result in column pointers "map/dictionary"
		if err := rows.Scan(columnPointers...); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(2)
		}

		//create map and retrieve the value of each column from the
		//columnPointers
		m := make(map[string]interface{}) //the actual key pair
		for i, colName := range dbColumns {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
			if valName, ok := (*val).(string); ok {
				tables = append(tables, valName)
			}
		}
		fmt.Fprintf(f, "%s", color.CyanString("=====================================================>"))
		//Dump to output
		str := pp.Sprint(m)
		fmt.Fprintf(f, "%s", str)
	}
	return tables
}

func gracefulCrash(err error, exit bool) bool {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		if exit {
			os.Exit(2)
		}
		return true
	}
	return false
}
