package config

const prodVS14 = "SET DATEFIRST 1; " +
	"SELECT MAX(DATEPART(hh,$dateTime)) AS hora, COUNT($ID) AS prod FROM $databaseInUse" +
	"WHERE $dateTime > '\".$dataInit.\" 00:00' \n\t\t\tAND $dateTime < '\".$dataInit.\" 16:30' " +
	"\n\t\t\t$paramModel $param \n\n\t\t\t" +
	"GROUP BY DATEPART(hh,$dateTime);"

const prodGEN3 = ""

const prodInv3 = ""

const prodYF = ""

const prodR744 = ""

const prodInv4 = ""

const prodInv4_2 = ""

const prodInv4_3 = ""
