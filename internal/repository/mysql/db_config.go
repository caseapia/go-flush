package mysql

/*
	The maximum number of log columns that
	will be returned to the user. Please
	do not specify more than 500 columns here.
*/
var LOGS_COLUMNS_LIMIT = 100

/*
	The maximum number of columns that will be returned to the user.
	This setting includes all data, but not logs and users. Please
	do not set more than 800 columns here.
*/
var COLUMNS_LIMIT = 30

/*
	Maximum number of user columns. Please do not set more than
	500 columns here.
*/
var USER_COLUMNS_LIMIT = 500
