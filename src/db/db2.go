package config

import (
	"database/sql"
)

var DB2 *sql.DB

const SqlGEN5 = `select count(R.ResultID) as prod, 'GEN5' as model from PalmelaPlant_GEN5_MES.dbo.Results R

where StationID = '9000' and Failed = 0 and R.Created >= @dataInicial and R.Created <= @dataFinal`

const ProductionGEN5 = `CREATE TABLE #prod
(
  Hour int, Timestamp datetime, A_Total_OK int, A_Total_NOK int, B_Total_OK int, B_Total_NOK int, C1_Total_OK int, C1_Total_NOK int, C2_Total_OK int, C2_Total_NOK int, D_Total_OK int, D_Total_NOK int
)

insert into #prod EXEC [dbo].[sp_HANON_CalculateProductionStatistics] @runMode = 4
select Hour as hora, Timestamp, A_Total_OK, A_Total_NOK, B_Total_OK, B_Total_NOK, C1_Total_OK, C1_Total_NOK,C2_Total_OK,C2_Total_NOK,D_Total_OK,D_Total_NOK from #prod
drop table #prod`
