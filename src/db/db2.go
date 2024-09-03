package config

import (
	"database/sql"
)

var DB2 *sql.DB

const SqlGEN5 = `select count(Results.ResultID) as prod, ISNULL(MAX(PN),'') as model from PalmelaPlant_GEN5_MES.dbo.Results

    inner join ProcessesResults on Results.ResultID = ProcessesResults.ResultID
    inner join Processes on ProcessesResults.ProcessID = Processes.ProcessID
    inner join ProcessesModels on Processes.ProcessID = ProcessesModels.ProcessID
    inner join Models on ProcessesModels.ModelID = Models.ModelID

where StationID = '9000' and Failed = 0 and CustomerPN is not null and ResultTimeStamp >= @dataInicial and ResultTimeStamp <= @dataFinal`
