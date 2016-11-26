@ECHO off

:: Batch file for
:: Go Group Project at GMIT 2016
:: Date:    Nov 17, 2016
:: Adapted from: Murach's ADO.NET 4 with Visual Basic 2010
:: 
:: Uses the SQLCMD utility to run a SQL script that creates
:: the Employee database.

ECHO Attempting to create the Maintenance database . . . 
sqlcmd -S localhost\SQLExpress -E /i CreateGoGroupDB.sql
ECHO.
ECHO If no error message is shown, then the database was created successfully.
ECHO.
PAUSE