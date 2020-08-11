package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/fatih/color"
	_ "github.com/godror/godror"
	"github.com/k0kubun/pp"
)

//var f *os.File = os.Stdout

//Currently there seems to be an issue with getting the output
//of this script into a file; however, it works well on the terminal
//provided your buffer lets you scroll that far back.
//The ouput is colour coded so that it is easy to read and can redirected
//into a file without the same errors I get when manually using the file stream
//perhaps a file descriptor would prove to be more effective
func main() {
	//create file
	f, err := os.OpenFile("oracleDB.audit", os.O_WRONLY|os.O_CREATE, 0644)
	//gracefulCrash(err, true)
	defer f.Close()
	//println(len(os.Args))
	//TODO let the script accept a connection string provided
	//as an argument
	var conn string
	if len(os.Args) == 2 {
		conn = os.Args[1]
	} else {
		conn = "sys/oracle@localhost AS SYSDBA"
	}
	db, err := sql.Open("godror", conn)
	gracefulCrash(err, true)
	defer db.Close()

	fmt.Fprintf(f, "Ensure All Default Passwords Are Changed\n")
	rows, err := db.Query("SELECT USERNAME FROM DBA_USERS_WITH_DEFPWD WHERE USERNAME NOT LIKE '%XS$NULL%'")
	gracefulCrash(err, true)
	prettyPrintQuery(rows, f)
	fmt.Fprintf(f, "Ensure All Sample Data And Users Have Been Removed\n")
	rows, err = db.Query("SELECT USERNAME FROM ALL_USERS WHERE USERNAME IN ('BI','HR','IX','OE','PM','SCOTT','SH')")
	gracefulCrash(err, true)
	prettyPrintQuery(rows, f)
	fmt.Fprintf(f, "Ensure 'AUDIT_SYS_OPERATIONS' Is Set to 'TRUE'\n")
	rows, err = db.Query("SELECT UPPER(VALUE) FROM V$PARAMETER WHERE UPPER(NAME) = 'AUDIT_SYS_OPERATIONS'")
	gracefulCrash(err, true)
	prettyPrintQuery(rows, f)
	fmt.Fprintf(f, "Ensure 'AUDIT_TRAIL' Is Set\n")
	rows, err = db.Query("SELECT UPPER(VALUE) FROM V$PARAMETER WHERE UPPER(NAME)='AUDIT_TRAIL'")
	gracefulCrash(err, true)
	prettyPrintQuery(rows, f)
	fmt.Fprintf(f, "Ensure 'GLOBAL_NAMES' Is Set to 'TRUE'\n")
	rows, err = db.Query("SELECT UPPER(VALUE) FROM V$PARAMETER WHERE UPPER(NAME)='GLOBAL_NAMES'")
	gracefulCrash(err, true)
	prettyPrintQuery(rows, f)
	fmt.Fprintf(f, "Ensure 'LOCAL_LISTENER' Is Set\n")
	rows, err = db.Query("SELECT UPPER(VALUE) FROM V$PARAMETER WHERE UPPER(NAME)='LOCAL_LISTENER'")
	gracefulCrash(err, true)
	prettyPrintQuery(rows, f)
	fmt.Fprintf(f, "Ensure 'OS_ROLES' Is Set to 'FALSE'\n")
	rows, err = db.Query("SELECT UPPER(VALUE) FROM V$PARAMETER WHERE UPPER(NAME)='OS_ROLES'")
	gracefulCrash(err, true)
	prettyPrintQuery(rows, f)
	fmt.Fprintf(f, "Ensure 'REMOTE_OS_ROLES' Is Set to 'FALSE'\n")
	rows, err = db.Query("SELECT UPPER(VALUE) FROM V$PARAMETER WHERE UPPER(NAME)='REMOTE_OS_ROLES'")
	gracefulCrash(err, true)
	prettyPrintQuery(rows, f)
	fmt.Fprintf(f, "Ensure 'SEC_CASE_SENSITIVE_LOGON' Is Set to 'TRUE'\n")
	rows, err = db.Query("SELECT UPPER(VALUE) FROM V$PARAMETER WHERE UPPER(NAME)='SEC_CASE_SENSITIVE_LOGON'")
	gracefulCrash(err, true)
	prettyPrintQuery(rows, f)
	fmt.Fprintf(f, "Ensure 'FAILED_LOGIN_ATTEMPTS' Is Less than or Equal to '5'\n")
	rows, err = db.Query("SELECT PROFILE, RESOURCE_NAME, LIMIT FROM DBA_PROFILES WHERE RESOURCE_NAME='FAILED_LOGIN_ATTEMPTS' AND (LIMIT = 'DEFAULT' OR LIMIT = 'UNLIMITED' OR LIMIT > 5)")
	gracefulCrash(err, true)
	prettyPrintQuery(rows, f)
	fmt.Fprintf(f, "Ensure 'PASSWORD_LIFE_TIME' Is Less than or Equal to '90'\n")
	rows, err = db.Query("SELECT PROFILE, RESOURCE_NAME, LIMIT FROM DBA_PROFILES WHERE RESOURCE_NAME='PASSWORD_LIFE_TIME' AND (LIMIT = 'DEFAULT' OR LIMIT = 'UNLIMITED' OR LIMIT > 90)")
	gracefulCrash(err, true)
	prettyPrintQuery(rows, f)
	fmt.Fprintf(f, "Ensure 'SESSIONS_PER_USER' Is Less than or Equal to '10'\n")
	rows, err = db.Query("SELECT PROFILE, RESOURCE_NAME, LIMIT FROM DBA_PROFILES WHERE RESOURCE_NAME='SESSIONS_PER_USER' AND (LIMIT = 'DEFAULT' OR LIMIT = 'UNLIMITED' OR LIMIT > 10)")
	gracefulCrash(err, true)
	prettyPrintQuery(rows, f)
	fmt.Fprintf(f, "Ensure 'EXECUTE' Is Revoked from 'PUBLIC' on 'DBMS_ADVISOR'\n")
	rows, err = db.Query("SELECT PRIVILEGE FROM DBA_TAB_PRIVS WHERE GRANTEE='PUBLIC' AND PRIVILEGE='EXECUTE' AND TABLE_NAME='DBMS_ADVISOR'")
	gracefulCrash(err, true)
	prettyPrintQuery(rows, f)
	fmt.Fprintf(f, "Ensure 'EXECUTE' Is Revoked from 'PUBLIC' on 'DBMS_JAVA'\n")
	rows, err = db.Query("SELECT PRIVILEGE FROM DBA_TAB_PRIVS WHERE GRANTEE='PUBLIC' AND PRIVILEGE='EXECUTE' AND TABLE_NAME='DBMS_JAVA'")
	gracefulCrash(err, true)
	prettyPrintQuery(rows, f)
	fmt.Fprintf(f, "Ensure 'EXECUTE' Is Revoked from 'PUBLIC' on 'DBMS_JAVA_TEST'\n")
	rows, err = db.Query("SELECT PRIVILEGE FROM DBA_TAB_PRIVS WHERE GRANTEE='PUBLIC' AND PRIVILEGE='EXECUTE' AND TABLE_NAME='DBMS_JAVA_TEST'")
	gracefulCrash(err, true)
	prettyPrintQuery(rows, f)
	fmt.Fprintf(f, "Ensure 'EXECUTE' Is Revoked from 'PUBLIC' on 'DBMS_BACKUP_RESTORE'\n")
	rows, err = db.Query("SELECT PRIVILEGE FROM DBA_TAB_PRIVS WHERE GRANTEE='PUBLIC' AND PRIVILEGE='EXECUTE' AND TABLE_NAME='DBMS_BACKUP_RESTORE'")
	gracefulCrash(err, true)
	prettyPrintQuery(rows, f)
	fmt.Fprintf(f, "Ensure 'DBA' Is Revoked from Unauthorized 'GRANTEE'\n")
	rows, err = db.Query("SELECT GRANTEE, GRANTED_ROLE FROM DBA_ROLE_PRIVS WHERE GRANTED_ROLE='DBA' AND GRANTEE NOT IN ('SYS','SYSTEM')")
	gracefulCrash(err, true)
	prettyPrintQuery(rows, f)
	fmt.Fprintf(f, "Enable 'USER' Audit Option\n")
	rows, err = db.Query("SELECT AUDIT_OPTION, SUCCESS, FAILURE FROM DBA_STMT_AUDIT_OPTS WHERE AUDIT_OPTION='USER' AND USER_NAME IS NULL AND PROXY_NAME IS NULL AND SUCCESS = 'BY ACCESS' AND FAILURE = 'BY ACCESS'")
	gracefulCrash(err, true)
	prettyPrintQuery(rows, f)
	fmt.Fprintf(f, "Enable 'ALTER	USER' Audit Option\n")
	rows, err = db.Query("SELECT AUDIT_OPTION, SUCCESS, FAILURE FROM DBA_STMT_AUDIT_OPTS WHERE AUDIT_OPTION='ALTER USER' AND USER_NAME IS NULL AND PROXY_NAME IS NULL AND SUCCESS = 'BY ACCESS' AND FAILURE = 'BY ACCESS'")
	gracefulCrash(err, true)
	prettyPrintQuery(rows, f)
	fmt.Fprintf(f, "Enable 'DROP USER' Audit Option\n")
	rows, err = db.Query("SELECT AUDIT_OPTION, SUCCESS, FAILURE FROM DBA_STMT_AUDIT_OPTS WHERE AUDIT_OPTION='DROP USER' AND USER_NAME IS NULL AND PROXY_NAME IS NULL AND SUCCESS = 'BY ACCESS' AND FAILURE = 'BY ACCESS'")
	gracefulCrash(err, true)
	prettyPrintQuery(rows, f)
	str := color.RedString("Script Completed Execution\n")
	fmt.Fprintf(f, str)
}

func prettyPrintQuery(rows *sql.Rows, f *os.File) {
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
		fmt.Fprintf(f, "%s\n", str)
	}
}

func prettyPrintTableQuery(rows *sql.Rows, f *os.File) []string {
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
		fmt.Fprintf(f, "%s\n", str)
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
